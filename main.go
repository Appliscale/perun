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

// A tool for CloudFormation template validation.
package main

import (
	"github.com/Appliscale/perun/checkingrequiredfiles"
	"github.com/Appliscale/perun/cliparser"
	"github.com/Appliscale/perun/configuration"
	"github.com/Appliscale/perun/configurator"
	"github.com/Appliscale/perun/context"
	"github.com/Appliscale/perun/estimatecost"
	"github.com/Appliscale/perun/linter"
	"github.com/Appliscale/perun/parameters"
	"github.com/Appliscale/perun/progress"
	"github.com/Appliscale/perun/stack"
	"github.com/Appliscale/perun/utilities"
	"github.com/Appliscale/perun/validator"
	"github.com/Appliscale/perun/validator/validators"
	"os"
)

func main() {
	ctx, err := context.GetContext(cliparser.ParseCliArguments, configuration.GetConfiguration, configuration.ReadInconsistencyConfiguration)
	checkingrequiredfiles.CheckingRequiredFiles(&ctx)

	if ctx.CliArguments.Lint != nil && *ctx.CliArguments.Lint {
		err = linter.CheckStyle(&ctx)
		if err != nil {
			os.Exit(1)
		}
	}

	if *ctx.CliArguments.Mode == cliparser.ValidateMode {
		ctx.InitializeAwsAPI()
		utilities.CheckFlagAndExit(validator.Validate(&ctx))
	}

	if *ctx.CliArguments.Mode == cliparser.ConfigureMode {
		configurator.CreateRequiredFilesInConfigureMode(&ctx)
		os.Exit(0)
	}

	if *ctx.CliArguments.Mode == cliparser.LintMode {
		err = linter.CheckStyle(&ctx)
		if err != nil {
			os.Exit(1)
		}
		os.Exit(0)
	}

	validationUnsuccessfullMsg := "To skip the validation part use the --no-validate flag"
	if *ctx.CliArguments.Mode == cliparser.CreateStackMode {
		ctx.InitializeAwsAPI()
		if *ctx.CliArguments.SkipValidation || (validator.Validate(&ctx) && validators.UserDecideGeneralRule(&ctx)) {
			utilities.CheckErrorCodeAndExit(stack.NewStack(&ctx))
		} else {
			ctx.Logger.Info(validationUnsuccessfullMsg)
		}
	}

	if *ctx.CliArguments.Mode == cliparser.DestroyStackMode {
		ctx.InitializeAwsAPI()
		utilities.CheckErrorCodeAndExit(stack.DestroyStack(&ctx))

	}

	if *ctx.CliArguments.Mode == cliparser.MfaMode {
		err := context.UpdateSessionToken(ctx.Config.DefaultProfile, ctx.Config.DefaultRegion, ctx.Config.DefaultDurationForMFA, &ctx)
		if err == nil {
			os.Exit(0)
		} else {
			ctx.Logger.Error(err.Error())
			os.Exit(1)
		}
	}

	if *ctx.CliArguments.Mode == cliparser.CreateChangeSetMode {
		ctx.InitializeAwsAPI()
		if *ctx.CliArguments.SkipValidation || (validator.Validate(&ctx) && validators.UserDecideGeneralRule(&ctx)) {
			err := stack.NewChangeSet(&ctx)
			if err != nil {
				ctx.Logger.Error(err.Error())
			}
		} else {
			ctx.Logger.Info(validationUnsuccessfullMsg)
		}
	}

	if *ctx.CliArguments.Mode == cliparser.DeleteChangeSetMode {
		ctx.InitializeAwsAPI()
		utilities.CheckErrorCodeAndExit(stack.DeleteChangeSet(&ctx))
	}

	if *ctx.CliArguments.Mode == cliparser.UpdateStackMode {
		ctx.InitializeAwsAPI()
		if *ctx.CliArguments.SkipValidation || (validator.Validate(&ctx) && validators.UserDecideGeneralRule(&ctx)) {
			utilities.CheckErrorCodeAndExit(stack.UpdateStack(&ctx))
		} else {
			ctx.Logger.Info(validationUnsuccessfullMsg)
		}

	}

	if *ctx.CliArguments.Mode == cliparser.SetupSinkMode {
		progress.ConfigureRemoteSink(&ctx)
		os.Exit(0)
	}

	if *ctx.CliArguments.Mode == cliparser.DestroySinkMode {
		progress.DestroyRemoteSink(&ctx)
		os.Exit(0)
	}

	if *ctx.CliArguments.Mode == cliparser.CreateParametersMode {
		parameters.ConfigureParameters(&ctx)
		os.Exit(0)
	}

	if *ctx.CliArguments.Mode == cliparser.SetStackPolicyMode {
		ctx.InitializeAwsAPI()
		if *ctx.CliArguments.DisableStackTermination || *ctx.CliArguments.EnableStackTermination {
			utilities.CheckErrorCodeAndExit(stack.SetTerminationProtection(&ctx))
		} else {
			utilities.CheckErrorCodeAndExit(stack.ApplyStackPolicy(&ctx))
		}
	}

	if *ctx.CliArguments.Mode == cliparser.EstimateCostMode {
		ctx.InitializeAwsAPI()
		estimatecost.EstimateCosts(&ctx)
	}
}
