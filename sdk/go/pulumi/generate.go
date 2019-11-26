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

// +build ignore

package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"text/template"
	"unicode"
)

type builtin struct {
	Name        string
	Type        string
	inputType   string
	Implements  []string
	elementType string
	Example     string
}

func (b builtin) DefineInputType() bool {
	return b.inputType == "" && b.Type != "AssetOrArchive"
}

func (b builtin) DefineInputMethods() bool {
	return b.Type != "AssetOrArchive"
}

func (b builtin) ElementType() string {
	if b.elementType != "" {
		return b.elementType
	}
	return b.Type
}

func (b builtin) InputType() string {
	if b.inputType != "" {
		return b.inputType
	}
	return b.Name
}

var builtins = makeBuiltins([]builtin{
	{Name: "Archive", Type: "*archive", inputType: "*archive", Implements: []string{"AssetOrArchive"}, Example: "NewFileArchive(\"foo.zip\")"},
	{Name: "Asset", Type: "*asset", inputType: "*asset", Implements: []string{"AssetOrArchive"}, Example: "NewFileAsset(\"foo.txt\")"},
	{Name: "AssetOrArchive", Type: "AssetOrArchive", Example: "NewFileArchive(\"foo.zip\")"},
	{Name: "Bool", Type: "bool", Example: "Bool(true)"},
	{Name: "Float32", Type: "float32", Example: "Float32(1.3)"},
	{Name: "Float64", Type: "float64", Example: "Float64(999.9)"},
	{Name: "ID", Type: "ID", inputType: "ID", Implements: []string{"String"}, Example: "ID(\"foo\")"},
	{Name: "Input", Type: "interface{}", Example: "String(\"any\")"},
	{Name: "Int", Type: "int", Example: "Int(42)"},
	{Name: "Int16", Type: "int16", Example: "Int16(33)"},
	{Name: "Int32", Type: "int32", Example: "Int32(24)"},
	{Name: "Int64", Type: "int64", Example: "Int64(15)"},
	{Name: "Int8", Type: "int8", Example: "Int8(6)"},
	{Name: "String", Type: "string", Example: "String(\"foo\")"},
	{Name: "URN", Type: "URN", inputType: "URN", Implements: []string{"String"}, Example: "URN(\"foo\")"},
	{Name: "Uint", Type: "uint", Example: "Uint(42)"},
	{Name: "Uint16", Type: "uint16", Example: "Uint16(33)"},
	{Name: "Uint32", Type: "uint32", Example: "Uint32(24)"},
	{Name: "Uint64", Type: "uint64", Example: "Uint64(15)"},
	{Name: "Uint8", Type: "uint8", Example: "Uint8(6)"},
})

var funcs = template.FuncMap{
	"Unexported": func(s string) string {
		runes := []rune(s)
		return string(append([]rune{unicode.ToLower(runes[0])}, runes[1:]...))
	},
}

func makeBuiltins(primitives []builtin) []builtin {
	// Augment primitives with array and map types.
	var builtins []builtin
	for _, p := range primitives {
		name := ""
		if p.Name != "Input" {
			builtins = append(builtins, p)
			name = p.Name
		}
		builtins = append(builtins, builtin{Name: name + "Array", Type: "[]" + name + "Input", elementType: "[]" + p.Type, Example: fmt.Sprintf("%sArray{%s}", name, p.Example)})
		builtins = append(builtins, builtin{Name: name + "Map", Type: "map[string]" + name + "Input", elementType: "map[string]" + p.Type, Example: fmt.Sprintf("%sMap{\"baz\": %s}", name, p.Example)})
	}
	return builtins
}

func main() {
	templates, err := template.New("templates").Funcs(funcs).ParseGlob("./templates/*")
	if err != nil {
		log.Fatalf("failed to parse templates: %v", err)
	}

	data := map[string]interface{}{
		"Builtins": builtins,
	}
	for _, t := range templates.Templates() {
		f, err := os.Create(t.Name())
		if err != nil {
			log.Fatalf("failed to create %v: %v", t.Name(), err)
		}
		if err := t.Execute(f, data); err != nil {
			log.Fatalf("failed to execute %v: %v", t.Name(), err)
		}
		f.Close()

		gofmt := exec.Command("gofmt", "-s", "-w", t.Name())
		stderr, err := gofmt.StderrPipe()
		if err != nil {
			log.Fatalf("failed to pipe stderr from gofmt: %v", err)
		}
		go func() {
			io.Copy(os.Stderr, stderr)
		}()
		if err := gofmt.Run(); err != nil {
			log.Fatalf("failed to gofmt %v: %v", t.Name(), err)
		}
	}
}
