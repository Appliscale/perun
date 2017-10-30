package cfofflinevalidator

import (
	"fmt"
	"encoding/json"
	"github.com/Appliscale/cftool/cfspecification"
	"io/ioutil"
	"path"
	"errors"
	"github.com/ghodss/yaml"
	"github.com/Appliscale/cftool/cfofflinevalidator/cftemplate"
	"github.com/Appliscale/cftool/cfofflinevalidator/cfvalidators"
	"github.com/Appliscale/cftool/cflogger"
)

var validators = map[string]interface{}{
	"AWS::EC2::VPC": cfvalidators.IsVpcValid,
}

func Validate(templatePath *string, configurationPath *string) {
	valid := false
	logger := cflogger.Logger{}
	defer printResult(&valid, &logger)

	specification, err := cfspecification.GetSpecification(*configurationPath)
	if err != nil {
		cflogger.LogError(&logger, err.Error())
		return
	}

	rawTemplate, err := ioutil.ReadFile(*templatePath)
	if err != nil {
		cflogger.LogError(&logger, err.Error())
		return
	}

	var template cftemplate.Template
	templateFileExtension := path.Ext(*templatePath)
	if templateFileExtension == ".json" {
		template, err = parseJSON(rawTemplate)
	} else if templateFileExtension == ".yaml" ||  templateFileExtension == ".yml" {
		template, err = parseYAML(rawTemplate)
	} else {
		err = errors.New("Invalid template file format.")
	}
	if err != nil {
		cflogger.LogError(&logger, err.Error())
		return
	}

	valid = validateResources(template.Resources, &specification, &logger)
}

func printResult(valid *bool, logger *cflogger.Logger) {
	cflogger.PrintErrors(logger)
	if !*valid {
		fmt.Println("Template is invalid!")
	} else {
		fmt.Println("Template is valid!")
	}
}

func validateResources(resources map[string]cftemplate.Resource, specification *cfspecification.Specification, logger *cflogger.Logger) (bool) {
	valid := true
	for resourceName, resourceValue := range resources {
		if resourceSpecification, ok := specification.ResourceTypes[resourceValue.Type]; ok {
			if !areRequiredPropertiesPresent(resourceSpecification, resourceValue, resourceName, logger) {
				valid = false
			}
		} else {
			cflogger.LogValidationError(logger, resourceName, "Type needs to be specified")
			valid = false
		}
		if validator, ok := validators[resourceValue.Type]; ok {
			if !validator.(func(string, cftemplate.Resource, *cflogger.Logger)(bool))(resourceName, resourceValue, logger) {
				valid = false
			}
		}
	}

	return valid
}
func areRequiredPropertiesPresent(resourceSpecification cfspecification.Resource, resourceValue cftemplate.Resource, resourceName string, logger *cflogger.Logger) bool {
	valid := true
	for propertyName, propertyValue := range resourceSpecification.Properties {
		if propertyValue.Required {
			if _, ok := resourceValue.Properties[propertyName]; !ok {
				cflogger.LogValidationError(logger, resourceName, "Property " + propertyName + " is required")
				valid = false
			}
		}
	}
	return valid
}

func parseJSON(templateFile []byte) (template cftemplate.Template, err error) {
	err = json.Unmarshal(templateFile, &template)
	if err != nil {
		return template, err
	}

	return template, nil
}

func parseYAML(templateFile []byte) (template cftemplate.Template, err error) {
	err = yaml.Unmarshal(templateFile, &template)
	if err != nil {
		return template, err
	}

	return template, nil
}