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

// UpdateStack prepares updateStackInput and updates stack.
func UpdateStack(context *context.Context) (err error) {
	template, stackName, err := getTemplateFromFile(context)
	if err != nil {
		return
	}
	templateStruct := updateStackInput(context, &template, &stackName)
	err = doUpdateStack(context, templateStruct)
	return
}

// doUpdateStack updates stack.
func doUpdateStack(context *context.Context, updateStackInput cloudformation.UpdateStackInput) error {
	if *context.CliArguments.Progress {
		conn, remoteSinkError := progress.GetRemoteSink(context)
		if remoteSinkError != nil {
			context.Logger.Error("Error getting remote sink configuration: " + remoteSinkError.Error())
			return remoteSinkError
		}
		updateStackInput.NotificationARNs = []*string{conn.TopicArn}
		_, updateError := context.CloudFormation.UpdateStack(&updateStackInput)
		if updateError != nil {
			context.Logger.Error("Error updating stack: " + updateError.Error())
			return updateError
		}
		conn.MonitorStackQueue()
	} else {
		_, updateError := context.CloudFormation.UpdateStack(&updateStackInput)
		if updateError != nil {
			context.Logger.Error("Error updating stack: " + updateError.Error())
			return updateError
		}
		context.Logger.Info("Stack update request successful")
	}
	return nil
}

// This function gets template and  name of stack. It creates "UpdateStackInput" structure.
func updateStackInput(context *context.Context, template *string, stackName *string) cloudformation.UpdateStackInput {
	params, err := parameters.ResolveParameters(context)
	if err != nil {
		context.Logger.Error(err.Error())
		return cloudformation.UpdateStackInput{}
	}
	rawCapabilities := *context.CliArguments.Capabilities
	capabilities := make([]*string, len(rawCapabilities))
	for i, capability := range rawCapabilities {
		capabilities[i] = &capability
	}
	templateStruct := cloudformation.UpdateStackInput{
		Parameters:   params,
		TemplateBody: template,
		StackName:    stackName,
		Capabilities: capabilities,
	}
	return templateStruct
}
