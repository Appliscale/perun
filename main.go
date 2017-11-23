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

// A tool for CloudFormation template validation and conversion.
package main

import (
	"github.com/Appliscale/perun/cliparser"
	"github.com/Appliscale/perun/converter"
	"github.com/Appliscale/perun/offlinevalidator"
	"github.com/Appliscale/perun/onlinevalidator"
	"github.com/Appliscale/perun/context"
	"os"
)

func main() {
	context, err := context.GetContext()
	if err != nil {
		os.Exit(1)
	}

	if *context.CliArguments.Mode == cliparser.ValidateMode {
		valid := onlinevalidator.ValidateAndEstimateCosts(&context)
		if valid {
			os.Exit(0)
		} else {
			os.Exit(1)
		}
	}

	if *context.CliArguments.Mode == cliparser.ConvertMode {
		err := converter.Convert(&context)
		if err == nil {
			os.Exit(0)
		} else {
			context.Logger.Error(err.Error())
			os.Exit(1)
		}
	}

	if *context.CliArguments.Mode == cliparser.OfflineValidateMode {
		valid := offlinevalidator.Validate(&context)
		if valid {
			os.Exit(0)
		} else {
			os.Exit(1)
		}
	}

	os.Exit(0)
}
