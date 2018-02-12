package intrinsicsolver

import (
	"strings"
)

/* Unfortunately the short-to-long-form function names exchange isn't solving the issue of YAML being ready for the YAML-JSON conversion.
In some cases the parser is misinterpretating function in it's long form with additional key and throws an error. We must enclose functions in curly braces. */
func fixLongFormCorrectness(line *string) {
	keyValue := strings.SplitAfterN(*line, ":", 2)
	if len(keyValue) == 2 && !strings.Contains(keyValue[0], "Fn:") {
		if strings.Contains(keyValue[1], "\"Fn::") && !strings.Contains(keyValue[0], "Fn") {
			*line = strings.Replace(*line, keyValue[1], (" {" + keyValue[1] + "}"), 1)
		} else if strings.Contains(keyValue[1], "\"Ref") && !strings.Contains(keyValue[0], "Ref") {
			*line = strings.Replace(*line, keyValue[1], (" {" + keyValue[1] + "}"), 1)
		}
	}
}
