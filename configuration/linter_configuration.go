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

package configuration

import (
	"github.com/Appliscale/perun/logger"
	"io/ioutil"
	"os"
)

// GetLinterConfigurationFile reads configuration from file.
func GetLinterConfigurationFile(linterFile *string, logger *logger.Logger) (rawLintConfiguration string) {
	if *linterFile != "" {
		bytesConfiguration, err := ioutil.ReadFile(*linterFile)
		if err != nil {
			logger.Error("Error reading linter configuration file from " + *linterFile)
			logger.Error(err.Error())
		}
		rawLintConfiguration = string(bytesConfiguration)
	} else {
		funcName(logger, "~/.config/perun/")
		conf, ok := getUserConfigFile(os.Stat, "style.yaml")
		if !ok {
			logger.Error("Error getting linter configuration file from ~/.config/perun/")
		}
		bytesConfiguration, err := ioutil.ReadFile(conf)
		if err != nil {
			logger.Error(err.Error())
		}
		rawLintConfiguration = string(bytesConfiguration)
	}
	return
}

func funcName(logger *logger.Logger, linterFile string) {
	logger.Info("Linter Configuration file from the following location will be used: " + linterFile)
}
