
---
title: "filter.proto"
weight: 5
---

<!-- Code generated by solo-kit. DO NOT EDIT. -->


### Package: `envoy.config.filter.http.aws.v2` 
##### Types:


- [LambdaPerRoute](#LambdaPerRoute)
- [LambdaProtocolExtension](#LambdaProtocolExtension)
  



##### Source File: [github.com/solo-io/gloo/projects/gloo/pkg/plugins/aws/filter.proto](https://github.com/solo-io/gloo/blob/master/projects/gloo/pkg/plugins/aws/filter.proto)





---
### <a name="LambdaPerRoute">LambdaPerRoute</a>

 
AWS Lambda contains the configuration necessary to perform transform regular http calls to
AWS Lambda invocations.

```yaml
"name": string
"qualifier": string
"async": bool

```

| Field | Type | Description | Default |
| ----- | ---- | ----------- |----------- | 
| `name` | `string` | The name of the function |  |
| `qualifier` | `string` | The qualifier of the function (defaults to $LATEST if not specified) |  |
| `async` | `bool` | Invocation type - async or regular. |  |




---
### <a name="LambdaProtocolExtension">LambdaProtocolExtension</a>



```yaml
"host": string
"region": string
"access_key": string
"secret_key": string

```

| Field | Type | Description | Default |
| ----- | ---- | ----------- |----------- | 
| `host` | `string` | The host header for AWS this cluster |  |
| `region` | `string` | The region for this cluster |  |
| `access_key` | `string` | The access_key for AWS this cluster |  |
| `secret_key` | `string` | The secret_key for AWS this cluster |  |





<!-- Start of HubSpot Embed Code -->
<script type="text/javascript" id="hs-script-loader" async defer src="//js.hs-scripts.com/5130874.js"></script>
<!-- End of HubSpot Embed Code -->
