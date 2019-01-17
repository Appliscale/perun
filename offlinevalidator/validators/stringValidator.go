package validators

import (
	"strconv"

	"github.com/Appliscale/perun/logger"
	"github.com/Appliscale/perun/offlinevalidator/template"
)

type funcArgKeyValueLogger func(string, string, *logger.ResourceValidation) bool

/*
StringValidator : takes the template Resources and checks all the `string` values using the certain, validation function and returns true if it would find any. Otherwise it returns false.
Assign this function result to a boolean variable and pass it to the warning function to print the check interpretation.
*/
func StringValidator(resourceContent template.Resource, resourceValidation *logger.ResourceValidation, fnKVL funcArgKeyValueLogger) bool {
	present := false // We assume that the template is well-formatted and thus the presence of the suspicious part is dubious.
	for propertyName, propertyValue := range resourceContent.Properties {
		initPath := []interface{}{propertyName}
		checkNested(propertyName, propertyValue, fnKVL, resourceValidation, initPath, &present)
	}
	return present
}

func checkNested(n interface{}, v interface{}, fnKVL funcArgKeyValueLogger, resourceValidation *logger.ResourceValidation, fullPath []interface{}, presence *bool) {
	if str, ok := v.(string); ok {
		where := ""
		for _, element := range fullPath {
			if stringElement, ok := element.(string); ok {
				if where != "" {
					where += ": " + stringElement
				} else {
					where = stringElement
				}
			} else if intElement, ok := element.(int); ok {
				where += "[" + strconv.Itoa(intElement) + "]"
			}
		}
		*presence = fnKVL(str, where, resourceValidation)
	} else if mp, ok := v.(map[string]interface{}); ok {
		for mpKey, mpValue := range mp {
			fullPath = append(fullPath, mpKey)
			checkNested(mpKey, mpValue, fnKVL, resourceValidation, fullPath, presence)
		}
	} else if slc, ok := v.([]interface{}); ok {
		for slcIdx, slcValue := range slc {
			fullPath = append(fullPath, slcIdx)
			checkNested(slcIdx, slcValue, fnKVL, resourceValidation, fullPath, presence)
		}
	}
}
