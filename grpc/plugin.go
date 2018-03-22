package grpc

import (
	"crypto/sha1"
	"fmt"
	"strings"

	envoyapi "github.com/envoyproxy/go-control-plane/envoy/api/v2"
	envoyroute "github.com/envoyproxy/go-control-plane/envoy/api/v2/route"
	"github.com/gogo/googleapis/google/api"
	"github.com/gogo/protobuf/proto"
	"github.com/gogo/protobuf/protoc-gen-gogo/descriptor"

	"github.com/gogo/protobuf/types"
	"github.com/pkg/errors"
	"github.com/solo-io/gloo-api/pkg/api/types/v1"
	"github.com/solo-io/gloo-plugins/transformation"
	"github.com/solo-io/gloo/pkg/log"
	"github.com/solo-io/gloo/pkg/plugin"
)

type Plugin struct {
	// map service names to their descriptors
	serviceDescriptors map[string]*descriptor.FileDescriptorSet
	// keep track of which service belongs to which upstream
	upstreamServices map[string]string
	transformation   *transformation.Plugin
}

const (
	filterName  = "envoy.grpc_json_transcoder"
	pluginStage = plugin.PreOutAuth

	ServiceTypeGRPC = "gRPC"
)

/*

for every request type, need to get all the top level fields
the field should come from either:
- body
- path
- query param

for each method:
  routes = aggregate all the routes for that method
  for each route:
    get the parameters (extractors).
    add an http rule to the destination function (method):

add_http_rule:
  look up the message type for the method's request
  for every top-level field in the message, we have to pick an extractor:
  - path
  - query param
  - body (currently just passthrough, nothing to do here)

*/

func (p *Plugin) GetDependencies(cfg *v1.Config) *plugin.Dependencies {
	deps := &plugin.Dependencies{}
	for _, us := range cfg.Upstreams {
		if !isOurs(us) {
			continue
		}
		serviceSpec, err := DecodeServiceProperties(us.ServiceInfo.Properties)
		if err != nil {
			log.Warnf("%v: error parsing service properties for upstream %v: %v",
				ServiceTypeGRPC, us.Name, err)
			continue
		}
		deps.FileRefs = append(deps.FileRefs, serviceSpec.DescriptorsFileRef)
	}
	return deps
}

func isOurs(in *v1.Upstream) bool {
	return in.ServiceInfo != nil && in.ServiceInfo.Type == ServiceTypeGRPC
}

func (p *Plugin) ProcessUpstream(params *plugin.UpstreamPluginParams, in *v1.Upstream, _ *envoyapi.Cluster) error {
	if !isOurs(in) {
		return nil
	}

	serviceProperties, err := DecodeServiceProperties(in.ServiceInfo.Properties)
	if err != nil {
		return errors.Wrap(err, "parsing service properties")
	}
	fileRef := serviceProperties.DescriptorsFileRef
	serviceName := serviceProperties.GRPCServiceName

	if fileRef == "" {
		return errors.New("service_properties.descriptors_file_ref cannot be empty")
	}
	if serviceName == "" {
		return errors.New("service_properties.service_name cannot be empty")
	}
	descriptorsFile := params.Files[fileRef]
	descriptors, err := convertProto(descriptorsFile)
	if err != nil {
		return errors.Wrapf(err, "parsing file %v as a proto descriptor set", fileRef)
	}

	if err := addHttpRulesToProto(in.Name, serviceName, descriptors); err != nil {
		return errors.Wrap(err, "failed to generate http rules for proto descriptors")
	}

	// cache the descriptors; we'll need then when we create our grpc filters
	p.serviceDescriptors[serviceName] = descriptors
	// keep track of which service belongs to which upstream
	p.upstreamServices[in.Name] = serviceName

	return nil
}

func convertProto(b []byte) (*descriptor.FileDescriptorSet, error) {
	var fileDescriptor descriptor.FileDescriptorSet
	err := proto.Unmarshal(b, &fileDescriptor)
	return &fileDescriptor, err
}

func (p *Plugin) ProcessRoute(pluginParams *plugin.RoutePluginParams, in *v1.Route, out *envoyroute.Route) error {
	switch {
	case in.SingleDestination != nil:
		err := p.processRouteForGRPC(in.SingleDestination, in.Extensions, out)
		if err != nil {
			return errors.Wrap(err, "processing route for gRPC destination")
		}

	case in.MultipleDestinations != nil:
		for _, dest := range in.MultipleDestinations {
			err := p.processRouteForGRPC(dest.Destination, in.Extensions, out)
			if err != nil {
				return errors.Wrap(err, "processing route for gRPC destination")
			}
		}
	}
	return nil
}

