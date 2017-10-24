package cfofflinevalidator

import (
	"fmt"
	"encoding/json"
	"github.com/Appliscale/cftool/cfspecification"
	"io/ioutil"
)

type Template struct {
	Resources map[string]Resource
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
	}
}

func validateResources(resources map[string]Resource, specification *cfspecification.Specification) (bool) {
	returnValue := true
	for resourceName, resourceValue := range resources {
		resourceSpecification := specification.ResourceTypes[resourceValue.Type]
		for propertyName, propertyValue := range resourceSpecification.Properties {
			if propertyValue.Required {
				if _, ok := resourceValue.Properties[propertyName]; !ok {
					fmt.Println("Property " + propertyName + " is required for resource " + resourceName)
					returnValue = false
				}
			}
		}
	}

	return returnValue
}

func parse(templateFile []byte) (template Template, err error) {
	err = json.Unmarshal(templateFile, &template)
	if err != nil {
		return template, err
	}

	return template, nil
}