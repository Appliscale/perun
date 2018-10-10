// Copyright 2018 Appliscale
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

// Package parameters provides tools for interactive creation of parameters file for aws
// cloud formation.
package parameters

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strings"

	"github.com/Appliscale/perun/context"
	"github.com/Appliscale/perun/helpers"
	"github.com/Appliscale/perun/validator/template"
	cloudformation2 "github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/awslabs/goformation/cloudformation"
)

// Parameter contains Key and Value. It describes how looks each parameter.
type Parameter struct {
	ParameterKey   string
	ParameterValue string
}

// GetJSONParameters gets parameters from context.CliArguments.TemplatePath.
func GetJSONParameters(context *context.Context) (resultString []byte, err error) {
	var parameters []*Parameter
	parameters, err = GetParameters(context)
	if err != nil {
		context.Logger.Error(err.Error())
		return
	}

	if *context.CliArguments.PrettyPrint {
		resultString, err = helpers.PrettyPrintJSON(parameters)
	} else {
		resultString, err = json.Marshal(parameters)
	}
	return
}

// ConfigureParameters allows to create file with parameters.
func ConfigureParameters(context *context.Context) error {
	resultString, err := GetJSONParameters(context)
	if err != nil {
		return err
	}
	if *context.CliArguments.OutputFilePath != "" {
		context.Logger.Info("Writing parameters configuration to file: " + *context.CliArguments.OutputFilePath)

		_, err = os.Stat(*context.CliArguments.OutputFilePath)
		if err == nil {
			context.Logger.Warning("File " + *context.CliArguments.OutputFilePath + " would be overriten by this action. Do you want to continue? [Y/N]")
			var answer string
			for answer != "n" && answer != "y" {
				fmt.Scanf("%s", &answer)
				answer = strings.ToLower(answer)
			}
			if answer == "n" {
				context.Logger.Info("Aborting..")
				return errors.New("user aborted")
			}
		}
		err = ioutil.WriteFile(*context.CliArguments.OutputFilePath, resultString, 0666)
		if err != nil {
			context.Logger.Error(err.Error())
		}
	} else {
		println(string(resultString))
	}
	return nil
}

// GetAwsParameters gets parameters from context and parses to AWS parameters.
func GetAwsParameters(context *context.Context) (parameters []*cloudformation2.Parameter, err error) {
	var params []*Parameter
	params, err = GetParameters(context)
	if err != nil {
		return
	}
	parameters = ParseParameterToAwsCompatible(params)
	return
}

// ParseParameterToAwsCompatible converts parameters from file to compatible with AWS.
func ParseParameterToAwsCompatible(params []*Parameter) (parameters []*cloudformation2.Parameter) {
	for paramnum := range params {
		parameters = append(parameters,
			&cloudformation2.Parameter{
				ParameterValue: &params[paramnum].ParameterValue,
				ParameterKey:   &params[paramnum].ParameterKey})
	}
	return
}

// Get the parameters - if parameters file provided - from file, else - interactively from user.
func ResolveParameters(context *context.Context) (params []*cloudformation2.Parameter, err error) {
	if *context.CliArguments.ParametersFile == "" {
		params, err = GetAwsParameters(context)
	} else {
		var parametersData []byte
		var readParameters []*Parameter
		parametersData, err = ioutil.ReadFile(*context.CliArguments.ParametersFile)
		if err != nil {
			return
		}
		err = json.Unmarshal(parametersData, &readParameters)
		if err != nil {
			return
		}
		params = ParseParameterToAwsCompatible(readParameters)
	}
	return
}

// GetParameters gets parameters from file, checks correctness and adds to Parameters.
func GetParameters(context *context.Context) (parameters []*Parameter, err error) {
	templateFile, err := parseTemplate(context)
	if err != nil {
		context.Logger.Error(err.Error())
		return nil, err
	}
	for parameterName, parameterSpec := range templateFile.Parameters {
		var parameterValid bool
		var parameterValue string
		if context.CliArguments.Parameters != nil {
			var exists bool
			parameterValue, exists = (*context.CliArguments.Parameters)[parameterName]
			if exists {
				parameterValid, err = checkParameterValid(parameterName, parameterSpec.(map[string]interface{}), parameterValue, context)
			}
		} else {
			parameterValid = false
		}
		for !parameterValid {
			print(parameterName, ": ")
			fmt.Scanf("%s", &parameterValue)
			parameterValid, err = checkParameterValid(parameterName, parameterSpec.(map[string]interface{}), parameterValue, context)
			if err != nil {
				context.Logger.Error(err.Error())
				return
			}
		}
		parameters = append(parameters, &Parameter{ParameterKey: parameterName, ParameterValue: parameterValue})
	}
	return
}

func checkParameterValid(parameterName string, parameterArgument map[string]interface{}, parameterValue string, context *context.Context) (bool, error) {
	if parameterArgument["AllowedValues"] != nil {
		allowedValues := getAllowedValues(parameterArgument)
		if !helpers.SliceContains(allowedValues, parameterValue) {
			context.Logger.Error("Value '" + parameterValue + "' is not allowed for Parameter " + parameterName + ". Value must be one of following: [" + strings.Join(allowedValues, ", ") + "]")
			return false, nil
		}
	}

	if parameterArgument["AllowedPattern"] != nil {
		allowedPattern := parameterArgument["AllowedPattern"].(string)
		matches, err := regexp.Match(fmt.Sprintf("^%s$", allowedPattern), []byte(parameterValue))
		if err != nil {
			return false, err
		}
		if !matches {
			context.Logger.Error("Value '" + parameterValue + "' does not match the required pattern: " + allowedPattern + " for Parameter " + parameterName)
			return false, nil
		}
	}
	return true, nil
}

func getAllowedValues(parameterArgument map[string]interface{}) (res []string) {
	list := parameterArgument["AllowedValues"].([]interface{})
	for _, val := range list {
		res = append(res, val.(string))
	}
	return
}

func parseTemplate(context *context.Context) (res cloudformation.Template, err error) {
	rawTemplate, err := ioutil.ReadFile(*context.CliArguments.TemplatePath)
	if err != nil {
		return
	}
	myTemplate := template.Template{}
	parser, err := helpers.GetParser(*context.CliArguments.TemplatePath)
	if err != nil {
		return
	}
	res, err = parser(rawTemplate, myTemplate, context.Logger)
	return
}
