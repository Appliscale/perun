// Copyright 2017 Appliscale
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

package offlinevalidator

import (
	"os"
	"testing"

	"github.com/Appliscale/perun/configuration"
	"github.com/Appliscale/perun/context"
	"github.com/Appliscale/perun/logger"
	"github.com/Appliscale/perun/offlinevalidator/template"
	"github.com/Appliscale/perun/specification"
	"github.com/stretchr/testify/assert"
)

var spec specification.Specification

var sink logger.Logger

var deadProp = make([]string, 0)
var deadRes = make([]string, 0)
var specInconsistency map[string]configuration.Property
var mockContext = context.Context{}

func setup() {
	var err error

	spec, err = specification.GetSpecificationFromFile("test_resources/test_specification.json")

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

	mockContext.Logger = &logger.Logger{}
	resources := make(map[string]template.Resource)
	resources["ExampleResource"] = createResourceWithOneProperty("ExampleResourceType", "ExampleProperty", "Property value")

	assert.True(t, validateResources(resources, &spec, deadProp, deadRes, specInconsistency, &mockContext), "This resource should be valid")
}

func TestInvalidResourceType(t *testing.T) {
	mockContext.Logger = &logger.Logger{}
	resources := make(map[string]template.Resource)
	resources["ExampleResource"] = createResourceWithOneProperty("InvalidType", "ExampleProperty", "Property value")

	assert.False(t, validateResources(resources, &spec, deadProp, deadRes, specInconsistency, &mockContext), "This resource should be invalid, it has invalid resource type")
}

func TestLackOfRequiredPropertyInResource(t *testing.T) {
	mockContext.Logger = &logger.Logger{}
	resources := make(map[string]template.Resource)
	resources["ExampleResource"] = createResourceWithOneProperty("ExampleResourceType", "SomeProperty", "Property value")

	assert.False(t, validateResources(resources, &spec, deadProp, deadRes, specInconsistency, &mockContext), "This resource should not be valid, it does not have required property")
}
func TestLackOfSubpropertyWithSpecification(t *testing.T) {
	mockContext.Logger = &logger.Logger{}
	resources := make(map[string]template.Resource)
	properties := map[string]interface{}{
		"Ec2KeyName": "SomeValue",
	}
	resources["cluster"] = createResourceWithNestedProperties("AWS::Nested3::Cluster", "SomeProperty", properties)

	assert.False(t, validateResources(resources, &spec, deadProp, deadRes, specInconsistency, &mockContext), "This resource should not be valid, it does not have property with specification")
}
func TestValidPrimitiveTypeInProperty(t *testing.T) {
	mockContext.Logger = &logger.Logger{}
	resources := make(map[string]template.Resource)
	properties := map[string]interface{}{
		"Ec2KeyName": "SomeValue",
	}
	resources["cluster"] = createResourceWithNestedProperties("AWS::Nested3::Cluster", "Instances", properties)

	assert.True(t, validateResources(resources, &spec, deadProp, deadRes, specInconsistency, &mockContext), "This resource should be valid")
}

func TestLackOfPrimitiveTypeInProperty(t *testing.T) {
	mockContext.Logger = &logger.Logger{}
	resources := make(map[string]template.Resource)
	properties := map[string]interface{}{
		"SomeProperty": "SomeValue",
	}
	resources["cluster"] = createResourceWithNestedProperties("AWS::Nested3::Cluster", "Instances", properties)

	assert.False(t, validateResources(resources, &spec, deadProp, deadRes, specInconsistency, &mockContext), "This resource shouldn't be valid")
}

func TestLackOfPrimitiveTypeInPropertyNestedInProperty(t *testing.T) {
	mockContext.Logger = &logger.Logger{}
	resources := make(map[string]template.Resource)
	properties := map[string]interface{}{
		"CoreInstanceGroup": map[string]interface{}{
			"AutoScalingPolicy": "SomePolicy",
			"DummyProperty":     "DummyProperty",
		},
	}
	resources["cluster"] = createResourceWithNestedProperties("AWS::Nested1::Cluster", "Instances", properties)

	assert.False(t, validateResources(resources, &spec, deadProp, deadRes, specInconsistency, &mockContext), "This resource shouldn't be valid, it lacks required property")
}
func TestLackOfRequiredSubproperty(t *testing.T) {
	mockContext.Logger = &logger.Logger{}
	resources := make(map[string]template.Resource)
	properties := map[string]interface{}{
		"DummySubproperty": map[string]interface{}{
			"DummyPrimitiveProperty": "SomeValue",
		},
	}
	resources["cluster"] = createResourceWithNestedProperties("AWS::Nested1::Cluster", "Instances", properties)

	assert.False(t, validateResources(resources, &spec, deadProp, deadRes, specInconsistency, &mockContext), "This resource shouldn't be valid, required subproperty is missing")
}
func TestLackOfRequiredPrimitiveTypeInNonrequiredSubproperty(t *testing.T) {
	mockContext.Logger = &logger.Logger{}
	resources := make(map[string]template.Resource)
	properties := map[string]interface{}{
		"ETag": "SomeEtagValue",
	}
	resources["ApiGatewayResource"] = createResourceWithNestedProperties("AWS::Nested2::RestApi", "BodyS3Location", properties)

	assert.False(t, validateResources(resources, &spec, deadProp, deadRes, specInconsistency, &mockContext), "This resource shouldn't be valid, required primitive property in nonrequired subproperty is missing")
}

