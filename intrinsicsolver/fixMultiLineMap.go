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

/* Function fixMultiLineMap detects if a function is of a multi-line map nature by checking what follows the function name.
At the moment the goformation library is inappropriately handling the case where the function name is in the same line as the key and the body of a function isn't in the same line.
There are many ways to solve this problem but the fastest is to move the function name to the next line, indent it and transform it to it's full name.
Other solutions include rewriting the whole function and it's body in one line but due the lack of knowledge of how nested the map internal structure is and where it ends,
this solution is not chosen. */
func fixMultiLineMap(line *string, lines *[]string, idx int, name string) {
	pLines := *lines
	short := shortForm(name)
	long := longForm(name)
	full := fullForm(long)
	if strings.Contains(*line, short) && !strings.Contains(*line, "|") {
		split := strings.Split(*line, short)
		if idx+1 < len(pLines) {
			if strings.Contains(pLines[idx+1], "-") && (len(split) == 1 || split[1] == "") {
				// If so - we have multiple-level function with a body created of a map elements as the hyphen-noted structures.
				if strings.Contains(*line, ":") {
					// If so - we have key and a function name in one line. We have to relocate the function name into the next line, indent it and change it to the long form.
					nextLineIndents := countLeadingSpaces(pLines[idx+1])
					fullIndents := strings.Repeat(" ", nextLineIndents)
					replacement := "\n" + fullIndents + full
					*line = strings.Replace(*line, short, replacement, -1)
				} else {
					// If so - we have function as the element of another map - we assume that it is well indented so we only change the form to the long one.
					*line = strings.Replace(*line, short, full, -1)
				}
			}
		}
	}
}