func (p *Plugin) processRouteForGRPC(dest *v1.Destination, extensions *types.Struct, out *envoyroute.Route) error {
	fnDest, ok := dest.DestinationType.(*v1.Destination_Function)
	if !ok {
		// not interested have a nice day
		return nil
	}
	upstreamName := fnDest.Function.UpstreamName
	serviceName, ok := p.upstreamServices[upstreamName]
	if !ok {
		// the upstream is not a grpc desintation
		return nil
	}

	// method name should be function name in this case. TODO: document in the api
	methodName := fnDest.Function.FunctionName

	// create the transformation for the route

	outPath := httpPath(upstreamName, serviceName, methodName)

	routeParams, err := transformation.DecodeRouteExtension(extensions)
	if err != nil {
		return errors.Wrap(err, "parsing route extensions")
	}

	return nil
}

func addHttpRulesToProto(upstreamName, serviceName string, set *descriptor.FileDescriptorSet) error {
	for _, file := range set.File {
		for _, svc := range file.Service {
			if *svc.Name == serviceName {
				for _, method := range svc.Method {
					extension, err := proto.GetExtension(method.Options, api.E_Http)
					if err != nil {
						return errors.Wrap(err, "getting http extensions from method.Options")
					}
					log.Printf("existing extension: %v", extension)
					if err := proto.SetExtension(method.Options, api.E_Http, &api.HttpRule{
						Pattern: &api.HttpRule_Post{
							Post: httpPath(upstreamName, serviceName, *method.Name),
						},
						Body: "*",
					}); err != nil {
						return errors.Wrap(err, "setting http extensions for method.Options")
					}
				}

			}
		}
	}
	return errors.Errorf("could not find match: %v/%v", upstreamName, serviceName)
}

func httpPath(upstreamName, serviceName, methodName string) string {
	h := sha1.New()
	h.Write([]byte(upstreamName + serviceName))
	return "/" + fmt.Sprintf("%x", h.Sum(nil))[:8] + "/" + upstreamName + "/" + serviceName + "/" + methodName
}

func FuncsForProto(serviceName string, set *descriptor.FileDescriptorSet) []*v1.Function {
	var funcs []*v1.Function
	for _, file := range set.File {
		for _, svc := range file.Service {
			if svc.Name == nil || *svc.Name != serviceName {
				continue
			}
			for _, method := range svc.Method {
				g, err := proto.GetExtension(method.Options, api.E_Http)
				if err != nil {
					log.Printf("missing http option on the extensions, skipping: %v", *method.Name)
					continue
				}
				httpRule, ok := g.(*api.HttpRule)
				if !ok {
					panic(g)
				}
				log.Printf("rule: %v", httpRule)
				verb, path := verbAndPathForRule(httpRule)
				fn := &v1.Function{
					Name: *method.Name,
					Spec: transformation.EncodeFunctionSpec(transformation.Template{
						Path:            toInjaTemplateFormat(path),
						Header:          map[string]string{":method": verb},
						PassthroughBody: true,
					}),
				}
				funcs = append(funcs, fn)
			}
		}
		log.Printf("%v", file.MessageType)
	}
	return funcs
}

func toInjaTemplateFormat(in string) string {
	in = strings.Replace(in, "{", "{{", -1)
	return strings.Replace(in, "}", "}}", -1)
}

func verbAndPathForRule(httpRule *api.HttpRule) (string, string) {
	switch rule := httpRule.Pattern.(type) {
	case *api.HttpRule_Get:
		return "GET", rule.Get
	case *api.HttpRule_Custom:
		return rule.Custom.Kind, rule.Custom.Path
	case *api.HttpRule_Delete:
		return "DELETE", rule.Delete
	case *api.HttpRule_Patch:
		return "PATCH", rule.Patch
	case *api.HttpRule_Post:
		return "POST", rule.Post
	case *api.HttpRule_Put:
		return "PUT", rule.Put
	}
	panic("unknown rule type")
}

func lookupMessageType(inputType string, messageTypes []*descriptor.DescriptorProto) *descriptor.DescriptorProto {
	for _, msg := range messageTypes {
		if *msg.Name == inputType {
			return msg
		}
	}
	return nil
}

//func (p *Plugin) HttpFilters(params *plugin.FilterPluginParams) []plugin.StagedFilter {
//
//	if len(p.CachedTransformations) == 0 {
//		return nil
//	}
//
//	filterConfig, err := protoutil.MarshalStruct(&Transformations{
//		Transformations: p.CachedTransformations,
//	})
//	if err != nil {
//		return nil
//	}
//
//	// clear cache
//	p.CachedTransformations = make(map[string]*Transformation)
//
//	return []plugin.StagedFilter{{HttpFilter: &envoyhttp.HttpFilter{
//		Name:   filterName,
//		Config: filterConfig,
//	}, Stage: pluginStage}}
//}
