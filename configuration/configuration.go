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

// Package configuration provides reader for perun configuration.
package configuration

import (
	"io/ioutil"
	"errors"
	"github.com/ghodss/yaml"
	"os"
	"github.com/Appliscale/perun/cliparser"
	"github.com/Appliscale/perun/logger"
)

// Perun configuration.
type Configuration struct {
	// AWS credentials profile.
	Profile string
	// AWS region (e.g. us-east-1).
	Region string
	// Map of resource specification CloudFront URL per region.
	SpecificationURL map[string]string
}

// Return URL to specification file. If there is no specification file for selected region, return error.
func (config Configuration) GetSpecificationFileURLForCurrentRegion() (string, error) {
	if url, ok := config.SpecificationURL[config.Region]; ok {
		return url + "/latest/gzip/CloudFormationResourceSpecification.json", nil
	}
	return "", errors.New("There is no specification file for region " + config.Region)
}

// Return perun configuration read from file.
func GetConfiguration(cliArguments cliparser.CliArguments, logger *logger.Logger) (config Configuration, err error) {
	configPath, err := getConfigurationPath(cliArguments, logger)
	if err != nil {
		return
	}
	rawConfiguration, err := ioutil.ReadFile(configPath)
	if err != nil {
		return
	}
	err = yaml.Unmarshal(rawConfiguration, &config)
	if err != nil {
		return
	}

	if config.Profile == "" {
		config.Profile = "default"
	}

	if config.Region == "" {
		config.Region = "us-east-1"
	}

	return
}

func getConfigurationPath(cliArguments cliparser.CliArguments, logger *logger.Logger) (configPath string, err error) {
	if _, err := os.Stat(*cliArguments.ConfigurationPath); err == nil {
		notifyUserAboutConfigurationFile(*cliArguments.ConfigurationPath, logger)
		return *cliArguments.ConfigurationPath, nil
	} else if path, ok := getConfigFileFromCurrentWorkingDirectory(os.Stat); ok {
		notifyUserAboutConfigurationFile(path, logger)
		return path, nil
	} else if path, ok := getUserConfigFile(os.Stat); ok {
		notifyUserAboutConfigurationFile(path, logger)
		return path, nil
	} else if path, ok := getGlobalConfigFile(os.Stat); ok {
		notifyUserAboutConfigurationFile(path, logger)
		return path, nil
	} else {
		return "", errors.New("Configuration file could not be read!")
	}
}

func notifyUserAboutConfigurationFile(configurationFilePath string, logger *logger.Logger) {
	logger.Info("Configuration file from the following location will be used: " + configurationFilePath)
}
