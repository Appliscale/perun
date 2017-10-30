package cfconfiguration

import (
	"io/ioutil"
	"errors"
	"github.com/ghodss/yaml"
	"os"
	"os/user"
	"fmt"
	"github.com/Appliscale/cftool/cfcontext"
)

type configuration struct {
	Profile string
	Region string
	SpecificationURL map[string]string
}

var config configuration
const globalConfigFile string = "/etc/.Appliscale/cftool/config.yaml"
const userConfigFile string = "/.config/cftool/config.yaml"

func GetSpecificationFileURL(context *cfcontext.Context) (string, error) {
	err := getConfiguration(*context.CliArguments.ConfigurationPath)
	if err != nil {
		return "", err
	}
	if url, ok := config.SpecificationURL[config.Region]; ok {
		return url + "/latest/gzip/CloudFormationResourceSpecification.json", nil
	}
	return "", errors.New("There is no specification file for region " + config.Region)
}

func GetRegion(context *cfcontext.Context) (string, error) {
	err := getConfiguration(*context.CliArguments.ConfigurationPath)
	if err != nil {
		return "", err
	}
	return config.Region, nil
}

func getConfiguration(configurationFilePath string) error {
	if len(config.SpecificationURL) == 0 {
		configPath, err := getConfigPath(configurationFilePath)
		if err != nil {
			config = configuration{}
			return err
		}
		rawConfiguration, err := ioutil.ReadFile(configPath)
		err = yaml.Unmarshal(rawConfiguration, &config)
		if err != nil {
			config = configuration{}
			return err
		}
	}

	return nil
}

func getConfigPath(configurationFilePath string) (configPath string, err error) {
	user, err := user.Current()
	usersHomeConfiguration := user.HomeDir + userConfigFile
	if _, err := os.Stat(configurationFilePath); err == nil {
		notifyUserAboutConfigurationFile(configurationFilePath)
		return configurationFilePath, nil
	} else if _, err := os.Stat(usersHomeConfiguration); err == nil {
		notifyUserAboutConfigurationFile(usersHomeConfiguration)
		return usersHomeConfiguration, nil
	} else if _, err := os.Stat(globalConfigFile); err == nil {
		notifyUserAboutConfigurationFile(globalConfigFile)
		return globalConfigFile, nil
	} else {
		return "", errors.New("There is no configuration file!")
	}
}
func notifyUserAboutConfigurationFile(configurationFilePath string) (int, error) {
	return fmt.Println("Configuration file from the following location will be use: " + configurationFilePath)
}