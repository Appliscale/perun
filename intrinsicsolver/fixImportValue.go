package intrinsicsolver

import "strings"

func fixImportValue(cache *string) {
	if strings.Contains(*cache, "!ImportValue") {
		getDivision := strings.Split(*cache, "!ImportValue")

		*cache = getDivision[0] + " { \"Fn::ImportValue\" : " + getDivision[1] + "}"

	}
}
