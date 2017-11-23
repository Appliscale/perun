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
	DefaultProfile string
	// AWS region (e.g. us-east-1).
	DefaultRegion string
	// Map of resource specification CloudFront URL per region.
	SpecificationURL map[string]string
	// Decision regarding if we use MFA token or not.
	DefaultDecisionForMFA bool
	// Duration for MFA token.
	DefaultDurationForMFA int64
	// Logger verbosity.
	DefaultVerbosity string
}

// Return URL to specification file. If there is no specification file for selected region, return error.
func (config Configuration) GetSpecificationFileURLForCurrentRegion() (string, error) {
	if url, ok := config.SpecificationURL[config.DefaultRegion]; ok {
		return url + "/latest/gzip/CloudFormationResourceSpecification.json", nil
	}

	return "", errors.New("There is no specification file for region " + config.DefaultRegion)
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

	postProcessing(&config, cliArguments, logger)

	return
}

func postProcessing(config *Configuration, cliArguments cliparser.CliArguments, logger *logger.Logger) {
	if config.DefaultProfile == "" {
		config.DefaultProfile = "default"
	}
	if config.DefaultRegion == "" {
		config.DefaultRegion = "us-east-1"
	}
	if config.DefaultVerbosity == "" {
		config.DefaultVerbosity = "INFO"
	}
	if *cliArguments.Verbosity != "" {
		config.DefaultVerbosity = *cliArguments.Verbosity
	}
	if *cliArguments.Region != "" {
		config.DefaultRegion = *cliArguments.Region
	}
	if *cliArguments.Profile != "" {
		config.DefaultProfile = *cliArguments.Profile
	}
	if *cliArguments.MFA != config.DefaultDecisionForMFA {
		config.DefaultDecisionForMFA = *cliArguments.MFA
	}
	if *cliArguments.DurationForMFA > 0 {
		config.DefaultDurationForMFA = *cliArguments.DurationForMFA
	}
}

func getConfigurationPath(cliArguments cliparser.CliArguments, logger *logger.Logger) (configPath string, err error) {
	if *cliArguments.Sandbox {
		return "", errors.New("No configuration file should be used.")
	}

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
