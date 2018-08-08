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
	"time"
)

var ValidateMode = "validate"
var ConfigureMode = "configure"
var CreateStackMode = "create-stack"
var DestroyStackMode = "delete-stack"
var UpdateStackMode = "update-stack"
var MfaMode = "mfa"
var SetupSinkMode = "setup-remote-sink"
var DestroySinkMode = "destroy-remote-sink"
var CreateParametersMode = "create-parameters"
var SetStackPolicyMode = "set-stack-policy"
var CreateChangeSetMode = "create-change-set"
var LintMode = "lint"

var ChangeSetDefaultName string

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
	Capabilities            *[]string
	PrettyPrint             *bool
	Progress                *bool
	ParametersFile          *string
	Block                   *bool
	Unblock                 *bool
	DisableStackTermination *bool
	EnableStackTermination  *bool
	ChangeSet               *string
	Lint                    *bool
	LinterConfiguration     *string
	EstimateCost            *bool
	SkipValidation          *bool
}

// Get and validate CLI arguments. Returns error if validation fails.
func ParseCliArguments(args []string) (cliArguments CliArguments, err error) {
	var (
		app = kingpin.New("Perun", "A command-line validation tool for AWS Cloud Formation that allows to conquer the cloud faster!")

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
		noValidate        = app.Flag("no-validate", "Disable validation before stack creation/update or creating Change Set.").Bool()

		validate                  = app.Command(ValidateMode, "Template Validation")
		validateTemplate          = validate.Arg("template", "A path to the template file.").Required().String()
		validateLint              = validate.Flag("lint", "Enable template linting").Bool()
		validateLintConfiguration = validate.Flag("lint-configuration", "A path to the configuration file").String()
		validateEstimateCost      = validate.Flag("estimate-cost", "Enable cost estimation during validation").Bool()
		validateParams            = validate.Flag("parameter", "list of parameters").StringMap()
		validateParametersFile    = validate.Flag("parameters-file", "filename with parameters").String()

		lint              = app.Command(LintMode, "Additional validation and template style checks")
		lintTemplate      = lint.Arg("template", "A path to the template file.").Required().String()
		lintConfiguration = lint.Flag("lint-configuration", "A path to the configuration file").String()

		configure = app.Command(ConfigureMode, "Create your own configuration mode")

		createStack                  = app.Command(CreateStackMode, "Creates a stack on aws")
		createStackName              = createStack.Arg("stack", "An AWS stack name.").Required().String()
		createStackTemplate          = createStack.Arg("template", "A path to the template file.").Required().String()
		createStackCapabilities      = createStack.Flag("capabilities", "Capabilities: CAPABILITY_IAM | CAPABILITY_NAMED_IAM").Enums("CAPABILITY_IAM", "CAPABILITY_NAMED_IAM")
		createStackParams            = createStack.Flag("parameter", "list of parameters").StringMap()
		createStackParametersFile    = createStack.Flag("parameters-file", "filename with parameters").String()
		createStackLint              = createStack.Flag("lint", "Enable template linting").Bool()
		createStackLintConfiguration = createStack.Flag("lint-configuration", "A path to the configuration file").String()
		createStackEstimateCost      = createStack.Flag("estimate-cost", "Enable cost estimation during validation").Bool()

		createChangeSet                  = app.Command(CreateChangeSetMode, "Creates a changeSet on aws")
		changeSetStackName               = createChangeSet.Arg("stack", "An AWS stack name").Required().String()
		changeSetTemplate                = createChangeSet.Arg("template", "A path to the template file").Required().String()
		createChangeSetName              = createChangeSet.Arg("changeSet", "An AWS Change Set name").String()
		createChangeSetParams            = createChangeSet.Flag("parameter", "list of parameters").StringMap()
		createChangeSetParametersFile    = createChangeSet.Flag("parameters-file", "filename with parameters").String()
		createChangeSetLint              = createChangeSet.Flag("lint", "Enable template linting").Bool()
		createChangeSetLintConfiguration = createChangeSet.Flag("lint-configuration", "A path to the configuration file").String()
		createChangeSetEstimateCost      = createChangeSet.Flag("estimate-cost", "Enable cost estimation during validation").Bool()

		deleteStack     = app.Command(DestroyStackMode, "Deletes a stack on aws")
		deleteStackName = deleteStack.Arg("stack", "An AWS stack name.").Required().String()

		updateStack                  = app.Command(UpdateStackMode, "Updates a stack on aws")
		updateStackName              = updateStack.Arg("stack", "An AWS stack name").Required().String()
		updateStackTemplate          = updateStack.Arg("template", "A path to the template file.").Required().String()
		updateStackCapabilities      = updateStack.Flag("capabilities", "Capabilities: CAPABILITY_IAM | CAPABILITY_NAMED_IAM").Enums("CAPABILITY_IAM", "CAPABILITY_NAMED_IAM")
		updateStackParams            = updateStack.Flag("parameter", "list of parameters").StringMap()
		updateStackParametersFile    = updateStack.Flag("parameters-file", "filename with parameters").String()
		updateStackLint              = updateStack.Flag("lint", "Enable template linting").Bool()
		updateStackLintConfiguration = updateStack.Flag("lint-configuration", "A path to the configuration file").String()
		updateStackEstimateCost      = updateStack.Flag("estimate-cost", "Enable cost estimation during validation").Bool()

		mfaCommand = app.Command(MfaMode, "Create temporary secure credentials with MFA.")

		setupSink = app.Command(SetupSinkMode, "Sets up resources required for progress report on stack events (SNS Topic, SQS Queue and SQS Queue Policy)")

		destroySink = app.Command(DestroySinkMode, "Destroys resources created with setup-remote-sink")

		createParameters                 = app.Command(CreateParametersMode, "Creates a JSON parameters configuration suitable for give cloud formation file")
		createParametersTemplate         = createParameters.Arg("template", "A path to the template file.").Required().String()
		createParametersParamsOutputFile = createParameters.Arg("output", "A path to file where parameters will be saved.").Required().String()
		createParametersParams           = createParameters.Flag("parameter", "list of parameters").StringMap()
		createParametersPrettyPrint      = createParameters.Flag("pretty-print", "Pretty printing JSON").Bool()

		setStackPolicy                  = app.Command(SetStackPolicyMode, "Set stack policy using JSON file.")
		setStackPolicyName              = setStackPolicy.Arg("stack", "An AWS stack name.").Required().String()
		setStackPolicyTemplate          = setStackPolicy.Arg("template", "A path to the template file.").Required().String()
		setDefaultBlockingStackPolicy   = setStackPolicy.Flag("block", "Blocking all actions.").Bool()
		setDefaultUnblockingStackPolicy = setStackPolicy.Flag("unblock", "Unblocking all actions.").Bool()
		setDisableStackTermination      = setStackPolicy.Flag("disable-stack-termination", "Allow to delete a stack.").Bool()
		setEnableStackTermination       = setStackPolicy.Flag("enable-stack-termination", "Protecting a stack from being deleted.").Bool()
	)

	app.HelpFlag.Short('h')
	app.Version(utilities.VersionStatus())

	switch kingpin.MustParse(app.Parse(args[1:])) {

	//online validate
	case validate.FullCommand():
		cliArguments.Mode = &ValidateMode
		cliArguments.TemplatePath = validateTemplate
		cliArguments.Lint = validateLint
		cliArguments.LinterConfiguration = validateLintConfiguration
		cliArguments.EstimateCost = validateEstimateCost
		cliArguments.Parameters = validateParams
		cliArguments.ParametersFile = validateParametersFile

		// configure
	case configure.FullCommand():
		cliArguments.Mode = &ConfigureMode

	case lint.FullCommand():
		cliArguments.Mode = &LintMode
		cliArguments.TemplatePath = lintTemplate
		cliArguments.LinterConfiguration = lintConfiguration

		// create Stack
	case createStack.FullCommand():
		cliArguments.Mode = &CreateStackMode
		cliArguments.Stack = createStackName
		cliArguments.TemplatePath = createStackTemplate
		cliArguments.Capabilities = createStackCapabilities
		cliArguments.Parameters = createStackParams
		cliArguments.ParametersFile = createStackParametersFile
		cliArguments.Lint = createStackLint
		cliArguments.LinterConfiguration = createStackLintConfiguration
		cliArguments.EstimateCost = createStackEstimateCost

		// delete Stack
	case deleteStack.FullCommand():
		cliArguments.Mode = &DestroyStackMode
		cliArguments.Stack = deleteStackName

		// generate MFA token
	case mfaCommand.FullCommand():
		cliArguments.Mode = &MfaMode

		// update Stack
	case updateStack.FullCommand():
		cliArguments.Mode = &UpdateStackMode
		cliArguments.Stack = updateStackName
		cliArguments.TemplatePath = updateStackTemplate
		cliArguments.Capabilities = updateStackCapabilities
		cliArguments.ParametersFile = updateStackParametersFile
		cliArguments.Parameters = updateStackParams
		cliArguments.Lint = updateStackLint
		cliArguments.LinterConfiguration = updateStackLintConfiguration
		cliArguments.EstimateCost = updateStackEstimateCost

		// create Parameters
	case createParameters.FullCommand():
		cliArguments.Mode = &CreateParametersMode
		cliArguments.TemplatePath = createParametersTemplate
		cliArguments.OutputFilePath = createParametersParamsOutputFile
		cliArguments.Parameters = createParametersParams
		cliArguments.PrettyPrint = createParametersPrettyPrint

		// set stack policy
	case setStackPolicy.FullCommand():
		cliArguments.Mode = &SetStackPolicyMode
		cliArguments.Block = setDefaultBlockingStackPolicy
		cliArguments.Unblock = setDefaultUnblockingStackPolicy
		cliArguments.Stack = setStackPolicyName
		cliArguments.TemplatePath = setStackPolicyTemplate
		cliArguments.DisableStackTermination = setDisableStackTermination
		cliArguments.EnableStackTermination = setEnableStackTermination

	case createChangeSet.FullCommand():
		cliArguments.Mode = &CreateChangeSetMode
		if *createChangeSetName != "" {
			cliArguments.ChangeSet = createChangeSetName
		} else {
			ChangeSetDefaultName = *changeSetStackName + time.Now().Format("-2006-01-02--15-04-05")
			cliArguments.ChangeSet = &ChangeSetDefaultName
		}
		cliArguments.TemplatePath = changeSetTemplate
		cliArguments.Stack = changeSetStackName
		cliArguments.Parameters = createChangeSetParams
		cliArguments.ParametersFile = createChangeSetParametersFile
		cliArguments.Lint = createChangeSetLint
		cliArguments.LinterConfiguration = createChangeSetLintConfiguration
		cliArguments.EstimateCost = createChangeSetEstimateCost

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
	cliArguments.SkipValidation = noValidate

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
