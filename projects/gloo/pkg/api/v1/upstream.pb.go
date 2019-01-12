// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: github.com/solo-io/gloo/projects/gloo/api/v1/upstream.proto

package v1 // import "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"

import proto "github.com/gogo/protobuf/proto"
import fmt "fmt"
import math "math"
import _ "github.com/gogo/protobuf/gogoproto"
import core "github.com/solo-io/solo-kit/pkg/api/v1/resources/core"

import bytes "bytes"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.GoGoProtoPackageIsVersion2 // please upgrade the proto package

//
// @solo-kit:resource.short_name=us
// @solo-kit:resource.plural_name=upstreams
// @solo-kit:resource.resource_groups=api.gloo.solo.io,discovery.gloo.solo.io,translator.supergloo.solo.io
//
// Upstreams represent destination for routing HTTP requests. Upstreams can be compared to
// [clusters](https://www.envoyproxy.io/docs/envoy/latest/api-v1/cluster_manager/cluster.html?highlight=cluster) in Envoy terminology.
// Each upstream in Gloo has a type. Supported types include `static`, `kubernetes`, `aws`, `consul`, and more.
// Each upstream type is handled by a corresponding Gloo plugin.
type Upstream struct {
	// Type-specific configuration. Examples include static, kubernetes, and aws.
	// The type-specific config for the upstream is called a spec.
	UpstreamSpec *UpstreamSpec `protobuf:"bytes,2,opt,name=upstream_spec,json=upstreamSpec" json:"upstream_spec,omitempty"`
	// Status indicates the validation status of the resource. Status is read-only by clients, and set by gloo during validation
	Status core.Status `protobuf:"bytes,6,opt,name=status" json:"status" testdiff:"ignore"`
	// Metadata contains the object metadata for this resource
	Metadata core.Metadata `protobuf:"bytes,7,opt,name=metadata" json:"metadata"`
	// Upstreams and their configuration can be automatically by Gloo Discovery
	// if this upstream is created or modified by Discovery, metadata about the operation will be placed here.
	DiscoveryMetadata    *DiscoveryMetadata `protobuf:"bytes,8,opt,name=discovery_metadata,json=discoveryMetadata" json:"discovery_metadata,omitempty"`
	XXX_NoUnkeyedLiteral struct{}           `json:"-"`
	XXX_unrecognized     []byte             `json:"-"`
	XXX_sizecache        int32              `json:"-"`
}

func (m *Upstream) Reset()         { *m = Upstream{} }
func (m *Upstream) String() string { return proto.CompactTextString(m) }
func (*Upstream) ProtoMessage()    {}
func (*Upstream) Descriptor() ([]byte, []int) {
	return fileDescriptor_upstream_745706e49a8c38ac, []int{0}
}
func (m *Upstream) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Upstream.Unmarshal(m, b)
}
func (m *Upstream) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Upstream.Marshal(b, m, deterministic)
}
func (dst *Upstream) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Upstream.Merge(dst, src)
}
func (m *Upstream) XXX_Size() int {
	return xxx_messageInfo_Upstream.Size(m)
}
func (m *Upstream) XXX_DiscardUnknown() {
	xxx_messageInfo_Upstream.DiscardUnknown(m)
}

var xxx_messageInfo_Upstream proto.InternalMessageInfo

func (m *Upstream) GetUpstreamSpec() *UpstreamSpec {
	if m != nil {
		return m.UpstreamSpec
	}
	return nil
}

func (m *Upstream) GetStatus() core.Status {
	if m != nil {
		return m.Status
	}
	return core.Status{}
}

func (m *Upstream) GetMetadata() core.Metadata {
	if m != nil {
		return m.Metadata
	}
	return core.Metadata{}
}

func (m *Upstream) GetDiscoveryMetadata() *DiscoveryMetadata {
	if m != nil {
		return m.DiscoveryMetadata
	}
	return nil
}

// created by discovery services
type DiscoveryMetadata struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *DiscoveryMetadata) Reset()         { *m = DiscoveryMetadata{} }
func (m *DiscoveryMetadata) String() string { return proto.CompactTextString(m) }
func (*DiscoveryMetadata) ProtoMessage()    {}
func (*DiscoveryMetadata) Descriptor() ([]byte, []int) {
	return fileDescriptor_upstream_745706e49a8c38ac, []int{1}
}
func (m *DiscoveryMetadata) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_DiscoveryMetadata.Unmarshal(m, b)
}
func (m *DiscoveryMetadata) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_DiscoveryMetadata.Marshal(b, m, deterministic)
}
func (dst *DiscoveryMetadata) XXX_Merge(src proto.Message) {
	xxx_messageInfo_DiscoveryMetadata.Merge(dst, src)
}
func (m *DiscoveryMetadata) XXX_Size() int {
	return xxx_messageInfo_DiscoveryMetadata.Size(m)
}
func (m *DiscoveryMetadata) XXX_DiscardUnknown() {
	xxx_messageInfo_DiscoveryMetadata.DiscardUnknown(m)
}

var xxx_messageInfo_DiscoveryMetadata proto.InternalMessageInfo

