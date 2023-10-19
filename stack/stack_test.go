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
package stack

import (
	"github.com/Appliscale/perun/stack/stack_mocks"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestGetTemplateFromFile(t *testing.T) {
	stackName := "StackName"
	templatePath := "./test_resources/test_template.yaml"
	context := stack_mocks.SetupContext(t, []string{"cmd", "create-stack", stackName, templatePath})

	templateBody := stack_mocks.ReadFile(t, templatePath)

	returnedTemplate, returnedStackName, err := getTemplateFromFile(context)
	if err != nil {
		t.Error(err.Error())
	}

	assert.Equal(t, returnedTemplate, templateBody)
	assert.Equal(t, returnedStackName, stackName)
}

func TestGetPath(t *testing.T) {
	templatePath := "./test_resources/test_template.yaml"
	stackName := "StackName"

	path := getPathForMode(t, []string{"cmd", "create-stack", stackName, templatePath})
	assert.Equal(t, path, templatePath)
	path = getPathForMode(t, []string{"cmd", "update-stack", stackName, templatePath})
	assert.Equal(t, path, templatePath)
	path = getPathForMode(t, []string{"cmd", "create-change-set", stackName, templatePath})
	assert.Equal(t, path, templatePath)

	path = getPathForMode(t, []string{"cmd", "set-stack-policy", stackName, templatePath, "--block"})
	assert.True(t, strings.HasSuffix(path, "/.config/perun/stack-policies/blocked.json"))

	path = getPathForMode(t, []string{"cmd", "set-stack-policy", stackName, templatePath, "--unblock"})
	assert.True(t, strings.HasSuffix(path, "/.config/perun/stack-policies/unblocked.json"))
}

func TestIsStackPolicyFileJSON(t *testing.T) {
	assert.False(t, isStackPolicyFileJSON("policyjson"))
	assert.True(t, isStackPolicyFileJSON("policy.json"))
	assert.False(t, isStackPolicyFileJSON("asd.yaml"))
}

func getPathForMode(t *testing.T, args []string) string {
	context := stack_mocks.SetupContext(t, args)
	path, err := getPath(context)
	if err != nil {
		t.Error(err.Error())
	}
	return path
}
