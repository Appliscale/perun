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
	"io/ioutil"
	"os"

	"github.com/Appliscale/perun/logger"
	"github.com/ghodss/yaml"
)

// InconsistencyConfiguration describes inconsistencies between specification and documentation.
type InconsistencyConfiguration struct {
	SpecificationInconsistency map[string]Property
}

// Property of inconsistency.
type Property map[string][]string

// ReadInconsistencyConfiguration gets configuration from file, if could not read return warning.
func ReadInconsistencyConfiguration(logger logger.LoggerInt) (config InconsistencyConfiguration) {
	if path, ok := getUserConfigFile(os.Stat, "specification_inconsistency.yaml"); ok {
		rawConfig, err := ioutil.ReadFile(path)
		if err != nil {
			logger.Warning("Could not read specification incosistencies configuration file")
			return
		}

		err = yaml.Unmarshal(rawConfig, &config)
		if err != nil {
			logger.Warning("Specification inconsistencies configuration file format is invalid")
			return
		}

		return
	}

	logger.Warning("Specification inconsistencies configuration file not found")
	return
}
