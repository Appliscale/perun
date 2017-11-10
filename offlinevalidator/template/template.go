// Copyright 2017 Appliscale
//
// Maintainers and Contributors:
//
//   - Piotr Figwer (piotr.figwer@appliscale.io)
//   - Wojciech Gawroński (wojciech.gawronski@appliscale.io)
//   - Kacper Patro (kacper.patro@appliscale.io)
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

type Template struct {
	AWSTemplateFormatVersion string
	Description string
	Metadata map[string]interface{}
	Parameters map[string]interface{}
	Mappings map[string]interface{}
	Conditions map[string]interface{}
	Transform map[string]interface{}
	Resources map[string]Resource
	Outputs map[string]interface{}
}

type Resource struct {
	Type string
	Properties map[string]interface{}
}
