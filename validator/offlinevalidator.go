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

// Package validator provides tools for offline CloudFormation template
// validation.
package validator

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"reflect"
	"strconv"
	"strings"

	"errors"

	"github.com/Appliscale/perun/configuration"
	"github.com/Appliscale/perun/context"
	"github.com/Appliscale/perun/helpers"
	"github.com/Appliscale/perun/logger"
	"github.com/Appliscale/perun/specification"
	"github.com/Appliscale/perun/validator/template"
	"github.com/Appliscale/perun/validator/validators"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/awslabs/goformation"
	"github.com/awslabs/goformation/cloudformation"
	"github.com/mitchellh/mapstructure"
)

var validatorsMap = map[string]interface{}{
	"AWS::EC2::VPC": validators.IsVpcValid,
}

func printResult(templateName string, valid *bool, logger logger.LoggerInt) {
	logger.PrintValidationErrors()
	if !*valid {
		logger.Error(fmt.Sprintf("Template %s is invalid!", templateName))
	} else {
		logger.Info(fmt.Sprintf("Template %s is valid!", templateName))
	}
}

// ValidateAndEstimateCost CloudFormation template.
func ValidateAndEstimateCost(ctx *context.Context) bool {
	return validateTemplateFile(*ctx.CliArguments.TemplatePath, *ctx.CliArguments.TemplatePath, ctx)
}

func validateTemplateFile(templatePath string, templateName string, context *context.Context) (valid bool) {
	valid = false
	defer printResult(templateName, &valid, context.Logger)

	resourceSpecification, err := specification.GetSpecification(context)

	if err != nil {
		context.Logger.Error(err.Error())
		return
	}

	rawTemplate, err := ioutil.ReadFile(templatePath)
	if err != nil {
		context.Logger.Error(err.Error())
		return
	}

	var perunTemplate template.Template
	var goFormationTemplate cloudformation.Template

	parser, err := helpers.GetParser(*context.CliArguments.TemplatePath)
	if err != nil {
		context.Logger.Error(err.Error())
		return
	}
	goFormationTemplate, err = parser(rawTemplate, perunTemplate, context.Logger)
	if err != nil {
		context.Logger.Error(err.Error())
		return
	}

	deNilizedTemplate, _ := nilNeutralize(goFormationTemplate, context.Logger)
	resources := obtainResources(deNilizedTemplate, perunTemplate, context.Logger)
	deadResources := getNilResources(resources)
	deadProperties := getNilProperties(resources)
	if hasAllowedValuesParametersValid(goFormationTemplate.Parameters, context.Logger) {
		valid = true
	} else {
		valid = false
		context.Logger.AddResourceForValidation("Parameters").AddValidationError("Allowed Values supports only Type String")
	}

	specInconsistency := context.InconsistencyConfig.SpecificationInconsistency

	templateBody := string(rawTemplate)
	valid = validateResources(resources, &resourceSpecification, deadProperties, deadResources, specInconsistency, context) && valid
	valid = awsValidate(context, &templateBody) && valid

	if *context.CliArguments.EstimateCost {
		estimateCosts(context, &templateBody)
	}

	return valid
}

// Looking for AllowedValues and checking what Type is it. If it finds Type other than String then it will return false.
func hasAllowedValuesParametersValid(parameters template.Parameters, logger logger.LoggerInt) bool {
	isType := false
	isAllovedValues := false
	for _, value := range parameters {
		valueof := reflect.ValueOf(value)
		isAllovedValues = false
		isType = false

		for _, key := range valueof.MapKeys() {

			keyValue := valueof.MapIndex(key)
			textType := "Type"
			keyString := key.Interface().(string)
			textValues := "AllowedValues"

			if textType == keyString {
				textString := "String"
				keyValueString := keyValue.Interface().(string)
				if textString != keyValueString {
					isType = true
				}
			} else if textValues == keyString {
				isAllovedValues = true
			}

			if isAllovedValues && isType {
				return false
			}
		}
	}
	return true
}

