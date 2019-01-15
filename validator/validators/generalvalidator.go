package validators

import (
	"github.com/Appliscale/perun/context"
	"github.com/Appliscale/perun/logger"
	"github.com/Appliscale/perun/validator/template"
	"strconv"
)

type Restrictor func(string) (bool, string)
var defaultRestrictor Restrictor = func(propertyName string) (valid bool, msg string) { return true, "Should pass"}

func GetRestrictor(key string, ctx *context.Context) Restrictor {
	return defaultRestrictor
}

func GeneralValidateResourceByName(resourceValue template.Resource, resourceValidation *logger.ResourceValidation, ctx *context.Context) {
	for propertyName, propertyValue := range resourceValue.Properties {
		analyzeObject(propertyName, propertyValue, resourceValidation, propertyName + " ", ctx)
	}
}

func analyzeObject(propertyName string, propertyValue interface{}, resourceValidation *logger.ResourceValidation, preMessage string, ctx *context.Context) {
	var propertyRestrictor Restrictor
	switch propertyValue.(type) {
	case string:
		propertyRestrictor = GetRestrictor(propertyName, ctx)
		if valid, msg := propertyRestrictor(propertyValue.(string)); !valid{
			resourceValidation.AddValidationError(preMessage + ": " + msg + ", but the value is: \"" + propertyValue.(string) + "\"")
		}
		break
	case []interface{}:
		if isStringList(propertyValue.([]interface{})) {
			propertyRestrictor = GetRestrictor(propertyName, ctx)
			for index, value := range propertyValue.([]interface{}) {
				if valid, msg := propertyRestrictor(value.(string)); !valid{
					resourceValidation.AddValidationError(preMessage + " -> [" + strconv.Itoa(index) + "]: " + msg + ", but the value is: \"" + value.(string) + "\"")
				}
			}
		} else {
			for index, value := range propertyValue.([]interface{}) {
				analyzeObject(propertyName, value,  resourceValidation, preMessage + " -> [" + strconv.Itoa(index) + "]", ctx)
			}
		}
		break
	case map[string]interface{}:
		for k, v := range propertyValue.(map[string]interface{}){
			analyzeObject(k, v, resourceValidation, preMessage + " -> " + k,ctx)
		}
		break
	default:
		//Do nothing - preparser is not ideal and skips some properties with intristic functions
	}
}


func isStringList(list []interface{}) bool{
	for _, v := range list {
		switch v.(type) {
		case string:
			break
		default:
			return false
			break
		}
	}
	return true
}
