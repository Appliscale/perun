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
	"encoding/json"
	"errors"
	"github.com/Appliscale/perun/cliparser"
	"github.com/Appliscale/perun/context"
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

	if *context.CliArguments.OutputFileFormat == cliparser.YAML {
		outputTemplate, err := toYAML(rawTemplate)
		if err != nil {
			return err
		}
		saveToFile(outputTemplate, *context.CliArguments.OutputFilePath, context.Logger)
	}

	if *context.CliArguments.OutputFileFormat == cliparser.JSON {
		var outputTemplate []byte
		if *context.CliArguments.PrettyPrint == false {
			outputTemplate, err = yamlToJSON(rawTemplate)
		} else if *context.CliArguments.PrettyPrint == true {
			outputTemplate, err = yamlToPrettyJSON(rawTemplate)
		}
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

func yamlToJSON(yamlTemplate []byte) ([]byte, error) {
	jsonTemplate, error := yaml.YAMLToJSON(yamlTemplate)
	if !govalidator.IsJSON(string(jsonTemplate)) {
		return nil, errors.New("This is not a valid YAML file")
	}
	return jsonTemplate, error
}

func yamlToPrettyJSON(yamlTemplate []byte) ([]byte, error) {
	var YAMLObj interface{}
	templateError := yaml.Unmarshal(yamlTemplate, &YAMLObj)

	jsonTemplate, templateError := json.MarshalIndent(YAMLObj, "", "    ")

	if !govalidator.IsJSON(string(jsonTemplate)) {
		return nil, errors.New("This is not a valid YAML file")
	}
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
