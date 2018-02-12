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

// Package offlinevalidator provides tools for offline CloudFormation template
// validation.
package offlinevalidator

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"path"
	"reflect"

	"github.com/Appliscale/perun/context"
	"github.com/Appliscale/perun/intrinsicsolver"
	"github.com/Appliscale/perun/logger"
	"github.com/Appliscale/perun/offlinevalidator/template"
	"github.com/Appliscale/perun/offlinevalidator/validators"
	"github.com/Appliscale/perun/specification"
	"github.com/awslabs/goformation"
	"github.com/awslabs/goformation/cloudformation"
	"github.com/ghodss/yaml"
	"github.com/mitchellh/mapstructure"
)

var validatorsMap = map[string]interface{}{
	"AWS::EC2::VPC": validators.IsVpcValid,
}

func printResult(valid *bool, logger *logger.Logger) {
	logger.PrintValidationErrors()
	if !*valid {
		logger.Error("Template is invalid!")
	} else {
		logger.Info("Template is valid!")
	}
}

// Validate CloudFormation template.
func Validate(context *context.Context) bool {
	valid := false
	defer printResult(&valid, context.Logger)

	specification, err := specification.GetSpecification(context)

	if err != nil {
		context.Logger.Error(err.Error())
		return false
	}

	rawTemplate, err := ioutil.ReadFile(*context.CliArguments.TemplatePath)
	if err != nil {
		context.Logger.Error(err.Error())
		return false
	}

	var perunTemplate template.Template
	var goFormationTemplate cloudformation.Template

	templateFileExtension := path.Ext(*context.CliArguments.TemplatePath)
	if templateFileExtension == ".json" {
		goFormationTemplate, err = parseJSON(rawTemplate, perunTemplate, context.Logger)
	} else if templateFileExtension == ".yaml" || templateFileExtension == ".yml" {
		goFormationTemplate, err = parseYAML(rawTemplate, perunTemplate, context.Logger)
	} else {
		err = errors.New("Invalid template file format.")
	}

	if err != nil {
		context.Logger.Error(err.Error())
		return false
	}

	resources := obtainResources(goFormationTemplate, perunTemplate)

	valid = validateResources(resources, &specification, context.Logger)
	return valid
}

func validateResources(resources map[string]template.Resource, specification *specification.Specification, sink *logger.Logger) bool {

	for resourceName, resourceValue := range resources {
		resourceValidation := sink.AddResourceForValidation(resourceName)

		if resourceSpecification, ok := specification.ResourceTypes[resourceValue.Type]; ok {
			for propertyName, propertyValue := range resourceSpecification.Properties {
				validateProperties(specification, resourceValue, propertyName, propertyValue, resourceValidation)
			}
		} else {
			resourceValidation.AddValidationError("Type needs to be specified")
		}
		if validator, ok := validatorsMap[resourceValue.Type]; ok {
			validator.(func(template.Resource, *logger.ResourceValidation) bool)(resourceValue, resourceValidation)
		}

	}
	return !sink.HasValidationErrors()
}

func validateProperties(
	specification *specification.Specification,
	resourceValue template.Resource,
	propertyName string,
	propertyValue specification.Property,
	resourceValidation *logger.ResourceValidation) {

	if _, ok := resourceValue.Properties[propertyName]; !ok {
		if propertyValue.Required {
			resourceValidation.AddValidationError("Property " + propertyName + " is required")
		}
	} else if len(propertyValue.Type) > 0 {
		if propertyValue.Type != "List" && propertyValue.Type != "Map" {
			checkNestedProperties(specification, resourceValue.Properties, resourceValue.Type, propertyName, propertyValue.Type, resourceValidation)
		} else if propertyValue.Type == "List" {
			checkListProperties(specification, resourceValue.Properties, resourceValue.Type, propertyName, propertyValue.ItemType, resourceValidation)
		} else if propertyValue.Type == "Map" {
			checkMapProperties(resourceValue.Properties, resourceValidation)
		}
	}
}

func checkListProperties(
	spec *specification.Specification,
	resourceProperties map[string]interface{},
	resourceValueType, propertyName, listItemType string,
	resourceValidation *logger.ResourceValidation) {

	if listItemType == "" {
		resourceSubproperties := toStringList(resourceProperties, propertyName)
		if reflect.TypeOf(resourceSubproperties).Kind() != reflect.Slice || len(resourceSubproperties) == 0 {
			resourceValidation.AddValidationError(propertyName + " must be a List")
		}
	} else if propertySpec, hasSpec := spec.PropertyTypes[resourceValueType+"."+listItemType]; hasSpec {

		resourceSubproperties := toMapList(resourceProperties, propertyName)
		for subpropertyName, subpropertyValue := range propertySpec.Properties {
			for _, listItem := range resourceSubproperties {
				if _, isPresent := listItem[subpropertyName]; !isPresent {
					if subpropertyValue.Required {
						resourceValidation.AddValidationError("Property " + subpropertyName + " is required in " + listItemType)
					}
				} else if isPresent {
					if subpropertyValue.IsSubproperty() {
						checkNestedProperties(spec, listItem, resourceValueType, subpropertyName, subpropertyValue.Type, resourceValidation)
					} else if subpropertyValue.Type == "Map" {
						checkMapProperties(listItem, resourceValidation)
					}
				}
			}
		}
	}
}

