package intrinsicsolver

import (
	"strings"
)

/* Function elongateForms is investigating for short-form functions and changes them for their long equivalent. */
func elongateForms(line *string, lines *[]string, idx int, name string) {
	var currentFunctions int
	pLines := *lines
	totalFunctions := strings.Count(*line, "!")
	for (currentFunctions != totalFunctions+1) && !strings.Contains(*line, "#!/bin/bash") && strings.Contains(*line, "!") {
		short := shortForm(name)
		long := longForm(name)
		full := fullForm(long)
		split := strings.Split(*line, short)
		if idx+1 < len(pLines) {
			if strings.Contains(*line, name) && strings.Contains(pLines[idx+1], "-") && (len(split) != 2) {
				// If so - we don't have to surround it with quotes.
				if strings.Contains(*line, short) && !strings.Contains(*line, "|") {
					*line = strings.Replace(*line, short, full, -1)
				} else if strings.Contains(*line, short) && strings.Contains(*line, "|") {
					*line = strings.Replace(*line, (short + " |"), full, -1)
				}
			} else if strings.Contains(*line, name) {
				if strings.Contains(*line, short) && !strings.Contains(*line, "|") {
					*line = strings.Replace(*line, short, ("\"" + long + "\":"), -1)
				} else if strings.Contains(*line, short) && strings.Contains(*line, "|") {
					*line = strings.Replace(*line, (short + " |"), ("\"" + long + "\":"), -1)
				} else if strings.Contains(*line, full) && !strings.Contains(*line, "|") {
					*line = strings.Replace(*line, full, ("\"" + long + "\":"), -1)
				} else if strings.Contains(*line, full) && strings.Contains(*line, "|") {
					*line = strings.Replace(*line, (full + " |"), ("\"" + long + "\":"), -1)
				}
			}
		}
		currentFunctions++
	}

}
