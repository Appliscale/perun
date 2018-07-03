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

// Package onlinevalidator provides tools for online CloudFormation template
// validation using AWS API.
package onlinevalidator

import (
	"github.com/Appliscale/perun/context"
	"github.com/Appliscale/perun/logger"
	"github.com/aws/aws-sdk-go/service/cloudformation"
	"io/ioutil"
)

// Validate template and get URL for cost estimation.
func ValidateAndEstimateCosts(context *context.Context) bool {
	valid := false
	defer printResult(&valid, context.Logger)

	context.InitializeAwsAPI()

	rawTemplate, err := ioutil.ReadFile(*context.CliArguments.TemplatePath)
	if err != nil {
		context.Logger.Error(err.Error())
		return false
	}

	template := string(rawTemplate)
	valid, err = isTemplateValid(context, &template)
	if err != nil {
		context.Logger.Error(err.Error())
		return false
	}

	estimateCosts(context, &template)

	return valid
}

func isTemplateValid(context *context.Context, template *string) (bool, error) {
	templateStruct := cloudformation.ValidateTemplateInput{
		TemplateBody: template,
	}
	_, err := context.CloudFormation.ValidateTemplate(&templateStruct)
	if err != nil {
		return false, err
	}

	return true, nil
}

func estimateCosts(context *context.Context, template *string) {
	templateCostInput := cloudformation.EstimateTemplateCostInput{
		TemplateBody: template,
	}
	output, err := context.CloudFormation.EstimateTemplateCost(&templateCostInput)

	if err != nil {
		context.Logger.Error(err.Error())
		return
	}

	context.Logger.Info("Costs estimation: " + *output.Url)
}

func printResult(valid *bool, logger *logger.Logger) {
	if !*valid {
		logger.Error("Template is invalid!")
	} else {
		logger.Info("Template is valid!")
	}
}
