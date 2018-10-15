package helpers

import (
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
