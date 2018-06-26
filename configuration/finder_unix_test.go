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

// +build !windows

package configuration

import (
	"errors"
	"os"
	"os/user"
	"testing"

	"github.com/stretchr/testify/assert"
)

func existStub(string) (os.FileInfo, error) {
	return nil, nil
}

func notExistStub(string) (os.FileInfo, error) {
	return nil, errors.New("")
}

func TestGetUserConfigFile(t *testing.T) {
	t.Run("File exist", func(t *testing.T) {
		path, ok := getUserConfigFile(existStub, "main.yaml")
		usr, _ := user.Current()
		assert.Equal(t, usr.HomeDir+"/.config/perun/main.yaml", path, "Should contain user home")
		assert.True(t, ok, "Should exist")
	})

	t.Run("File does not exist", func(t *testing.T) {
		_, ok := getUserConfigFile(notExistStub, "main.yaml")
		assert.False(t, ok, "Should not exist")
	})
}

func TestGetGlobalConfigFile(t *testing.T) {
	t.Run("File exist", func(t *testing.T) {
		path, ok := getGlobalConfigFile(existStub)
		assert.Equal(t, "/etc/perun/main.yaml", path, "Should contain /etc")
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
		assert.Equal(t, dir+"/.perun", path, "Should contain current working directory")
		assert.True(t, ok, "Should exist")
	})

	t.Run("File does not exist", func(t *testing.T) {
		_, ok := getConfigFileFromCurrentWorkingDirectory(notExistStub)
		assert.False(t, ok, "Should not exist")
	})
}
