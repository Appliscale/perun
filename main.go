package main

import (
	"github.com/Appliscale/cftool/cfcliparser"
	"github.com/Appliscale/cftool/cfconverter"
	"github.com/Appliscale/cftool/cfofflinevalidator"
	"github.com/Appliscale/cftool/cfonlinevalidator"
	"github.com/Appliscale/cftool/cfcontext"
	"os"
)

func main() {
	context, err := cfcontext.GetContext()
	if err != nil {
		os.Exit(1)
	}

	if *context.CliArguments.Mode == cfcliparser.ValidateMode {
		valid := cfonlinevalidator.ValidateAndEstimateCosts(&context)
		if valid {
			os.Exit(0)
		} else {
			os.Exit(1)
		}
	}

	if *context.CliArguments.Mode == cfcliparser.ConvertMode {
		err := cfconverter.Convert(&context)
		if err == nil {
			os.Exit(0)
		} else {
			context.Logger.Error(err.Error())
			os.Exit(1)
		}
	}

	if *context.CliArguments.Mode == cfcliparser.OfflineValidateMode {
		valid := cfofflinevalidator.Validate(&context)
		if valid {
			os.Exit(0)
		} else {
			os.Exit(1)
		}
	}

	os.Exit(0)
}