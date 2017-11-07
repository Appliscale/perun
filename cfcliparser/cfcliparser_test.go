package cfcliparser

import (
	"testing"
	"os"
	"github.com/stretchr/testify/assert"
)

func TestModeNotSpecified(t *testing.T) {
	assert.Equal(t, "You should specify what you want to do with --mode flag",
		parseCliArguments([]string{"cmd"}).Error())
}

func TestInvalidMode(t *testing.T) {
	assert.Equal(t, "Invalid mode. Use validate, validate_offline or convert",
		parseCliArguments([]string{"cmd", "-m=some_mode"}).Error())
}


func TestNoTemplatePath(t *testing.T) {
	assert.Equal(t, "You should specify a source of the template file with --template flag",
		parseCliArguments([]string{"cmd", "--mode=validate"}).Error())
}

func TestNoOutputPathInConvertMode(t *testing.T) {
	assert.Equal(t, "You should specify a output file path with --output flag",
		parseCliArguments([]string{"cmd", "--mode=convert", "--template=some_path"}).Error())
}

func TestNoOutputFormatInConvertMode(t *testing.T) {
	assert.Equal(t, "You should specify a output file format with --format flag",
		parseCliArguments([]string{"cmd", "--mode=convert", "--template=some_path", "--output=some_path"}).Error())
}

func TestInvalidOutputFormatInConvertMode(t *testing.T) {
	assert.Equal(t, "Invalid output file format. Use JSON or YAML",
		parseCliArguments([]string{"cmd", "--mode=convert", "--template=some_path", "--output=some_path", "--format=wrong_format"}).Error())
}

func TestValidArgs(t *testing.T) {
	assert.Nil(t, parseCliArguments([]string{"cmd", "--mode=validate_offline", "--template=some_path"}))
}

func parseCliArguments(args []string) error {
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()
	os.Args = args
	_, err := ParseCliArguments()
	return err
}