func checkNestedProperties(
	spec *specification.Specification,
	resourceProperties map[string]interface{},
	resourceValueType, propertyName, propertyType string,
	resourceValidation *logger.ResourceValidation) {

	if propertySpec, hasSpec := spec.PropertyTypes[resourceValueType+"."+propertyType]; hasSpec {
		resourceSubproperties := toMap(resourceProperties, propertyName)
		for subpropertyName, subpropertyValue := range propertySpec.Properties {
			if _, isPresent := resourceSubproperties[subpropertyName]; !isPresent {
				if subpropertyValue.Required {
					resourceValidation.AddValidationError("Property " + subpropertyName + " is required" + "in " + propertyName)
				}
			} else if isPresent {
				if subpropertyValue.IsSubproperty() {
					checkNestedProperties(spec, resourceSubproperties, resourceValueType, subpropertyName, subpropertyValue.Type, resourceValidation)
				} else if subpropertyValue.Type == "List" {
					checkListProperties(spec, resourceSubproperties, resourceValueType, subpropertyName, subpropertyValue.ItemType, resourceValidation)
				} else if subpropertyValue.Type == "Map" {
					checkMapProperties(resourceSubproperties, resourceValidation)
				}
			}
		}
	}
}

func checkMapProperties(
	resourceProperties map[string]interface{},
	resourceValidation *logger.ResourceValidation) {

	for subpropertyName, subpropertyValue := range resourceProperties {
		if reflect.TypeOf(subpropertyValue).Kind() != reflect.Map {
			resourceValidation.AddValidationError(subpropertyName + " must be a Map")
		}
	}
}

func parseJSON(templateFile []byte, refTemplate template.Template, logger *logger.Logger) (template cloudformation.Template, err error) {

	err = json.Unmarshal(templateFile, &refTemplate)
	if err != nil {
		return template, err
	}

	tempJSON, err := goformation.ParseJSON(templateFile)
	if err != nil {
		logger.Error(err.Error())
	}

	returnTemplate := *tempJSON

	return returnTemplate, nil
}

func parseYAML(templateFile []byte, refTemplate template.Template, logger *logger.Logger) (template cloudformation.Template, err error) {

	err = yaml.Unmarshal(templateFile, &refTemplate)
	if err != nil {
		return template, err
	}

	preprocessed, preprocessingError := intrinsicsolver.FixFunctions(templateFile, logger, "multiline")
	if preprocessingError != nil {
		logger.Error(preprocessingError.Error())
	}
	tempYAML, err := goformation.ParseYAML(preprocessed)
	if err != nil {
		logger.Error(err.Error())
	}

	returnTemplate := *tempYAML

	return returnTemplate, nil
}

func obtainResources(goformationTemplate cloudformation.Template, perunTemplate template.Template) map[string]template.Resource {
	perunResources := perunTemplate.Resources
	goformationResources := goformationTemplate.Resources

	errDecode := mapstructure.Decode(goformationResources, &perunResources)
	if errDecode != nil {
		/*
			Printing errDecode would log:

			ERROR error(s) decoding:
			[template.Resource name] expected a map, got 'bool'

			whenever a value of a property would be a boolean value (e.g. evaluated by !Equals intrinsic function; or e.g. 'got string', 'got float' etc. in other options).
			But after logging all the decoding errors, it would log if template is valid or not and eventually log the missing property as it should do
			and the error doesn't stand as obstacle of validation.
		*/
	}

	return perunResources
}

func toMapList(resourceProperties map[string]interface{}, propertyName string) []map[string]interface{} {
	subproperties, ok := resourceProperties[propertyName].([]interface{})
	if !ok {
		return []map[string]interface{}{}
	}
	mapList := make([]map[string]interface{}, len(subproperties))
	for index, value := range subproperties {
		mapList[index] = value.(map[string]interface{})
	}
	return mapList
}

func toStringList(resourceProperties map[string]interface{}, propertyName string) []string {
	subproperties, ok := resourceProperties[propertyName].([]interface{})
	if !ok {
		return nil
	}
	list := make([]string, len(subproperties))
	for index, value := range subproperties {
		list[index] = value.(string)
	}
	return list
}

func toMap(resourceProperties map[string]interface{}, propertyName string) map[string]interface{} {
	subproperties, ok := resourceProperties[propertyName].(map[string]interface{})
	if !ok {
		return map[string]interface{}{}
	}
	return subproperties
}
