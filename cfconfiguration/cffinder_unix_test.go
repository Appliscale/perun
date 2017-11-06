// +build !windows

package cfconfiguration

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"os"
	"os/user"
	"testing"
)

func TestMain(m *testing.M) {
	retCode := m.Run()
	os.Exit(retCode)
}

func TestGetUserConfigFile(t *testing.T) {
	t.Run("File exist", func(t *testing.T) {
		tmp := checkFileExistence
		defer func() { checkFileExistence = tmp }()
		checkFileExistence = func(name string) (os.FileInfo, error) {
			return nil, nil
		}

		path, ok := getUserConfigFile()
		usr, _ := user.Current()
		assert.Equal(t, usr.HomeDir + "/.config/cftool/config.yaml", path, "Should contain user home")
		assert.True(t, ok, "Should exist")
	})

	t.Run("File do not exist", func(t *testing.T) {
		tmp := checkFileExistence
		defer func() { checkFileExistence = tmp }()
		checkFileExistence = func(name string) (os.FileInfo, error) {
			return nil, errors.New("")
		}

		_, ok := getConfigFileFromProjectRoot()
		assert.False(t, ok, "Should not exist")
	})
}

func TestGetGlobalConfigFile(t *testing.T) {
	t.Run("File exist", func(t *testing.T) {
		tmp := checkFileExistence
		defer func() { checkFileExistence = tmp }()
		checkFileExistence = func(name string) (os.FileInfo, error) {
			return nil, nil
		}

		path, ok := getGlobalConfigFile()
		assert.Equal(t, "/etc/.Appliscale/cftool/config.yaml", path, "Should contain /etc")
		assert.True(t, ok, "Should exist")
	})

	t.Run("File do not exist", func(t *testing.T) {
		tmp := checkFileExistence
		defer func() { checkFileExistence = tmp }()
		checkFileExistence = func(name string) (os.FileInfo, error) {
			return nil, errors.New("")
		}

		_, ok := getConfigFileFromProjectRoot()
		assert.False(t, ok, "Should not exist")
	})
}

func TestGetConfigFileFromProjectRoot(t *testing.T) {
	t.Run("File exist", func(t *testing.T) {
		tmp := checkFileExistence
		defer func() { checkFileExistence = tmp }()
		checkFileExistence = func(name string) (os.FileInfo, error) {
			return nil, nil
		}

		_, ok := getConfigFileFromProjectRoot()
		assert.True(t, ok, "Should exist")
	})

	t.Run("File do not exist", func(t *testing.T) {
		tmp := checkFileExistence
		defer func() { checkFileExistence = tmp }()
		checkFileExistence = func(name string) (os.FileInfo, error) {
			return nil, errors.New("")
		}

		_, ok := getConfigFileFromProjectRoot()
		assert.False(t, ok, "Should not exist")
	})
}