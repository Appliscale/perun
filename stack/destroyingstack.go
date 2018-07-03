package stack

import (
	"github.com/Appliscale/perun/context"
	"github.com/Appliscale/perun/mysession"
	"github.com/Appliscale/perun/progress"
	"github.com/aws/aws-sdk-go/service/cloudformation"
)

// DestroyStack bases on "DeleteStackInput" structure and destroys stack. It uses "StackName" to choose which stack will be destroy. Before that it creates session.
func DestroyStack(context *context.Context) error {
	delStackInput := deleteStackInput(context)
	currentSession := mysession.InitializeSession(context)
	api := cloudformation.New(currentSession)

	var err error = nil
	if *context.CliArguments.Progress {
		conn, err := progress.GetRemoteSink(context, currentSession)
		if err != nil {
			context.Logger.Error("Error getting remote sink configuration: " + err.Error())
			return err
		}
		_, err = api.DeleteStack(&delStackInput)
		conn.MonitorStackQueue()
	} else {
		_, err = api.DeleteStack(&delStackInput)
	}
	if err != nil {
		context.Logger.Error(err.Error())
		return err
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
