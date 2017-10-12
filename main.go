package main

import (
	"github.com/appliscale/cftool/cfonlinevalidator"
	"fmt"
	"github.com/appliscale/cftool/cfcliparser"
	"github.com/appliscale/cftool/cfconverter"
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