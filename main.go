package main

import (
	"github.com/Appliscale/cftool/cfcliparser"
	"github.com/Appliscale/cftool/cfconverter"
	"github.com/Appliscale/cftool/cfofflinevalidator"
	"github.com/Appliscale/cftool/cfonlinevalidator"
	"github.com/Appliscale/cftool/cfcontext"
)

func main() {
	context, err := cfcontext.GetContext()
	if err != nil {
		return
	}
	defer context.Logger.PrintErrors()


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