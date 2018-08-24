package checkingrequiredfiles

import (
	"bufio"
	"fmt"
	"github.com/Appliscale/perun/cliparser"
	"github.com/Appliscale/perun/configuration"
	"github.com/Appliscale/perun/configurator"
	"github.com/Appliscale/perun/context"
	"github.com/Appliscale/perun/logger"
	"github.com/Appliscale/perun/myuser"
	"github.com/go-ini/ini"
	"io"
	"net/http"
	"os"
	"strings"
)

//CheckingRequiredFiles looks for required and default files and if doesn't find will create these.
func CheckingRequiredFiles(ctx *context.Context) {
	myLogger := logger.CreateDefaultLogger()

	mainYAMLexists, mainError := isMainYAMLPresent(myLogger)
	if mainError != nil {
		myLogger.Error(mainError.Error())
	}

	configAWSExists, configError := isAWSConfigPresent(myLogger)
	if configError != nil {
		myLogger.Error(configError.Error())
	}

	credentialsExists, credentialsError := isCredentialsPresent(myLogger)
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
			profilesInConfig := getProfilesFromFile(homePath+"/.aws/config", myLogger)
			isDefaultProfile := findProfile(profilesInConfig, profile)

			if isDefaultProfile {
				var answer string
				myLogger.GetInput("Default profile exists, do you want to use it *Y* or create your own *N*?", &answer)

				if strings.ToUpper(answer) == "Y" {
					region = findRegionForProfile(profile, homePath+"/.aws/config")
					con := configurator.CreateMainYaml(ctx, profile, region)
					configuration.SaveToFile(con, homePath+"/.config/perun/main.yaml", &myLogger)
					*ctx, _ = context.GetContext(cliparser.ParseCliArguments, configuration.GetConfiguration, configuration.ReadInconsistencyConfiguration)

				} else if strings.ToUpper(answer) == "N" {
					configurator.CreateRequiredFilesInConfigureMode(ctx)
					configurator.CreateAWSCredentialsFile(ctx, profile)
				}

			} else { // isDefaultProfile == false
				myLogger.Always("Available profiles from config:")
				for _, prof := range profilesInConfig {
					myLogger.Always(prof)
				}
				myLogger.GetInput("Which profile should perun use as a default?", &profile)
				isUserProfile := findProfile(profilesInConfig, profile)
				for !isUserProfile {
					myLogger.GetInput("I cannnot find this profile, try again", &profile)
					isUserProfile = findProfile(profilesInConfig, profile)
				}
				region = findRegionForProfile(profile, homePath+"/.aws/config")
				con := configurator.CreateMainYaml(ctx, profile, region)
				configuration.SaveToFile(con, homePath+"/.config/perun/main.yaml", &myLogger)
				*ctx, _ = context.GetContext(cliparser.ParseCliArguments, configuration.GetConfiguration, configuration.ReadInconsistencyConfiguration)

			}
			if !credentialsExists {
				configurator.CreateAWSCredentialsFile(ctx, profile)
			}
		} else { //configAWSExists == false
			var answer string
			myLogger.GetInput("Config doesn't exist, create default *Y* or new *N*?", &answer)

			if strings.ToUpper(answer) == "N" {
				profile, region = configurator.GetRegionAndProfile(ctx)
				configurator.CreateAWSConfigFile(ctx, profile, region)
				con := configurator.CreateMainYaml(ctx, profile, region)
				configuration.SaveToFile(con, homePath+"/.config/perun/main.yaml", &myLogger)
				*ctx, _ = context.GetContext(cliparser.ParseCliArguments, configuration.GetConfiguration, configuration.ReadInconsistencyConfiguration)
				configurator.CreateAWSCredentialsFile(ctx, profile)
				profilesInCredentials := getProfilesFromFile(homePath+"/.aws/credentials", myLogger)
				profilesInConfig := getProfilesFromFile(homePath+"/.aws/config", myLogger)
				profiles := findNewProfileInCredentials(profilesInCredentials, profilesInConfig)

				if len(profiles) > 0 {
					for _, prof := range profiles {
						myLogger.Always("I found profile " + prof + " in credentials, but not in config. \nCreating new profile in config.")
						myLogger.GetInput("Region", &region)
						configurator.CreateAWSConfigFile(ctx, prof, region)
					}
				}
			} else if strings.ToUpper(answer) == "Y" {
				configurator.CreateAWSConfigFile(ctx, profile, region)
				con := configurator.CreateMainYaml(ctx, profile, region)
				configuration.SaveToFile(con, homePath+"/.config/perun/main.yaml", &myLogger)
				*ctx, _ = context.GetContext(cliparser.ParseCliArguments, configuration.GetConfiguration, configuration.ReadInconsistencyConfiguration)
				configurator.CreateAWSCredentialsFile(ctx, profile)

			}
		}

		if credentialsExists {
			isProfileInPresent := isProfileInCredentials(profile, homePath+"/.aws/credentials")

			if !isProfileInPresent {
				configurator.CreateAWSCredentialsFile(ctx, profile)
			}
		}

	} else { //mainYAMLexists == true

		if configAWSExists {
			if !credentialsExists {
				myLogger.Always("Profile from main.yaml: " + ctx.Config.DefaultProfile)
				configurator.CreateAWSCredentialsFile(ctx, ctx.Config.DefaultProfile)
			} else {
				isProfileInPresent := isProfileInCredentials(ctx.Config.DefaultProfile, homePath+"/.aws/credentials")
				if !isProfileInPresent {
					myLogger.Always("Profile from main.yaml: " + ctx.Config.DefaultProfile)
					configurator.CreateAWSCredentialsFile(ctx, ctx.Config.DefaultProfile)
				}
			}

		} else { //configAWSExists ==false
			var answer string
			myLogger.GetInput("Config doesn't exist, create default *Y* or new *N*?", &answer)
			if strings.ToUpper(answer) == "Y" {
				configurator.CreateAWSConfigFile(ctx, ctx.Config.DefaultProfile, ctx.Config.DefaultRegion)
			} else if strings.ToUpper(answer) == "N" {
				profile, region = configurator.GetRegionAndProfile(ctx)
				con := configurator.CreateMainYaml(ctx, profile, region)
				configuration.SaveToFile(con, homePath+"/.config/perun/main.yaml", &myLogger)
				configurator.CreateAWSCredentialsFile(ctx, profile)
				configurator.CreateAWSConfigFile(ctx, profile, region)
				*ctx, _ = context.GetContext(cliparser.ParseCliArguments, configuration.GetConfiguration, configuration.ReadInconsistencyConfiguration)

			}
			profilesInCredentials := getProfilesFromFile(homePath+"/.aws/credentials", myLogger)
			profilesInConfig := getProfilesFromFile(homePath+"/.aws/config", myLogger)
			profiles := findNewProfileInCredentials(profilesInCredentials, profilesInConfig)
			if len(profiles) > 0 {
				for _, prof := range profiles {
					myLogger.Always("I found profile " + prof + " in credentials, but not in config. \nCreating new profile in config.")
					myLogger.GetInput("Region", &region)
					configurator.CreateAWSConfigFile(ctx, prof, region)
				}
			}
		}
	}

	downloadError := downloadDefaultFile()
	if downloadError != nil {
		myLogger.Error(downloadError.Error())
	}
}

