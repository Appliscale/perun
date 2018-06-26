package stack

import (
	"github.com/Appliscale/perun/context"
	"github.com/Appliscale/perun/progress"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudformation"
)

// This function gets template and  name of stack. It creates "CreateStackInput" structure.
func createStackInput(template *string, stackName *string, context *context.Context) (cloudformation.CreateStackInput, error) {
	params, err := getParameters(context)
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

// This function uses CreateStackInput variable to create Stack.
func createStack(templateStruct cloudformation.CreateStackInput, session *session.Session) (err error) {
	api := cloudformation.New(session)
	_, err = api.CreateStack(&templateStruct)

	return
}

// NewStack uses all functions above and session to create Stack.
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

	currentSession, sessionError := prepareSession(context)
	if sessionError == nil {

		if *context.CliArguments.Progress {
			conn, remoteSinkError := progress.GetRemoteSink(context, currentSession)
			if remoteSinkError != nil {
				context.Logger.Error("Error getting remote sink configuration: " + remoteSinkError.Error())
				return remoteSinkError
			}
			templateStruct.NotificationARNs = []*string{conn.TopicArn}
			creationError := createStack(templateStruct, currentSession)
			if creationError != nil {
				context.Logger.Error("Error creating stack: " + creationError.Error())
				return creationError
			}
			conn.MonitorQueue()
		} else {
			creationError := createStack(templateStruct, currentSession)
			if creationError != nil {
				context.Logger.Error("Error creating stack: " + creationError.Error())
				return creationError
			}
		}
	} else {
		context.Logger.Error(sessionError.Error())
		return sessionError
	}
	return nil
}
