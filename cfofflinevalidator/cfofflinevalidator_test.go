package cfofflinevalidator

import (
	"testing"
	"os"
	"github.com/Appliscale/cftool/cfspecification"
	"github.com/stretchr/testify/assert"
)

var specification cfspecification.Specification

func setup() {
	var err error
	specification, err = cfspecification.GetSpecificationFromFile("resources/test_specification.json")
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
	resources := make(map[string]Resource)
	resources["ExampleResource"] = createResourceWithOneProperty("ExampleResourceType", "ExampleProperty", "Property value")

	assert.True(t, validateResources(resources, &specification), "This resource should be valid")
}

func TestInvalidResourceType(t *testing.T) {
	resources := make(map[string]Resource)
	resources["ExampleResource"] = createResourceWithOneProperty("InvalidType", "ExampleProperty", "Property value")

	assert.False(t, validateResources(resources, &specification), "This resource should be valid, it has invalid resource type")
}

func TestLackOfRequiredPropertyInResource(t *testing.T) {
	resources := make(map[string]Resource)
	resources["ExampleResource"] = createResourceWithOneProperty("ExampleResourceType", "SomeProperty", "Property value")

	assert.False(t, validateResources(resources, &specification), "This resource should not be valid, it do not have required property")
}

func createResourceWithOneProperty(resourceType string, propertyName string, propertyValue string) (Resource) {
	resource := Resource{}
	resource.Type = resourceType
	resource.Properties = make(map[string]interface{})
	resource.Properties[propertyName] = propertyValue

	return resource
}