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
	"github.com/Appliscale/perun/progress"
	"github.com/aws/aws-sdk-go/service/cloudformation"
)

// DestroyStack bases on "DeleteStackInput" structure and destroys stack. It uses "StackName" to choose which stack will be destroy. Before that it creates session.
func DestroyStack(context *context.Context) error {
	delStackInput := deleteStackInput(context)

	var err error = nil
	if *context.CliArguments.Progress {
		conn, err := progress.GetRemoteSink(context)
		if err != nil {
			context.Logger.Error("Error getting remote sink configuration: " + err.Error())
			return err
		}
		_, err = context.CloudFormation.DeleteStack(&delStackInput)
		if err != nil {
			context.Logger.Error(err.Error())
			return err
		}
		conn.MonitorStackQueue()
	} else {
		_, err = context.CloudFormation.DeleteStack(&delStackInput)
		if err != nil {
			context.Logger.Error(err.Error())
			return err
		}
		context.Logger.Info("Stack deletion request successful")
	}
	return nil
}

// This function gets "StackName" from Stack in CliArguments and creates "DeleteStackInput" structure.
func deleteStackInput(context *context.Context) cloudformation.DeleteStackInput {
	name := *context.CliArguments.Stack
	templateStruct := cloudformation.DeleteStackInput{
		StackName: &name,
	}
	return templateStruct
}
