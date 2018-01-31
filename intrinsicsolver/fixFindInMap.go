package intrinsicsolver

import (
	"regexp"
	"strings"
)

func fixFindInMap(cache *string) {
	if strings.Contains(*cache, "!FindInMap") {

		findInMap := strings.SplitAfter(*cache, "!FindInMap") //get the parameters
		divideEnd := strings.SplitAfter(findInMap[1], "]")    //!FindInMap sometimes is nested so we are trying to get only the tag without anything extra

		endFields := splitter(findInMap[1], ",][")

		isAlphaBeta, _ := regexp.Compile(`[A-Za-z]`)

		for i := 0; i < len(endFields); i++ {

			if !strings.Contains(endFields[i], "!") && isAlphaBeta.MatchString(endFields[i]) {
				endFields[i] = "\"" + strings.Replace(endFields[i], " ", "", -1) + "\""
			}

		}

		var ending string

		if divideEnd[1] != "" {
			ending = divideEnd[1]
		} else {
			ending = ""
		}

		if ending != "" {
			ending = ending + "}"
		}

		*cache = strings.Replace(findInMap[0], "!FindInMap", " { \"Fn::FindInMap", -1) + "\" : " + "[" + endFields[1] + ", " + endFields[2] + ", " + endFields[3] + "]" + " }" + ending

	}
}