func TestLackOfRequiredPropertyInNonRequiredProperty(t *testing.T) {
	mockContext.Logger = &logger.Logger{}

	resources := make(map[string]template.Resource)
	properties := map[string]interface{}{
		"Location": map[string]interface{}{
			"DummyValue": "",
		},
	}
	resources["ExampleResource"] = createResourceWithNestedProperties("AWS::Nested4::Method", "Definition", properties)

	assert.True(t, validateResources(resources, &spec, deadProp, deadRes, specInconsistency, &mockContext), "This resource should be valid")
}

func TestLackOfRequiredNestedPrimitivePropertyInListItem(t *testing.T) {
	mockContext.Logger = &logger.Logger{}
	resources := make(map[string]template.Resource)
	properties := []interface{}{
		0: map[string]interface{}{
			"DummyProperty": "SomeValue",
		},
		1: map[string]interface{}{
			"DummyProperty": "SomeValue",
		},
	}
	resource := template.Resource{}
	resource.Type = "AWS::List1::Cluster"
	resource.Properties = make(map[string]interface{})
	resource.Properties["BootstrapActions"] = properties
	resources["ExampleResource"] = resource

	assert.False(t, validateResources(resources, &spec, deadProp, deadRes, specInconsistency, &mockContext), "This resource should not be valid, List is empty")
}

func TestLackOfRequiredListItemSubpropertyInList(t *testing.T) {
	mockContext.Logger = &logger.Logger{}
	resources := make(map[string]template.Resource)
	properties := map[string]interface{}{
		"RoutingRules": []interface{}{
			0: map[string]interface{}{
				"SomeSubproperty": map[string]interface{}{
					"SomeDummyValue": "dummy1.example.com",
				},
			},
		},
	}
	resources["ExampleResource"] = createResourceWithNestedProperties("AWS::List2::Bucket", "WebsiteConfiguration", properties)

	assert.False(t, validateResources(resources, &spec, deadProp, deadRes, specInconsistency, &mockContext), "This resource should not be valid, It must contain RedirectRule property")
}

func TestLackOfRequiredPrimitiveTypeListItemInList(t *testing.T) {
	mockContext.Logger = &logger.Logger{}
	resources := make(map[string]template.Resource)
	properties := map[string]interface{}{
		"RoutingRules": []interface{}{
			0: map[string]interface{}{
				"RedirectRule": map[string]interface{}{
					"SomeProperty": "SomeValue1",
					"HostName":     "dummy1.example.com",
				},
			},
			1: map[string]interface{}{
				"RedirectRule": map[string]interface{}{
					"HttpRedirectCode": "SomeValue1",
					"HostName":         "dummy2.example.com",
				},
			},
		},
	}
	resources["ExampleResource"] = createResourceWithNestedProperties("AWS::List2::Bucket", "WebsiteConfiguration", properties)

	assert.False(t, validateResources(resources, &spec, deadProp, deadRes, specInconsistency, &mockContext), "This resource should not be valid, RedirectRule must contain HostName and HttpRedirectCode")
}

func TestValidRequiredPrimitiveTypeListItemInList(t *testing.T) {
	mockContext.Logger = &logger.Logger{}
	resources := make(map[string]template.Resource)
	properties := map[string]interface{}{
		"RoutingRules": []interface{}{
			0: map[string]interface{}{
				"RedirectRule": map[string]interface{}{
					"HostName":         "dummy1.example.com",
					"HttpRedirectCode": "SomeValue1",
				},
			},
		},
	}
	resources["ExampleResource"] = createResourceWithNestedProperties("AWS::List2::Bucket", "WebsiteConfiguration", properties)

	assert.True(t, validateResources(resources, &spec, deadProp, deadRes, specInconsistency, &mockContext), "This resource should be valid")
}

func TestLackOfNonRequiredNestedListItemProperty(t *testing.T) {
	mockContext.Logger = &logger.Logger{}

	resources := make(map[string]template.Resource)
	properties := map[string]interface{}{
		"Rules": []interface{}{
			0: map[string]interface{}{
				"Id":               "SomeValue",
				"Status":           "Enabled",
				"ExpirationInDays": 60,
			},
		},
	}
	resources["ExampleResource"] = createResourceWithNestedProperties("AWS::List3::Bucket", "LifecycleConfiguration", properties)

	assert.True(t, validateResources(resources, &spec, deadProp, deadRes, specInconsistency, &mockContext), "This resource should be valid")
}

