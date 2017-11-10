// +build !windows

package cfconfiguration

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"os"
	"os/user"
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
		usr, _ := user.Current()
		assert.Equal(t, usr.HomeDir + "/.config/cftool/main.yaml", path, "Should contain user home")
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
		assert.Equal(t, "/etc/cftool/main.yaml", path, "Should contain /etc")
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
		assert.Equal(t, dir + "/.cftool", path, "Should contain current working directory")
		assert.True(t, ok, "Should exist")
	})

	t.Run("File does not exist", func(t *testing.T) {
		_, ok := getConfigFileFromCurrentWorkingDirectory(notExistStub)
		assert.False(t, ok, "Should not exist")
	})
}