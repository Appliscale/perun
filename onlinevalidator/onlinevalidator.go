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
	"github.com/Appliscale/perun/mysession"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudformation"
	"io/ioutil"
)

// Validate template and get URL for cost estimation.
func ValidateAndEstimateCosts(context *context.Context) bool {
	valid := false
	defer printResult(&valid, context.Logger)

	if context.Config.DefaultDecisionForMFA {
		err := mysession.UpdateSessionToken(context.Config.DefaultProfile, context.Config.DefaultRegion, context.Config.DefaultDurationForMFA, context)
		if err != nil {
			context.Logger.Error(err.Error())
			return false
		}
	}

	session, err := mysession.CreateSession(context, context.Config.DefaultProfile, &context.Config.DefaultRegion)
	if err != nil {
		context.Logger.Error(err.Error())
		return false
	}

	rawTemplate, err := ioutil.ReadFile(*context.CliArguments.TemplatePath)
	if err != nil {
		context.Logger.Error(err.Error())
		return false
	}

	template := string(rawTemplate)
	valid, err = isTemplateValid(session, &template)
	if err != nil {
		context.Logger.Error(err.Error())
		return false
	}

	estimateCosts(session, &template, context.Logger)

	return valid
}

func isTemplateValid(session *session.Session, template *string) (bool, error) {
	api := cloudformation.New(session)
	templateStruct := cloudformation.ValidateTemplateInput{
		TemplateBody: template,
	}
	_, error := api.ValidateTemplate(&templateStruct)
	if error != nil {
		return false, error
	}

	return true, nil
}

func estimateCosts(session *session.Session, template *string, logger *logger.Logger) {
	api := cloudformation.New(session)
	templateCostInput := cloudformation.EstimateTemplateCostInput{
		TemplateBody: template,
	}
	output, err := api.EstimateTemplateCost(&templateCostInput)

	if err != nil {
		logger.Error(err.Error())
		return
	}

	logger.Info("Costs estimation: " + *output.Url)
}

func printResult(valid *bool, logger *logger.Logger) {
	if !*valid {
		logger.Error("Template is invalid!")
	} else {
		logger.Info("Template is valid!")
	}
}
