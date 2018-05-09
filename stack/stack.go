package stack

import (
	"encoding/json"
	"errors"
	"github.com/Appliscale/perun/context"
	"github.com/Appliscale/perun/mysession"
	"github.com/Appliscale/perun/myuser"
	"github.com/Appliscale/perun/parameters"
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

// Looking for path to user/default template.
func getPath(context *context.Context) (path string, err error) {
	homePath, pathError := myuser.GetUserHomeDir()
	if pathError != nil {
		context.Logger.Error(pathError.Error())
		return "", pathError
	}

	if *context.CliArguments.Mode == "create-stack" {
		path = *context.CliArguments.TemplatePath
	} else if *context.CliArguments.Mode == "set-stack-policy" {
		if *context.CliArguments.Unblock {
			path = homePath + "/.config/perun/stack-policies/unblocked.json"
		} else if *context.CliArguments.Block {
			path = homePath + "/.config/perun/stack-policies/blocked.json"
		} else if len(*context.CliArguments.TemplatePath) > 0 && isStackPolicyFileJSON(*context.CliArguments.TemplatePath) {
			path = *context.CliArguments.TemplatePath
		} else {
			return "", errors.New("Incorrect path")
		}

	}

	return
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

// Checking is file type JSON.
func isStackPolicyFileJSON(filename string) bool {
	templateFileExtension := path.Ext(filename)
	if templateFileExtension == ".json" {
		return true
	}

	return false
}
