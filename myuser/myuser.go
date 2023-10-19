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

// Package myuser provides function which ease work with user information e.g path to home directory.
package myuser

import (
	"os/user"
)

// GetUserHomeDir gets path to user's home directory. It's used when perun checks if configuration files exists.
func GetUserHomeDir() (string, error) {
	user, userError := user.Current()
	if userError != nil {
		return "", userError
	}
	path := user.HomeDir

	return path, nil
}
