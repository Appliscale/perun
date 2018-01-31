package intrinsicsolver

import "strings"

func fixIf(cache *string, lines []string, idx int) {
	pLines := &lines
	if strings.Contains(*cache, "!If") {
		var s int
		var appendResult, joined string
		if string(lines[idx][len(lines[idx])-1]) == "[" {
			//then multi-line !If tag

			ifSlice := make([]string, 0)
			for !(strings.Contains((lines[idx+s]), "]") && !strings.Contains((lines[idx+s]), ",")) {
				var t int
				for string(lines[idx+s][t]) == " " {
					t++
				}
				matched := strings.Replace((lines[idx+s]), " ", "", t)
				ifSlice = append(ifSlice, matched)
				s++
			}

			joined = strings.Join(ifSlice, "")

		} else {
			//then single-line !If tag
			var t int
			for string(lines[idx][t]) == " " {
				t++
			}
			joined = strings.Replace((lines[idx]), " ", "", t)

		}

		joinedChars := strings.Split(joined, "")

		var ignore int
		externalCommaSlice := make([]int, 0)
		commaFixed := make([]string, 0)

		//get indexes of char's that are commas INSIDE !If parameters (not commas that are separators of those parameters)
		for i, d := range joinedChars {
			if d == "{" {
				ignore = 1
			}
			if d == "}" {
				ignore = 0
			}
			if ignore == 1 {
				if d == "," {
					externalCommaSlice = append(externalCommaSlice, i)
				}
			}
		}

		for _, d := range externalCommaSlice {
			joinedChars[d] = "#" //replace those commas with #
		}

		joinedANew := strings.Join(joinedChars, "")
		threeElemesSplit := strings.Split(joinedANew, ",") //split only the parameters (avoid splitting by parameters internal commas)

		for i, d := range threeElemesSplit {
			if i == 0 {
				d = strings.SplitAfter(d, "[")[1]
			}
			commaFixed = append(commaFixed, strings.Replace(d, "#", ",", -1)) //after all replace #'es with commas to get the original construction
		}

		lineSlice := make([]string, 0)
		for _, d := range commaFixed {
			//for every of three elements
			if strings.Contains(d, "{") || strings.Contains(d, "}") {
				//then it is complex expression in curly brackets

				//look for nested maps
				chars := strings.Split(d, "")

				var ignore int
				externalSlice := make([]int, 0)

				//get indexes of char's that are commas INSIDE parameters maps to be preserved
				for i, d := range chars {
					if d == "[" {
						ignore = 1
					}
					if d == "]" {
						ignore = 0
					}
					if ignore == 1 {
						if d == "," {
							externalSlice = append(externalSlice, i)
						}
					}
				}

				for _, d := range externalSlice {
					chars[d] = "#"
				}

				joinedANew := strings.Join(chars, "")

				subElementsSlice := make([]string, 0)
				d = strings.Replace(strings.Replace(joinedANew, "}", "", -1), "{", "", -1)
				commaSplit := strings.Split(d, ",")
				for _, d := range commaSplit {
					//for every comma-separated expression
					if strings.Count(d, ":") == 1 && !strings.Contains(d, "!") {
						//then it is key : value without intrinsic functions inside
						colonSplit := strings.Split(d, ":")
						key := "\"" + strings.Replace(colonSplit[0], " ", "", -1) + "\""
						var value string
						if isNumeric(strings.Replace(colonSplit[1], " ", "", -1)) {
							//if value is numeric
							value = strings.Replace(colonSplit[1], " ", "", -1)
						} else if strings.Contains(colonSplit[1], "[") && strings.Contains(colonSplit[1], "]") {
							//if value is a map
							value = colonSplit[1]
						} else {
							value = "\"" + strings.Replace(colonSplit[1], " ", "", -1) + "\""
						}
						toAppend := key + " : " + value
						subElementsSlice = append(subElementsSlice, toAppend)
					} else {
						var value string
						//then it is key : value with function inside
						colonSplit := strings.Split(d, ":")
						key := "\"" + strings.Replace(colonSplit[0], " ", "", -1) + "\""
						value = colonSplit[1]
						toAppend := key + " : " + value
						subElementsSlice = append(subElementsSlice, toAppend)
					}
				}
				connected := "{" + strings.Join(subElementsSlice, ",") + "}"
				lineSlice = append(lineSlice, connected)
			} else if !strings.Contains(d, ":") && !strings.Contains(d, "!") {
				//then it is only value
				if string(d[len(d)-1]) == "]" {
					d = strings.Replace(d, "]", "", -1)
				}
				value := "\"" + d + "\""
				lineSlice = append(lineSlice, value)

			} else {
				//then it is expression with function
				lineSlice = append(lineSlice, d)
			}
		}

		appendResult = strings.Replace(strings.Join(lineSlice, ", "), "#", ",", -1) //exchange the #'es with commas inside paameters internal structures

		keyValue := strings.SplitAfter(*cache, "!If")

		for i := 0; i < s; i++ {
			deleteFromSlice((*pLines), idx, 1)
		}

		*cache = strings.Replace(keyValue[0], "!If", "{ \"Fn::If\" : ", -1) + "[" + appendResult + "]" + "}"

	}

}
