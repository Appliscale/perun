package main

import (
	"github.com/Appliscale/cftool/cfonlinevalidator"
	"fmt"
	"github.com/Appliscale/cftool/cfcliparser"
	"github.com/Appliscale/cftool/cfconverter"
	"github.com/Appliscale/cftool/cfofflinevalidator"
	"github.com/Appliscale/cftool/cfspecification"
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

	if *cliArguments.Mode == cfcliparser.OfflineValidateMode {
		specification, err := cfspecification.GetEmbeddedSpecification(
			"CloudFormationResourceSpecification.json")
		if err != nil {
			fmt.Println(err)
			return
		}
		cfofflinevalidator.Validate(cliArguments.FilePath, &specification)
		return
	}
}