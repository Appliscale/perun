// Copyright 2017 Appliscale
//
// Maintainers and contributors are listed in README file inside repository.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package cliparser provides tools and structures for parsing and validating
// perun CLI arguments.
package cliparser

import (
	"errors"
	"github.com/Appliscale/perun/logger"
	"github.com/Appliscale/perun/utilities"
	"gopkg.in/alecthomas/kingpin.v2"
)

var ValidateMode = "validate"
var ConvertMode = "convert"
var OfflineValidateMode = "validate_offline"
var ConfigureMode = "configure"
var CreateStackMode = "create-stack"
var DestroyStackMode = "delete-stack"
var SetupSinkMode = "setup-remote-sink"
var DestroySinkMode = "destroy-remote-sink"
var CreateParametersMode = "create-parameters"
var SetStackPolicyMode = "set-stack-policy"

const JSON = "json"
const YAML = "yaml"

type CliArguments struct {
	Mode                    *string
	TemplatePath            *string
	Parameters              *map[string]string
	OutputFilePath          *string
	ConfigurationPath       *string
	Quiet                   *bool
	Yes                     *bool
	Verbosity               *string
	MFA                     *bool
	DurationForMFA          *int64
	Profile                 *string
	Region                  *string
	Sandbox                 *bool
	Stack                   *string
	PrettyPrint             *bool
	Progress                *bool
	ParametersFile          *string
	Block                   *bool
	Unblock                 *bool
	DisableStackTermination *bool
	EnableStackTermination  *bool
}

