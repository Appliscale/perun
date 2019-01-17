package validators

import (
	"encoding/json"
	"fmt"

	"github.com/davecgh/go-spew/spew"

	"github.com/asaskevich/govalidator"
	"github.com/ghodss/yaml"

	"github.com/Appliscale/perun/intrinsicsolver"
	"github.com/Appliscale/perun/logger"
	"github.com/Appliscale/perun/offlinevalidator/template"
)

type fnStringInterface func(string, interface{})

func MapScanner(templateFile []byte, refTemplate template.Template, logger *logger.Logger) (bool, string, error) {

	matching := false

	if govalidator.IsJSON(string(templateFile)) {
		err := json.Unmarshal(templateFile, &refTemplate)
		if err != nil {
			return false, "", err
		}
	} else {
		preprocessed, preprocessingError := intrinsicsolver.FixFunctions(templateFile, logger, "multiline", "elongate", "correctlong")
		if preprocessingError != nil {
			logger.Error(preprocessingError.Error())
		}
		err := yaml.Unmarshal(preprocessed, &refTemplate)
		if err != nil {
			return false, "", err
		}
	}

	resources := refTemplate.Resources
	maps := make([]string, 0)
	initPath := make([]interface{}, 0)
	for _, resourceContent := range resources {
		for propertyName, propertyContent := range resourceContent.Properties {
			initPath = []interface{}{propertyName}
			findFindInMap(propertyName, propertyContent, &maps, initPath)
		}
	}

	mappings := refTemplate.Mappings
	whatMatches := ""
	fmt.Println()
	fmt.Println("MAPPINGS:")
	spew.Dump(mappings)
	fmt.Println()
	fmt.Println("RESOURCES:")
	spew.Dump(resources)
	fmt.Println()
	for mappingName, _ := range mappings {
		matching, whatMatches = isNotPresent(maps, mappingName)
	}

	return matching, whatMatches, nil
}

func findFindInMap(i interface{}, content interface{}, mapSlice *[]string, fullPath []interface{}) {
	if key, ok := i.(string); ok {
		if m, ok := content.(map[string]interface{}); ok {
			for mKey, mContent := range m {
				fullPath = append(fullPath, mKey)
				findFindInMap(mKey, mContent, mapSlice, fullPath)
			}
		} else if key == "Fn::FindInMap" {
			if args, ok := content.([]interface{}); ok {
				if mapName, ok := args[0].(string); ok {
					*mapSlice = append(*mapSlice, mapName)
				}
			}
		}
	} else if _, ok := i.(int); ok {
		if s, ok := content.([]interface{}); ok {
			for sIndex, sContent := range s {
				fullPath = append(fullPath, sIndex)
				findFindInMap(sIndex, sContent, mapSlice, fullPath)
			}
		}
	}
}

func isNotPresent(slice []string, match string) (bool, string) {
	for _, s := range slice {
		if s != match {
			return true, s
		}
	}
	return false, ""
}

func isPresent(slice []string, match string) bool {
	for _, s := range slice {
		if s == match {
			return true
		}
	}
	return false
}
