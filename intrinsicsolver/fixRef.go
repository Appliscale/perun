package intrinsicsolver

import "strings"

func fixRef(cache *string) {
	if strings.Contains(*cache, "!Ref") {
		letters := strings.Split(*cache, "") //split line into char's
		counters := make([]int, 0)
		enders := make([]int, 0)

		//get the index of letter where the !Ref tag starts to be
		for o := 0; o < len(letters); o++ {
			if string(letters[o]) == "!" && string(letters[o+1]) == "R" {
				counters = append(counters, o)
			}
		}

		//then count the indexes form the position where the !Ref tag starts
		for _, d := range counters {
			c := 0
			for q := d; q < len(letters); q++ {

				if string(letters[q]) == "\"" || string(letters[q]) == "," { //if it encounters the !Ref parameter or it's end, count +1
					c++
				}

				if c == 2 { //if count == 2, then we found the parameter begining (+1) and the end (+1) == 2; we can register the position
					enders = append(enders, q)
					break
				}

			}
		}

		sl := make([]string, 0)

		for j := 0; j < len(counters); j++ {
			begin := counters[j]
			end := enders[j]
			sl = append(sl, strings.Join(letters[begin:end], "")) //get the extract of '!Reg PARAMETER' from the whole line
		}

		for _, d := range sl {
			if string(d[len(d)-1]) != "\"" && strings.Contains(d, "\"") {
				d = d + "\""
			}

			var value string

			if strings.Contains(d, "\"") {
				value = strings.Split(d, " ")[1]
			} else {
				value = "\"" + strings.SplitN(d, " ", 2)[1] + "\""
			}

			var endTag string

			if strings.Contains(value, ",") { //if there was some lost comma by accident, replace it after the parameter
				value = strings.Replace(value, ",", "", -1)
				endTag = ","
			}

			exchange := "{ \"Ref\"" + " : " + value + " }" + endTag

			*cache = strings.Replace(*cache, d, exchange, -1)

		}
	}

}
