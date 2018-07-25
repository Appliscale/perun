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

package template

type TemplateWithDetails struct {
	AWSTemplateFormatVersion *TemplateElement
	Description              *TemplateElement
	Metadata                 *TemplateElement
	Parameters               *TemplateElement
	Mappings                 *TemplateElement
	Conditions               *TemplateElement
	Transform                *TemplateElement
	Resources                *TemplateElement
	Outputs                  *TemplateElement
}

type TemplateElement struct {
	Name     string
	Value    interface{}
	Type     TemplateElementValueType
	Children interface{}
	Line     int
	Column   int
}

type TemplateElementValueType int

const (
	NotExist = TemplateElementValueType(iota)
	String
	Number
	Object
	Array
	Boolean
	Null
	Unknown
)

func (te *TemplateElement) GetChildrenMap() map[string]*TemplateElement {
	return te.Children.(map[string]*TemplateElement)
}

func (te *TemplateElement) GetChildrenSlice() []*TemplateElement {
	return *te.Children.(*[]*TemplateElement)
}

func (te *TemplateElement) Traverse(iterator func(element *TemplateElement, parent *TemplateElement, depth int)) {
	if te != nil {
		te.traverse(iterator, nil, 0)
	}
}

func (te *TemplateElement) traverse(iterator func(element *TemplateElement, parent *TemplateElement, depth int), parent *TemplateElement, depth int) {
	iterator(te, parent, depth)
	if te.Type == Object {
		for _, v := range te.GetChildrenMap() {
			v.traverse(iterator, te, depth+1)
		}
	} else if te.Type == Array {
		for _, v := range te.GetChildrenSlice() {
			v.traverse(iterator, te, depth+1)
		}
	}
}

func (twd TemplateWithDetails) Traverse(iterator func(element *TemplateElement, parent *TemplateElement, depth int)) {
	twd.AWSTemplateFormatVersion.Traverse(iterator)
	twd.Description.Traverse(iterator)
	twd.Metadata.Traverse(iterator)
	twd.Parameters.Traverse(iterator)
	twd.Mappings.Traverse(iterator)
	twd.Conditions.Traverse(iterator)
	twd.Transform.Traverse(iterator)
	twd.Resources.Traverse(iterator)
	twd.Outputs.Traverse(iterator)
}
