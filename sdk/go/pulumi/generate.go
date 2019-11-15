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
	"log"
	"os"
	"os/exec"
	"strings"
	"text/template"
)

type builtin struct {
	Name       string
	Type       string
	inputType  string
	Implements []string
}

func (b builtin) DefineInputType() bool {
	return b.inputType == "" && b.Type != "AssetOrArchive"
}

func (b builtin) DefineInputMethods() bool {
	return b.Type != "AssetOrArchive"
}

func (b builtin) InputType() string {
	if b.inputType != "" {
		return b.inputType
	}
	return b.Name
}

func (b builtin) ExportedType() string {
	switch b.Type {
	case "*archive":
		return "Archive"
	case "*asset":
		return "Asset"
	}
	return b.Type
}

var builtins = []builtin{
	{Name: "Any", Type: "interface{}", inputType: "anyInput"},
	{Name: "Archive", Type: "*archive", inputType: "*archive", Implements: []string{"AssetOrArchive"}},
	{Name: "Array", Type: "[]interface{}"},
	{Name: "Asset", Type: "*asset", inputType: "*asset", Implements: []string{"AssetOrArchive"}},
	{Name: "AssetOrArchive", Type: "AssetOrArchive"},
	{Name: "Bool", Type: "bool"},
	{Name: "Float32", Type: "float32"},
	{Name: "Float64", Type: "float64"},
	{Name: "ID", Type: "ID", inputType: "ID", Implements: []string{"String"}},
	{Name: "Int", Type: "int"},
	{Name: "Int16", Type: "int16"},
	{Name: "Int32", Type: "int32"},
	{Name: "Int64", Type: "int64"},
	{Name: "Int8", Type: "int8"},
	{Name: "Map", Type: "map[string]interface{}"},
	{Name: "String", Type: "string"},
	{Name: "URN", Type: "URN", inputType: "URN", Implements: []string{"String"}},
	{Name: "Uint", Type: "uint"},
	{Name: "Uint16", Type: "uint16"},
	{Name: "Uint32", Type: "uint32"},
	{Name: "Uint64", Type: "uint64"},
	{Name: "Uint8", Type: "uint8"},
}

var funcs = template.FuncMap{
	"ToLower": func(s string) string {
		return strings.ToLower(s)
	},
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

		if err := exec.Command("gofmt", "-s", "-w", t.Name()).Run(); err != nil {
			log.Fatalf("failed to gofmt %v: %v", t.Name(), err)
		}
	}
}
