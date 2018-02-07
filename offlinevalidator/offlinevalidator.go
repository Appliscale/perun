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

func printResult(valid *bool, logger *logger.Logger) {
	logger.PrintValidationErrors()
	if !*valid {
		logger.Error("Template is invalid!")
	} else {
		logger.Info("Template is valid!")
	}
}

func validateResources(resources map[string]template.Resource, specification *specification.Specification, sink *logger.Logger) bool {
	valid := true
	for resourceName, resourceValue := range resources {
		if resourceSpecification, ok := specification.ResourceTypes[resourceValue.Type]; ok {
			if !areRequiredPropertiesPresent(resourceSpecification, resourceValue, resourceName, sink) {
				valid = false
			}
		} else {
			sink.ValidationError(resourceName, "Type needs to be specified")
			valid = false
		}
		if validator, ok := validatorsMap[resourceValue.Type]; ok {
			if !validator.(func(string, template.Resource, *logger.Logger) bool)(resourceName, resourceValue, sink) {
				valid = false
			}
		}
	}

	return valid
}
func areRequiredPropertiesPresent(resourceSpecification specification.Resource, resourceValue template.Resource, resourceName string, logger *logger.Logger) bool {
	valid := true
	for propertyName, propertyValue := range resourceSpecification.Properties {
		if propertyValue.Required {
			if _, ok := resourceValue.Properties[propertyName]; !ok {
				logger.ValidationError(resourceName, "Property "+propertyName+" is required")
				valid = false
			}
		}
	}
	return valid
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

	preprocessed, preprocessingError := intrinsicsolver.FixFunctions(templateFile, logger)
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
