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
	"strings"
)

var ValidateMode = "validate"
var ConvertMode = "convert"
var OfflineValidateMode = "validate_offline"
var ConfigureMode = "configure"
var CreateStackMode = "create-stack"
var DestroyStackMode = "delete-stack"

const JSON = "json"
const YAML = "yaml"

type CliArguments struct {
	Mode              *string
	TemplatePath      *string
	OutputFilePath    *string
	OutputFileFormat  *string
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
	PrettyPrint       *bool
}

func availableFormats() []string {
	return []string{JSON, YAML}
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

		onlineValidate         = app.Command(ValidateMode, "Online template Validation")
		onlineValidateTemplate = onlineValidate.Arg("template", "A path to the template file.").Required().String()

		offlineValidate         = app.Command(OfflineValidateMode, "Offline Template Validation")
		offlineValidateTemplate = offlineValidate.Arg("template", "A path to the template file.").Required().String()

		convert             = app.Command(ConvertMode, "Convertion between JSON and YAML of template files")
		convertTemplate     = convert.Arg("template", "A path to the template file.").Required().String()
		convertOutputFile   = convert.Arg("output", "A path where converted file will be saved.").Required().String()
		convertOutputFormat = convert.Arg("format", "Output format: "+strings.ToUpper(JSON)+" | "+strings.ToUpper(YAML)+".").HintAction(availableFormats).Required().String()

		configure = app.Command(ConfigureMode, "Create your own configuration mode")

		createStack     = app.Command(CreateStackMode, "Creates a stack on aws")
		createStackName = createStack.Arg("stack", "An AWS stack name.").Required().String()

		deleteStack     = app.Command(DestroyStackMode, "Deletes a stack on aws")
		deleteStackName = deleteStack.Arg("stack", "An AWS stack name.").Required().String()
	)
	app.HelpFlag.Short('h')
	app.Version(utilities.VersionStatus())

	switch kingpin.MustParse(app.Parse(args[1:])) {

	//online validate
	case onlineValidate.FullCommand():
		cliArguments.Mode = &ValidateMode
		cliArguments.TemplatePath = onlineValidateTemplate

		// offline validation
	case offlineValidate.FullCommand():
		cliArguments.Mode = &OfflineValidateMode
		cliArguments.TemplatePath = offlineValidateTemplate

		// convert
	case convert.FullCommand():
		cliArguments.Mode = &ConvertMode
		cliArguments.TemplatePath = convertTemplate
		cliArguments.OutputFilePath = convertOutputFile
		cliArguments.OutputFileFormat = convertOutputFormat

		// configure
	case configure.FullCommand():
		cliArguments.Mode = &ConfigureMode

		// create Stack

	case createStack.FullCommand():
		cliArguments.Mode = &CreateStackMode
		cliArguments.Stack = createStackName

		// delete Stack
	case deleteStack.FullCommand():
		cliArguments.Mode = &DestroyStackMode
		cliArguments.Stack = deleteStackName
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

	if *cliArguments.Mode == ConvertMode {
		*cliArguments.OutputFileFormat = strings.ToLower(*cliArguments.OutputFileFormat)
		if *cliArguments.OutputFileFormat != JSON && *cliArguments.OutputFileFormat != YAML {
			err = errors.New("Invalid output file format. Use JSON or YAML")
			return
		}

	}

	return
}
