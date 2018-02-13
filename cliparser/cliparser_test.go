// Copyright 2017 Appliscale
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

package cliparser

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestInvalidVerbosity(t *testing.T) {
	assert.Equal(t, "You specified invalid value for --verbosity flag",
		parseCliArguments([]string{"cmd", "validate", "some_path", "--verbosity=TEST"}).Error())
}

func TestTooSmallDurationForMFA(t *testing.T) {
	assert.Equal(t, "You should specify value for duration of MFA token greater than zero",
		parseCliArguments([]string{"cmd", "validate", "some_path", "--duration=-1"}).Error())
}

func TestTooBigDurationForMFA(t *testing.T) {
	assert.Equal(t, "You should specify value for duration of MFA token smaller than 129600 (3 hours)",
		parseCliArguments([]string{"cmd", "validate", "some_path", "--duration=50000000"}).Error())
}

func TestValidArgs(t *testing.T) {
	assert.Nil(t, parseCliArguments([]string{"cmd", "validate_offline", "some_path"}))
}

func parseCliArguments(args []string) error {
	_, err := ParseCliArguments(args)
	return err
}
