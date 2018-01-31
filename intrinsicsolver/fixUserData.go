package intrinsicsolver

import (
	"regexp"
	"strings"
)

func fixUserData(cache *string, lines []string, idx int) {
	pLines := &lines
	if strings.Contains(*cache, "\"Fn::Base64\":") || strings.Contains(*cache, "!Base64") || strings.Contains(*cache, "Fn::Base64:") {
		if strings.Contains(lines[idx+1], "!Sub |") {

			every := strings.SplitAfter(*cache, "")

			var i, spacesCount, s, t, start, end int

			for every[i] == " " {
				spacesCount++
				i++
			}

			spaces := strings.Repeat(" ", spacesCount)

			reg, _ := regexp.Compile(`[A-Za-z/#!]+`)
			userDataSlice := make([]string, 0)

			for reg.MatchString(lines[idx+s]) {
				for string(lines[idx+s][t]) == " " {
					t++
				}
				matched := "\"" + strings.Replace((lines[idx+s]), " ", "", t) + "\""
				userDataSlice = append(userDataSlice, matched)
				s++
			}

			for i, l := range userDataSlice {

				if strings.Contains(l, "#cloud-config") || strings.Contains(l, "#!/bin/bash") {
					start = i
				}

				if string(l[0:2]) == "\"/" {
					end = i
					break
				}

			}

			userScripts := strings.Join(userDataSlice[start:end], ",")
			cfnInit := userDataSlice[end]
			cfnSignal := userDataSlice[end+1]

			for i := 0; i < (end + 1); i++ {
				deleteFromSlice((*pLines), idx, 1)
			}

			*cache = spaces + "{\"Fn::Base64\" : " + "{ \"Fn::Join\" : [\"\n\", [" + userScripts + "," + "{ \"Fn::Sub\" : " + cfnInit + "}," + "{ \"Fn::Sub\" : " + cfnSignal + "}]]}" + "}"

		}
	}

}