func init() {
	proto.RegisterType((*Upstream)(nil), "gloo.solo.io.Upstream")
	proto.RegisterType((*DiscoveryMetadata)(nil), "gloo.solo.io.DiscoveryMetadata")
}
func (this *Upstream) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*Upstream)
	if !ok {
		that2, ok := that.(Upstream)
		if ok {
			that1 = &that2
		} else {
			return false
		}
	}
	if that1 == nil {
		return this == nil
	} else if this == nil {
		return false
	}
	if !this.UpstreamSpec.Equal(that1.UpstreamSpec) {
		return false
	}
	if !this.Status.Equal(&that1.Status) {
		return false
	}
	if !this.Metadata.Equal(&that1.Metadata) {
		return false
	}
	if !this.DiscoveryMetadata.Equal(that1.DiscoveryMetadata) {
		return false
	}
	if !bytes.Equal(this.XXX_unrecognized, that1.XXX_unrecognized) {
		return false
	}
	return true
}
func (this *DiscoveryMetadata) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*DiscoveryMetadata)
	if !ok {
		that2, ok := that.(DiscoveryMetadata)
		if ok {
			that1 = &that2
		} else {
			return false
		}
	}
	if that1 == nil {
		return this == nil
	} else if this == nil {
		return false
	}
	if !bytes.Equal(this.XXX_unrecognized, that1.XXX_unrecognized) {
		return false
	}
	return true
}

func init() {
	proto.RegisterFile("github.com/solo-io/gloo/projects/gloo/api/v1/upstream.proto", fileDescriptor_upstream_745706e49a8c38ac)
}

var fileDescriptor_upstream_745706e49a8c38ac = []byte{
	// 325 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x94, 0x90, 0x4d, 0x4e, 0x02, 0x31,
	0x1c, 0xc5, 0x85, 0x18, 0x24, 0x15, 0x17, 0x8c, 0xc4, 0x20, 0x0b, 0x31, 0xb3, 0x72, 0x63, 0x2b,
	0x9a, 0x18, 0x83, 0x0b, 0x13, 0x62, 0xe2, 0x4a, 0x17, 0x43, 0xdc, 0xb8, 0x21, 0xa5, 0x53, 0x6a,
	0xe5, 0xe3, 0xdf, 0xb4, 0x1d, 0x12, 0x2f, 0xe3, 0xda, 0xa3, 0x78, 0x0a, 0x16, 0x1e, 0xc1, 0x13,
	0x98, 0x29, 0xed, 0x04, 0xd4, 0x05, 0xae, 0x66, 0x3a, 0xef, 0xfd, 0xde, 0xf4, 0x3d, 0x74, 0x2d,
	0xa4, 0x7d, 0xce, 0x86, 0x98, 0xc1, 0x94, 0x18, 0x98, 0xc0, 0xa9, 0x04, 0x22, 0x26, 0x00, 0x44,
	0x69, 0x78, 0xe1, 0xcc, 0x9a, 0xe5, 0x89, 0x2a, 0x49, 0xe6, 0x1d, 0x92, 0x29, 0x63, 0x35, 0xa7,
	0x53, 0xac, 0x34, 0x58, 0x88, 0x6a, 0xb9, 0x86, 0x73, 0x0c, 0x4b, 0x68, 0x35, 0x04, 0x08, 0x70,
	0x02, 0xc9, 0xdf, 0x96, 0x9e, 0x56, 0xe7, 0x8f, 0x1f, 0xb8, 0xe7, 0x58, 0xda, 0x10, 0x3b, 0xe5,
	0x96, 0xa6, 0xd4, 0x52, 0x8f, 0x90, 0x0d, 0x10, 0x63, 0xa9, 0xcd, 0x8c, 0x07, 0xba, 0xff, 0x2a,
	0xa1, 0x26, 0x99, 0x90, 0x33, 0xcf, 0xc6, 0x6f, 0x65, 0x54, 0x7d, 0xf4, 0xb5, 0xa2, 0x1b, 0xb4,
	0x17, 0x2a, 0x0e, 0x8c, 0xe2, 0xac, 0x59, 0x3e, 0x2e, 0x9d, 0xec, 0x9e, 0xb7, 0xf0, 0x6a, 0x51,
	0x1c, 0xec, 0x7d, 0xc5, 0x59, 0x52, 0xcb, 0x56, 0x4e, 0xd1, 0x1d, 0xaa, 0x2c, 0x6f, 0xd6, 0xac,
	0x38, 0xb2, 0x81, 0x19, 0x68, 0x5e, 0x90, 0x7d, 0xa7, 0xf5, 0x0e, 0x3f, 0x16, 0xed, 0xad, 0xaf,
	0x45, 0xbb, 0x6e, 0xb9, 0xb1, 0xa9, 0x1c, 0x8d, 0xba, 0xb1, 0x14, 0x33, 0xd0, 0x3c, 0x4e, 0x3c,
	0x1e, 0x5d, 0xa1, 0x6a, 0x58, 0xa5, 0xb9, 0xe3, 0xa2, 0x0e, 0xd6, 0xa3, 0xee, 0xbd, 0xda, 0xdb,
	0xce, 0xc3, 0x92, 0xc2, 0x1d, 0x3d, 0xa0, 0x28, 0x95, 0x86, 0xc1, 0x9c, 0xeb, 0xd7, 0x41, 0x91,
	0x51, 0x75, 0x19, 0xed, 0xf5, 0x22, 0xb7, 0xc1, 0x17, 0xc2, 0x92, 0x7a, 0xfa, 0xf3, 0x53, 0xbc,
	0x8f, 0xea, 0xbf, 0x7c, 0xbd, 0xcb, 0xf7, 0xcf, 0xa3, 0xd2, 0xd3, 0xd9, 0x66, 0xbb, 0xab, 0xb1,
	0xf0, 0xdb, 0x0f, 0x2b, 0x6e, 0xf4, 0x8b, 0xef, 0x00, 0x00, 0x00, 0xff, 0xff, 0xa0, 0x26, 0x75,
	0x6c, 0x77, 0x02, 0x00, 0x00,
}
