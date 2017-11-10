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

package specification

import (
	"testing"
	"os"
	"github.com/stretchr/testify/assert"
)

var spec Specification

func setup() {
	var err error
	spec, err = GetSpecificationFromFile("test_resources/test_specification.json")
	if err != nil {
		panic(err)
	}
}

func TestMain(m *testing.M) {
	setup()
	retCode := m.Run()
	os.Exit(retCode)
}

func TestResouceSpecificationVersion(t *testing.T) {
	assert.Equal(t, "1.2.3", spec.ResourceSpecificationVersion, "Invalid ResouceSpecificationVersion")
}

func TestPropertyTypes(t *testing.T) {

	assert.Equal(t, "Documentation string", spec.PropertyTypes["ExampleProperties"].Documentation,
		getErrorMessage("PropertyTypes:ExampleProperties:Documentation"))

	assert.Equal(t, "Documentation string", spec.PropertyTypes["ExampleProperties"].Properties["ExampleProperty"].Documentation,
		getErrorMessage("PropertyTypes:ExampleProperties:ExampleProperty:Documentation"))

	assert.Equal(t, true, spec.PropertyTypes["ExampleProperties"].Properties["ExampleProperty"].Required,
		getErrorMessage("PropertyTypes:ExampleProperties:ExampleProperty:Required"))

	assert.Equal(t, "String", spec.PropertyTypes["ExampleProperties"].Properties["ExampleProperty"].PrimitiveType,
		getErrorMessage("PropertyTypes:ExampleProperties:ExampleProperty:PrimitiveType"))

	assert.Equal(t, "Mutable", spec.PropertyTypes["ExampleProperties"].Properties["ExampleProperty"].UpdateType,
		getErrorMessage("PropertyTypes:ExampleProperties:ExampleProperty:UpdateType"))

	assert.Equal(t, true, spec.PropertyTypes["ExampleProperties"].Properties["AnotherExampleProperty"].DuplicatesAllowed,
		getErrorMessage("PropertyTypes:ExampleProperties:AnotherExampleProperty:DuplicatesAllowed"))

	assert.Equal(t, "List", spec.PropertyTypes["ExampleProperties"].Properties["AnotherExampleProperty"].Type,
		getErrorMessage("PropertyTypes:ExampleProperties:AnotherExampleProperty:Type"))

	assert.Equal(t, "Tag", spec.PropertyTypes["ExampleProperties"].Properties["AnotherExampleProperty"].ItemType,
		getErrorMessage("PropertyTypes:ExampleProperties:AnotherExampleProperty:ItemType"))

	assert.Equal(t, "String", spec.PropertyTypes["ExampleProperties"].Properties["AnotherExampleProperty"].PrimitiveItemType,
		getErrorMessage("PropertyTypes:ExampleProperties:AnotherExampleProperty:PrimitiveItemType"))
}

func TestResourceTypes(t *testing.T) {

	assert.Equal(t, "Documentation string", spec.ResourceTypes["ExampleResourceType"].Documentation,
		getErrorMessage("ResourceTypes:ExampleResourceType:Documentation"))

	assert.Equal(t, "String", spec.ResourceTypes["ExampleResourceType"].Attributes["ExampleAttribute"].PrimitiveItemType,
		getErrorMessage("ResourceTypes:ExampleResourceType:ExampleAttribute:PrimitiveType"))

	assert.Equal(t, "List", spec.ResourceTypes["ExampleResourceType"].Attributes["ExampleAttribute"].Type,
		getErrorMessage("ResourceTypes:ExampleResourceType:ExampleAttribute:Type"))

	assert.Equal(t, "Tag", spec.ResourceTypes["ExampleResourceType"].Attributes["AnotherExampleAttribute"].ItemType,
		getErrorMessage("ResourceTypes:ExampleResourceType:AnotherExampleAttribute:ItemType"))

	assert.Equal(t, "String", spec.ResourceTypes["ExampleResourceType"].Attributes["AnotherExampleAttribute"].PrimitiveType,
		getErrorMessage("ResourceTypes:ExampleResourceType:AnotherExampleAttribute:PrimitiveType"))

	assert.Equal(t, "Documentation string", spec.ResourceTypes["ExampleResourceType"].Properties["ExampleProperty"].Documentation,
		getErrorMessage("ResourceTypes:ExampleResourceType:ExampleProperty:Documentation"))

	assert.Equal(t, "String", spec.ResourceTypes["ExampleResourceType"].Properties["ExampleProperty"].PrimitiveType,
		getErrorMessage("ResourceTypes:ExampleResourceType:ExampleProperty:PrimitiveType"))

	assert.Equal(t, true, spec.ResourceTypes["ExampleResourceType"].Properties["ExampleProperty"].Required,
		getErrorMessage("ResourceTypes:ExampleResourceType:ExampleProperty:Required"))

	assert.Equal(t, "Immutable", spec.ResourceTypes["ExampleResourceType"].Properties["ExampleProperty"].UpdateType,
		getErrorMessage("ResourceTypes:ExampleResourceType:ExampleProperty:UpdateType"))
}

func getErrorMessage(fieldName string) string {
	return "Invalid field " + fieldName
}

