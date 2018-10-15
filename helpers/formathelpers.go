// Copyright 2018 Appliscale
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

// Package helpers has some useful functions to choose parser and ease scan maps and slices.
package helpers

import (
	"encoding/json"
	"errors"
	"path"
	"strconv"

	"github.com/Appliscale/perun/intrinsicsolver"
	"github.com/Appliscale/perun/logger"
	"github.com/Appliscale/perun/validator/template"
	"github.com/awslabs/goformation"
	"github.com/awslabs/goformation/cloudformation"
	"github.com/ghodss/yaml"
	"regexp"
)

// GetParser chooses parser based on file extension.
func GetParser(filename string) (func([]byte, template.Template, logger.LoggerInt) (cloudformation.Template, error), error) {
	templateFileExtension := path.Ext(filename)
	if templateFileExtension == ".json" {
		return ParseJSON, nil
	} else if templateFileExtension == ".yaml" || templateFileExtension == ".yml" {
		return ParseYAML, nil
	} else {
		return nil, errors.New("Invalid template file format.")
	}
}

// ParseJSON parses JSON template file to cloudformation template.
func ParseJSON(templateFile []byte, refTemplate template.Template, logger logger.LoggerInt) (template cloudformation.Template, err error) {
	err = json.Unmarshal(templateFile, &refTemplate)
	if err != nil {
		if syntaxError, isSyntaxError := err.(*json.SyntaxError); isSyntaxError {
			syntaxOffset := int(syntaxError.Offset)
			line, character := lineAndCharacter(string(templateFile), syntaxOffset)
			logger.Error("Syntax error at line " + strconv.Itoa(line) + ", column " + strconv.Itoa(character))
		} else if typeError, isTypeError := err.(*json.UnmarshalTypeError); isTypeError {
			typeOffset := int(typeError.Offset)
			line, character := lineAndCharacter(string(templateFile), typeOffset)
			logger.Error("Type error at line " + strconv.Itoa(line) + ", column " + strconv.Itoa(character))
		}
		return template, err
	}

	tempJSON, err := goformation.ParseJSON(templateFile)
	if err != nil {
		logger.Error(err.Error())
	}

	returnTemplate := *tempJSON

	return returnTemplate, nil
}

// ParseYAML parses YAML template file to cloudformation template.
func ParseYAML(templateFile []byte, refTemplate template.Template, logger logger.LoggerInt) (template cloudformation.Template, err error) {
	err = yaml.Unmarshal(templateFile, &refTemplate)
	if err != nil {
		return template, err
	}

	for resource := range refTemplate.Resources {
		var validDeletionPolicy = regexp.MustCompile("(^$)|(Delete)$|(Retain)$|(Snapshot)$")
		if !validDeletionPolicy.MatchString(refTemplate.Resources[resource].DeletionPolicy) {
			err = errors.New("Deletion Policy in resource: " + resource + " has to be a string literal, cannot be parametrized")
		}
	}
	preprocessed, preprocessingError := intrinsicsolver.FixFunctions(templateFile, logger, "multiline", "elongate", "correctlong")
	if preprocessingError != nil {
		logger.Error(preprocessingError.Error())
	}

	tempYAML, parseError := goformation.ParseYAML(preprocessed)
	if parseError != nil {
		return *cloudformation.NewTemplate(), parseError
	}

	returnTemplate := *tempYAML

	return returnTemplate, err
}

// PrettyPrintJSON prepares JSON file with indent to ease reading it.
func PrettyPrintJSON(toPrint interface{}) ([]byte, error) {
	return json.MarshalIndent(toPrint, "", "    ")
}

func lineAndCharacter(input string, offset int) (line int, character int) {
	lf := rune(0x0A)

	if offset > len(input) || offset < 0 {
		return 0, 0
	}

	line = 1

	for i, b := range input {
		if b == lf {
			if i < offset {
				line++
				character = 0
			}
		} else {
			character++
		}
		if i == offset {
			break
		}
	}
	return line, character
}

// CountLeadingSpaces counts leading spaces. It's used in checkYamlIndentation() to find indentation error in template.
func CountLeadingSpaces(line string) int {
	i := 0
	for _, runeValue := range line {
		if runeValue == ' ' {
			i++
		} else {
			break
		}
	}
	return i
}
