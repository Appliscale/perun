package intrinsicsolver

import (
	"regexp"
	"strings"
)

func fixEquals(cache *string) {
	if strings.Contains(*cache, "!Equals") {
		keyValue := strings.SplitAfter(*cache, ":")                 //get slice of key : value
		equalsAndRest := strings.SplitAfter(keyValue[1], "!Equals") //split after "!Equals" tag to get the tag's parameters
		regB := regexp.MustCompile(`\[([^()]*)\]`)
		insideBrackets := strings.Join(regB.FindAllString(equalsAndRest[1], 1), "") //retrieve everything inside []
		recut := insideBrackets[1 : len(insideBrackets)-1]                          //get rid of the brackets (we have onlu parameters)

		commaSeparate := strings.SplitN(recut, ",", 2)

		reg := regexp.MustCompile(`}|{|!|:`) //checks if there are plain values or e.g. functions inside
		value1 := strings.Replace(commaSeparate[0], " ", "", -1)
		value2 := commaSeparate[1]

		if !reg.MatchString(value1) {
			if !strings.Contains(value1, "\"") {
				value1 = "\"" + strings.Replace(value1, " ", "", 1) + "\""
			}
		}

		if !reg.MatchString(value2) {
			if !strings.Contains(value2, "\"") {
				value2 = "\"" + strings.Replace(value2, " ", "", 1) + "\""
			}
		}

		*cache = keyValue[0] + " {\"Fn::Equals\" : [" + value1 + ", " + value2 + "]" + "}"

	}
}
