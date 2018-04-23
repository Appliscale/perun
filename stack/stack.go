package stack

import (
	"encoding/json"
	"errors"
	"github.com/Appliscale/perun/context"
	"github.com/Appliscale/perun/mysession"
	"github.com/Appliscale/perun/parameters"
	"github.com/Appliscale/perun/progress"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudformation"
	"io/ioutil"
	"path"
)

// This function reads "StackName" from Stack in CliArguments and file from TemplatePath in CliArguments. It converts these to type string.
func getTemplateFromFile(context *context.Context) (string, string, error) {
	var rawTemplate []byte
	var readFileError error
	path, pathError := getPath(context)
	if pathError != nil {
		return "", "", pathError
	}

	rawTemplate, readFileError = ioutil.ReadFile(path)
	if readFileError != nil {
		context.Logger.Error(readFileError.Error())
		return "", "", readFileError
	}

	rawStackName := *context.CliArguments.Stack
	template := string(rawTemplate)
	stackName := rawStackName
	return template, stackName, nil
}

func getPath(context *context.Context) (path string, err error) {
	if *context.CliArguments.Mode == "create-stack" {
		path = *context.CliArguments.TemplatePath
	} else if *context.CliArguments.Mode == "set-stack-policy" {
		if *context.CliArguments.Unblock {
			path = "./stack/defaultstackpolicy/unblocked.json"
		} else if *context.CliArguments.Block {
			path = "./stack/defaultstackpolicy/blocked.json"
		} else if len(*context.CliArguments.TemplatePath) > 0 {
			path = *context.CliArguments.TemplatePath
		} else {
			return "", errors.New("Incorrect path.")
		}

	}

	return
}

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
func NewStack(context *context.Context) {
	template, stackName, incorrectPath := getTemplateFromFile(context)
	if incorrectPath != nil {
		context.Logger.Error(incorrectPath.Error())
		return
	}
	templateStruct, templateError := createStackInput(&template, &stackName, context)
	if templateError != nil {
		context.Logger.Error(templateError.Error())
		return
	}

	currentSession, sessionError := prepareSession(context)
	if sessionError == nil {

		if *context.CliArguments.Progress {
			conn, err := progress.GetRemoteSink(context, currentSession)
			if err != nil {
				context.Logger.Error("Error getting remote sink configuration: " + err.Error())
				return
			}
			templateStruct.NotificationARNs = []*string{conn.TopicArn}
			err = createStack(templateStruct, currentSession)
			if err != nil {
				context.Logger.Error("Error creating stack: " + err.Error())
				return
			}
			conn.MonitorQueue()
		} else {
			err := createStack(templateStruct, currentSession)
			if err != nil {
				context.Logger.Error("Error creating stack: " + err.Error())
				return
			}
		}
	}
}

// DestroyStack bases on "DeleteStackInput" structure and destroys stack. It uses "StackName" to choose which stack will be destroy. Before that it creates session.
func DestroyStack(context *context.Context) {
	delStackInput := deleteStackInput(context)
	currentSession, sessionError := prepareSession(context)
	if sessionError == nil {
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
}

// This function gets "StackName" from Stack in CliArguments and creates "DeleteStackInput" structure.
func deleteStackInput(context *context.Context) cloudformation.DeleteStackInput {
	name := *context.CliArguments.Stack
	templateStruct := cloudformation.DeleteStackInput{
		StackName: &name,
	}
	return templateStruct
}

// Get the parameters - if parameters file provided - from file, else - interactively from user
func getParameters(context *context.Context) (params []*cloudformation.Parameter, err error) {
	if *context.CliArguments.ParametersFile == "" {
		params, err = parameters.GetAwsParameters(context)
	} else {
		var parametersData []byte
		var readParameters []*parameters.Parameter
		parametersData, err = ioutil.ReadFile(*context.CliArguments.ParametersFile)
		if err != nil {
			return
		}
		err = json.Unmarshal(parametersData, &readParameters)
		if err != nil {
			return
		}
		params = parameters.ParseParameterToAwsCompatible(readParameters)
	}
	return
}

// NewStackPolicy creates StackPolicy from JSON file.
func NewStackPolicy(context *context.Context) {
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
func createStackPolicy(templateStruct cloudformation.SetStackPolicyInput, session *session.Session) (err error) {
	api := cloudformation.New(session)
	_, err = api.SetStackPolicy(&templateStruct)

	return err
}

// This function gets template and  name of stack. It creates "CreateStackInput" structure.
func createStackPolicyInput(template *string, stackName *string, context *context.Context) cloudformation.SetStackPolicyInput {
	templateStruct := cloudformation.SetStackPolicyInput{
		StackPolicyBody: template,
		StackName:       stackName,
	}
	return templateStruct
}

// Creating SessionToken and Session.
func prepareSession(context *context.Context) (*session.Session, error) {
	tokenError := mysession.UpdateSessionToken(context.Config.DefaultProfile, context.Config.DefaultRegion, context.Config.DefaultDurationForMFA, context)
	if tokenError != nil {
		context.Logger.Error(tokenError.Error())
	}
	currentSession, createSessionError := mysession.CreateSession(context, context.Config.DefaultProfile, &context.Config.DefaultRegion)
	if createSessionError != nil {
		context.Logger.Error(createSessionError.Error())
	}
	return currentSession, createSessionError
}

// Checking if file is type JSON.
func isStackPolicyFileJSON(filename string) bool {
	templateFileExtension := path.Ext(filename)
	if templateFileExtension == ".json" {
		return true
	}
	return false

}
