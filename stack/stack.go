package stack

import (
	"github.com/Appliscale/perun/context"
	"github.com/Appliscale/perun/mysession"
	//"github.com/Appliscale/perun/notificationservice"
	"io/ioutil"
	"os"

	"github.com/Appliscale/perun/progress"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudformation"
)

// This function gets template and  name of stack. It creates "CreateStackInput" structure.
func createStackInput(context *context.Context, template *string, stackName *string) cloudformation.CreateStackInput {

	rawCapabilities := *context.CliArguments.Capabilities
	capabilities := make([]*string, len(rawCapabilities))
	for i, capability := range rawCapabilities {
		capabilities[i] = &capability
	}

	templateStruct := cloudformation.CreateStackInput{
		TemplateBody: template,
		StackName:    stackName,
		Capabilities: capabilities,
	}
	return templateStruct
}

// This function gets template and  name of stack. It creates "CreateStackInput" structure.
func updateStackInput(context *context.Context, template *string, stackName *string) cloudformation.UpdateStackInput {
	rawCapabilities := *context.CliArguments.Capabilities
	capabilities := make([]*string, len(rawCapabilities))
	for i, capability := range rawCapabilities {
		capabilities[i] = &capability
	}
	templateStruct := cloudformation.UpdateStackInput{
		TemplateBody: template,
		StackName:    stackName,
		Capabilities: capabilities,
	}
	return templateStruct
}

// This function reads "StackName" from Stack in CliArguments and file from TemplatePath in CliArguments. It converts these to type string.
func getTemplateFromFile(context *context.Context) (string, string) {

	rawTemplate, readFileError := ioutil.ReadFile(*context.CliArguments.TemplatePath)
	if readFileError != nil {
		context.Logger.Error(readFileError.Error())
	}

	rawStackName := *context.CliArguments.Stack
	template := string(rawTemplate)
	stackName := rawStackName
	return template, stackName
}

// This function uses CreateStackInput variable to create Stack.
func createStack(context *context.Context, templateStruct cloudformation.CreateStackInput, session *session.Session) {
	api := cloudformation.New(session)
	_, err := api.CreateStack(&templateStruct)
	if err != nil {
		context.Logger.Error("Error creating stack: " + err.Error())
	}
}

// This function uses all functions above and session to create Stack.
func NewStack(context *context.Context) {
	template, stackName := getTemplateFromFile(context)

	templateStruct := createStackInput(context, &template, &stackName)
	currentSession := mysession.InitializeSession(context)

	if *context.CliArguments.Progress {
		conn, err := progress.GetRemoteSink(context, currentSession)
		if err != nil {
			context.Logger.Error("Error getting remote sink configuration: " + err.Error())
			return
		}
		templateStruct.NotificationARNs = []*string{conn.TopicArn}
		createStack(context, templateStruct, currentSession)
		conn.MonitorQueue()
	} else {
		createStack(context, templateStruct, currentSession)
	}
}

// This function bases on "DeleteStackInput" structure and destroys stack. It uses "StackName" to choose which stack will be destroy. Before that it creates session.
func DestroyStack(context *context.Context) {
	delStackInput := deleteStackInput(context)
	currentSession := mysession.InitializeSession(context)
	api := cloudformation.New(currentSession)

	var err error = nil
	if *context.CliArguments.Progress {
		conn, err := progress.GetRemoteSink(context, currentSession)
		if err != nil {
			context.Logger.Error("Error getting remote sink configuration: " + err.Error())
			return
		}
		_, err = api.DeleteStack(&delStackInput)
		conn.MonitorQueue()
	} else {
		_, err = api.DeleteStack(&delStackInput)
	}
	if err != nil {
		context.Logger.Error(err.Error())
	}
}

func UpdateStack(context *context.Context) {
	template, stackName := getTemplateFromFile(context)
	templateStruct := updateStackInput(context, &template, &stackName)
	session := mysession.InitializeSession(context)
	err := updateStack(templateStruct, session)
	if err != nil {
		context.Logger.Error(err.Error())
		os.Exit(1)
	}
}

func updateStack(updateStackInput cloudformation.UpdateStackInput, session *session.Session) error {
	api := cloudformation.New(session)
	_, err := api.UpdateStack(&updateStackInput)
	return err
}

// This function gets "StackName" from Stack in CliArguments and creates "DeleteStackInput" structure.
func deleteStackInput(context *context.Context) cloudformation.DeleteStackInput {
	name := *context.CliArguments.Stack
	templateStruct := cloudformation.DeleteStackInput{
		StackName: &name,
	}
	return templateStruct
}
