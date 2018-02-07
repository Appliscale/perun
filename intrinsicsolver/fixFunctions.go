package intrinsicsolver

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

// FixFunctions : takes []byte file and firstly converts all single quotation marks to double ones (anything between single ones is treated as the rune in GoLang), then deconstructs file into lines, checks for intrinsic functions of a map nature where the function name is located in one line and it's body (map elements) are located in the following lines (if this would be not fixed an error would be thrown: `json: unsupported type: map[interface {}]interface {}`). The function changes the notation by putting function name in the next line with proper indentation and saves the result to temporary file, then opens it and returns []byte array.
func FixFunctions(template []byte) []byte {
	var quotationProcessed, temporaryResult []string
	preLines := parseFileIntoLines(template)

	// All single quotation marks are transformed to double ones.
	for _, e := range preLines {
		var fixed string
		if strings.Contains(e, "'") {
			fixed = strings.Replace(e, "'", "\"", -1)
		} else {
			fixed = e
		}

		quotationProcessed = append(quotationProcessed, fixed)
	}

	lines := quotationProcessed

	// These are the YAML short names of a functions which take the arguments in a form of a map.
	multiLiners := []string{"!FindInMap", "!Join", "!Select", "!Split", "!Sub", "!And", "!Equals", "!If", "!Not", "!Or"}

	for idx, d := range lines {
		for _, ml := range multiLiners {
			fixMultiLineMap(&d, &lines, idx, ml)
		}

		temporaryResult = append(temporaryResult, d)
	}

	// Function writeLines saves the processed result to a file (if there would be any errors, it could be investigated there).
	if err := writeLines(temporaryResult, "preprocessed.yml"); err != nil {
		log.Fatalf("writeLines: %s", err)
	}

	// Then the temporary result is opened and returned as a []byte.
	preprocessedTemplate, err := ioutil.ReadFile("preprocessed.yml")
	if err != nil {
		fmt.Println(err)
	}

	return preprocessedTemplate
}

// Function parseFileIntoLines is reading the []byte file and returns it line by line as []string slice.
func parseFileIntoLines(template []byte) []string {
	bytesReader := bytes.NewReader(template)
	lines := make([]string, 0)
	scanner := bufio.NewScanner(bytesReader)

	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return lines
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
