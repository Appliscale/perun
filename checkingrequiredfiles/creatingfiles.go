package checkingrequiredfiles

import (
	"github.com/Appliscale/perun/cliparser"
	"github.com/Appliscale/perun/configuration"
	"github.com/Appliscale/perun/configurator"
	"github.com/Appliscale/perun/context"
	"github.com/Appliscale/perun/helpers"
	"github.com/Appliscale/perun/logger"
	"strings"
)

type CreatingFiles interface {
	createNewMainYaml(profile string, homePath string, ctx *context.Context, myLogger logger.LoggerInt) context.Context
	useProfileFromConfig(profilesInConfig []string, profile string, myLogger *logger.Logger) string
	addNewProfileFromCredentialsToConfig(profile string, homePath string, ctx *context.Context, myLogger *logger.Logger)
	addProfileToCredentials(profile string, homePath string, ctx *context.Context, myLogger *logger.Logger)
	configIsPresent(profile string, homePath string, ctx *context.Context, myLogger logger.Logger) (string, context.Context)
	newConfigFile(profile string, region string, homePath string, ctx *context.Context, myLogger *logger.Logger) (string, string, context.Context)
	createCredentials(profile string, homePath string, ctx *context.Context, myLogger *logger.Logger)
}

// Creating main.yaml.
func createNewMainYaml(profile string, homePath string, ctx *context.Context, myLogger logger.LoggerInt) context.Context {
	region := findRegionForProfile(profile, homePath+"/.aws/config", myLogger)
	con := configurator.CreateMainYaml(ctx, profile, region)
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
				configurator.CreateAWSConfigFile(ctx, prof, region)
			}
		}
	}
}

// Checking if profile is in .aws/credentials.
func addProfileToCredentials(profile string, homePath string, ctx *context.Context, myLogger *logger.Logger) {
	profilesInCredentials := getProfilesFromFile(homePath+"/.aws/credentials", myLogger)
	temp := helpers.SliceContains(profilesInCredentials, profile)
	if !temp {
		configurator.CreateAWSCredentialsFile(ctx, profile)
	} else {
		myLogger.Always("Profile " + profile + " has already credentials")
	}
}

// Creating main.yaml based on .aws/config or in configure mode.
func configIsPresent(profile string, homePath string, ctx *context.Context, myLogger logger.Logger) (string, context.Context) {
	profilesInConfig := getProfilesFromFile(homePath+"/.aws/config", &myLogger)
	isDefaultProfile := helpers.SliceContains(profilesInConfig, profile)
	if isDefaultProfile {
		var answer string
		myLogger.GetInput("Default profile exists, do you want to use it *Y* or create your own *N*?", &answer)
		if strings.ToUpper(answer) == "Y" {
			*ctx = createNewMainYaml(profile, homePath, ctx, &myLogger)
		} else if strings.ToUpper(answer) == "N" {
			configurator.CreateRequiredFilesInConfigureMode(ctx)

		}
	} else { // isDefaultProfile == false
		profile = useProfileFromConfig(profilesInConfig, profile, &myLogger)
		*ctx = createNewMainYaml(profile, homePath, ctx, &myLogger)
	}

	return profile, *ctx
}

// Creating new .aws/config and main.yaml for profile.
func newConfigFile(profile string, region string, homePath string, ctx *context.Context, myLogger *logger.Logger) (string, string, context.Context) {
	profile, region = configurator.GetRegionAndProfile(myLogger)
	configurator.CreateAWSConfigFile(ctx, profile, region)
	*ctx = createNewMainYaml(profile, homePath, ctx, myLogger)
	return profile, region, *ctx
}

// Creating credentials for all present profiles.
func createCredentials(profile string, homePath string, ctx *context.Context, myLogger *logger.Logger) {
	isProfileInPresent := isProfileInCredentials(profile, homePath+"/.aws/credentials", myLogger)

	if !isProfileInPresent {
		var answer string
		myLogger.GetInput("I found profile "+profile+" in .aws/config without credentials, add? Y/N", &answer)
		if strings.ToUpper(answer) == "Y" {
			configurator.CreateAWSCredentialsFile(ctx, profile)
		}
	}

}
