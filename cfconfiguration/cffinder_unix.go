// +build !windows

package cfconfiguration

import "go/build"
import "os"
import "os/user"

var checkFileExistence = func (name string) (os.FileInfo, error) { return os.Stat(name) }

func getUserConfigFile() (val string, ok bool) {
	const relativeUserConfigPath = "/.config/cftool/config.yaml"

	usr, err := user.Current()
	if err != nil {
		return "", false
	}

	userConfigPath := usr.HomeDir + relativeUserConfigPath

	if _, err := checkFileExistence(userConfigPath); err != nil {
		return "", false
	}

	return userConfigPath, true
}

func getGlobalConfigFile() (val string, ok bool) {
	const globalConfigPath = "/etc/.Appliscale/cftool/config.yaml"

	if _, err := checkFileExistence(globalConfigPath); err != nil {
		return "", false
	}

	return globalConfigPath, true
}

func getConfigFileFromProjectRoot() (val string, ok bool) {
	const relativeProjectRoot = "github.com/Appliscale/cftool"

	goPath := build.Default.GOPATH

	configPath := goPath + "/src/" + relativeProjectRoot + "/config.yaml"

	if _, err := checkFileExistence(configPath); err != nil {
		return "", false
	}

	return configPath, true
}