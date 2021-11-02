package codegen

const tempFile = `// Code generated by protoc-gen-natsrpc. DO NOT EDIT.
package {{.PackageName}}

import (
	"context"

	"github.com/byebyebruce/natsrpc"
	"github.com/golang/protobuf/proto"
	"github.com/nats-io/nats.go"
)
var _ proto.Message
{{- range .ServiceList}}
{{$serviceAsync := .ServiceAsync}}
{{$clientAsync := .ClientAsync}}

// {{ .ServiceName }}
type {{ .ServiceName }} interface {
{{- range .MethodList }}
// {{ .MethodName }}
	{{- if eq .Publish false }}
		{{- if eq $serviceAsync true }}
		{{ .MethodName }}(ctx context.Context, req *{{ .InputTypeName }}, cb func(*{{ .OutputTypeName }}, error))	
		{{- else }}
		{{ .MethodName }}(ctx context.Context, req *{{ .InputTypeName }}) (*{{ .OutputTypeName }}, error)
		{{- end }}
	{{- else }}
		{{ .MethodName }}(ctx context.Context, req *{{ .InputTypeName }})
	{{- end }}
{{- end }}
}

// Register{{ .ServiceName }}
func Register{{ .ServiceName }}(server *natsrpc.Server, s {{ .ServiceName }}, opts ...natsrpc.ServiceOption) (natsrpc.IService, error) {
	return server.Register("{{ $.GoPackageName }}.{{ .ServiceName }}", s, opts...)
}

{{- $clientName := .ServiceName}}


// {{ $clientName }}Client
type {{ $clientName }}Client struct {
	c *natsrpc.Client
}

// New{{$clientName}}Client
func New{{ $clientName }}Client(enc *nats.EncodedConn, opts ...natsrpc.ClientOption) (*{{ $clientName }}Client, error) {
	c, err := natsrpc.NewClient(enc, "{{ $.GoPackageName }}.{{ $clientName }}", opts...)
	if err != nil {
		return nil, err
	}
	ret := &{{ $clientName }}Client{
		c:c,
	}
	return ret, nil
}

{{- range .MethodList }}
// {{ .MethodName }}
	{{- if eq .Publish false }}
		{{- if eq $clientAsync true }}
		func (c *{{ $clientName }}Client) {{ .MethodName }}(ctx context.Context, req *{{ .InputTypeName }}})(*{{ .OutputTypeName }}, error) {
			rep := &{{ .OutputTypeName }}{}
			err := c.c.Request(ctx, "{{ .MethodName }}", req, rep)
			cb(rep, err)
		}
		{{- else }}
		func (c *{{ $clientName }}Client) {{ .MethodName }}(ctx context.Context, req *{{ .InputTypeName }}}, cb func(*{{ .OutputTypeName }}, error) {
			rep := &{{ .OutputTypeName }}{}
			err := c.c.Request(ctx, "{{ .MethodName }}", req, rep)
			return rep, err 
		}
		{{- end }}
	{{- else }}
		func (c *{{ $clientName }}Client) {{ .MethodName }}(ctx context.Context, notify *{{ .InputTypeName }}}) error {
			return c.c.Publish("{{ .MethodName }}", notify)
		}
	{{- end }}
{{- end }}

{{- end }}
`

/*
{{ range .MethodList}}
	{{- if eq .MethodType 0 -}}
// {{ .MethodName }}Async
func (c *{{ $clientName }}Client) {{ .MethodName }}Async(ctx context.Context, req *{{ .InputTypeName }}, cb func(*{{ .OutputTypeName }}, error)){
	rep := &{{ .OutputTypeName }}{}
	f := func(_ proto.Message, err error) {
		cb(rep, err)
	}
	c.c.AsyncRequest("{{ .MethodName }}", req, rep, f)
}
	{{- end}}

	{{- if eq .MethodType 1 -}}
// {{ .MethodName }}
func (c *{{ $clientName }}Client) {{ .MethodName }}(ctx context.Context, req *{{ .InputTypeName }}) (*{{ .OutputTypeName }}, error) {
	rep := &{{ .OutputTypeName }}{}
	err := c.c.Request(ctx, "{{ .MethodName }}", req, rep)
	return rep, err
}
	{{- end}}

	{{- if eq .MethodType 2 -}}
// Publish{{ .MethodName }}
func (c *{{ $clientName }}Client) Publish{{ .MethodName }}(ctx context.Context, notify *{{ .InputTypeName }}) error {
	return c.c.Publish("{{ .MethodName }}", notify)
}
	{{- end}}
{{ end }}

{{- end}}
*/

const serviceTmpl = `{{- range .ServiceList}}
{{$serviceAsync := .ServiceAsync}}
{{$clientAsync := .ClientAsync}}

// {{ .ServiceName }}
type {{ .ServiceName }} interface {
{{- range .MethodList }}
// {{ .MethodName }}
	{{- if eq .Publish false }}
		{{- if eq $serviceAsync true }}
		{{ .MethodName }}(ctx context.Context, req *{{ .InputTypeName }}, cb func(*{{ .OutputTypeName }}, error))	
		{{- else }}
		{{ .MethodName }}(ctx context.Context, req *{{ .InputTypeName }}) (*{{ .OutputTypeName }}, error)
		{{- end }}
	{{- else }}
		{{ .MethodName }}(ctx context.Context, req *{{ .InputTypeName }})
	{{- end }}
{{- end }}
}

// Register{{ .ServiceName }}
func Register{{ .ServiceName }}(server *natsrpc.Server, s {{ .ServiceName }}, opts ...natsrpc.ServiceOption) (natsrpc.IService, error) {
	return server.Register("{{ $.GoPackageName }}.{{ .ServiceName }}", s, opts...)
}
{{- end }}
`