// Looking for main.yaml.
func isMainYAMLPresent(myLogger logger.Logger) (bool, error) {
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

// Looking for .aws.config.
func isAWSConfigPresent(myLogger logger.Logger) (bool, error) {
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
func isCredentialsPresent(myLogger logger.Logger) (bool, error) {
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
func getProfilesFromFile(path string, mylogger logger.Logger) []string {
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
			profiles = append(profiles, profile)
		}
	}
	return profiles
}

// Looking for user's profile in credentials or config.
func isProfileInCredentials(profile string, path string) bool {
	credentials, credentialsError := os.Open(path)
	if credentialsError != nil {
		fmt.Println(credentialsError.Error())
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

// Looking for profile in profiles.
func findProfile(profiles []string, myProfile string) bool {
	for _, profile := range profiles {
		if myProfile == profile {
			return true
		}
	}
	return false
}

// Looking for region for profile.
func findRegionForProfile(profile string, path string) string {
	configuration, loadError := ini.Load(path)
	if loadError != nil {
		fmt.Println(loadError.Error())
	}
	section, sectionError := configuration.GetSection(profile)
	if sectionError != nil {
		fmt.Println(sectionError.Error())
		return ""
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
		isProfileHere := findProfile(config, cred)
		if !isProfileHere {
			profiles = append(profiles, cred)
			return profiles
		}
	}
	return []string{}
}

// Downloadind other files.
func downloadDefaultFile() error {
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

		_, openError := os.Open(homePath + file)
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
