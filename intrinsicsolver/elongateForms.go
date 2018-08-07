package intrinsicsolver

import (
	"strings"
	"regexp"
)

/* Function elongateForms is investigating for short-form functions and changes them for their long equivalent. */
func elongateForms(line *string, lines *[]string, idx int, name string) {
	var currentFunctions int
	pLines := *lines
	totalFunctions := strings.Count(*line, "!")
	for (currentFunctions != totalFunctions+1) && !strings.Contains(*line, "#!") && strings.Contains(*line, "!") {
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

				line = addQuotes(short, split, line)

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

func addQuotes(short string, split []string, line *string) *string {
	// Function !Sub can take only a string in its short form - It has to be marked as string
	if short == "!Sub" {
		whiteSpaceTrimmed := strings.TrimSpace(split[1])
		if !regexp.MustCompile(`".*"`).MatchString(whiteSpaceTrimmed) && !strings.Contains(*line, "|") {
			*line = strings.Replace(*line, whiteSpaceTrimmed, ("\"" + whiteSpaceTrimmed + "\""), -1)
		}
	}
	return line
}
