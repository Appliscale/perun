// Copyright 2018 Appliscale
//
// Maintainers and contributors are listed in README file inside repository.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package checkingrequiredfiles checks if .aws/config, .aws/credentials main.yaml and other configuration files exist.
package checkingrequiredfiles

import (
	"bufio"
	"github.com/Appliscale/perun/configurator"
	"github.com/Appliscale/perun/context"
	"github.com/Appliscale/perun/helpers"
	"github.com/Appliscale/perun/logger"
	"github.com/Appliscale/perun/myuser"
	"github.com/go-ini/ini"
	"io"
	"net/http"
	"os"
	"strings"
)

var EnvironmentVariables = map[string]string{
	"profile": os.Getenv("AWS_PROFILE"),
	"id":      os.Getenv("AWS_ACCESS_KEY_ID"),
	"key":     os.Getenv("AWS_SECRET_ACCESS_KEY"),
	"region":  os.Getenv("AWS_DEFAULT_REGION"),
	"output":  os.Getenv("AWS_DEFAULT_OUTPUT"),
	"token":   os.Getenv("AWS_SESSION_TOKEN"),
}

//CheckingRequiredFiles looks for required and default files and if doesn't find will create these.
func CheckingRequiredFiles(ctx *context.Context) (offline bool) {
	myLogger := logger.CreateDefaultLogger()
	offline = false
	mainYAMLexists, mainError := isMainYAMLPresent(&myLogger)
	if mainError != nil {
		myLogger.Error(mainError.Error())
	}

	configAWSExists, configError := isAWSConfigPresent(&myLogger)
	if configError != nil {
		myLogger.Error(configError.Error())
	}

	credentialsExists, credentialsError := isCredentialsPresent(&myLogger)
	if credentialsError != nil {
		myLogger.Error(credentialsError.Error())
	}

	homePath, pathError := myuser.GetUserHomeDir()
	if pathError != nil {
		myLogger.Error(pathError.Error())
	}

	profile := "default"
	region := "us-east-1"

	if !mainYAMLexists {
		if configAWSExists {
			profile, *ctx = configIsPresent(profile, homePath, ctx, &myLogger)
			if !credentialsExists {
				//dane z zmiennych srodowiskowych
				// var answer string
				// ctx.Logger.GetInput("Creating aws/credentials based on environment variables? Y/N", &answer)
				// if strings.ToUpper(answer) == "Y" {
				// 	createCredentialsBasedOnEnvironmentVariables(EnvironmentVariables, ctx.Logger)
				// 	profile = EnvironmentVariables["profile"]
				// 	region = EnvironmentVariables["region"]

				// } else if *ctx.CliArguments.Mode == cliparser.ValidateMode {
				// 	var answer string
				// 	myLogger.GetInput("You haven't got credentials file, run only offline validation? Y/N", &answer)
				// 	if strings.ToUpper(answer) == "Y" {
				// 		return true //offline
				// 	}
				createCredentials(profile, homePath, ctx, &myLogger)
				//}
			}
		} else { //configAWSExists == false
			var answer string
			myLogger.GetInput("Config doesn't exist, create default *Y* or new *N*?", &answer)
			if strings.ToUpper(answer) == "N" {
				profile, region, *ctx = newConfigFile(profile, region, homePath, ctx, &myLogger)
				addProfileToCredentials(profile, homePath, ctx, &myLogger)
				addNewProfileFromCredentialsToConfig(profile, homePath, ctx, &myLogger)

			} else if strings.ToUpper(answer) == "Y" {
				configurator.CreateAWSConfigFile(&myLogger, profile, region)
				*ctx = createNewMainYaml(profile, homePath, ctx, &myLogger)
				configurator.CreateAWSCredentialsFile(ctx, profile)
			}

			if credentialsExists {
				createCredentials(profile, homePath, ctx, &myLogger)
			}
		}
	} else { //mainYAMLexists == true
		if configAWSExists {
			if !credentialsExists {
				offline, ctx.Config.DefaultProfile, ctx.Config.DefaultRegion = checkingCredentials(ctx, profile, region)
				if offline == true {
					return offline
				}
				myLogger.Always("Profile from main.yaml: " + ctx.Config.DefaultProfile)
				addProfileToCredentials(ctx.Config.DefaultProfile, homePath, ctx, ctx.Logger)
				//configurator.CreateAWSCredentialsFile(ctx, ctx.Config.DefaultProfile)
			} else {
				isProfileInPresent := isProfileInCredentials(ctx.Config.DefaultProfile, homePath+"/.aws/credentials", &myLogger)
				if !isProfileInPresent {
					myLogger.Always("Profile from main.yaml: " + ctx.Config.DefaultProfile)
					configurator.CreateAWSCredentialsFile(ctx, ctx.Config.DefaultProfile)
				}
			}
		} else { //configAWSExists ==false
			var answer string
			myLogger.GetInput("Config doesn't exist, create default - "+ctx.Config.DefaultProfile+" *Y* or new *N*?", &answer)
			if strings.ToUpper(answer) == "Y" {
				configurator.CreateAWSConfigFile(ctx.Logger, ctx.Config.DefaultProfile, ctx.Config.DefaultRegion)
			} else if strings.ToUpper(answer) == "N" {
				profile, region, *ctx = newConfigFile(profile, region, homePath, ctx, &myLogger)
				addProfileToCredentials(profile, homePath, ctx, &myLogger)
			}
			addNewProfileFromCredentialsToConfig(ctx.Config.DefaultProfile, homePath, ctx, &myLogger)

			if credentialsExists {
				createCredentials(ctx.Config.DefaultProfile, homePath, ctx, &myLogger)

			}
		}
	}
	downloadError := downloadDefaultFiles()
	if downloadError != nil {
		myLogger.Error(downloadError.Error())
	}
	return false
}

