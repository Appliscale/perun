// Package cfconfiguration provides reader for PerunCloud configuration.
package cfconfiguration

import (
	"io/ioutil"
	"errors"
	"github.com/ghodss/yaml"
	"os"
	"github.com/Appliscale/cftool/cfcliparser"
	"github.com/Appliscale/cftool/cflogger"
)

// PerunCloud configuration.
type Configuration struct {
	// AWS credentials profile.
	Profile string
	// AWS region (e.g. us-east-1).
	Region string
	// Map of resource specification cloudfront URL per region.
	SpecificationURL map[string]string
}

// Return URL to specification file. If there is no specification file for selected region, return error.
func (config Configuration) GetSpecificationFileURLForCurrentRegion() (string, error) {
	if url, ok := config.SpecificationURL[config.Region]; ok {
		return url + "/latest/gzip/CloudFormationResourceSpecification.json", nil
	}
	return "", errors.New("There is no specification file for region " + config.Region)
}

// Return PerunCloud configuration readed from file.
func GetConfiguration(cliArguments cfcliparser.CliArguments, logger *cflogger.Logger) (config Configuration, err error) {
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

func getConfigurationPath(cliArguments cfcliparser.CliArguments, logger *cflogger.Logger) (configPath string, err error) {
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

func notifyUserAboutConfigurationFile(configurationFilePath string, logger *cflogger.Logger) {
	logger.Info("Configuration file from the following location will be used: " + configurationFilePath)
}