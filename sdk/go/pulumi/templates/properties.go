// Copyright 2016-2018, Pulumi Corporation.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// nolint: lll
package pulumi

import (
	"context"
	"reflect"
	"sync"

	"github.com/pkg/errors"
)

const (
	outputPending = iota
	outputResolved
	outputRejected
)

// Output helps encode the relationship between resources in a Pulumi application. Specifically an output property
// holds onto a value and the resource it came from. An output value can then be provided when constructing new
// resources, allowing that new resource to know both the value as well as the resource the value came from.  This
// allows for a precise "dependency graph" to be created, which properly tracks the relationship between resources.
type Output struct {
	*outputState // protect against value aliasing.
}

// outputState is a heap-allocated block of state for each output property, in case of aliasing.
type outputState struct {
	mutex sync.Mutex
	cond  *sync.Cond

	state uint32 // one of output{Pending,Resolved,Rejected}

	value interface{} // the value of this output if it is resolved.
	err   error       // the error associated with this output if it is rejected.
	known bool        // true if this output's value is known.

	deps []Resource // the dependencies associated with this output property.
}

func (o *outputState) dependencies() []Resource {
	if o == nil {
		return nil
	}
	return o.deps
}

func (o *outputState) fulfill(value interface{}, known bool, err error) {
	if o == nil {
		return
	}

	o.mutex.Lock()
	defer func() {
		o.mutex.Unlock()
		o.cond.Broadcast()
	}()

	if o.state != outputPending {
		return
	}

	if err != nil {
		o.state, o.err, o.known = outputRejected, err, true
	} else {
		o.state, o.value, o.known = outputResolved, value, known
	}
}

func (o *outputState) resolve(value interface{}, known bool) {
	o.fulfill(value, known, nil)
}

func (o *outputState) reject(err error) {
	o.fulfill(nil, true, err)
}

func (o *outputState) await(ctx context.Context) (interface{}, bool, error) {
	for {
		if o == nil {
			// If the state is nil, treat its value as resolved and unknown.
			return nil, false, nil
		}

		o.mutex.Lock()
		for o.state == outputPending {
			if ctx.Err() != nil {
				return nil, true, ctx.Err()
			}
			o.cond.Wait()
		}
		o.mutex.Unlock()

		if !o.known || o.err != nil {
			return nil, o.known, o.err
		}

		ov, ok := isOutput(o.value)
		if !ok {
			return o.value, true, nil
		}
		o = ov.outputState
	}
}

func newOutput(deps ...Resource) Output {
	out := Output{
		&outputState{
			deps: deps,
		},
	}
	out.outputState.cond = sync.NewCond(&out.outputState.mutex)
	return out
}

var outputType = reflect.TypeOf(Output{})

func isOutput(v interface{}) (Output, bool) {
	if v != nil {
		rv := reflect.ValueOf(v)
		if rv.Type().ConvertibleTo(outputType) {
			return rv.Convert(outputType).Interface().(Output), true
		}
	}
	return Output{}, false
}

// NewOutput returns an output value that can be used to rendezvous with the production of a value or error.  The
// function returns the output itself, plus two functions: one for resolving a value, and another for rejecting with an
// error; exactly one function must be called. This acts like a promise.
func NewOutput() (Output, func(interface{}), func(error)) {
	out := newOutput()

	resolve := func(v interface{}) {
		out.resolve(v, true)
	}
	reject := func(err error) {
		out.reject(err)
	}

	return out, resolve, reject
}

// ApplyWithContext transforms the data of the output property using the applier func. The result remains an output
// property, and accumulates all implicated dependencies, so that resources can be properly tracked using a DAG.
// This function does not block awaiting the value; instead, it spawns a Goroutine that will await its availability.
func (out Output) Apply(applier func(v interface{}) (interface{}, error)) Output {
	return out.ApplyWithContext(context.Background(), func(ctx context.Context, v interface{}) (interface{}, error) {
		return applier(v)
	})
}

