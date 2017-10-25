package cfconfiguration

import (
	"io/ioutil"
	"errors"
	"github.com/ghodss/yaml"
)

type configuration struct {
	SpecificationURL map[string]string
}

var config configuration

func GetSpecificationFileURL(region string) (string, error) {
	err := getConfiguration()
	if err != nil {
		return "", err
	}
	if url, ok := config.SpecificationURL[region]; ok {
		return url + "/latest/gzip/CloudFormationResourceSpecification.json", nil
	}
	return "", errors.New("There is no specification file for region " + region)
}

func getConfiguration() (error) {
	if len(config.SpecificationURL) == 0 {
		rawConfiguration, err := ioutil.ReadFile("/etc/.Appliscale/cftool/config.yaml")
		if err != nil {
			return err
		}
		err = yaml.Unmarshal(rawConfiguration, &config)
		if err != nil {
			config = configuration{}
			return err
		}
	}

	return nil
}