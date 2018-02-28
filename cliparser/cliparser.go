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
var UpdateStackMode = "update-stack"
var MfaMode = "mfa"
var SetupSinkMode = "setup-remote-sink"
var DestroySinkMode = "destroy-remote-sink"

const JSON = "json"
const YAML = "yaml"

type CliArguments struct {
	Mode              *string
	TemplatePath      *string
	OutputFilePath    *string
	ConfigurationPath *string
	Quiet             *bool
	Yes               *bool
	Verbosity         *string
	MFA               *bool
	DurationForMFA    *int64
	Profile           *string
	Region            *string
	Sandbox           *bool
	Stack             *string
	Capabilities      *[]string
	PrettyPrint       *bool
	Progress          *bool
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
		prettyPrint          = convert.Flag("pretty-print", "Pretty printing JSON").Bool()

		configure = app.Command(ConfigureMode, "Create your own configuration mode")

		createStack             = app.Command(CreateStackMode, "Creates a stack on aws")
		createStackName         = createStack.Arg("stack", "An AWS stack name.").String()
		createStackTemplate     = createStack.Arg("template", "A path to the template file.").String()
		createStackImpName      = createStack.Flag("stack", "Sn AWS stack name.").String()
		createStackImpTemplate  = createStack.Flag("template", "A path to the template file.").String()
		createStackCapabilities = createStack.Flag("capabilities", "Capabilities: CAPABILITY_IAM | CAPABILITY_NAMED_IAM").Enums("CAPABILITY_IAM", "CAPABILITY_NAMED_IAM")

		deleteStack        = app.Command(DestroyStackMode, "Deletes a stack on aws")
		deleteStackName    = deleteStack.Arg("stack", "An AWS stack name.").String()
		deleteStackImpName = deleteStack.Flag("stack", "An AWS stack name.").String()

		updateStack             = app.Command(UpdateStackMode, "Updates a stack on aws")
		updateStackName         = updateStack.Arg("stack", "An AWS stack name").String()
		updateStackTemplate     = updateStack.Arg("template", "A path to the template file.").String()
		updateStackImpName      = updateStack.Flag("stack", "Sn AWS stack name.").String()
		updateStackImpTemplate  = updateStack.Flag("template", "A path to the template file.").String()
		updateStackCapabilities = updateStack.Flag("capabilities", "Capabilities: CAPABILITY_IAM | CAPABILITY_NAMED_IAM").Enums("CAPABILITY_IAM", "CAPABILITY_NAMED_IAM")

		mfaCommand = app.Command(MfaMode, "Create temporary secure credentials with MFA.")

		setupSink = app.Command(SetupSinkMode, "Sets up resources required for progress report on stack events (SNS Topic, SQS Queue and SQS Queue Policy)")

		destroySink = app.Command(DestroySinkMode, "Destroys resources created with setup-remote-sink")
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
		cliArguments.PrettyPrint = prettyPrint

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
		cliArguments.Capabilities = createStackCapabilities

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

		// generate MFA token
	case mfaCommand.FullCommand():
		cliArguments.Mode = &MfaMode

		// update Stack
	case updateStack.FullCommand():
		cliArguments.Mode = &UpdateStackMode
		cliArguments.Capabilities = updateStackCapabilities
		if len(*updateStackImpTemplate) > 0 && len(*updateStackImpName) > 0 {
			cliArguments.Stack = updateStackImpName
			cliArguments.TemplatePath = updateStackImpTemplate
		} else if len(*updateStackName) > 0 && len(*updateStackTemplate) > 0 {
			cliArguments.Stack = updateStackName
			cliArguments.TemplatePath = updateStackTemplate
		} else if len(*updateStackName) > 0 && len(*updateStackImpTemplate) > 0 {
			cliArguments.Stack = updateStackName
			cliArguments.TemplatePath = updateStackImpTemplate
		} else {
			err = errors.New("You have to specify stack name and template file, try --help")
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
