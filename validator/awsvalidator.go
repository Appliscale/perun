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

// Package awsvalidator provides tools for online CloudFormation template
// validation using AWS API.
package validator

import (
	"github.com/Appliscale/perun/context"
	"github.com/aws/aws-sdk-go/service/cloudformation"
	"strings"
)

func awsValidate(ctx *context.Context, templateBody *string) bool {
	valid, err := isTemplateValid(ctx, templateBody)
	if err != nil {
		ctx.Logger.Error(err.Error())
		return false
	}
	return valid
}

func isTemplateValid(context *context.Context, template *string) (bool, error) {
	templateStruct := cloudformation.ValidateTemplateInput{
		TemplateBody: template,
	}
	_, err := context.CloudFormation.ValidateTemplate(&templateStruct)
	if err != nil {
		if strings.Contains(err.Error(), "ExpiredToken:") {
			context.Logger.Error(err.Error())
			return true, nil
		}
		return false, err
	}
	return true, nil
}
