package cfofflinevalidator

import (
	"fmt"
	"encoding/json"
	"github.com/Appliscale/cftool/cfspecification"
	"io/ioutil"
)

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

func Validate(templatePath *string, specification *cfspecification.Specification) {

	rawTemplate, err := ioutil.ReadFile(*templatePath)
	if err != nil {
		fmt.Println(err)
	}

	template, err := parse(rawTemplate)
	if err != nil {
		fmt.Println(err)
	}

	valid := validateResources(template.Resources, specification)
	if !valid {
		fmt.Println("Template is invalid!")
	} else {
		fmt.Println("Template is valid!")
	}
}

func validateResources(resources map[string]Resource, specification *cfspecification.Specification) (bool) {
	valid := true
	for resourceName, resourceValue := range resources {
		if resourceSpecification, ok := specification.ResourceTypes[resourceValue.Type]; ok {
			if !areResourcePropertiesValid(resourceSpecification, resourceValue, resourceName) {
				valid = false
			}
		} else {
			fmt.Println("Type needs to be specified for resource " + resourceName)
			valid = false
		}
	}

	return valid
}
func areResourcePropertiesValid(resourceSpecification cfspecification.Resource, resourceValue Resource, resourceName string) bool {
	valid := true
	for propertyName, propertyValue := range resourceSpecification.Properties {
		if propertyValue.Required {
			if _, ok := resourceValue.Properties[propertyName]; !ok {
				fmt.Println("Property " + propertyName + " is required for resource " + resourceName)
				valid = false
			}
		}
	}
	return valid
}

func parse(templateFile []byte) (template Template, err error) {
	err = json.Unmarshal(templateFile, &template)
	if err != nil {
		return template, err
	}

	return template, nil
}