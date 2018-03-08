// Copyright 2017 Appliscale
//
// Maintainers and contributors are listed in README file inside repository.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package converter provides tools for JSON/YAML CloudFormation templates
// conversion.
package converter

import (
	"errors"
	"github.com/Appliscale/perun/context"
	"github.com/Appliscale/perun/helpers"
	"github.com/Appliscale/perun/intrinsicsolver"
	"github.com/Appliscale/perun/logger"
	"github.com/asaskevich/govalidator"
	"github.com/ghodss/yaml"
	"io/ioutil"
	"os"
)

// Read template from the file, convert it and check if it has valid structure.
// Then save converted template to file.
func Convert(context *context.Context) error {
	rawTemplate, err := ioutil.ReadFile(*context.CliArguments.TemplatePath)
	if err != nil {
		return err
	}
	format := detectFormatFromContent(rawTemplate)
	var outputTemplate []byte

	// If input type file is JSON convert to YAML.
	if format == "JSON" {
		outputTemplate, err = jsonToYaml(rawTemplate)
		if err != nil {
			return err
		}
		saveToFile(outputTemplate, *context.CliArguments.OutputFilePath, context.Logger)

		// If input type file is YAML, check all functions and create JSON (with or not --pretty-print flag).
	} else if format == "YAML" {
		preprocessed, preprocessingError := intrinsicsolver.FixFunctions(rawTemplate, context.Logger, "multiline", "elongate", "correctlong")
		if preprocessingError != nil {
			context.Logger.Error(preprocessingError.Error())
		}
		if *context.CliArguments.PrettyPrint == false {
			outputTemplate, err = yamlToJson(preprocessed)
		} else if *context.CliArguments.PrettyPrint == true {
			outputTemplate, err = yamlToPrettyJson(preprocessed)
		}
		if err != nil {
			return err
		}
		err = saveToFile(outputTemplate, *context.CliArguments.OutputFilePath, context.Logger)
		if err != nil {
			return err
		}
	} else {
		context.Logger.Always(format)
		return nil
	}

	return nil
}

func jsonToYaml(jsonTemplate []byte) ([]byte, error) {
	if !govalidator.IsJSON(string(jsonTemplate)) {
		return nil, errors.New("This is not a valid JSON file")
	}

	yamlTemplate, err := yaml.JSONToYAML(jsonTemplate)

	return yamlTemplate, err
}

func yamlToJson(yamlTemplate []byte) ([]byte, error) {
	jsonTemplate, err := yaml.YAMLToJSON(yamlTemplate)
	return jsonTemplate, err
}

func yamlToPrettyJson(yamlTemplate []byte) ([]byte, error) {
	var YAMLObj interface{}
	templateError := yaml.Unmarshal(yamlTemplate, &YAMLObj)

	jsonTemplate, templateError := helpers.PrettyPrintJSON(YAMLObj)

	return jsonTemplate, templateError

}

func saveToFile(template []byte, path string, logger *logger.Logger) error {
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

func detectFormatFromContent(rawTemplate []byte) (format string) {
	_, errorYAML := jsonToYaml(rawTemplate)
	_, errorJSON := yamlToJson(rawTemplate)

	if errorYAML == nil {
		return "JSON"
	} else if errorJSON == nil {
		return "YAML"
	}
	return "Unsupported file format. The input file must be either a valid JSON or YAML file."
}