func TestInvalidList(t *testing.T) {
	mockContext.Logger = &logger.Logger{}

	resources := make(map[string]template.Resource)
	properties := "DummyValue"

	resources["ExampleResource"] = createResourceWithOneProperty("AWS::List4::DBSubnetGroup", "SubnetIds", properties)

	assert.False(t, validateResources(resources, &spec, deadProp, deadRes, specInconsistency, &mockContext), "This resource should be valid")
}

func TestValidList(t *testing.T) {
	mockContext.Logger = &logger.Logger{}

	resources := make(map[string]template.Resource)

	resource := template.Resource{}
	resource.Type = "AWS::List4::DBSubnetGroup"
	resource.Properties = make(map[string]interface{})
	resource.Properties["SubnetIds"] = []interface{}{
		0: "subnet-33333",
	}
	resources["ExampleResource"] = resource

	assert.True(t, validateResources(resources, &spec, deadProp, deadRes, specInconsistency, &mockContext), "This resource should be valid")
}

func TestValidIfMapInNestedPropertyIsMap(t *testing.T) {
	mockContext.Logger = &logger.Logger{}

	resources := make(map[string]template.Resource)
	properties := map[string]interface{}{
		"Attributes": map[string]interface{}{
			"MapProperty1": "MapValue1",
		},
	}
	resources["ExampleResource"] = createResourceWithNestedProperties("AWS::Map2::Thing", "AttributePayload", properties)

	assert.True(t, validateResources(resources, &spec, deadProp, deadRes, specInconsistency, &mockContext), "This resource should be valid")
}

func TestInvalidNestedNonMapProperty(t *testing.T) {
	mockContext.Logger = &logger.Logger{}

	resources := make(map[string]template.Resource)
	properties := map[string]interface{}{
		"Attributes": "DummyValue",
	}
	resources["ExampleResource"] = createResourceWithNestedProperties("AWS::Map2::Thing", "AttributePayload", properties)

	assert.False(t, validateResources(resources, &spec, deadProp, deadRes, specInconsistency, &mockContext), "This resource shouldn't be valid - Attributes should be a Map")
}

func TestValidMapProperty(t *testing.T) {
	mockContext.Logger = &logger.Logger{}

	resources := make(map[string]template.Resource)
	resource := template.Resource{}
	resource.Type = "AWS::Map3::DBParameterGroup"
	resource.Properties = make(map[string]interface{})
	resource.Properties["Parameters"] = map[string]interface{}{
		"general_log":     1,
		"long_query_time": 10,
		"slow_query_log":  1,
	}
	resource.Properties["Family"] = "mysql5.6"
	resources["ExampleResource"] = resource

	assert.True(t, validateResources(resources, &spec, deadProp, deadRes, specInconsistency, &mockContext), "This resource should be valid")
}

func TestInvalidMapProperty(t *testing.T) {
	mockContext.Logger = &logger.Logger{}

	resources := make(map[string]template.Resource)
	resource := template.Resource{}
	resource.Type = "AWS::Map3::DBParameterGroup"
	resource.Properties = make(map[string]interface{})
	resource.Properties["Parameters"] = "DummyValue"
	resource.Properties["Family"] = "mysql5.6"
	resources["ExampleResource"] = resource

	assert.False(t, validateResources(resources, &spec, deadProp, deadRes, specInconsistency, &mockContext), "This resource should be valid")
}

func TestHasAllowedValuesParametersValid(t *testing.T) {
	sink = logger.Logger{}
	data := make(map[string]interface{})

	data["AllowedValues"] = ""
	data["Type"] = "String"
	parameters := createParameters("Correct", data)

	assert.True(t, hasAllowedValuesParametersValid(parameters, &sink), "This template has AllowedValues with Type String")
}

func TestHasAllowedValuesParametersInvalid(t *testing.T) {
	sink = logger.Logger{}
	data := make(map[string]interface{})

	data["AllowedValues"] = ""
	data["Type"] = "AWS::EC2::VPC::Id"
	parameters := createParameters("Incorrect", data)

	assert.False(t, hasAllowedValuesParametersValid(parameters, &sink), "This template has AllowedValues with Type other than String")
}

func createResourceWithNestedProperties(resourceType string, propertyName string, nestedPropertyValue map[string]interface{}) template.Resource {
	resource := template.Resource{}
	resource.Type = resourceType
	resource.Properties = make(map[string]interface{})
	resource.Properties[propertyName] = nestedPropertyValue

	return resource
}

func createResourceWithOneProperty(resourceType string, propertyName string, propertyValue string) template.Resource {
	resource := template.Resource{}
	resource.Type = resourceType
	resource.Properties = make(map[string]interface{})
	resource.Properties[propertyName] = propertyValue

	return resource
}

func createParameters(name string, value map[string]interface{}) map[string]interface{} {
	parameters := make(map[string]interface{})
	parameters[name] = value

	return parameters
}
