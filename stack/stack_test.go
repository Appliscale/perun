package stack

import (
	"github.com/Appliscale/perun/stack/mocks"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestGetTemplateFromFile(t *testing.T) {
	stackName := "StackName"
	templatePath := "./test_resources/test_template.yaml"
	context := mocks.SetupContext(t, []string{"cmd", "create-stack", stackName, templatePath})

	templateBody := mocks.ReadFile(t, templatePath)

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
	context := mocks.SetupContext(t, args)
	path, err := getPath(context)
	if err != nil {
		t.Error(err.Error())
	}
	return path
}