// ApplyWithContext transforms the data of the output property using the applier func. The result remains an output
// property, and accumulates all implicated dependencies, so that resources can be properly tracked using a DAG.
// This function does not block awaiting the value; instead, it spawns a Goroutine that will await its availability.
// The provided context can be used to reject the output as canceled.
func (out Output) ApplyWithContext(ctx context.Context,
	applier func(ctx context.Context, v interface{}) (interface{}, error)) Output {

	result := newOutput(out.dependencies()...)
	go func() {
		v, known, err := out.await(ctx)
		if err != nil || !known {
			result.fulfill(nil, known, err)
			return
		}

		// If we have a known value, run the applier to transform it.
		u, err := applier(ctx, v)
		if err != nil {
			result.reject(err)
			return
		}

		// Fulfill the result.
		result.fulfill(u, true, nil)
	}()
	return result 
}
{{range .Builtins}}
// Apply{{.Name}} is like Apply, but returns a {{.Name}}Output.
func (out Output) Apply{{.Name}}(applier func (interface{}) ({{.ExportedType}}, error)) {{.Name}}Output {
	return out.Apply{{.Name}}WithContext(context.Background(), func (_ context.Context, v interface{}) ({{.ExportedType}}, error) {
		return applier(v)
	})
}

// Apply{{.Name}}WithContext is like ApplyWithContext, but returns a {{.Name}}Output.
func (out Output) Apply{{.Name}}WithContext(ctx context.Context, applier func(context.Context, interface{}) ({{.ExportedType}}, error)) {{.Name}}Output {
	return {{.Name}}Output(out.ApplyWithContext(ctx, func (ctx context.Context, v interface{}) (interface{}, error) {
		return applier(ctx, v)
	}))
}
{{end}}
// Input is the type of a generic input value for a Pulumi resource. This type is used in conjunction with Output
// to provide polymorphism over strongly-typed input values.
//
// The intended pattern for nested Pulumi value types is to define an input interface and a plain, input, and output
// variant of the value type that implement the input interface.
//
// For example, given a nested Pulumi value type with the following shape:
//
//     type Nested struct {
//         Foo int
//         Bar string
//     }
//
// We would define the following:
//
//     var nestedType = reflect.TypeOf((*Nested)(nil))
//
//     type NestedInputType interface {
//         pulumi.Input
//
//         isNested()
//     }
//
//     type Nested struct {
//         Foo int `pulumi:"foo"`
//         Bar string `pulumi:"bar"`
//     }
//
//     func (*Nested) ElementType() reflect.Type {
//         return nestedType
//     }
//
//     func (*Nested) isNested() {}
//
//     type NestedInput struct {
//         Foo pulumi.IntInput `pulumi:"foo"`
//         Bar pulumi.StringInput `pulumi:"bar"`
//     }
//
//     func (*NestedInput) ElementType() reflect.Type {
//         return nestedType
//     }
//
//     func (*NestedInput) isNested() {}
//
//     type NestedOutput pulumi.Output
//
//     func (NestedOutput) ElementType() reflect.Type {
//         return nestedType
//     }
//
//     func (NestedOutput) isNested() {}
//
//     func (out NestedOutput) Apply(applier func(*Nested) (interface{}, error)) {
//         return out.ApplyWithContext(context.Background(), func(_ context.Context, v *Nested) (interface{}, error) {
//             return applier(v)
//         })
//     }
//
//     func (out NestedOutput) ApplyWithContext(ctx context.Context, applier func(context.Context, *Nested) (interface{}, error) {
//         return pulumi.Output(out).ApplyWithContext(ctx, func(ctx context.Context, v interface{}) (interface{}, error) {
//             return applier(ctx, v.(*Nested))
//         })
//     }
//
type Input interface {
	ElementType() reflect.Type
}

type anyInput struct {
	v interface{}
}

func Any(v interface{}) AnyInput {
	return anyInput{v: v}
}

{{with $builtins := .Builtins}}
{{range $builtins}}
var {{.Name | ToLower}}Type = reflect.TypeOf((*{{.Type}})(nil)).Elem()

