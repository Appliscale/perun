package cfcliparser

import (
	"errors"
	"strings"
	"gopkg.in/alecthomas/kingpin.v2"
)

const ValidateMode string = "validate"
const ConvertMode string = "convert"
const OfflineValidateMode = "validate_offline"
const JSON string = "json"
const YAML string = "yaml"

type CliArguments struct {
	Mode              *string
	TemplatePath      *string
	OutputFilePath    *string
	OutputFileFormat  *string
	ConfigurationPath *string
	Quiet             *bool
	Yes               *bool
	Verbosity         *string
}

func ParseCliArguments() (cliArguments CliArguments, err error) {

	cliArguments.Mode = kingpin.Flag("mode", ValidateMode+"|"+OfflineValidateMode+"|"+ConvertMode).Short('f').String()
	cliArguments.TemplatePath = kingpin.Flag("template", "A path to the template").Short('t').String()
	cliArguments.OutputFilePath = kingpin.Flag("output", "A path, where converted file will be saved").Short('o').String()
	cliArguments.OutputFileFormat = kingpin.Flag("format", "Output format: " + strings.ToUpper(JSON)+ "|"+ strings.ToUpper(YAML)).Short('x').String()
	cliArguments.ConfigurationPath = kingpin.Flag("config", "A path to the configuration file").Short('c').String()
	cliArguments.Quiet = kingpin.Flag("quiet", "No console output, just return code").Short('q').Bool()
	cliArguments.Yes = kingpin.Flag("yes", "Always say yes").Short('y').Bool()
	cliArguments.Verbosity = kingpin.Flag("verbosity", "TRACE|DEBUG|INFO|ERROR").Short('v').String()

	kingpin.Parse()

	if *cliArguments.Mode == "" {
		err = errors.New("You should specify what you want to do with -mode flag")
		return
	}

	if *cliArguments.Mode != ValidateMode && *cliArguments.Mode != ConvertMode && *cliArguments.Mode != OfflineValidateMode {
		err = errors.New("Invalid mode. Use validate or convert")
		return
	}

	if *cliArguments.TemplatePath == "" {
		err = errors.New("You should specify a source of the template file with -file flag")
		return
	}

	if *cliArguments.Mode == ConvertMode {
		if *cliArguments.OutputFilePath == "" {
			err = errors.New("You should specify a output file path with -output flag")
			return
		}

		if *cliArguments.OutputFileFormat == "" {
			err = errors.New("You should specify a output file format with -format flag")
			return
		}

		*cliArguments.OutputFileFormat = strings.ToLower(*cliArguments.OutputFileFormat)
		if *cliArguments.OutputFileFormat != JSON && *cliArguments.OutputFileFormat != YAML {
			err = errors.New("Invalid mode. Use validate or convert")
			return
		}
	}

	return
}
