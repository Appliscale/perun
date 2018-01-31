package intrinsicsolver

import "strings"

func fixSelect(cache *string, lines []string, idx int) {
	pLines := &lines
	var getIndex, getList, rest string
	if strings.Contains(*cache, "!Select") {

		getDivision := strings.SplitAfter(*cache, "!Select")

		if getDivision[1] == "" || getDivision[1] == " " {
			//then multi-line
			fullIndex := lines[idx+1]
			fullList := lines[idx+2]
			getIndex = strings.Replace(strings.SplitAfter(fullIndex, "- ")[1], " ", "", 1)
			getList = strings.Replace(strings.SplitAfter(fullList, "- ")[1], " ", "", 1)

			deleteFromSlice((*pLines), idx, 1)
			deleteFromSlice((*pLines), idx, 1)

			rest = "[" + getIndex + ", " + getList + "]" + "}"

		} else {
			//then single-line
			rest = getDivision[1] + "}"
		}

		*cache = strings.Replace(getDivision[0], "!Select", "{ \"Fn::Select\" :", -1) + rest

	}
}
