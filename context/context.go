// Copyright 2018 Appliscale
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

// Package context provides context for perun.
package context

import (
	"os"

	"github.com/Appliscale/perun/awsapi"
	"github.com/Appliscale/perun/cliparser"
	"github.com/Appliscale/perun/configuration"
	"github.com/Appliscale/perun/logger"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudformation"
)

// Context contains perun's logger, configuration, information about inconsistency
// between specification and documentation, and session.
type Context struct {
	CliArguments        cliparser.CliArguments
	Logger              *logger.Logger
	Config              configuration.Configuration
	InconsistencyConfig configuration.InconsistencyConfiguration
	CloudFormation      awsapi.CloudFormationAPI
	CurrentSession      *session.Session
}

type cliArgumentsParser func(args []string) (cliparser.CliArguments, error)
type configurationReader func(cliparser.CliArguments, *logger.Logger) (configuration.Configuration, error)
type inconsistenciesReader func(*logger.Logger) configuration.InconsistencyConfiguration

// GetContext creates CLI context. Creating logger and config and checking inconsistency.
func GetContext(cliArgParser cliArgumentsParser, confReader configurationReader, inconsistReader inconsistenciesReader) (context Context, err error) {
	myLogger := logger.CreateDefaultLogger()

	cliArguments, err := cliArgParser(os.Args)
	if err != nil {
		myLogger.Error(err.Error())
		return
	}

	if cliArguments.Quiet != nil {
		myLogger.Quiet = *cliArguments.Quiet
	}

	if cliArguments.Yes != nil {
		myLogger.Yes = *cliArguments.Yes
	}

	config, err := confReader(cliArguments, &myLogger)
	if err != nil {
		myLogger.Error(err.Error())
		return
	}

	myLogger.SetVerbosity(config.DefaultVerbosity)

	iconsistenciesConfig := inconsistReader(&myLogger)

	context = Context{
		CliArguments:        cliArguments,
		Logger:              &myLogger,
		Config:              config,
		InconsistencyConfig: iconsistenciesConfig,
	}
	return
}

// InitializeAwsAPI creates session.
func (context *Context) InitializeAwsAPI() {
	context.CurrentSession = InitializeSession(context)
	context.CloudFormation = awsapi.NewAWSCloudFormation(cloudformation.New(context.CurrentSession))
}
