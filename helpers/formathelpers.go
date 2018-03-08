package helpers

import (
	"encoding/json"
	"errors"
	"github.com/Appliscale/perun/intrinsicsolver"
	"github.com/Appliscale/perun/logger"
	"github.com/Appliscale/perun/offlinevalidator/template"
	"github.com/awslabs/goformation"
	"github.com/awslabs/goformation/cloudformation"
	"github.com/ghodss/yaml"
	"path"
	"strconv"
)

func GetParser(filename string) (func([]byte, template.Template, *logger.Logger) (cloudformation.Template, error), error) {
	templateFileExtension := path.Ext(filename)
	if templateFileExtension == ".json" {
		return ParseJSON, nil
	} else if templateFileExtension == ".yaml" || templateFileExtension == ".yml" {
		return ParseYAML, nil
	} else {
		return nil, errors.New("Invalid template file format.")
	}
}

func ParseJSON(templateFile []byte, refTemplate template.Template, logger *logger.Logger) (template cloudformation.Template, err error) {

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

func ParseYAML(templateFile []byte, refTemplate template.Template, logger *logger.Logger) (template cloudformation.Template, err error) {

	err = yaml.Unmarshal(templateFile, &refTemplate)
	if err != nil {
		return template, err
	}

	preprocessed, preprocessingError := intrinsicsolver.FixFunctions(templateFile, logger, "multiline", "elongate", "correctlong")
	if preprocessingError != nil {
		logger.Error(preprocessingError.Error())
	}
	tempYAML, err := goformation.ParseYAML(preprocessed)
	if err != nil {
		logger.Error(err.Error())
	}

	returnTemplate := *tempYAML

	return returnTemplate, nil
}

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
