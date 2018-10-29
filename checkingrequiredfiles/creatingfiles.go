package checkingrequiredfiles

import (
	"github.com/Appliscale/perun/cliparser"
	"github.com/Appliscale/perun/configuration"
	"github.com/Appliscale/perun/configurator"
	"github.com/Appliscale/perun/context"
	"github.com/Appliscale/perun/helpers"
	"github.com/Appliscale/perun/logger"
	"github.com/Appliscale/perun/myuser"
	"strings"
)

// Creating main.yaml.
func createNewMainYaml(profile string, homePath string, ctx *context.Context, myLogger logger.LoggerInt) context.Context {
	region := findRegionForProfile(profile, homePath+"/.aws/config", myLogger)
	con := configurator.CreateMainYaml(myLogger, profile, region)
	configuration.SaveToFile(con, homePath+"/.config/perun/main.yaml", myLogger)
	*ctx, _ = context.GetContext(cliparser.ParseCliArguments, configuration.GetConfiguration, configuration.ReadInconsistencyConfiguration)
	return *ctx
}

// If config exists, use profile from .aws/config.
func useProfileFromConfig(profilesInConfig []string, profile string, myLogger logger.LoggerInt) string {
	myLogger.Always("Available profiles from config:")
	for _, prof := range profilesInConfig {
		myLogger.Always(prof)
	}
	myLogger.GetInput("Which profile should perun use as a default?", &profile)
	isUserProfile := helpers.SliceContains(profilesInConfig, profile)
	for !isUserProfile {
		myLogger.GetInput("I cannot find this profile, try again", &profile)
		isUserProfile = helpers.SliceContains(profilesInConfig, profile)
	}
	return profile
}

// If profile exists in .aws/credentials, but not in aws/config, add profile.
func addNewProfileFromCredentialsToConfig(profile string, homePath string, ctx *context.Context, myLogger logger.LoggerInt) {
	profilesInCredentials := getProfilesFromFile(homePath+"/.aws/credentials", myLogger)
	profilesInConfig := getProfilesFromFile(homePath+"/.aws/config", myLogger)
	profiles := findNewProfileInCredentials(profilesInCredentials, profilesInConfig)
	if len(profiles) > 0 {
		for _, prof := range profiles {
			var answer string
			myLogger.GetInput("I found profile "+prof+" in credentials, but not in config. \nCreate new profile in config? Y/N", &answer)
			if strings.ToUpper(answer) == "Y" {
				var region string
				myLogger.GetInput("Region", &region)
				configurator.CreateAWSConfigFile(ctx.Logger, prof, region)
			}
		}
	}
}

// Checking if profile is in .aws/credentials.
func addProfileToCredentials(profile string, homePath string, ctx *context.Context, myLogger logger.LoggerInt) {
	profilesInCredentials := getProfilesFromFile(homePath+"/.aws/credentials", myLogger)
	temp := helpers.SliceContains(profilesInCredentials, profile)
	if !temp {
		configurator.CreateAWSCredentialsFile(ctx, profile)
	} else {
		myLogger.Always("Profile " + profile + " has already credentials")
	}
}

// Creating main.yaml based on .aws/config or in configure mode.
func configIsPresent(profile string, homePath string, ctx *context.Context, myLogger logger.LoggerInt) (string, context.Context) {
	profilesInConfig := getProfilesFromFile(homePath+"/.aws/config", myLogger)
	isDefaultProfile := helpers.SliceContains(profilesInConfig, profile)
	if isDefaultProfile {
		var answer string
		myLogger.GetInput("Default profile exists, do you want to use it *Y* or create your own *N*?", &answer)
		if strings.ToUpper(answer) == "Y" {
			*ctx = createNewMainYaml(profile, homePath, ctx, myLogger)
		} else if strings.ToUpper(answer) == "N" {
			configurator.CreateRequiredFilesInConfigureMode(ctx)

		}
	} else { // isDefaultProfile == false
		profile = useProfileFromConfig(profilesInConfig, profile, myLogger)
		*ctx = createNewMainYaml(profile, homePath, ctx, myLogger)
	}

	return profile, *ctx
}

// Creating new .aws/config and main.yaml for profile.
func newConfigFile(profile string, region string, homePath string, ctx *context.Context, myLogger *logger.Logger) (string, string, context.Context) {
	profile, region = configurator.GetRegionAndProfile(myLogger)
	configurator.CreateAWSConfigFile(ctx.Logger, profile, region)
	*ctx = createNewMainYaml(profile, homePath, ctx, myLogger)
	return profile, region, *ctx
}

// Creating credentials for all present profiles.
func createCredentials(profile string, homePath string, ctx *context.Context, myLogger logger.LoggerInt) {
	isProfileInPresent := isProfileInCredentials(profile, homePath+"/.aws/credentials", myLogger)

	if !isProfileInPresent {
		var answer string
		myLogger.GetInput("I found profile "+profile+" in .aws/config without credentials, add? Y/N", &answer)
		if strings.ToUpper(answer) == "Y" {
			configurator.CreateAWSCredentialsFile(ctx, profile)
		}
	}

}

func createCredentialsBasedOnEnvironmentVariables(envVar map[string]string, myLogger logger.LoggerInt) {
	if len(envVar["profile"]) > 0 {
		myLogger.Always("Creating .aws/credentials file based on environmentVariables")
		homePath, pathError := myuser.GetUserHomeDir()
		if pathError != nil {
			myLogger.Error(pathError.Error())
		}
		path := homePath + "/.aws/credentials"
		line := "[" + envVar["profile"] + "-long-term" + "]\n"
		configurator.AppendStringToFile(path, line)
		line = "aws_access_key_id" + " = " + envVar["id"] + "\n"
		configurator.AppendStringToFile(path, line)
		line = "aws_secret_access_key" + " = " + envVar["key"] + "\n"
		configurator.AppendStringToFile(path, line)
		line = "[" + envVar["profile"] + "]\n"
		configurator.AppendStringToFile(path, line)
		line = "aws_session_token" + " = " + envVar["token"] + "\n"
		configurator.AppendStringToFile(path, line)
	}
}

func checkingCredentials(ctx *context.Context, profile string, region string) (bool, string, string) {
	if EnvironmentVariables["profile"] != "" {
		var answer string
		ctx.Logger.GetInput("Creating aws/credentials based on environment variables? Y/N", &answer)
		if strings.ToUpper(answer) == "Y" {
			createCredentialsBasedOnEnvironmentVariables(EnvironmentVariables, ctx.Logger)
			profile = EnvironmentVariables["profile"]
			region = EnvironmentVariables["region"]
			return false, profile, region
		} else if *ctx.CliArguments.Mode == cliparser.ValidateMode {
			var answer string //offline walidacja
			ctx.Logger.GetInput("You haven't got credentials file, run only offline validation? Y/N", &answer)
			if strings.ToUpper(answer) == "Y" {
				return true, "", "" //offline
			}
		}
	}
	return false, profile, region

}
