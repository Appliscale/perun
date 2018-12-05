package checkingrequiredfiles

import (
	"github.com/Appliscale/perun/cliparser"
	"github.com/Appliscale/perun/configuration"
	"github.com/Appliscale/perun/configurator"
	"github.com/Appliscale/perun/context"
	"github.com/Appliscale/perun/helpers"
	"github.com/Appliscale/perun/logger"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/ec2metadata"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"os"
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
	profilesInCredentials := getProfilesFromFile(homePath + "/.aws/credentials")
	profilesInConfig := getProfilesFromFile(homePath + "/.aws/config")
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
	profilesInCredentials := getProfilesFromFile(homePath + "/.aws/credentials")
	temp := helpers.SliceContains(profilesInCredentials, profile)
	if !temp {
		configurator.CreateAWSCredentialsFile(ctx, profile)
	} else {
		myLogger.Always("Profile " + profile + " has already credentials")
	}
}

// Creating main.yaml based on .aws/config or in configure mode.
func configIsPresent(profile string, homePath string, ctx *context.Context, myLogger logger.LoggerInt) (string, context.Context) {
	profilesInConfig := getProfilesFromFile(homePath + "/.aws/config")
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
	configurator.CreateAWSConfigFile(myLogger, profile, region)
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

func getIamInstanceProfileAssociations(myLogger logger.LoggerInt, region string) (*ec2.DescribeIamInstanceProfileAssociationsOutput, error) {
	// Create a Session with a custom region
	sess := session.Must(session.NewSession(&aws.Config{
		Region: &region,
	}))
	svc := ec2.New(sess)
	input := &ec2.DescribeIamInstanceProfileAssociationsInput{}

	result, err := svc.DescribeIamInstanceProfileAssociations(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			default:
				myLogger.Error(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			myLogger.Error(err.Error())
		}
		return result, err
	}

	return result, nil
}

// Get AWS region and check if perun is running on EC2.
func getRegion() (string, bool, error) {
	svc := ec2metadata.New(session.New())
	region, err := svc.Region()
	if err != nil {
		return "", false, err
	}
	return region, true, nil
}

// Get IAM Instance profile name to use it as profile name.
func getInstanceProfileName(output *ec2.DescribeIamInstanceProfileAssociationsOutput) string {
	arn := output.IamInstanceProfileAssociations[0].IamInstanceProfile.Arn
	name := strings.SplitAfter(*arn, "/")
	return name[len(name)-1]
}

// Getting information about EC2 and prepare to run perun there.
func workingOnEC2(myLogger logger.LoggerInt) (profile string, region string, err error) {
	region, _, regionError := getRegion()
	myLogger.Info("Running on EC2")
	if regionError != nil {
		myLogger.Error(regionError.Error())
		return "", "", regionError
	}
	instanceProfileAssociations, instanceError := getIamInstanceProfileAssociations(myLogger, region)
	if instanceError != nil {
		myLogger.Error(instanceError.Error())
		return "", "", instanceError
	}
	instanceProfileName := getInstanceProfileName(instanceProfileAssociations)
	return instanceProfileName, region, nil
}

// Create context and main.yaml if perun is running on EC2.
func createEC2context(profile string, homePath string, region string, ctx *context.Context, myLogger logger.LoggerInt) context.Context {
	con := configurator.CreateMainYaml(myLogger, profile, region)
	_, err := os.Stat(homePath + "/.config/perun")
	if os.IsNotExist(err) {
		mkdirError := os.MkdirAll(homePath+"/.config/perun", 0755)
		if mkdirError != nil {
			myLogger.Error(mkdirError.Error())
		}
	}
	configuration.SaveToFile(con, homePath+"/.config/perun/main.yaml", myLogger)
	*ctx, _ = context.GetContext(cliparser.ParseCliArguments, configuration.GetConfiguration, configuration.ReadInconsistencyConfiguration)
	return *ctx
}
