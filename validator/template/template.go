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

// Package template provides struct which describes AWS Template.
package template

// Template describes AWS Template.
type Template struct {
	AWSTemplateFormatVersion string                 `yaml:"AWSTemplateFormatVersion"`
	Description              string                 `yaml:"Description"`
	Metadata                 map[string]interface{} `yaml:"Metadata"`
	Parameters               map[string]interface{} `yaml:"Parameters"`
	Mappings                 map[string]interface{} `yaml:"Mappings"`
	Conditions               map[string]interface{} `yaml:"Conditions"`
	Transform                map[string]interface{} `yaml:"Transform"`
	Resources                map[string]Resource    `yaml:"Resources"`
	Outputs                  map[string]interface{} `yaml:"Outputs"`
}

// Resource describes structure of Resources in Template.
type Resource struct {
	Type           string                 `yaml:"Type"`
	Properties     map[string]interface{} `yaml:"Properties"`
	DeletionPolicy string                 `yaml:"DeletionPolicy"`
}

// Parameters describes structure of Parameters in Template.
type Parameters map[string]interface{}
