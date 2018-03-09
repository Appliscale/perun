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
	"reflect"
	"strconv"
	"strings"

	"github.com/Appliscale/perun/context"
	"github.com/Appliscale/perun/helpers"
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

	parser, err := helpers.GetParser(*context.CliArguments.TemplatePath)
	if err != nil {
		context.Logger.Error(err.Error())
		return false
	}
	goFormationTemplate, err = parser(rawTemplate, perunTemplate, context.Logger)
	if err != nil {
		context.Logger.Error(err.Error())
		return false
	}

	deNilizedTemplate, _ := nilNeutralize(goFormationTemplate, context.Logger)
	resources := obtainResources(deNilizedTemplate, perunTemplate, context.Logger)
	deadResources := getNilResources(resources)
	deadProperties := getNilProperties(resources)

	valid = validateResources(resources, &specification, context.Logger, deadProperties, deadResources)
	return valid
}

func validateResources(resources map[string]template.Resource, specification *specification.Specification, sink *logger.Logger, deadProp []string, deadRes []string) bool {

	for resourceName, resourceValue := range resources {
		if deadResource := helpers.SliceContains(deadRes, resourceName); !deadResource {
			resourceValidation := sink.AddResourceForValidation(resourceName)

			if resourceSpecification, ok := specification.ResourceTypes[resourceValue.Type]; ok {
				for propertyName, propertyValue := range resourceSpecification.Properties {
					if deadProperty := helpers.SliceContains(deadProp, propertyName); !deadProperty {
						validateProperties(specification, resourceValue, propertyName, propertyValue, resourceValidation)
					}
				}
			} else {
				resourceValidation.AddValidationError("Type needs to be specified")
			}
			if validator, ok := validatorsMap[resourceValue.Type]; ok {
				validator.(func(template.Resource, *logger.ResourceValidation) bool)(resourceValue, resourceValidation)
			}

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
			checkMapProperties(resourceValue.Properties, propertyName, resourceValidation)
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
						checkMapProperties(listItem, propertyName, resourceValidation)
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
		resourceSubproperties, _ := toMap(resourceProperties, propertyName)
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
					checkMapProperties(resourceSubproperties, subpropertyName, resourceValidation)
				}
			}
		}
	}
}

func checkMapProperties(
	resourceProperties map[string]interface{},
	propertyName string,
	resourceValidation *logger.ResourceValidation) {

	_, err := toMap(resourceProperties, propertyName)
	if err != nil {
		resourceValidation.AddValidationError(err.Error())
	}
	for subpropertyName, subpropertyValue := range resourceProperties {
		if reflect.TypeOf(subpropertyValue).Kind() != reflect.Map {
			resourceValidation.AddValidationError(subpropertyName + " must be a Map")
		}
	}
}

