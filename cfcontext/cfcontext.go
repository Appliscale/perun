// Copyright 2017 Appliscale
//
// Maintainers and Contributors:
//
//   - Piotr Figwer (piotr.figwer@appliscale.io)
//   - Wojciech Gawroński (wojciech.gawronski@appliscale.io)
//   - Kacper Patro (kacper.patro@appliscale.io)
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

// Package cfcontext provides context for PerunCloud.
package cfcontext

import (
	"github.com/Appliscale/cftool/cfcliparser"
	"github.com/Appliscale/cftool/cflogger"
	"github.com/Appliscale/cftool/cfconfiguration"
)

type Context struct {
	CliArguments cfcliparser.CliArguments
	Logger* cflogger.Logger
	Config cfconfiguration.Configuration
}

// Create PerunCloud context.
func GetContext() (context Context, err error) {
	logger := cflogger.CreateDefaultLogger()

	cliArguments, err := cfcliparser.ParseCliArguments()
	if err != nil {
		logger.Error(err.Error())
		return
	}

	logger.Quiet = *cliArguments.Quiet
	logger.Yes = *cliArguments.Yes
	logger.SetVerbosity(*cliArguments.Verbosity)

	config, err := cfconfiguration.GetConfiguration(cliArguments, &logger)
	if err != nil {
		logger.Error(err.Error())
		return
	}

	context = Context{
		CliArguments: cliArguments,
		Logger: &logger,
		Config: config,
	}

	return
}
