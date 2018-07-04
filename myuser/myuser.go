package myuser

import (
	"os/user"
)

// GetUserHomeDir gets path to home directory.
func GetUserHomeDir() (string, error) {
	user, userError := user.Current()
	if userError != nil {
		return "", userError
	}
	path := user.HomeDir

	return path, nil
}
