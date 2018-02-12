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
	"io/ioutil"
	"os"
	"testing"

	"github.com/Appliscale/perun/logger"
	"github.com/stretchr/testify/assert"
)

var sink logger.Logger

func setup() {
	sink = logger.Logger{}
}

func TestMain(m *testing.M) {
	setup()
	retCode := m.Run()
	os.Exit(retCode)
}

func TestIndentations(t *testing.T) {
	line := "                Key: Value       "
	lineIndent := indentations(line)
	firstLetter := string(line[lineIndent])
	assert.Equal(t, 16, lineIndent, "MSG")
	assert.Equal(t, "K", firstLetter, "MSG")
}

func TestMultiline(t *testing.T) {
	rawTemplate, _ := ioutil.ReadFile("./test_resources/test_map.yaml")
	expectedTemplate, _ := ioutil.ReadFile("./test_resources/manual_test_map.yaml")
	fixed, _ := FixFunctions(rawTemplate, &sink, "multiline")
	expected, _ := parseFileIntoLines(expectedTemplate, &sink)
	actual, _ := parseFileIntoLines(fixed, &sink)

	if string(actual[len(actual)-1]) == "" {
		actual = actual[:(len(actual) - 1)]
	}

	assert.Equal(t, expected, actual, "MSG")
}

func TestElongate(t *testing.T) {
	rawTemplate, _ := ioutil.ReadFile("./test_resources/test_elongate.yaml")
	expectedTemplate, _ := ioutil.ReadFile("./test_resources/manual_test_elongate.yaml")
	fixed, _ := FixFunctions(rawTemplate, &sink, "elongate")
	expected, _ := parseFileIntoLines(expectedTemplate, &sink)
	actual, _ := parseFileIntoLines(fixed, &sink)

	if string(actual[len(actual)-1]) == "" {
		actual = actual[:(len(actual) - 1)]
	}

	assert.Equal(t, expected, actual, "MSG")
}

func TestCorrectLong(t *testing.T) {
	rawTemplate, _ := ioutil.ReadFile("./test_resources/manual_test_elongate.yaml")
	expectedTemplate, _ := ioutil.ReadFile("./test_resources/manual_test_correctlong.yaml")
	fixed, _ := FixFunctions(rawTemplate, &sink, "correctlong")
	expected, _ := parseFileIntoLines(expectedTemplate, &sink)
	actual, _ := parseFileIntoLines(fixed, &sink)

	if string(actual[len(actual)-1]) == "" {
		actual = actual[:(len(actual) - 1)]
	}

	assert.Equal(t, expected, actual, "MSG")
}