// Looking for main.yaml.
func isMainYAMLPresent(myLogger *logger.Logger) (bool, error) {
	homePath, pathError := myuser.GetUserHomeDir()
	if pathError != nil {
		myLogger.Error(pathError.Error())
		return false, pathError
	}
	_, mainError := os.Open(homePath + "/.config/perun/main.yaml")
	if mainError != nil {
		_, mainError = os.Open(homePath + "/etc/perun/main.yaml")
		if mainError != nil {
			return false, pathError
		}
		return true, pathError
	}
	return true, pathError
}

// Looking for .aws/config.
func isAWSConfigPresent(myLogger *logger.Logger) (bool, error) {
	homePath, pathError := myuser.GetUserHomeDir()
	if pathError != nil {
		myLogger.Error(pathError.Error())
		return false, pathError
	}
	_, credentialsError := os.Open(homePath + "/.aws/config")
	if credentialsError != nil {
		return false, credentialsError
	}
	return true, nil

}

// Looking for .aws/credentials.
func isCredentialsPresent(myLogger *logger.Logger) (bool, error) {
	homePath, pathError := myuser.GetUserHomeDir()
	if pathError != nil {
		myLogger.Error(pathError.Error())
		return false, pathError
	}
	_, credentialsError := os.Open(homePath + "/.aws/credentials")
	if credentialsError != nil {
		return false, credentialsError
	}
	return true, nil
}

// Looking for [profiles] in credentials or config and return all.
func getProfilesFromFile(path string, mylogger logger.LoggerInt) []string {
	credentials, credentialsError := os.Open(path)
	if credentialsError != nil {
		mylogger.Error(credentialsError.Error())
		return []string{}
	}
	defer credentials.Close()
	profiles := make([]string, 0)
	scanner := bufio.NewScanner(credentials)
	for scanner.Scan() {
		if strings.Contains(scanner.Text(), "[") {
			profile := strings.TrimPrefix(scanner.Text(), "[")
			profile = strings.TrimSuffix(profile, "]")
			if strings.Contains(profile, "profile ") {
				profile = strings.TrimPrefix(profile, "profile ")
			}
			if strings.Contains(profile, "-long-term") {
				profile = strings.TrimSuffix(profile, "-long-term")
			}
			profiles = append(profiles, profile)
		}
	}
	return profiles
}

// Looking for user's profile in credentials or config.
func isProfileInCredentials(profile string, path string, mylogger logger.LoggerInt) bool {
	credentials, credentialsError := os.Open(path)
	if credentialsError != nil {
		mylogger.Error(credentialsError.Error())
	}
	defer credentials.Close()
	scanner := bufio.NewScanner(credentials)
	for scanner.Scan() {
		if strings.Contains(scanner.Text(), "["+profile+"]") || strings.Contains(scanner.Text(), "["+profile+"-long-term]") {
			return true
		}
	}
	return false

}

// Looking for region for profile.
func findRegionForProfile(profile string, path string, mylogger logger.LoggerInt) string {
	configuration, loadError := ini.Load(path)
	if loadError != nil {
		mylogger.Error(loadError.Error())
	}
	section, sectionError := configuration.GetSection(profile)
	if sectionError != nil {
		section, sectionError = configuration.GetSection("profile " + profile)
		if sectionError != nil {
			mylogger.Error(sectionError.Error())
			return ""
		}
	}
	region := section.Key("region").Value()
	return region

}

// Getting profiles from credentials and config, if credentials has new profile, add to config.
func findNewProfileInCredentials(credentials []string, config []string) []string {
	profiles := make([]string, 0)
	for i, cred := range credentials {
		if strings.Contains(cred, "-long-term") {
			cred = strings.TrimSuffix(cred, "-long-term")
			credentials[i] = cred
		}
	}
	for _, cred := range credentials {
		isProfileHere := helpers.SliceContains(config, cred)
		if !isProfileHere {
			profiles = append(profiles, cred)
			return profiles
		}
	}
	return []string{}
}

// Downloading other files.
func downloadDefaultFiles() error {
	urls := make(map[string]string)
	urls["blocked.json"] = "https://s3.amazonaws.com/perun-default-file/blocked.json"
	urls["unblocked.json"] = "https://s3.amazonaws.com/perun-default-file/unblocked.json"
	urls["style.yaml"] = "https://s3.amazonaws.com/perun-default-file/style.yaml"
	urls["specification_inconsistency.yaml"] = "https://s3.amazonaws.com/perun-default-file/specification_inconsistency.yaml"

	for file, url := range urls {
		homePath, _ := myuser.GetUserHomeDir()
		homePath += "/.config/perun/"

		if strings.Contains(file, "blocked") {
			homePath += "stack-policies/"
		}

		_, err := os.Stat(homePath)
		if os.IsNotExist(err) {
			os.Mkdir(homePath, 0755)
		}

		_, openError := os.Open(homePath + file) //checking if file exists
		if openError != nil {
			out, creatingFileError := os.Create(homePath + file)

			if creatingFileError != nil {
				return creatingFileError
			}
			defer out.Close()

			resp, httpGetError := http.Get(url)
			if httpGetError != nil {
				return httpGetError
			}
			defer resp.Body.Close()

			_, copyError := io.Copy(out, resp.Body)
			if copyError != nil {
				return copyError
			}
		}
	}
	return nil
}
