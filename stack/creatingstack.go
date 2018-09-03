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
