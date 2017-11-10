package cfconfiguration

import "os"

func getConfigFileFromCurrentWorkingDirectory_(existenceChecker myStat, relativePath string) (val string, ok bool) {
	var err error
	var dir string

	dir, err = os.Getwd()
	if err != nil {
		return "", false
	}

	configPath := dir + relativePath

	_, err = existenceChecker(configPath)
	if err != nil {
		return "", false
	}

	return configPath, true
}