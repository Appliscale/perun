package configuration

import (
	"io/ioutil"
	"os"

	"github.com/Appliscale/perun/logger"
	"github.com/ghodss/yaml"
)

type InconsistencyConfiguration struct {
	// Inconsistencies between specification and documentation.
	SpecificationInconsistency map[string]Property
}

type Property map[string][]string

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
