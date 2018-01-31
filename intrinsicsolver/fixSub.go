package intrinsicsolver

import (
	"strings"
)

func fixSub(cache *string, lines []string, idx int) {
	pLines := &lines
	if strings.Contains(*cache, "!Sub") && !strings.Contains(*cache, "|") {
		slice := strings.Fields(*cache)
		every := strings.SplitAfter(*cache, "")

		var i, spacesCount int

		for every[i] == " " {
			spacesCount++
			i++
		}

		spaces := strings.Repeat(" ", spacesCount)

		var key string

		subKeyValue := strings.SplitAfter(*cache, ":")
		stringSub := strings.Split(lines[idx+1], "- ")[1]
		valuesSub := strings.Split(lines[idx+2], "- ")[1]
		keyValue := strings.Split(valuesSub, ":")

		if strings.Contains(keyValue[0], "{") {
			key = strings.Replace(strings.SplitAfter(keyValue[0], "{")[1], " ", "", 1)
		} else {
			key = keyValue[0]
		}

		if string(keyValue[1][len(keyValue[1])-1]) == "}" {
			keyValue[1] = keyValue[1][:(len(keyValue[1]) - 1)]
		}

		deleteFromSlice((*pLines), idx, 1)
		deleteFromSlice((*pLines), idx, 1)

		var prefix, endParens string

		if slice[0] == "-" {
			if slice[1] != "!Sub" {
				prefix = subKeyValue[0] + " "
				endParens = "]}"
			} else {
				prefix = spaces + "- "
				endParens = "}]"
			}

		} else {
			prefix = subKeyValue[0] + " "
			endParens = "]}"
		}

		*cache = prefix + "{ \"Fn::Sub\" : [ " + strings.Replace(stringSub, "'", "\"", -1) + ", {" + "\"" + key + "\" :" + keyValue[1] + "}" + endParens

	}
}
