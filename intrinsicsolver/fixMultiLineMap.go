package intrinsicsolver

import (
	"strings"
)

// Function fixMultiLineMap detects if a function is of a multi-line map nature by checking what follows the function name. At the moment the goformation library is inappropriately handling the case where the function name is in the same line as the key and the body of a function isn't in the same line. There are many ways to solve this problem but the fastest is to move the function name to the next line, indent it and transform it to it's full name. Other solutions include rewriting the whole function and it's body in one line but due the lack of knowledge of how nested the map internal structure is and where it ends, this solution is not chosen.
func fixMultiLineMap(line *string, lines *[]string, idx int, tagName string) {
	pLines := *lines
	longName := "Fn::" + strings.Split(tagName, "!")[1] + ":"
	if strings.Contains(*line, tagName) && !strings.Contains(*line, "|") {
		split := strings.Split(*line, tagName)
		if strings.Contains(pLines[idx+1], "-") && split[1] == "" {
			// If so - we have multiple-level function with a body created of a map elements as the hyphen-noted structures.
			if strings.Contains(*line, ":") {
				// If so - we have key and a function name in one line. We have to relocate the function name into the next line, indent it and change it to the long form.
				nextLineIndents := indentations(pLines[idx+1])
				fullIndents := strings.Repeat(" ", nextLineIndents)
				replacement := "\n" + fullIndents + longName
				*line = strings.Replace(*line, tagName, replacement, -1)
			} else {
				// If so - we have function as the element of another map - we assume that it is well indented so we only change the form to the long one.
				*line = strings.Replace(*line, tagName, longName, -1)
			}
		}
	}
}
