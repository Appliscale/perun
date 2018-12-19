package helpers

import (
	"github.com/awslabs/goformation/cloudformation"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

func TestGetParser(t *testing.T) {
	parser, err := GetParser("myfile.json")
	assert.Equal(t, reflect.ValueOf(ParseJSON), reflect.ValueOf(parser))
	assert.Nil(t, err)

	parser, err = GetParser("myfile.yaml")
	assert.Equal(t, reflect.ValueOf(ParseYAML), reflect.ValueOf(parser))
	assert.Nil(t, err)

	parser, err = GetParser("myfile.yml")
	assert.Equal(t, reflect.ValueOf(ParseYAML), reflect.ValueOf(parser))
	assert.Nil(t, err)

	parser, err = GetParser("myfile.alamakota")
	assert.NotEqual(t, reflect.ValueOf(ParseYAML), reflect.ValueOf(parser))
	assert.NotEqual(t, reflect.ValueOf(ParseJSON), reflect.ValueOf(parser))
	assert.NotNil(t, err)
}

func TestCountLeadingSpaces(t *testing.T) {
	assert.Equal(t, 4, CountLeadingSpaces("    dskjhfmasjkd"))
	assert.Equal(t, 0, CountLeadingSpaces("ajksdasd"))
	assert.Equal(t, 0, CountLeadingSpaces("a    sdkajsd"))
	assert.Equal(t, 0, CountLeadingSpaces(",    "))
}

func TestLineAndCharacter(t *testing.T) {
	line, character := lineAndCharacter("0123456789asd", 8)
	assert.Equal(t, 1, line)
	assert.Equal(t, 9, character)

	line, character = lineAndCharacter("0123456789\nasd", 11)
	assert.Equal(t, 2, line)
	assert.Equal(t, 1, character)
}

func TestFindFnImportValue(t *testing.T) {
	var templateFile = make([]byte, 0)
	tempYAML := cloudformation.Template{}
	err := findFnImportValue(templateFile, &tempYAML)
	assert.Nilf(t, err, "Error should be nil")
}

func TestReplaceImportValue(t *testing.T) {
	path := []string{"Resorce", "Property", "Name", "Value"}
	cfTemplate := cloudformation.Template{}

	err := replaceImportValue(path, &cfTemplate)
	assert.NotNilf(t, err, "Error should not be nil")

}

func TestAddToPathAndReplace(t *testing.T) {
	path := []string{}
	name := "name"
	value := "value"
	tempYAML := cloudformation.Template{}
	startPath := []string{}

	err := addToPathAndReplace(path, name, value, &tempYAML, startPath)
	assert.Nilf(t, err, "Error should be nil")

}
