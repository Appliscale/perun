package validators

import (
	"github.com/Appliscale/perun/logger"
)

// IsLonger64 : function argument to the `validators.StringValidator`. Checks if the string is longer than 64 characters.
func IsLonger64(value string, where string, resourceValidation *logger.ResourceValidation) bool {
	if len(value) > 64 {
		resourceValidation.AddValidationError(where + " value longer than 64 characters.")
	}
	return true
}

// Longer64Warning : takes the result of the `validators.StringValidator` and if it is `true`, it warns that the strings exceeding the 64 character could crash the CloudFormation console.
func Longer64Warning(result bool, sink *logger.Logger) {
	if result {
		sink.Warning("Encountered strings longer than 64 characters. There is some probability that this could make this template rejected by the CloudFormation.")
	}
}
