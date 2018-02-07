package intrinsicsolver

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/Appliscale/perun/logger"
)

/*
FixFunctions : takes []byte file and firstly converts all single quotation marks to double ones (anything between single ones is treated as the rune in GoLang),
then deconstructs file into lines, checks for intrinsic functions of a map nature where the function name is located in one line and it's body (map elements)
are located in the following lines (if this would be not fixed an error would be thrown: `json: unsupported type: map[interface {}]interface {}`).
The function changes the notation by putting function name in the next line with proper indentation and saves the result to temporary file,
then opens it and returns []byte array.
*/
func FixFunctions(template []byte, logger *logger.Logger) ([]byte, error) {
	var quotationProcessed, temporaryResult []string
	preLines, err := parseFileIntoLines(template, logger)

	// All single quotation marks are transformed to double ones.
	for _, line := range preLines {
		var fixed string
		if strings.Contains(line, "'") {
			fixed = strings.Replace(line, "'", "\"", -1)
		} else {
			fixed = line
		}

		quotationProcessed = append(quotationProcessed, fixed)
	}

	lines := quotationProcessed

	// These are the YAML short names of a functions which take the arguments in a form of a map.
	multiLiners := []string{"!FindInMap", "!Join", "!Select", "!Split", "!Sub", "!And", "!Equals", "!If", "!Not", "!Or"}

	for idx, d := range lines {
		for _, function := range multiLiners {
			fixMultiLineMap(&d, &lines, idx, function)
		}

		temporaryResult = append(temporaryResult, d)
	}

	// Function writeLines saves the processed result to a file (if there would be any errors, it could be investigated there).
	if err := writeLines(temporaryResult, "preprocessed.yml"); err != nil {
		logger.Error(err.Error())
		return nil, err
	}

	// Then the temporary result is opened and returned as a []byte.
	preprocessedTemplate, err := ioutil.ReadFile("preprocessed.yml")
	if err != nil {
		logger.Error(err.Error())
		return preprocessedTemplate, err
	}

	return preprocessedTemplate, nil
}

// Function parseFileIntoLines is reading the []byte file and returns it line by line as []string slice.
func parseFileIntoLines(template []byte, logger *logger.Logger) ([]string, error) {
	bytesReader := bytes.NewReader(template)
	lines := make([]string, 0)
	scanner := bufio.NewScanner(bytesReader)

	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		logger.Error(err.Error())
		return nil, err
	}

	return lines, nil
}

// Function writeLines takes []string slice and writes it element by element as line by line file
func writeLines(lines []string, path string) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	w := bufio.NewWriter(file)
	for _, line := range lines {
		fmt.Fprintln(w, line)
	}

	return w.Flush()
}
