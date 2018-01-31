package intrinsicsolver

import "strings"

func fixGetAtt(cache *string) {
	if strings.Contains(*cache, "!GetAtt") {
		keyValue := strings.SplitAfter(*cache, "!GetAtt")
		getValue := strings.Replace(keyValue[1], " ", "", -1)
		divideGetValue := strings.Split(getValue, ".")

		*cache = strings.Replace(keyValue[0], "!GetAtt", "{ \"Fn::GetAtt", -1) + "\" : " + "[" + "\"" + strings.Replace(divideGetValue[0], "\"", "", -1) + "\"" + ", " + "\"" + strings.Replace(divideGetValue[1], "\"", "", -1) + "\"" + " ] " + "}"

	}
}