func ParseJSON(templateFile []byte, refTemplate template.Template, logger *logger.Logger) (template cloudformation.Template, err error) {

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

func ParseYAML(templateFile []byte, refTemplate template.Template, logger *logger.Logger) (template cloudformation.Template, err error) {

	err = yaml.Unmarshal(templateFile, &refTemplate)
	if err != nil {
		return template, err
	}

	preprocessed, preprocessingError := intrinsicsolver.FixFunctions(templateFile, logger, "multiline", "elongate", "correctlong")
	if preprocessingError != nil {
		logger.Error(preprocessingError.Error())
	}
	tempYAML, err := goformation.ParseYAML(preprocessed)
	return *tempYAML, err
}

func obtainResources(goformationTemplate cloudformation.Template, perunTemplate template.Template, logger *logger.Logger) map[string]template.Resource {
	perunResources := perunTemplate.Resources
	goformationResources := goformationTemplate.Resources

	mapstructure.Decode(goformationResources, &perunResources)

	for propertyName, propertyContent := range perunResources {
		if propertyContent.Properties == nil {
			logger.Warning(propertyName + " <--- is nil.")
		} else {
			for element, elementValue := range propertyContent.Properties {
				initPath := []interface{}{element} // The path from the Property name to the <nil> element.
				var discarded interface{}          // Container which stores the encountered nodes that aren't on the path.
				checkWhereIsNil(element, elementValue, propertyName, logger, initPath, &discarded)
			}
		}
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
		if _, ok := value.(map[string]interface{}); ok {
			mapList[index] = value.(map[string]interface{})
		}
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
		if value != nil {
			list[index] = value.(string)
		}
	}
	return list
}

func toMap(resourceProperties map[string]interface{}, propertyName string) (map[string]interface{}, error) {
	subproperties, ok := resourceProperties[propertyName].(map[string]interface{})
	if !ok {
		return nil, errors.New(propertyName + " must be a Map")
	}
	return subproperties, nil
}

// There is a possibility that a hash map inside the template would have one of it's element's being an intrinsic function designed to output `key : value` pair.
// If this function would be unresolved, it would output a standalone <nil> of type interface{}. It would be an alien element in a hash map.
// To prevent the parser from breaking, we wipe out the entire, expected hash map element.
func nilNeutralize(template cloudformation.Template, logger *logger.Logger) (output cloudformation.Template, err error) {
	bytes, initErr := json.Marshal(template)
	if initErr != nil {
		logger.Error(err.Error())
	}
	byteSlice := string(bytes)

	var info int
	var check1, check2, check3 string
	if strings.Contains(byteSlice, ",null,") {
		check1 = strings.Replace(byteSlice, ",null,", ",", -1)
		info++
	} else {
		check1 = byteSlice
	}
	if strings.Contains(check1, "[null,") {
		check2 = strings.Replace(check1, "[null,", "[", -1)
		info++
	} else {
		check2 = check1
	}
	if strings.Contains(check2, ",null]") {
		check3 = strings.Replace(check2, ",null]", "]", -1)
		info++
	} else {
		check3 = check2
	}

	byteSliceCorrected := []byte(check3)

	tempJSON, err := goformation.ParseJSON(byteSliceCorrected)
	if err != nil {
		logger.Error(err.Error())
	}

	infoOpening, link, part, occurences, elements, a, t := "", "", "", "", "", "", ""
	if info > 0 {
		if info == 1 {
			elements = "element"
			t = "this "
			a = "a"
			infoOpening = "is an intrinsic function "
			link = "is"
			part = "part"
		} else {
			elements = "elements"
			t = "those "
			occurences = strconv.Itoa(info)
			infoOpening = "are " + occurences + " intrinsic functions "
			link = "are"
			part = "parts"
		}
		logger.Info("There " + infoOpening + "which would output `key : value` pair but " + link + " unresolved and " + link + " evaluated to <nil>. As " + t + elements + " of a template should be " + a + " hash table " + elements + ", " + t + "standalone <nil> " + link + " deleted completely. It is recommended to investigate " + t + part + " of a template manually.")
	}

	returnTemplate := *tempJSON

	return returnTemplate, nil
}

func getNilProperties(resources map[string]template.Resource) []string {
	list := make([]string, 0)
	for _, resourceContent := range resources {
		properties := resourceContent.Properties
		for propertyName, propertyContent := range properties {
			if propertyContent == nil {
				list = append(list, propertyName)
			}
		}
	}
	return list
}

func getNilResources(resources map[string]template.Resource) []string {
	list := make([]string, 0)
	for resourceName, resourceContent := range resources {
		if resourceContent.Properties == nil {
			list = append(list, resourceName)
		}

	}
	return list
}

func checkWhereIsNil(n interface{}, v interface{}, baseLevel string, logger *logger.Logger, fullPath []interface{}, dsc *interface{}) {
	if v == nil { // Value we encountered is nil - this is the end of investigation.
		where := ""
		for _, element := range fullPath {
			if stringElement, ok := element.(string); ok {
				if where != "" {
					where += ": " + stringElement
				} else {
					where = stringElement
				}
			} else if intElement, ok := element.(int); ok {
				where += "[" + strconv.Itoa(intElement) + "]"
			}
		}
		logger.Warning(baseLevel + ": " + where + " <--- is nil.")
	} else if mp, ok := v.(map[string]interface{}); ok { // Value we encountered is a map.
		if helpers.IsPlainMap(mp) { // Check is it plain, non-nil map.
			// It is - we shouldn't dive into.
			*dsc = n // The name is stored in the `discarded` container as the name of the blind alley.
		} else {
			for kmp, vmp := range mp {
				if helpers.IsNonStringFloatBool(vmp) {
					fullPath = append(fullPath, kmp)
					fullPath = helpers.Discard(fullPath, *dsc) // If the output path would be different, it seems that we've encountered some node which is not on the way to the <nil>. It will be discarded from the path. Otherwise the paths are the same and we hit the point.
					checkWhereIsNil(kmp, vmp, baseLevel, logger, fullPath, dsc)
				}
			}
		}
	} else if slc, ok := v.([]interface{}); ok { // The same flow as above.
		if helpers.IsPlainSlice(slc) {
			*dsc = n
		} else {
			for islc, vslc := range slc {
				if helpers.IsNonStringFloatBool(vslc) {
					fullPath = append(fullPath, islc)
					fullPath = helpers.Discard(fullPath, *dsc)
					checkWhereIsNil(islc, vslc, baseLevel, logger, fullPath, dsc)
				}
			}
		}
	}
}
