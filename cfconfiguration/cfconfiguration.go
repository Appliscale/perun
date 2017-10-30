package cfconfiguration

import (
	"io/ioutil"
	"errors"
	"github.com/ghodss/yaml"
	"os"
	"os/user"
	"fmt"
	"github.com/Appliscale/cftool/cfcliparser"
)

type Configuration struct {
	Profile string
	Region string
	SpecificationURL map[string]string
}

const GlobalConfigFile string = "/etc/.Appliscale/cftool/config.yaml"
const UserConfigFile string = "/.config/cftool/config.yaml"

func (config Configuration) GetSpecificationFileURLForCurrentRegion() (string, error) {
	if url, ok := config.SpecificationURL[config.Region]; ok {
		return url + "/latest/gzip/CloudFormationResourceSpecification.json", nil
	}
	return "", errors.New("There is no specification file for region " + config.Region)
}

func GetConfiguration(cliArguments cfcliparser.CliArguments) (config Configuration, err error) {
	configPath, err := getConfigurationPath(cliArguments)
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

	return
}

func getConfigurationPath(cliArguments cfcliparser.CliArguments) (configPath string, err error) {
	user, err := user.Current()
	usersHomeConfiguration := user.HomeDir + UserConfigFile
	if _, err := os.Stat(*cliArguments.ConfigurationPath); err == nil {
		notifyUserAboutConfigurationFile(*cliArguments.ConfigurationPath)
		return *cliArguments.ConfigurationPath, nil
	} else if _, err := os.Stat(usersHomeConfiguration); err == nil {
		notifyUserAboutConfigurationFile(usersHomeConfiguration)
		return usersHomeConfiguration, nil
	} else if _, err := os.Stat(GlobalConfigFile); err == nil {
		notifyUserAboutConfigurationFile(GlobalConfigFile)
		return GlobalConfigFile, nil
	} else {
		return "", errors.New("There is no configuration file!")
	}
}
func notifyUserAboutConfigurationFile(configurationFilePath string) (int, error) {
	return fmt.Println("Configuration file from the following location will be use: " + configurationFilePath)
}