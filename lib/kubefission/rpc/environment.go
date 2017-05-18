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

/* RPC stubs for Environment resource provider */

// EnvironmentToken is the type token corresponding to the Environment package type.
const EnvironmentToken = tokens.Type("kubefission:environment:Environment")

// EnvironmentProviderOps is a pluggable interface for Environment-related management functionality.
type EnvironmentProviderOps interface {
    Check(ctx context.Context, obj *Environment) ([]mapper.FieldError, error)
    Create(ctx context.Context, obj *Environment) (resource.ID, error)
    Get(ctx context.Context, id resource.ID) (*Environment, error)
    InspectChange(ctx context.Context,
        id resource.ID, old *Environment, new *Environment, diff *resource.ObjectDiff) ([]string, error)
    Update(ctx context.Context,
        id resource.ID, old *Environment, new *Environment, diff *resource.ObjectDiff) error
    Delete(ctx context.Context, id resource.ID) error
}

// EnvironmentProvider is a dynamic gRPC-based plugin for managing Environment resources.
type EnvironmentProvider struct {
    ops EnvironmentProviderOps
}

// NewEnvironmentProvider allocates a resource provider that delegates to a ops instance.
func NewEnvironmentProvider(ops EnvironmentProviderOps) lumirpc.ResourceProviderServer {
    contract.Assert(ops != nil)
    return &EnvironmentProvider{ops: ops}
}

func (p *EnvironmentProvider) Check(
    ctx context.Context, req *lumirpc.CheckRequest) (*lumirpc.CheckResponse, error) {
    contract.Assert(req.GetType() == string(EnvironmentToken))
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

func (p *EnvironmentProvider) Name(
    ctx context.Context, req *lumirpc.NameRequest) (*lumirpc.NameResponse, error) {
    contract.Assert(req.GetType() == string(EnvironmentToken))
    obj, _, decerr := p.Unmarshal(req.GetProperties())
    if decerr != nil {
        return nil, decerr
    }
    if obj.Name == "" {
        return nil, errors.New("Name property cannot be empty")
    }
    return &lumirpc.NameResponse{Name: obj.Name}, nil
}

func (p *EnvironmentProvider) Create(
    ctx context.Context, req *lumirpc.CreateRequest) (*lumirpc.CreateResponse, error) {
    contract.Assert(req.GetType() == string(EnvironmentToken))
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

func (p *EnvironmentProvider) Get(
    ctx context.Context, req *lumirpc.GetRequest) (*lumirpc.GetResponse, error) {
    contract.Assert(req.GetType() == string(EnvironmentToken))
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

func (p *EnvironmentProvider) InspectChange(
    ctx context.Context, req *lumirpc.ChangeRequest) (*lumirpc.InspectChangeResponse, error) {
    contract.Assert(req.GetType() == string(EnvironmentToken))
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

func (p *EnvironmentProvider) Update(
    ctx context.Context, req *lumirpc.ChangeRequest) (*pbempty.Empty, error) {
    contract.Assert(req.GetType() == string(EnvironmentToken))
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

func (p *EnvironmentProvider) Delete(
    ctx context.Context, req *lumirpc.DeleteRequest) (*pbempty.Empty, error) {
    contract.Assert(req.GetType() == string(EnvironmentToken))
    id := resource.ID(req.GetId())
    if err := p.ops.Delete(ctx, id); err != nil {
        return nil, err
    }
    return &pbempty.Empty{}, nil
}

func (p *EnvironmentProvider) Unmarshal(
    v *pbstruct.Struct) (*Environment, resource.PropertyMap, mapper.DecodeError) {
    var obj Environment
    props := resource.UnmarshalProperties(v)
    result := mapper.MapIU(props.Mappable(), &obj)
    return &obj, props, result
}

/* Marshalable Environment structure(s) */

// Environment is a marshalable representation of its corresponding IDL type.
type Environment struct {
    Name string `json:"name"`
    RunContainerImageURL string `json:"runContainerImageURL"`
}

// Environment's properties have constants to make dealing with diffs and property bags easier.
const (
    Environment_Name = "name"
    Environment_RunContainerImageURL = "runContainerImageURL"
)


