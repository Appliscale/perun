// +build windows

package cfconfiguration

import "os"

type myStat func(string) (os.FileInfo, error)

func getUserConfigFile(existenceChecker myStat) (val string, ok bool) {
	const envVar = "LOCALAPPDATA"
	const relativeUserConfigPath = "\\cftool\\main.yaml"

	return checkConfigExistence(envVar, relativeUserConfigPath, existenceChecker)
}

func getGlobalConfigFile(existenceChecker myStat) (val string, ok bool) {
	const envVar = "ALLUSERSPROFILE"
	const relativeGlobalConfigPath = "\\cftool\\main.yaml"

	return checkConfigExistence(envVar, relativeGlobalConfigPath, existenceChecker)
}

func checkConfigExistence(envVar string, relativeConfigPath string, existenceChecker myStat) (val string, ok bool) {
	absoluteConfigPath, ok := buildAbsolutePath(envVar, relativeConfigPath)
	if !ok {
		return "", false
	}

	_, err := existenceChecker(absoluteConfigPath)
	if err != nil {
		return "", false
	}

	return absoluteConfigPath, true
}

func buildAbsolutePath(envVar string, relativeConfigPath string) (val string, ok bool) {
	envVal, ok := os.LookupEnv(envVar)
	if !ok {
		return "", false
	}

	absoluteConfigPath := envVal + relativeConfigPath

	return absoluteConfigPath, true
}

func getConfigFileFromCurrentWorkingDirectory(existenceChecker myStat) (val string, ok bool) {
	return getConfigFileFromCurrentWorkingDirectory_(existenceChecker, "\\.cftool")
}