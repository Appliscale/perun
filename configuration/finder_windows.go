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

// +build windows

package configuration

import "os"

type myStat func(string) (os.FileInfo, error)

func getUserConfigFile(existenceChecker myStat) (val string, ok bool) {
	const envVar = "LOCALAPPDATA"
	const relativeUserConfigPath = "\\perun\\main.yaml"

	return checkConfigExistence(envVar, relativeUserConfigPath, existenceChecker)
}

func getGlobalConfigFile(existenceChecker myStat) (val string, ok bool) {
	const envVar = "ALLUSERSPROFILE"
	const relativeGlobalConfigPath = "\\perun\\main.yaml"

	return checkConfigExistence(envVar, relativeGlobalConfigPath, existenceChecker)
}

func checkConfigExistence(envVar string, relativeConfigPath string, existenceChecker myStat) (val string, ok bool) {
	absoluteConfigPath, ok := buildAbsolutePath(envVar, relativeConfigPath)
	if !ok {
		return "", false
	}

	_, err := existenceChecker(absoluteConfigPath)
	if err != nil {
		return "", false
	}

	return absoluteConfigPath, true
}

func buildAbsolutePath(envVar string, relativeConfigPath string) (val string, ok bool) {
	envVal, ok := os.LookupEnv(envVar)
	if !ok {
		return "", false
	}

	absoluteConfigPath := envVal + relativeConfigPath

	return absoluteConfigPath, true
}

func getConfigFileFromCurrentWorkingDirectory(existenceChecker myStat) (val string, ok bool) {
	return getConfigFileFromCurrentWorkingDirectory_(existenceChecker, "\\.perun")
}
