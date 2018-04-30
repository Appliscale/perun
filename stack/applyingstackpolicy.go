package stack

import (
	"github.com/Appliscale/perun/context"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudformation"
)

// ApplyStackPolicy creates StackPolicy from JSON file.
func ApplyStackPolicy(context *context.Context) {
	template, stackName, incorrectPath := getTemplateFromFile(context)
	if incorrectPath != nil {
		context.Logger.Error(incorrectPath.Error())
		return
	}
	templateStruct := createStackPolicyInput(&template, &stackName, context)

	currentSession, sessionError := prepareSession(context)
	if sessionError == nil {

		err := createStackPolicy(templateStruct, currentSession)
		if err != nil {
			context.Logger.Error("Error creating stack policy: " + err.Error())
			return
		}
	}
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
