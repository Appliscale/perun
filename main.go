package main

import (
	"fmt"
	"github.com/Appliscale/cftool/cfcliparser"
	"github.com/Appliscale/cftool/cfconverter"
	"github.com/Appliscale/cftool/cfofflinevalidator"
	"github.com/Appliscale/cftool/cfonlinevalidator"
)

func main() {
	cliArguments, error := cfcliparser.ParseCliArguments()

	if error != nil {
		fmt.Println(error)
		return
	}

	if *cliArguments.Mode == cfcliparser.ValidateMode {
		cfonlinevalidator.ValidateAndEstimateCosts(cliArguments.FilePath, cliArguments.ConfigurationPath)
		return
	}

	if *cliArguments.Mode == cfcliparser.ConvertMode {
		cfconverter.Convert(cliArguments.FilePath, cliArguments.OutputFilePath, cliArguments.OutputFileFormat)
		return
	}

	if *cliArguments.Mode == cfcliparser.OfflineValidateMode {
		cfofflinevalidator.Validate(cliArguments.FilePath, cliArguments.ConfigurationPath)
		return
	}
}