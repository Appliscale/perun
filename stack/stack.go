package stack

import (
	"errors"
	"path"

	"io/ioutil"

	"github.com/Appliscale/perun/context"
	"github.com/Appliscale/perun/myuser"
)

// This function reads "StackName" from Stack in CliArguments and file from TemplatePath in CliArguments. It converts these to type string.
func getTemplateFromFile(context *context.Context) (string, string, error) {
	var rawTemplate []byte
	var readFileError error
	templatePath, pathError := getPath(context)
	if pathError != nil {
		return "", "", pathError
	}

	rawTemplate, readFileError = ioutil.ReadFile(templatePath)
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

	if *context.CliArguments.Mode == "set-stack-policy" {
		if *context.CliArguments.Unblock {
			path = homePath + "/.config/perun/stack-policies/unblocked.json"
		} else if *context.CliArguments.Block {
			path = homePath + "/.config/perun/stack-policies/blocked.json"
		} else if len(*context.CliArguments.TemplatePath) > 0 && isStackPolicyFileJSON(*context.CliArguments.TemplatePath) {
			path = *context.CliArguments.TemplatePath
		} else {
			return "", errors.New("Incorrect path")
		}
	} else if len(*context.CliArguments.TemplatePath) > 0 {
		path = *context.CliArguments.TemplatePath
	}
	return
}

// Checking is file type JSON.
func isStackPolicyFileJSON(filename string) bool {
	templateFileExtension := path.Ext(filename)
	if templateFileExtension == ".json" {
		return true
	}

	return false
}
