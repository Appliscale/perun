package stack

import (
	"github.com/Appliscale/perun/context"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudformation"
)

// ApplyStackPolicy creates StackPolicy from JSON file.
func ApplyStackPolicy(context *context.Context) error {
	template, stackName, incorrectPath := getTemplateFromFile(context)
	if incorrectPath != nil {
		context.Logger.Error(incorrectPath.Error())
		return incorrectPath
	}
	templateStruct := createStackPolicyInput(&template, &stackName, context)

	currentSession, sessionError := prepareSession(context)
	if sessionError == nil {

		creationError := createStackPolicy(templateStruct, currentSession)
		if creationError != nil {
			context.Logger.Error("Error creating stack policy: " + creationError.Error())
			return creationError
		}
	} else {
		context.Logger.Error(sessionError.Error())
		return sessionError
	}
	return nil
}

// Getting template from file and setting StackPolicy.
func createStackPolicy(templateStruct cloudformation.SetStackPolicyInput, session *session.Session) (apiError error) {
	api := cloudformation.New(session)
	_, apiError = api.SetStackPolicy(&templateStruct)

	return apiError
}

// This function gets template and  name of stack. It creates "CreateStackInput" structure.
func createStackPolicyInput(template *string, stackName *string, context *context.Context) cloudformation.SetStackPolicyInput {
	templateStruct := cloudformation.SetStackPolicyInput{
		StackPolicyBody: template,
		StackName:       stackName,
	}

	return templateStruct
}
