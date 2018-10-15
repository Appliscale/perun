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
	"errors"
	"github.com/Appliscale/perun/context"
	"github.com/aws/aws-sdk-go/service/cloudformation"
)

// Getting stackName and flag which describes stack termination protection to create struct.
func createUpdateTerminationProtectionInput(stackName string, isProtectionEnable bool) cloudformation.UpdateTerminationProtectionInput {
	templateStruct := cloudformation.UpdateTerminationProtectionInput{
		EnableTerminationProtection: &isProtectionEnable,
		StackName:                   &stackName,
	}

	return templateStruct
}

// Checking flag and settting protection.
func isProtectionEnable(context *context.Context) (bool, error) {
	if *context.CliArguments.DisableStackTermination {
		return true, nil
	} else if *context.CliArguments.EnableStackTermination {
		return false, nil
	}

	return false, errors.New("Incorrect StackTerminationProtection flag")
}

// SetTerminationProtection turn off or on stack protection from being deleted.
func SetTerminationProtection(context *context.Context) error {
	stackName := context.CliArguments.Stack
	isProtectionEnable, stackTerminationError := isProtectionEnable(context)
	if stackTerminationError != nil {
		return stackTerminationError
	}
	templateStruct := createUpdateTerminationProtectionInput(*stackName, isProtectionEnable)
	_, apiError := context.CloudFormation.UpdateTerminationProtection(&templateStruct)
	if apiError != nil {
		context.Logger.Error("Error setting stack termination protection: " + apiError.Error())
		return apiError
	}

	if isProtectionEnable {
		context.Logger.Info("Terminaction Protection Enabled successfully")
	} else {
		context.Logger.Info("Termination Protection Disabled successfully")
	}

	return nil
}
