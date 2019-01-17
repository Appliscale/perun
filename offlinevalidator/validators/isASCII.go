package validators

import (
	"github.com/Appliscale/perun/logger"
	"github.com/asaskevich/govalidator"
)

// IsASCII : function argument to the `validators.StringValidator`. Checks if the string is contains non-ASCII characters.
func IsASCII(value string, where string, resourceValidation *logger.ResourceValidation) bool {
	if !govalidator.IsASCII(value) {
		resourceValidation.AddValidationError(where + " value contains non-ASCII characters.")
	}
	return true
}

// ASCIIWarning : takes the result of the `validators.StringValidator` and if it is `true`, it warns that the strings exceeding the 64 character could crash the CloudFormation console.
func ASCIIWarning(result bool, sink *logger.Logger) {
	if result {
		sink.Warning("Encountered strings containing non-ASCII characters. There is some probability that this could make this template rejected by the CloudFormation.")
	}
}
