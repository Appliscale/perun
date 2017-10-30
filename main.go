package main

import (
	"github.com/Appliscale/cftool/cfcliparser"
	"github.com/Appliscale/cftool/cfconverter"
	"github.com/Appliscale/cftool/cfofflinevalidator"
	"github.com/Appliscale/cftool/cfonlinevalidator"
	"github.com/Appliscale/cftool/cfcontext"
	"github.com/Appliscale/cftool/cflogger"
)

func main() {
	logger := cflogger.Logger{}
	defer logger.PrintErrors()

	cliArguments, err := cfcliparser.ParseCliArguments()
	if err != nil {
		logger.LogError(err.Error())
		return
	}

	context := cfcontext.Context{
		CliArguments: cliArguments,
		Logger: &logger,
	}

	if *context.CliArguments.Mode == cfcliparser.ValidateMode {
		cfonlinevalidator.ValidateAndEstimateCosts(&context)
		return
	}

	if *context.CliArguments.Mode == cfcliparser.ConvertMode {
		cfconverter.Convert(&context)
		return
	}

	if *context.CliArguments.Mode == cfcliparser.OfflineValidateMode {
		cfofflinevalidator.Validate(&context)
		return
	}
}