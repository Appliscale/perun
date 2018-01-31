package intrinsicsolver

import "strings"

func fixSplit(cache *string) {
	if strings.Contains(*cache, "!Split") {
		getDivision := strings.Split(*cache, "!Split")

		if string(getDivision[1][0]) == " " {
			getDivision[1] = strings.Replace(getDivision[1], " ", "", 1)
		}

		*cache = "{ \"Fn::Split\" : " + getDivision[1] + " }"

	}
}
