// *** WARNING: this file was generated by the Lumi IDL Compiler (LUMIDL). ***
// *** Do not edit by hand unless you're certain you know what you are doing! ***

package kubefission

import (
    "errors"

    pbempty "github.com/golang/protobuf/ptypes/empty"
    pbstruct "github.com/golang/protobuf/ptypes/struct"
    "golang.org/x/net/context"

    "github.com/pulumi/lumi/pkg/resource"
    "github.com/pulumi/lumi/pkg/tokens"
    "github.com/pulumi/lumi/pkg/util/contract"
    "github.com/pulumi/lumi/pkg/util/mapper"
    "github.com/pulumi/lumi/sdk/go/pkg/lumirpc"
)

/* RPC stubs for HTTPTrigger resource provider */

// HTTPTriggerToken is the type token corresponding to the HTTPTrigger package type.
const HTTPTriggerToken = tokens.Type("kubefission:httptrigger:HTTPTrigger")

// HTTPTriggerProviderOps is a pluggable interface for HTTPTrigger-related management functionality.
type HTTPTriggerProviderOps interface {
    Check(ctx context.Context, obj *HTTPTrigger) ([]mapper.FieldError, error)
    Create(ctx context.Context, obj *HTTPTrigger) (resource.ID, error)
    Get(ctx context.Context, id resource.ID) (*HTTPTrigger, error)
    InspectChange(ctx context.Context,
        id resource.ID, old *HTTPTrigger, new *HTTPTrigger, diff *resource.ObjectDiff) ([]string, error)
    Update(ctx context.Context,
        id resource.ID, old *HTTPTrigger, new *HTTPTrigger, diff *resource.ObjectDiff) error
    Delete(ctx context.Context, id resource.ID) error
}

// HTTPTriggerProvider is a dynamic gRPC-based plugin for managing HTTPTrigger resources.
type HTTPTriggerProvider struct {
    ops HTTPTriggerProviderOps
}

// NewHTTPTriggerProvider allocates a resource provider that delegates to a ops instance.
func NewHTTPTriggerProvider(ops HTTPTriggerProviderOps) lumirpc.ResourceProviderServer {
    contract.Assert(ops != nil)
    return &HTTPTriggerProvider{ops: ops}
}

func (p *HTTPTriggerProvider) Check(
    ctx context.Context, req *lumirpc.CheckRequest) (*lumirpc.CheckResponse, error) {
    contract.Assert(req.GetType() == string(HTTPTriggerToken))
    obj, _, decerr := p.Unmarshal(req.GetProperties())
    if decerr == nil || len(decerr.Failures()) == 0 {
        failures, err := p.ops.Check(ctx, obj)
        if err != nil {
            return nil, err
        }
        if len(failures) > 0 {
            decerr = mapper.NewDecodeErr(failures)
        }
    }
    return resource.NewCheckResponse(decerr), nil
}

func (p *HTTPTriggerProvider) Name(
    ctx context.Context, req *lumirpc.NameRequest) (*lumirpc.NameResponse, error) {
    contract.Assert(req.GetType() == string(HTTPTriggerToken))
    obj, _, decerr := p.Unmarshal(req.GetProperties())
    if decerr != nil {
        return nil, decerr
    }
    if obj.Name == "" {
        return nil, errors.New("Name property cannot be empty")
    }
    return &lumirpc.NameResponse{Name: obj.Name}, nil
}

func (p *HTTPTriggerProvider) Create(
    ctx context.Context, req *lumirpc.CreateRequest) (*lumirpc.CreateResponse, error) {
    contract.Assert(req.GetType() == string(HTTPTriggerToken))
    obj, _, decerr := p.Unmarshal(req.GetProperties())
    if decerr != nil {
        return nil, decerr
    }
    id, err := p.ops.Create(ctx, obj)
    if err != nil {
        return nil, err
    }
    return &lumirpc.CreateResponse{
        Id:   string(id),
    }, nil
}

func (p *HTTPTriggerProvider) Get(
    ctx context.Context, req *lumirpc.GetRequest) (*lumirpc.GetResponse, error) {
    contract.Assert(req.GetType() == string(HTTPTriggerToken))
    id := resource.ID(req.GetId())
    obj, err := p.ops.Get(ctx, id)
    if err != nil {
        return nil, err
    }
    return &lumirpc.GetResponse{
        Properties: resource.MarshalProperties(
            nil, resource.NewPropertyMap(obj), resource.MarshalOptions{}),
    }, nil
}

func (p *HTTPTriggerProvider) InspectChange(
    ctx context.Context, req *lumirpc.ChangeRequest) (*lumirpc.InspectChangeResponse, error) {
    contract.Assert(req.GetType() == string(HTTPTriggerToken))
    id := resource.ID(req.GetId())
    old, oldprops, decerr := p.Unmarshal(req.GetOlds())
    if decerr != nil {
        return nil, decerr
    }
    new, newprops, decerr := p.Unmarshal(req.GetNews())
    if decerr != nil {
        return nil, decerr
    }
    var replaces []string
    diff := oldprops.Diff(newprops)
    if diff != nil {
        if diff.Changed("name") {
            replaces = append(replaces, "name")
        }
    }
    more, err := p.ops.InspectChange(ctx, id, old, new, diff)
    if err != nil {
        return nil, err
    }
    return &lumirpc.InspectChangeResponse{
        Replaces: append(replaces, more...),
    }, err
}

func (p *HTTPTriggerProvider) Update(
    ctx context.Context, req *lumirpc.ChangeRequest) (*pbempty.Empty, error) {
    contract.Assert(req.GetType() == string(HTTPTriggerToken))
    id := resource.ID(req.GetId())
    old, oldprops, err := p.Unmarshal(req.GetOlds())
    if err != nil {
        return nil, err
    }
    new, newprops, err := p.Unmarshal(req.GetNews())
    if err != nil {
        return nil, err
    }
    diff := oldprops.Diff(newprops)
    if err := p.ops.Update(ctx, id, old, new, diff); err != nil {
        return nil, err
    }
    return &pbempty.Empty{}, nil
}

func (p *HTTPTriggerProvider) Delete(
    ctx context.Context, req *lumirpc.DeleteRequest) (*pbempty.Empty, error) {
    contract.Assert(req.GetType() == string(HTTPTriggerToken))
    id := resource.ID(req.GetId())
    if err := p.ops.Delete(ctx, id); err != nil {
        return nil, err
    }
    return &pbempty.Empty{}, nil
}

func (p *HTTPTriggerProvider) Unmarshal(
    v *pbstruct.Struct) (*HTTPTrigger, resource.PropertyMap, mapper.DecodeError) {
    var obj HTTPTrigger
    props := resource.UnmarshalProperties(v)
    result := mapper.MapIU(props.Mappable(), &obj)
    return &obj, props, result
}

/* Marshalable HTTPTrigger structure(s) */

// HTTPTrigger is a marshalable representation of its corresponding IDL type.
type HTTPTrigger struct {
    Name string `json:"name"`
    URLPattern string `json:"urlPattern"`
    Method string `json:"method"`
    Function resource.ID `json:"function"`
}

// HTTPTrigger's properties have constants to make dealing with diffs and property bags easier.
const (
    HTTPTrigger_Name = "name"
    HTTPTrigger_URLPattern = "urlPattern"
    HTTPTrigger_Method = "method"
    HTTPTrigger_Function = "function"
)


