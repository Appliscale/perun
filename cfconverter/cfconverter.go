// Copyright 2017 Appliscale
//
// Maintainers and Contributors:
//
//   - Piotr Figwer (piotr.figwer@appliscale.io)
//   - Wojciech Gawroński (wojciech.gawronski@appliscale.io)
//   - Kacper Patro (kacper.patro@appliscale.io)
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

// Package cfconverter provides tools for JSON/YAML cloudformation templates conversion.
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

// Read template from the file, convert it and check if it has valid structure. Then save converted template to file.
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
