package cfconfiguration

import (
	"io/ioutil"
	"errors"
	"github.com/ghodss/yaml"
	"os"
	"os/user"
	"github.com/Appliscale/cftool/cfcliparser"
	"github.com/Appliscale/cftool/cflogger"
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

	return
}

func getConfigurationPath(cliArguments cfcliparser.CliArguments, logger *cflogger.Logger) (configPath string, err error) {
	user, err := user.Current()
	usersHomeConfiguration := user.HomeDir + UserConfigFile
	if _, err := os.Stat(*cliArguments.ConfigurationPath); err == nil {
		notifyUserAboutConfigurationFile(*cliArguments.ConfigurationPath, logger)
		return *cliArguments.ConfigurationPath, nil
	} else if _, err := os.Stat(usersHomeConfiguration); err == nil {
		notifyUserAboutConfigurationFile(usersHomeConfiguration, logger)
		return usersHomeConfiguration, nil
	} else if _, err := os.Stat(GlobalConfigFile); err == nil {
		notifyUserAboutConfigurationFile(GlobalConfigFile, logger)
		return GlobalConfigFile, nil
	} else {
		return "", errors.New("There is no configuration file!")
	}
}
func notifyUserAboutConfigurationFile(configurationFilePath string, logger *cflogger.Logger) {
	logger.Info("Configuration file from the following location will be use: " + configurationFilePath)
}