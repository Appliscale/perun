package cfconverter

import (
	"github.com/ghodss/yaml"
	"github.com/asaskevich/govalidator"
	"errors"
	"io/ioutil"
	"github.com/Appliscale/cftool/cfcliparser"
	"os"
	"github.com/Appliscale/cftool/cflogger"
	"github.com/Appliscale/cftool/cfcontext"
)

func Convert(context *cfcontext.Context) {
	defer context.Logger.PrintErrors()

	rawTemplate, err := ioutil.ReadFile(*context.CliArguments.FilePath)
	if err != nil {
		context.Logger.LogError(err.Error())
		return
	}

	if *context.CliArguments.OutputFileFormat == cfcliparser.YAML {
		outputTemplate, err := toYAML(rawTemplate)
		if err != nil {
			context.Logger.LogError(err.Error())
			return
		}
		saveToFile(outputTemplate, *context.CliArguments.OutputFilePath, context.Logger)
	}

	if *context.CliArguments.OutputFileFormat == cfcliparser.JSON {
		outputTemplate, err := toJSON(rawTemplate)
		if err != nil {
			context.Logger.LogError(err.Error())
			return
		}
		saveToFile(outputTemplate, *context.CliArguments.OutputFilePath, context.Logger)
	}
}

func toYAML(jsonTemplate []byte) ([]byte, error) {
	if !govalidator.IsJSON(string(jsonTemplate)) {
		return nil, errors.New("This is not a valid JSON file")
	}

	yamlTemplate, error := yaml.JSONToYAML(jsonTemplate)

	return yamlTemplate, error
}

func toJSON(yamlTemplate []byte) ([]byte, error) {
	jsonTemplate, error := yaml.YAMLToJSON(yamlTemplate)

	if !govalidator.IsJSON(string(jsonTemplate)) {
		return nil, errors.New("This is not a valid YAML file")
	}

	return jsonTemplate, error
}

func saveToFile(template []byte, path string, logger *cflogger.Logger) {
	outputFile, err := os.Create(path)
	if err != nil {
		logger.LogError(err.Error())
		return
	}

	defer outputFile.Close()

	_, err = outputFile.Write(template)
	if err != nil {
		logger.LogError(err.Error())
		return
	}
}