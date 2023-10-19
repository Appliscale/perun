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

package estimatecost

import (
	"github.com/Appliscale/perun/context"
	"github.com/Appliscale/perun/parameters"
	"github.com/aws/aws-sdk-go/service/cloudformation"
	"io/ioutil"
)

func estimateCosts(context *context.Context, template *string) (err error) {
	templateParameters, err := parameters.ResolveParameters(context)
	if err != nil {
		context.Logger.Error(err.Error())
		return
	}
	templateCostInput := cloudformation.EstimateTemplateCostInput{
		TemplateBody: template,
		Parameters:   templateParameters,
	}
	output, err := context.CloudFormation.EstimateTemplateCost(&templateCostInput)
	if err != nil {
		context.Logger.Error(err.Error())
		return
	}
	context.Logger.Info("Costs estimation: " + *output.Url)
	return
}

func EstimateCosts(ctx *context.Context) {
	templatePath := *ctx.CliArguments.TemplatePath
	rawTemplate, err := ioutil.ReadFile(templatePath)
	if err != nil {
		ctx.Logger.Error(err.Error())
		return
	}
	templateBody := string(rawTemplate)
	estimateError := estimateCosts(ctx, &templateBody)
	if estimateError != nil {
		ctx.Logger.Error(estimateError.Error())
	}
}
