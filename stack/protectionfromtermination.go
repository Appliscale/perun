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
	return nil
}
