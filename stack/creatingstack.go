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

package stack

import (
	"github.com/Appliscale/perun/context"
	"github.com/Appliscale/perun/parameters"
	"github.com/Appliscale/perun/progress"
	"github.com/aws/aws-sdk-go/service/cloudformation"
)

// This function gets template and  name of stack. It creates "CreateStackInput" structure.
func createStackInput(template *string, stackName *string, context *context.Context) (cloudformation.CreateStackInput, error) {
	params, err := parameters.ResolveParameters(context)
	if err != nil {
		context.Logger.Error(err.Error())
		return cloudformation.CreateStackInput{}, err
	}

	templateStruct := cloudformation.CreateStackInput{
		Parameters:   params,
		TemplateBody: template,
		StackName:    stackName,
	}
	return templateStruct, nil
}

// NewStack create Stack. It's get template from context.CliArguments.TemplatePath.
func NewStack(context *context.Context) error {
	template, stackName, incorrectPath := getTemplateFromFile(context)
	if incorrectPath != nil {
		context.Logger.Error(incorrectPath.Error())
		return incorrectPath
	}
	templateStruct, templateError := createStackInput(&template, &stackName, context)
	if templateError != nil {
		context.Logger.Error(templateError.Error())
		return templateError
	}

	if *context.CliArguments.Progress {
		conn, remoteSinkError := progress.GetRemoteSink(context)
		if remoteSinkError != nil {
			context.Logger.Error("Error getting remote sink configuration: " + remoteSinkError.Error())
			return remoteSinkError
		}
		templateStruct.NotificationARNs = []*string{conn.TopicArn}
		_, creationError := context.CloudFormation.CreateStack(&templateStruct)
		if creationError != nil {
			context.Logger.Error("Error creating stack: " + creationError.Error())
			return creationError
		}
		conn.MonitorStackQueue()
	} else {
		_, creationError := context.CloudFormation.CreateStack(&templateStruct)
		if creationError != nil {
			context.Logger.Error("Error creating stack: " + creationError.Error())
			return creationError
		}
		context.Logger.Info("Stack creation request successful")
	}

	return nil
}
