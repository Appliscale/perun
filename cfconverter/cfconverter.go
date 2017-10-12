package cfconverter

import (
	"github.com/ghodss/yaml"
	"github.com/asaskevich/govalidator"
	"errors"
	"io/ioutil"
	"fmt"
	"github.com/Appliscale/cftool/cfcliparser"
	"os"
)

func Convert(sourceFilePath *string, destinationFilePath *string, format *string) {
	rawTemplate, error := ioutil.ReadFile(*sourceFilePath)
	if error != nil {
		fmt.Println(error)
		return
	}

	if *format == cfcliparser.YAML {
		outputTemplate, error := toYAML(rawTemplate)
		if error != nil {
			fmt.Println(error)
			return
		}
		saveToFile(outputTemplate, destinationFilePath)
	}

	if *format == cfcliparser.JSON {
		outputTemplate, error := toJSON(rawTemplate)
		if error != nil {
			fmt.Println(error)
			return
		}
		saveToFile(outputTemplate, destinationFilePath)
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

func saveToFile(template []byte, path *string) {
	outputFile, error := os.Create(*path)
	if error != nil {
		fmt.Println(error)
		return
	}

	defer outputFile.Close()

	_, error = outputFile.Write(template)
	if error != nil {
		fmt.Println(error)
		return
	}
}