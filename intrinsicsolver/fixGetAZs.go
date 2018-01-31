package intrinsicsolver

import "strings"

func fixGetAZs(cache *string) {
	if strings.Contains(*cache, "!GetAZs") {
		getDivision := strings.SplitAfter(*cache, "!GetAZs")

		if string(getDivision[1][0]) == " " {
			getDivision[1] = strings.Replace(getDivision[1], " ", "", 1)
		}

		if !strings.Contains(getDivision[1], "!") || !strings.Contains(getDivision[1], "{") || !strings.Contains(getDivision[1], "}") {
			getDivision[1] = "\"" + getDivision[1] + "\""
		}

		*cache = strings.Replace(getDivision[0], "!GetAZs", "{ \"Fn::GetAZs\" : ", -1) + getDivision[1] + " }"

	}
}
