package cfconverter

import (
	"github.com/ghodss/yaml"
	"github.com/asaskevich/govalidator"
	"errors"
	"io/ioutil"
	"github.com/Appliscale/cftool/cfcliparser"
	"os"
	"github.com/Appliscale/cftool/cflogger"
)

func Convert(sourceFilePath *string, destinationFilePath *string, format *string) {
	logger := cflogger.Logger{}
	defer cflogger.PrintErrors(&logger)

	rawTemplate, err := ioutil.ReadFile(*sourceFilePath)
	if err != nil {
		cflogger.LogError(&logger, err.Error())
		return
	}

	if *format == cfcliparser.YAML {
		outputTemplate, err := toYAML(rawTemplate)
		if err != nil {
			cflogger.LogError(&logger, err.Error())
			return
		}
		saveToFile(outputTemplate, destinationFilePath, &logger)
	}

	if *format == cfcliparser.JSON {
		outputTemplate, err := toJSON(rawTemplate)
		if err != nil {
			cflogger.LogError(&logger, err.Error())
			return
		}
		saveToFile(outputTemplate, destinationFilePath, &logger)
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

func saveToFile(template []byte, path *string, logger *cflogger.Logger) {
	outputFile, err := os.Create(*path)
	if err != nil {
		cflogger.LogError(logger, err.Error())
		return
	}

	defer outputFile.Close()

	_, err = outputFile.Write(template)
	if err != nil {
		cflogger.LogError(logger, err.Error())
		return
	}
}