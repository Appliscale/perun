package main

import (
	"github.com/Appliscale/cftool/cfonlinevalidator"
	"fmt"
	"github.com/Appliscale/cftool/cfcliparser"
	"github.com/Appliscale/cftool/cfconverter"
)

func main() {
	cliArguments, error := cfcliparser.ParseCliArguments()

	if error != nil {
		fmt.Println(error)
		return
	}

	if *cliArguments.Mode == cfcliparser.ValidateMode {
		cfonlinevalidator.ValidateAndEstimateCosts(cliArguments.FilePath, cliArguments.Region)
		return
	}

	if *cliArguments.Mode == cfcliparser.ConvertMode {
		cfconverter.Convert(cliArguments.FilePath, cliArguments.OutputFilePath, cliArguments.OutputFileFormat)
		return
	}
}