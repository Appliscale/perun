// Copyright 2018 Appliscale
//
// Maintainers and contributors are listed in README file inside repository.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package intrinsicsolver

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"strings"

	"github.com/Appliscale/perun/logger"
)

var Functions = []string{"Base64", "GetAtt", "GetAZs", "ImportValue", "Ref", "FindInMap", "Join", "Select", "Split", "Sub", "And", "Equals", "If", "Not", "Or"}
var mapNature = Functions[5:]

/*
FixFunctions : takes []byte file and firstly converts all single quotation marks to double ones (anything between single ones is treated as the rune in GoLang),
then deconstructs file into lines, checks for intrinsic functions. The FixFunctions has modes: `multiline`, `elongate` and `correctlong`.
Mode `multiline` looks for functions of a map nature where the function name is located in one line and it's body (map elements)
are located in the following lines (if this would be not fixed an error would be thrown: `json: unsupported type: map[interface {}]interface {}`).
The function changes the notation by putting function name in the next line with proper indentation.
Mode `elongate` exchanges the short function names into their proper, long equivalent.
Mode `correctlong` prepares the file for conversion into JSON. If the file is a YAML with every line being solicitously indented, there is no problem and the `elongate` mode is all we need.
But if there is any mixed notation (e.g. indented maps along with one-line maps, functions in one line with the key), parsing must be preceded with some additional operations.
The result is returned as a []byte array.
*/
func FixFunctions(template []byte, logger logger.LoggerInt, mode ...string) ([]byte, error) {
	var quotationProcessed, temporaryResult []string
	preLines, err := parseFileIntoLines(template, logger)
	if err != nil {
		logger.Error(err.Error())
		return nil, err
	}

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

	// In case the intrinsic function is in the last line and the the next line is investigated in search for it's multi-line body, we have to add one, blank line.
	quotationProcessed = append(quotationProcessed, "")

	lines := quotationProcessed

	for idx, d := range lines {
		for _, m := range mode {
			if m == "multiline" {
				for _, function := range mapNature {
					fixMultiLineMap(&d, &lines, idx, function)
				}
			}
			if m == "elongate" {
				for _, function := range Functions {
					elongateForms(&d, &lines, idx, function)
				}
			}
			if m == "correctlong" {
				fixLongFormCorrectness(&d)
			}
		}

		temporaryResult = append(temporaryResult, d)
	}

	stringStream := strings.Join(temporaryResult, "\n")
	output := []byte(stringStream)

	return output, nil
}

// Expands the function name to it's long form without a colon. For example - Fn::FindInMap.
func longForm(name string) string {
	var fullName string
	if name != "Ref" {
		fullName = "Fn::" + name
	} else {
		fullName = name
	}
	return fullName
}

/* Expands the function name by adding a colon. For example - Fn::FindInMap:.
It is crucial to pass here the output from the longForm function.*/
func fullForm(name string) string {
	return (name + ":")
}

// Expands the function name to it's short form. For example - !FindInMap.
func shortForm(name string) string {
	return ("!" + name)
}

// Function parseFileIntoLines is reading the []byte file and returns it line by line as []string slice.
func parseFileIntoLines(template []byte, logger logger.LoggerInt) ([]string, error) {
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
