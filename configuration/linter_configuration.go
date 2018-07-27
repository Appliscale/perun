package configuration

import (
	"github.com/Appliscale/perun/logger"
	"io/ioutil"
	"os"
)

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
