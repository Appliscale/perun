// +build windows

package cfconfiguration

import "errors"
import "os"

const relativeUserConfigPath string = "\\.config\\cftool\\config.yaml"
const relativeGlobalConfigPath string = "\\.Appliscale\\cftool\\config.yaml"

func GetUserConfigFile() (val string, err error) {
	adp, ok := os.LookupEnv("APPDATA")
	if !ok {
		return "", errors.New("missed env lookup")
	}

	userConfigPath := adp + relativeUserConfigPath

	if _, err := os.Stat(userConfigPath); err != nil {
		return "", err
	}

	return userConfigPath, nil
}

func GetGlobalConfigFile() (val string, err error) {
	pfp, ok := os.LookupEnv("PROGRAMFILES")
	if !ok {
		return "", errors.New("missed env lookup")
	}

	globalConfigPath := pfp + relativeGlobalConfigPath

	if _, err := os.Stat(globalConfigPath); err != nil {
		return "", err
	}

	return globalConfigPath, nil
}