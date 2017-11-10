// +build windows

package cfconfiguration

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func existStub(string) (os.FileInfo, error) {
	return nil, nil
}

func notExistStub(string) (os.FileInfo, error) {
	return nil, errors.New("")
}

func TestGetUserConfigFile(t *testing.T) {
	t.Run("File exist", func(t *testing.T) {
		path, ok := getUserConfigFile(existStub)
		envVal, _ := os.LookupEnv("LOCALAPPDATA")
		assert.Equal(t, envVal + "\\cftool\\main.yaml", path, "Should contain Local")
		assert.True(t, ok, "Should exist")
	})

	t.Run("File does not exist", func(t *testing.T) {
		_, ok := getUserConfigFile(notExistStub)
		assert.False(t, ok, "Should not exist")
	})
}

func TestGetGlobalConfigFile(t *testing.T) {
	t.Run("File exist", func(t *testing.T) {
		path, ok := getGlobalConfigFile(existStub)
		envVal, _ := os.LookupEnv("ALLUSERSPROFILE")
		assert.Equal(t, envVal + "\\cftool\\main.yaml", path, "Should contain ProgramData")
		assert.True(t, ok, "Should exist")
	})

	t.Run("File does not exist", func(t *testing.T) {
		_, ok := getGlobalConfigFile(notExistStub)
		assert.False(t, ok, "Should not exist")
	})
}

func TestGetConfigFileFromCurrentWorkingDirectory(t *testing.T) {
	t.Run("File exist", func(t *testing.T) {
		path, ok := getConfigFileFromCurrentWorkingDirectory(existStub)
		dir, _ := os.Getwd()
		assert.Equal(t, dir + "\\.cftool", path, "Should contain current working directory")
		assert.True(t, ok, "Should exist")
	})

	t.Run("File does not exist", func(t *testing.T) {
		_, ok := getConfigFileFromCurrentWorkingDirectory(notExistStub)
		assert.False(t, ok, "Should not exist")
	})
}