func validateResources(resources map[string]template.Resource, specification *specification.Specification, deadProp []string, deadRes []string, specInconsistency map[string]configuration.Property, ctx *context.Context) bool {
	sink := ctx.Logger
	for resourceName, resourceValue := range resources {
		if deadResource := helpers.SliceContains(deadRes, resourceName); !deadResource {
			resourceValidation := sink.AddResourceForValidation(resourceName)
			processNestedTemplates(resourceValue.Properties, ctx)
			if resourceSpecification, ok := specification.ResourceTypes[resourceValue.Type]; ok {
				for propertyName, propertyValue := range resourceSpecification.Properties {
					if deadProperty := helpers.SliceContains(deadProp, propertyName); !deadProperty {
						validateProperties(specification, resourceValue, propertyName, propertyValue, resourceValidation, specInconsistency, sink)
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
	resourceValidation *logger.ResourceValidation,
	specInconsistency map[string]configuration.Property,
	logger logger.LoggerInt) {

	warnAboutSpecificationInconsistencies(propertyName, specInconsistency[resourceValue.Type], logger)
	if _, ok := resourceValue.Properties[propertyName]; !ok {
		if propertyValue.Required {
			resourceValidation.AddValidationError("Property " + propertyName + " is required")
		}
	} else if len(propertyValue.Type) > 0 {
		if propertyValue.Type != "List" && propertyValue.Type != "Map" {
			checkNestedProperties(specification, resourceValue.Properties, resourceValue.Type, propertyName, propertyValue.Type, resourceValidation, specInconsistency, logger)
		} else if propertyValue.Type == "List" {
			checkListProperties(specification, resourceValue.Properties, resourceValue.Type, propertyName, propertyValue.ItemType, resourceValidation, specInconsistency, logger)
		} else if propertyValue.Type == "Map" {
			checkMapProperties(resourceValue.Properties, propertyName, resourceValidation)
		}
	}
}

// check should be before validate, someone might add property because he thought it is required and here he would not get notified about inconsistency...
func warnAboutSpecificationInconsistencies(subpropertyName string, specInconsistentProperty configuration.Property, logger logger.LoggerInt) {
	if specInconsistentProperty[subpropertyName] != nil {
		for _, inconsistentPropertyName := range specInconsistentProperty[subpropertyName] {
			if inconsistentPropertyName == "Required" {
				logger.Warning(subpropertyName + "->" + inconsistentPropertyName + " in documentation is not consistent with specification")
			}
		}
	}
}

func checkListProperties(
	spec *specification.Specification,
	resourceProperties map[string]interface{},
	resourceValueType, propertyName, listItemType string,
	resourceValidation *logger.ResourceValidation,
	specInconsistency map[string]configuration.Property,
	logger logger.LoggerInt) {

	if listItemType == "" {
		resourceSubproperties := toStringList(resourceProperties, propertyName)
		if reflect.TypeOf(resourceSubproperties).Kind() != reflect.Slice || len(resourceSubproperties) == 0 {
			resourceValidation.AddValidationError(propertyName + " must be a List")
		}
	} else if propertySpec, hasSpec := spec.PropertyTypes[resourceValueType+"."+listItemType]; hasSpec {
		resourceSubproperties := toMapList(resourceProperties, propertyName)
		for subpropertyName, subpropertyValue := range propertySpec.Properties {
			for _, listItem := range resourceSubproperties {
				warnAboutSpecificationInconsistencies(subpropertyName, specInconsistency[resourceValueType+"."+listItemType], logger)
				if _, isPresent := listItem[subpropertyName]; !isPresent {
					if subpropertyValue.Required {
						resourceValidation.AddValidationError("Property " + subpropertyName + " is required in " + listItemType)
					}
				} else if isPresent {
					if subpropertyValue.IsSubproperty() {
						checkNestedProperties(spec, listItem, resourceValueType, subpropertyName, subpropertyValue.Type, resourceValidation, specInconsistency, logger)
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
	resourceValidation *logger.ResourceValidation,
	specInconsistency map[string]configuration.Property,
	logger logger.LoggerInt) {

	if propertySpec, hasSpec := spec.PropertyTypes[resourceValueType+"."+propertyType]; hasSpec {
		resourceSubproperties, _ := toMap(resourceProperties, propertyName)
		for subpropertyName, subpropertyValue := range propertySpec.Properties {
			warnAboutSpecificationInconsistencies(subpropertyName, specInconsistency[resourceValueType+"."+propertyName], logger)
			if _, isPresent := resourceSubproperties[subpropertyName]; !isPresent {
				if subpropertyValue.Required {
					resourceValidation.AddValidationError("Property " + subpropertyName + " is required " + "in " + propertyName)
				}
			} else if isPresent {
				if subpropertyValue.IsSubproperty() {
					checkNestedProperties(spec, resourceSubproperties, resourceValueType, subpropertyName, subpropertyValue.Type, resourceValidation, specInconsistency, logger)
				} else if subpropertyValue.Type == "List" {
					checkListProperties(spec, resourceSubproperties, resourceValueType, subpropertyName, subpropertyValue.ItemType, resourceValidation, specInconsistency, logger)
				} else if subpropertyValue.Type == "Map" {
					checkMapProperties(resourceSubproperties, subpropertyName, resourceValidation)
				}
			}
		}
	}
}

func processNestedTemplates(properties map[string]interface{}, ctx *context.Context) {
	if rawTemplateURL, ok := properties["TemplateURL"]; ok {
		if templateURL, ok := rawTemplateURL.(string); ok {
			err := validateNestedTemplate(templateURL, ctx)
			if err != nil {
				ctx.Logger.Error(err.Error())
				os.Exit(1)
			}
		}
	}
}

func validateNestedTemplate(templateURL string, ctx *context.Context) error {
	err := context.UpdateSessionToken(ctx.Config.DefaultProfile, ctx.Config.DefaultRegion, ctx.Config.DefaultDurationForMFA, ctx)
	if err != nil {
		return err
	}

	tempfile, err := ioutil.TempFile(ctx.Config.DefaultTemporaryFilesDirectory, "")
	if err != nil {
		return err
	}
	defer os.Remove(tempfile.Name())

	if err := downloadTemplateFromBucket(templateURL, tempfile, ctx); err != nil {
		return err
	}

	validateTemplateFile(tempfile.Name(), templateURL, ctx)

	if err = tempfile.Close(); err != nil {
		return err
	}

	return nil
}

func downloadTemplateFromBucket(templateURL string, file io.WriterAt, ctx *context.Context) error {
	region, bucket, key := fetchBucketDataFromURL(templateURL)

	session, err := context.CreateSession(ctx, ctx.Config.DefaultProfile, &region)
	if err != nil {
		return err
	}

	downloader := s3manager.NewDownloader(session)

	_, err = downloader.Download(file, &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return err
	}

	return nil
}

func fetchBucketDataFromURL(url string) (region string, bucket string, key string) {
	path := strings.SplitN(url, "/", 5)
	host := strings.Split(path[2], ".")

	region = host[1]
	bucket = path[3]
	key = path[4]
	return
}

func checkMapProperties(
	resourceProperties map[string]interface{},
	propertyName string,
	resourceValidation *logger.ResourceValidation) {

	_, err := toMap(resourceProperties, propertyName)
	if err != nil {
		resourceValidation.AddValidationError(err.Error())
	}
}

func obtainResources(goformationTemplate cloudformation.Template, perunTemplate template.Template, logger logger.LoggerInt) map[string]template.Resource {
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
func nilNeutralize(template cloudformation.Template, logger logger.LoggerInt) (output cloudformation.Template, err error) {
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

func checkWhereIsNil(n interface{}, v interface{}, baseLevel string, logger logger.LoggerInt, fullPath []interface{}, dsc *interface{}) {
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
