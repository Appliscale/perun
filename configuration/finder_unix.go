// Copyright 2017 Appliscale
//
// Maintainers and Contributors:
//
//   - Piotr Figwer (piotr.figwer@appliscale.io)
//   - Wojciech Gawroński (wojciech.gawronski@appliscale.io)
//   - Kacper Patro (kacper.patro@appliscale.io)
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

import "os"
import "os/user"

type myStat func(string) (os.FileInfo, error)

func getUserConfigFile(existenceChecker myStat) (val string, ok bool) {
	const relativeUserConfigPath = "/.config/perun/main.yaml"

	var err error
	var usr *user.User

	usr, err = user.Current()
	if err != nil {
		return "", false
	}

	userConfigPath := usr.HomeDir + relativeUserConfigPath

	_, err = existenceChecker(userConfigPath)
	if err != nil {
		return "", false
	}

	return userConfigPath, true
}

func getGlobalConfigFile(existenceChecker myStat) (val string, ok bool) {
	const globalConfigPath = "/etc/perun/main.yaml"

	_, err := existenceChecker(globalConfigPath)
	if err != nil {
		return "", false
	}

	return globalConfigPath, true
}

func getConfigFileFromCurrentWorkingDirectory(existenceChecker myStat) (val string, ok bool) {
	return getConfigFileFromCurrentWorkingDirectory_(existenceChecker, "/.perun")
}