// Get and validate CLI arguments. Returns error if validation fails.
func ParseCliArguments(args []string) (cliArguments CliArguments, err error) {
	var (
		app = kingpin.New("Perun", "Swiss army knife for AWS CloudFormation templates - validation, conversion, generators and other various stuff.")

		quiet             = app.Flag("quiet", "No console output, just return code.").Short('q').Bool()
		yes               = app.Flag("yes", "Always say yes.").Short('y').Bool()
		verbosity         = app.Flag("verbosity", "Logger verbosity: TRACE | DEBUG | INFO | ERROR.").Short('v').String()
		mfa               = app.Flag("mfa", "Enable AWS MFA.").Bool()
		DurationForMFA    = app.Flag("duration", "Duration for AWS MFA token (seconds value from range [1, 129600]).").Short('d').Int64()
		profile           = app.Flag("profile", "An AWS profile name.").Short('p').String()
		region            = app.Flag("region", "An AWS region to use.").Short('r').String()
		sandbox           = app.Flag("sandbox", "Do not use configuration files hierarchy.").Bool()
		configurationPath = app.Flag("config", "A path to the configuration file").Short('c').String()
		showProgress      = app.Flag("progress", "Show progress of stack creation. Option available only after setting up a remote sink").Bool()

		onlineValidate            = app.Command(ValidateMode, "Online template Validation")
		onlineValidateTemplate    = onlineValidate.Arg("template", "A path to the template file.").String()
		onlineValidateImpTemplate = onlineValidate.Flag("template", "A path to the template file.").String()

		offlineValidate            = app.Command(OfflineValidateMode, "Offline Template Validation")
		offlineValidateTemplate    = offlineValidate.Arg("template", "A path to the template file.").String()
		offlineValidateImpTemplate = offlineValidate.Flag("template", "A path to the template file.").String()

		convert              = app.Command(ConvertMode, "Convertion between JSON and YAML of template files")
		convertTemplate      = convert.Arg("template", "A path to the template file.").String()
		convertOutputFile    = convert.Arg("output", "A path where converted file will be saved.").String()
		convertImpTemplate   = convert.Flag("from", "A path to the template file.").String()
		convertImpOutputFile = convert.Flag("to", "A path where converted file will be saved.").String()
		convertPrettyPrint   = convert.Flag("pretty-print", "Pretty printing JSON").Bool()

		configure = app.Command(ConfigureMode, "Create your own configuration mode")

		createStack               = app.Command(CreateStackMode, "Creates a stack on aws")
		createStackName           = createStack.Arg("stack", "An AWS stack name.").String()
		createStackTemplate       = createStack.Arg("template", "A path to the template file.").String()
		createStackImpName        = createStack.Flag("stack", "An AWS stack name.").String()
		createStackImpTemplate    = createStack.Flag("template", "A path to the template file.").String()
		createStackParams         = createStack.Flag("parameter", "list of parameters").StringMap()
		createStackParametersFile = createStack.Flag("parameters-file", "filename with parameters").String()

		deleteStack        = app.Command(DestroyStackMode, "Deletes a stack on aws")
		deleteStackName    = deleteStack.Arg("stack", "An AWS stack name.").String()
		deleteStackImpName = deleteStack.Flag("stack", "An AWS stack name.").String()

		setupSink = app.Command(SetupSinkMode, "Sets up resources required for progress report on stack events (SNS Topic, SQS Queue and SQS Queue Policy)")

		destroySink = app.Command(DestroySinkMode, "Destroys resources created with setup-remote-sink")

		createParameters                 = app.Command(CreateParametersMode, "Creates a JSON parameters configuration suitable for give cloud formation file")
		createParametersTemplate         = createParameters.Arg("template", "A path to the template file.").String()
		createParametersImpTemplate      = createParameters.Flag("template", "A path to the template file.").String()
		createParametersParamsOutputFile = createParameters.Flag("output", "A path to file where parameters will be saved.").String()
		createParametersParams           = createParameters.Flag("parameter", "list of parameters").StringMap()
		createParametersPrettyPrint      = createParameters.Flag("pretty-print", "Pretty printing JSON").Bool()

		setStackPolicy                  = app.Command(SetStackPolicyMode, "Set stack policy using JSON file.")
		setStackPolicyName              = setStackPolicy.Arg("stack", "An AWS stack name.").String()
		setStackPolicyImpName           = setStackPolicy.Flag("stack", "An AWS stack name.").String()
		setStackPolicyTemplate          = setStackPolicy.Arg("template", "A path to the template file.").String()
		setStackPolicyImpTemplate       = setStackPolicy.Flag("template", "A path to the template file.").String()
		setDefaultBlockingStackPolicy   = setStackPolicy.Flag("block", "Blocking all actions.").Bool()
		setDefaultUnblockingStackPolicy = setStackPolicy.Flag("unblock", "Unblocking all actions.").Bool()
		setDisableStackTermination      = setStackPolicy.Flag("disable-stack-termination", "Allow to delete a stack.").Bool()
		setEnableStackTermination       = setStackPolicy.Flag("enable-stack-termination", "Protecting a stack from being deleted.").Bool()
	)
	app.HelpFlag.Short('h')
	app.Version(utilities.VersionStatus())

	switch kingpin.MustParse(app.Parse(args[1:])) {

	//online validate
	case onlineValidate.FullCommand():
		cliArguments.Mode = &ValidateMode
		if len(*onlineValidateTemplate) > 0 {
			cliArguments.TemplatePath = onlineValidateTemplate
		} else if len(*onlineValidateImpTemplate) > 0 {
			cliArguments.TemplatePath = onlineValidateImpTemplate
		} else {
			err = errors.New("You have to specify the template, try --help")
			return
		}

		// offline validation
	case offlineValidate.FullCommand():
		cliArguments.Mode = &OfflineValidateMode
		if len(*offlineValidateTemplate) > 0 {
			cliArguments.TemplatePath = offlineValidateTemplate
		} else if len(*offlineValidateImpTemplate) > 0 {
			cliArguments.TemplatePath = offlineValidateImpTemplate
		} else {
			err = errors.New("You have to specify the template, try --help")
			return
		}

		// convert
	case convert.FullCommand():
		cliArguments.Mode = &ConvertMode
		cliArguments.PrettyPrint = convertPrettyPrint

		if len(*convertImpOutputFile) > 0 && len(*convertImpTemplate) > 0 {
			cliArguments.TemplatePath = convertImpTemplate
			cliArguments.OutputFilePath = convertImpOutputFile
		} else if len(*convertOutputFile) > 0 && len(*convertTemplate) > 0 {
			cliArguments.TemplatePath = convertTemplate
			cliArguments.OutputFilePath = convertOutputFile
		} else if len(*convertTemplate) > 0 && len(*convertImpOutputFile) > 0 {
			cliArguments.TemplatePath = convertTemplate
			cliArguments.OutputFilePath = convertImpOutputFile
		} else {
			err = errors.New("You have to specify the template and the output file, try --help")
			return
		}

		// configure
	case configure.FullCommand():
		cliArguments.Mode = &ConfigureMode

		// create Stack
	case createStack.FullCommand():
		cliArguments.Mode = &CreateStackMode
		cliArguments.Parameters = createStackParams
		cliArguments.ParametersFile = createStackParametersFile
		if len(*createStackImpTemplate) > 0 && len(*createStackImpName) > 0 {
			cliArguments.Stack = createStackImpName
			cliArguments.TemplatePath = createStackImpTemplate
		} else if len(*createStackName) > 0 && len(*createStackTemplate) > 0 {
			cliArguments.Stack = createStackName
			cliArguments.TemplatePath = createStackTemplate
		} else if len(*createStackName) > 0 && len(*createStackImpTemplate) > 0 {
			cliArguments.Stack = createStackName
			cliArguments.TemplatePath = createStackImpTemplate
		} else {
			err = errors.New("You have to specify stack name and template file, try --help")
			return
		}

		// delete Stack
	case deleteStack.FullCommand():
		cliArguments.Mode = &DestroyStackMode
		if len(*deleteStackName) > 0 {
			cliArguments.Stack = deleteStackName
		} else if len(*deleteStackImpName) > 0 {
			cliArguments.Stack = deleteStackImpName
		} else {
			err = errors.New("You have to specify the stack name, try --help")
			return
		}

		// create Parameters
	case createParameters.FullCommand():
		cliArguments.Mode = &CreateParametersMode
		if len(*createParametersTemplate) > 0 {
			cliArguments.TemplatePath = createParametersTemplate
		} else if len(*createParametersImpTemplate) > 0 {
			cliArguments.TemplatePath = createParametersImpTemplate
		} else {
			err = errors.New("You have to specify the cloud formation template, try --help")
			return
		}
		cliArguments.OutputFilePath = createParametersParamsOutputFile
		cliArguments.Parameters = createParametersParams
		cliArguments.PrettyPrint = createParametersPrettyPrint

		// set stack policy
	case setStackPolicy.FullCommand():
		cliArguments.Mode = &SetStackPolicyMode
		cliArguments.Block = setDefaultBlockingStackPolicy
		cliArguments.Unblock = setDefaultUnblockingStackPolicy
		cliArguments.DisableStackTermination = setDisableStackTermination
		cliArguments.EnableStackTermination = setEnableStackTermination
		if len(*setStackPolicyImpTemplate) > 0 && len(*setStackPolicyImpName) > 0 {
			cliArguments.Stack = setStackPolicyImpName
			cliArguments.TemplatePath = setStackPolicyImpTemplate
		} else if len(*setStackPolicyName) > 0 && len(*setStackPolicyTemplate) > 0 {
			cliArguments.Stack = setStackPolicyName
			cliArguments.TemplatePath = setStackPolicyTemplate
		} else if len(*setStackPolicyName) > 0 && len(*setStackPolicyImpTemplate) > 0 {
			cliArguments.Stack = setStackPolicyName
			cliArguments.TemplatePath = setStackPolicyImpTemplate
		} else if len(*setStackPolicyName) > 0 && *setDefaultBlockingStackPolicy {
			cliArguments.Stack = setStackPolicyName
			cliArguments.Block = setDefaultBlockingStackPolicy
		} else if len(*setStackPolicyImpName) > 0 && *setDefaultBlockingStackPolicy {
			cliArguments.Stack = setStackPolicyImpName
			cliArguments.Block = setDefaultBlockingStackPolicy
		} else if len(*setStackPolicyName) > 0 && *setDefaultUnblockingStackPolicy {
			cliArguments.Stack = setStackPolicyName
			cliArguments.Unblock = setDefaultUnblockingStackPolicy
		} else if len(*setStackPolicyImpName) > 0 && *setDefaultUnblockingStackPolicy {
			cliArguments.Stack = setStackPolicyImpName
			cliArguments.Unblock = setDefaultUnblockingStackPolicy
		} else if len(*setStackPolicyName) > 0 && *setDisableStackTermination {
			cliArguments.Stack = setStackPolicyName
			cliArguments.DisableStackTermination = setDisableStackTermination
		} else if len(*setStackPolicyImpName) > 0 && *setDisableStackTermination {
			cliArguments.Stack = setStackPolicyImpName
			cliArguments.DisableStackTermination = setDisableStackTermination
		} else if len(*setStackPolicyName) > 0 && *setEnableStackTermination {
			cliArguments.Stack = setStackPolicyName
			cliArguments.EnableStackTermination = setEnableStackTermination
		} else if len(*setStackPolicyImpName) > 0 && *setEnableStackTermination {
			cliArguments.Stack = setStackPolicyImpName
			cliArguments.EnableStackTermination = setEnableStackTermination
		} else {
			err = errors.New("You have to specify stack name and template file or flag, try --help")
			return
		}

		// set up remote sink
	case setupSink.FullCommand():
		cliArguments.Mode = &SetupSinkMode

		// destroy remote sink
	case destroySink.FullCommand():
		cliArguments.Mode = &DestroySinkMode
	}

	// OTHER FLAGS
	cliArguments.Quiet = quiet
	cliArguments.Yes = yes
	cliArguments.Verbosity = verbosity
	cliArguments.MFA = mfa
	cliArguments.DurationForMFA = DurationForMFA
	cliArguments.Profile = profile
	cliArguments.Region = region
	cliArguments.Sandbox = sandbox
	cliArguments.ConfigurationPath = configurationPath
	cliArguments.Progress = showProgress

	if *cliArguments.DurationForMFA < 0 {
		err = errors.New("You should specify value for duration of MFA token greater than zero")
		return
	}

	if *cliArguments.DurationForMFA > 129600 {
		err = errors.New("You should specify value for duration of MFA token smaller than 129600 (3 hours)")
		return
	}

	if *cliArguments.Verbosity != "" && !logger.IsVerbosityValid(*cliArguments.Verbosity) {
		err = errors.New("You specified invalid value for --verbosity flag")
		return
	}

	return
}