// {{.Name}}Input is an input type that accepts {{.Name}} and {{.Name}}Output values.
type {{.Name}}Input interface {
	Input

	// nolint: unused
	is{{.Name}}()
}
{{if .DefineInputType}}
// {{.Name}} is an input type for {{.ExportedType}} values.
type {{.Name}} {{.Type}} 
{{end}}
{{if .DefineInputMethods}}
// ElementType returns the element type of this Input ({{.Type}}).
func ({{.InputType}}) ElementType() reflect.Type {
	return {{.Name | ToLower}}Type
}

func ({{.InputType}}) is{{.Name}}() {}

{{with $builtin := .}}
{{range $t := .Implements}}
func ({{$builtin.InputType}}) is{{$t}}() {}
{{end}}
{{end}}
{{end}}
// {{.Name}}Output is an Output that returns {{.ExportedType}} values.
type {{.Name}}Output Output

// ElementType returns the element type of this Output ({{.ExportedType}}).
func ({{.Name}}Output) ElementType() reflect.Type {
	return {{.Name | ToLower}}Type
}

func ({{.Name}}Output) is{{.Name}}() {}

{{with $builtin := .}}
{{range $t := .Implements}}
func ({{$builtin.Name}}Output) is{{$t}}() {}
{{end}}
{{end}}
// Apply applies a transformation to the {{.Name | ToLower}} value when it is available.
func (out {{.Name}}Output) Apply(applier func({{.ExportedType}}) (interface{}, error)) Output {
	return out.ApplyWithContext(context.Background(), func(_ context.Context, v {{.ExportedType}}) (interface{}, error) {
		return applier(v)
	})
}

// ApplyWithContext applies a transformation to the {{.Name | ToLower}} value when it is available.
func (out {{.Name}}Output) ApplyWithContext(ctx context.Context, applier func(context.Context, {{.ExportedType}}) (interface{}, error)) Output {
	return Output(out).ApplyWithContext(ctx, func(ctx context.Context, v interface{}) (interface{}, error) {
		return applier(ctx, {{if eq .Type "interface{}"}}v{{else}}convert(v, {{.Name | ToLower}}Type).({{.Type}}){{end}})
	})
}
{{with $me := .}}
{{range $builtins}}
// Apply{{.Name}} is like Apply, but returns a {{.Name}}Output.
func (out {{$me.Name}}Output) Apply{{.Name}}(applier func (v {{$me.ExportedType}}) ({{.ExportedType}}, error)) {{.Name}}Output {
	return out.Apply{{.Name}}WithContext(context.Background(), func (_ context.Context, v {{$me.ExportedType}}) ({{.ExportedType}}, error) {
		return applier(v)
	})
}

// Apply{{.Name}}WithContext is like ApplyWithContext, but returns a {{.Name}}Output.
func (out {{$me.Name}}Output) Apply{{.Name}}WithContext(ctx context.Context, applier func(context.Context, {{$me.ExportedType}}) ({{.ExportedType}}, error)) {{.Name}}Output {
	return {{.Name}}Output(Output(out).ApplyWithContext(ctx, func (ctx context.Context, v interface{}) (interface{}, error) {
		return applier(ctx, {{if eq $me.Type "interface{}"}}v{{else}}convert(v, {{$me.Name | ToLower}}Type).({{$me.Type}}){{end}})
	}))
}
{{end}}
{{end}}
{{end}}
{{end}}

func (out IDOutput) await(ctx context.Context) (ID, bool, error) {
	id, known, err := Output(out).await(ctx)
	if !known || err != nil {
		return "", known, err
	}
	return ID(convert(id, stringType).(string)), true, nil
}

func (out URNOutput) await(ctx context.Context) (URN, bool, error) {
	id, known, err := Output(out).await(ctx)
	if !known || err != nil {
		return "", known, err
	}
	return URN(convert(id, stringType).(string)), true, nil
}

func convert(v interface{}, to reflect.Type) interface{} {
	rv := reflect.ValueOf(v)
	if !rv.Type().ConvertibleTo(to) {
		panic(errors.Errorf("cannot convert output value of type %s to %s", rv.Type(), to))
	}
	return rv.Convert(to).Interface()
}
