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

package offlinevalidator

import (
	"testing"
	"os"
	"github.com/Appliscale/perun/specification"
	"github.com/stretchr/testify/assert"
	"github.com/Appliscale/perun/offlinevalidator/template"
	"github.com/Appliscale/perun/logger"
)

var specification specification.Specification
var logger logger.Logger

func setup() {
	var err error
	logger = logger.Logger{}
	specification, err = specification.GetSpecificationFromFile("test_resources/test_specification.json")
	if err != nil {
		panic(err)
	}
}

func TestMain(m *testing.M) {
	setup()
	retCode := m.Run()
	os.Exit(retCode)
}

func TestValidResource(t *testing.T) {
	resources := make(map[string]template.Resource)
	resources["ExampleResource"] = createResourceWithOneProperty("ExampleResourceType", "ExampleProperty", "Property value")

	assert.True(t, validateResources(resources, &specification, &logger), "This resource should be valid")
}

func TestInvalidResourceType(t *testing.T) {
	resources := make(map[string]template.Resource)
	resources["ExampleResource"] = createResourceWithOneProperty("InvalidType", "ExampleProperty", "Property value")

	assert.False(t, validateResources(resources, &specification, &logger), "This resource should be valid, it has invalid resource type")
}

func TestLackOfRequiredPropertyInResource(t *testing.T) {
	resources := make(map[string]template.Resource)
	resources["ExampleResource"] = createResourceWithOneProperty("ExampleResourceType", "SomeProperty", "Property value")

	assert.False(t, validateResources(resources, &specification, &logger), "This resource should not be valid, it do not have required property")
}

func createResourceWithOneProperty(resourceType string, propertyName string, propertyValue string) (template.Resource) {
	resource := template.Resource{}
	resource.Type = resourceType
	resource.Properties = make(map[string]interface{})
	resource.Properties[propertyName] = propertyValue

	return resource
}
