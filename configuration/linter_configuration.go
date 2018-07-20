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
			logger.Error(err.Error())
		}
		rawLintConfiguration = string(bytesConfiguration)
	} else {
		conf, ok := getUserConfigFile(os.Stat, "style.yaml")
		if !ok {
			logger.Error("Error getting linter configuration")
		}
		bytesConfiguration, err := ioutil.ReadFile(conf)
		if err != nil {
			logger.Error(err.Error())
		}
		rawLintConfiguration = string(bytesConfiguration)
	}
	return
}
