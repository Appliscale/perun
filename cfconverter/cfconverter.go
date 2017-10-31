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

func Convert(context *cfcontext.Context) error {
	rawTemplate, err := ioutil.ReadFile(*context.CliArguments.TemplatePath)
	if err != nil {
		return err
	}

	if *context.CliArguments.OutputFileFormat == cfcliparser.YAML {
		outputTemplate, err := toYAML(rawTemplate)
		if err != nil {
			return err
		}
		saveToFile(outputTemplate, *context.CliArguments.OutputFilePath, context.Logger)
	}

	if *context.CliArguments.OutputFileFormat == cfcliparser.JSON {
		outputTemplate, err := toJSON(rawTemplate)
		if err != nil {
			return err
		}
		err = saveToFile(outputTemplate, *context.CliArguments.OutputFilePath, context.Logger)
		if err != nil {
			return err
		}
	}

	return nil
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

func saveToFile(template []byte, path string, logger *cflogger.Logger) error {
	outputFile, err := os.Create(path)
	if err != nil {
		return err
	}

	defer outputFile.Close()

	_, err = outputFile.Write(template)
	if err != nil {
		return err
	}

	return nil
}