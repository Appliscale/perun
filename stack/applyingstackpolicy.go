package stack

import (
	"github.com/Appliscale/perun/context"
	"github.com/aws/aws-sdk-go/service/cloudformation"
)

// ApplyStackPolicy creates StackPolicy from JSON file.
func ApplyStackPolicy(context *context.Context) error {
	template, stackName, incorrectPath := getTemplateFromFile(context)
	if incorrectPath != nil {
		context.Logger.Error(incorrectPath.Error())
		return incorrectPath
	}
	templateStruct := createStackPolicyInput(&template, &stackName)

	_, creationError := context.CloudFormation.SetStackPolicy(&templateStruct)
	if creationError != nil {
		context.Logger.Error("Error creating stack policy: " + creationError.Error())
		return creationError
	}

	return nil
}

// This function gets template and  name of stack. It creates "CreateStackInput" structure.
func createStackPolicyInput(template *string, stackName *string) cloudformation.SetStackPolicyInput {
	templateStruct := cloudformation.SetStackPolicyInput{
		StackPolicyBody: template,
		StackName:       stackName,
	}
	return templateStruct
}
