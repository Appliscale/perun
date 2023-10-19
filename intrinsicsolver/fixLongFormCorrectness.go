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
