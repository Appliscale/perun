package cfspecification

import (
	"testing"
	"os"
	"github.com/stretchr/testify/assert"
)

var specification Specification

func setup() {
	var err error
	specification, err = GetSpecificationFromFile("test_resources/test_specification.json")
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
	assert.Equal(t,"1.2.3", specification.ResourceSpecificationVersion, "Invalid ResouceSpecificationVersion")
}

func TestPropertyTypes(t *testing.T) {

	assert.Equal(t, "Documentation string", specification.PropertyTypes["ExampleProperties"].Documentation,
		getErrorMessage("PropertyTypes:ExampleProperties:Documentation"))

	assert.Equal(t, "Documentation string", specification.PropertyTypes["ExampleProperties"].Properties["ExampleProperty"].Documentation,
		getErrorMessage("PropertyTypes:ExampleProperties:ExampleProperty:Documentation"))

	assert.Equal(t, true, specification.PropertyTypes["ExampleProperties"].Properties["ExampleProperty"].Required,
		getErrorMessage("PropertyTypes:ExampleProperties:ExampleProperty:Required"))

	assert.Equal(t, "String", specification.PropertyTypes["ExampleProperties"].Properties["ExampleProperty"].PrimitiveType,
		getErrorMessage("PropertyTypes:ExampleProperties:ExampleProperty:PrimitiveType"))

	assert.Equal(t, "Mutable", specification.PropertyTypes["ExampleProperties"].Properties["ExampleProperty"].UpdateType,
		getErrorMessage("PropertyTypes:ExampleProperties:ExampleProperty:UpdateType"))

	assert.Equal(t, true, specification.PropertyTypes["ExampleProperties"].Properties["AnotherExampleProperty"].DuplicatesAllowed,
		getErrorMessage("PropertyTypes:ExampleProperties:AnotherExampleProperty:DuplicatesAllowed"))

	assert.Equal(t, "List", specification.PropertyTypes["ExampleProperties"].Properties["AnotherExampleProperty"].Type,
		getErrorMessage("PropertyTypes:ExampleProperties:AnotherExampleProperty:Type"))

	assert.Equal(t, "Tag", specification.PropertyTypes["ExampleProperties"].Properties["AnotherExampleProperty"].ItemType,
		getErrorMessage("PropertyTypes:ExampleProperties:AnotherExampleProperty:ItemType"))

	assert.Equal(t, "String", specification.PropertyTypes["ExampleProperties"].Properties["AnotherExampleProperty"].PrimitiveItemType,
		getErrorMessage("PropertyTypes:ExampleProperties:AnotherExampleProperty:PrimitiveItemType"))
}

func TestResourceTypes(t *testing.T) {

	assert.Equal(t, "Documentation string", specification.ResourceTypes["ExampleResourceType"].Documentation,
		getErrorMessage("ResourceTypes:ExampleResourceType:Documentation"))

	assert.Equal(t, "String", specification.ResourceTypes["ExampleResourceType"].Attributes["ExampleAttribute"].PrimitiveItemType,
		getErrorMessage("ResourceTypes:ExampleResourceType:ExampleAttribute:PrimitiveType"))

	assert.Equal(t, "List", specification.ResourceTypes["ExampleResourceType"].Attributes["ExampleAttribute"].Type,
		getErrorMessage("ResourceTypes:ExampleResourceType:ExampleAttribute:Type"))

	assert.Equal(t, "Tag", specification.ResourceTypes["ExampleResourceType"].Attributes["AnotherExampleAttribute"].ItemType,
		getErrorMessage("ResourceTypes:ExampleResourceType:AnotherExampleAttribute:ItemType"))

	assert.Equal(t, "String", specification.ResourceTypes["ExampleResourceType"].Attributes["AnotherExampleAttribute"].PrimitiveType,
		getErrorMessage("ResourceTypes:ExampleResourceType:AnotherExampleAttribute:PrimitiveType"))

	assert.Equal(t, "Documentation string", specification.ResourceTypes["ExampleResourceType"].Properties["ExampleProperty"].Documentation,
		getErrorMessage("ResourceTypes:ExampleResourceType:ExampleProperty:Documentation"))

	assert.Equal(t, "String", specification.ResourceTypes["ExampleResourceType"].Properties["ExampleProperty"].PrimitiveType,
		getErrorMessage("ResourceTypes:ExampleResourceType:ExampleProperty:PrimitiveType"))

	assert.Equal(t, true, specification.ResourceTypes["ExampleResourceType"].Properties["ExampleProperty"].Required,
		getErrorMessage("ResourceTypes:ExampleResourceType:ExampleProperty:Required"))

	assert.Equal(t, "Immutable", specification.ResourceTypes["ExampleResourceType"].Properties["ExampleProperty"].UpdateType,
		getErrorMessage("ResourceTypes:ExampleResourceType:ExampleProperty:UpdateType"))
}

func getErrorMessage(fieldName string) string {
	return "Invalid field " + fieldName
}

