// +build !windows

package cfconfiguration

import "os"
import "os/user"

type myStat func(string) (os.FileInfo, error)

func getUserConfigFile(existenceChecker myStat) (val string, ok bool) {
	const relativeUserConfigPath = "/.config/cftool/main.yaml"

	var err error
	var usr *user.User

	usr, err = user.Current()
	if err != nil {
		return "", false
	}

	userConfigPath := usr.HomeDir + relativeUserConfigPath

	_, err = existenceChecker(userConfigPath)
	if err != nil {
		return "", false
	}

	return userConfigPath, true
}

func getGlobalConfigFile(existenceChecker myStat) (val string, ok bool) {
	const globalConfigPath = "/etc/cftool/main.yaml"

	_, err := existenceChecker(globalConfigPath)
	if err != nil {
		return "", false
	}

	return globalConfigPath, true
}

func getConfigFileFromCurrentWorkingDirectory(existenceChecker myStat) (val string, ok bool) {
	var err error
	var dir string

	dir, err = os.Getwd()
	if err != nil {
		return "", false
	}

	configPath := dir + "/.cftool"

	_, err = existenceChecker(configPath)
	if err != nil {
		return "", false
	}

	return configPath, true
}