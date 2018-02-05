// Copyright 2018 Appliscale
//
// Maintainers and contributors are listed in README file inside repository.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package parsers

import (
	"testing"
	"os"
	"io/ioutil"
	template "github.com/Appliscale/perun/offlinevalidator/template"
	"strings"
	"github.com/stretchr/testify/assert"
	"fmt"
)

func TestMain(m *testing.M) {
	retCode := m.Run()
	os.Exit(retCode)
}

func TestJsonParser(t *testing.T) {
	fileName := "test_resources/sample_template.json"
	fileContents, err := ioutil.ReadFile(fileName)
	if err != nil {
		t.Fatal("Error while loading file: " + fileName, err)
	}
	tmpl := template.TemplateWithDetails{}
	err = ParseJson(fileContents, &tmpl)
	if err != nil {
		t.Fatal("Error while parsing file: ", err)
	}
	element := tmpl.Mappings.GetChildrenMap()["RegionMap"].GetChildrenMap()["us-east-1"].GetChildrenMap()["32"]
	assert.NotNil(t, element)
	assert.Equal(t, template.String, element.Type)
	assert.Equal(t, "32", element.Name)
	assert.Equal(t, "ami-6411e20d", element.Value)
	assert.Equal(t, 5, element.Line)
	assert.Equal(t, 23, element.Column)
	assert.Nil(t, element.Children)

	element = tmpl.Resources.GetChildrenMap()["myEC2Instance"].GetChildrenMap()["Properties"].GetChildrenMap()["ImageId"].GetChildrenMap()["Fn::FindInMap"].GetChildrenSlice()[0]
	assert.NotNil(t, element)
	assert.Equal(t, template.String, element.Type)
	assert.Equal(t, "[0]", element.Name)
	assert.Equal(t, "RegionMap", element.Value)
	assert.Equal(t, 17, element.Line)
	assert.Equal(t, 43, element.Column)
	assert.Nil(t, element.Children)

	element = tmpl.Resources.GetChildrenMap()["myEC2Instance"].GetChildrenMap()["Properties"].GetChildrenMap()["ImageId"].GetChildrenMap()["Fn::FindInMap"].GetChildrenSlice()[1]
	assert.NotNil(t, element)
	assert.Equal(t, template.Object, element.Type)
	assert.Equal(t, "[1]", element.Name)
	assert.Nil(t, element.Value)
	assert.Equal(t, 17, element.Line)
	assert.Equal(t, 56, element.Column)
	assert.NotNil(t, element.Children)

	tmpl.Traverse(func(element *template.TemplateElement, parent *template.TemplateElement, depth int) {
		prefix := strings.Repeat("\t", depth)
		fmt.Printf("%s%s @%d:%d\n", prefix, element.Name, element.Line, element.Column)
		if element.Type != template.Object && element.Type != template.Array {
			fmt.Printf("%sValue: %v\n", prefix, element.Value)
		}
	})

}
