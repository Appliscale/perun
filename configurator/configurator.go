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

// Package configurator allows to create configuration file main.yaml.
package configurator

import (
	"github.com/Appliscale/perun/cliparser"
	"github.com/Appliscale/perun/configuration"
	"github.com/Appliscale/perun/context"
	"github.com/Appliscale/perun/logger"
	"github.com/Appliscale/perun/myuser"
	"os"
	"strconv"
	"strings"
)

// ResourceSpecificationURL contains links to Resource Specification for all regions.
var ResourceSpecificationURL = map[string]string{
	"us-east-2":      "https://dnwj8swjjbsbt.cloudfront.net",
	"us-east-1":      "https://d1uauaxba7bl26.cloudfront.net",
	"us-west-1":      "https://d68hl49wbnanq.cloudfront.net",
	"us-west-2":      "https://d201a2mn26r7lk.cloudfront.net",
	"ap-south-1":     "https://d2senuesg1djtx.cloudfront.net",
	"ap-northeast-2": "https://d1ane3fvebulky.cloudfront.net",
	"ap-southeast-1": "https://doigdx0kgq9el.cloudfront.net",
	"ap-southeast-2": "https://d2stg8d246z9di.cloudfront.net",
	"ap-northeast-1": "https://d33vqc0rt9ld30.cloudfront.net",
	"ca-central-1":   "https://d2s8ygphhesbe7.cloudfront.net",
	"eu-central-1":   "https://d1mta8qj7i28i2.cloudfront.net",
	"eu-west-1":      "https://d3teyb21fexa9r.cloudfront.net",
	"eu-west-2":      "https://d1742qcu2c1ncx.cloudfront.net",
	"sa-east-1":      "https://d3c9jyj3w509b0.cloudfront.net",
}

// CreateRequiredFilesInConfigureMode creates main.yaml and .aws/credentials in configure mode.
func CreateRequiredFilesInConfigureMode(ctx *context.Context) {
	homePath, pathError := myuser.GetUserHomeDir()
	if pathError != nil {
		ctx.Logger.Error(pathError.Error())
	}
	myLogger := logger.CreateDefaultLogger()
	homePath += "/.config/perun"
	ctx.Logger.Always("Configure file could be in \n  " + homePath + "\n  /etc/perun")
	var yourPath string
	var yourName string
	ctx.Logger.GetInput("Your path", &yourPath)
	ctx.Logger.GetInput("Filename", &yourName)
	myProfile, myRegion := GetRegionAndProfile(&myLogger)
	createConfigurationFile(yourPath+"/"+yourName, ctx, myProfile, myRegion)
	*ctx, _ = context.GetContext(cliparser.ParseCliArguments, configuration.GetConfiguration, configuration.ReadInconsistencyConfiguration)
	var answer string
	ctx.Logger.GetInput("Do you want to create .aws/credentials for this profile? Y/N", &answer)
	if strings.ToUpper(answer) == "Y" {
		CreateAWSCredentialsFile(ctx, myProfile)
	}
}

// Creating main.yaml in user's path.
func createConfigurationFile(path string, context *context.Context, myProfile string, myRegion string) {
	context.Logger.Always("File will be created in " + path)
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		con := CreateMainYaml(context, myProfile, myRegion)
		configuration.SaveToFile(con, path, context.Logger)
	} else {
		var answer string
		context.Logger.GetInput("File already exists in this path. Do you want to overwrite this file? Y/N", &answer)
		if strings.ToUpper(answer) == "Y" {
			con := CreateMainYaml(context, myProfile, myRegion)
			configuration.SaveToFile(con, path, context.Logger)

		}
	}
}

//List of all available regions.
func showRegions(myLogger logger.LoggerInt) {
	regions := makeArrayRegions()
	myLogger.Always("Regions:")
	for i := 0; i < len(regions); i++ {
		pom := strconv.Itoa(i)
		myLogger.Always("Number " + pom + " region " + regions[i])
	}
}

// Choosing one region.
func setRegions(myLogger logger.LoggerInt) (region string, err bool) {
	var numberRegion int
	myLogger.GetInput("Choose region", &numberRegion)
	regions := makeArrayRegions()
	if numberRegion >= 0 && numberRegion < 14 {
		region = regions[numberRegion]
		myLogger.Always("Your region is: " + region)
		err = true
	} else {
		myLogger.Error("Invalid region")
		err = false
	}
	return
}

// Choosing one profile.
func setProfile(myLogger logger.LoggerInt) (profile string, err bool) {
	myLogger.GetInput("Input name of profile", &profile)
	if profile != "" {
		myLogger.Always("Your profile is: " + profile)
		err = true
	} else {
		myLogger.Error("Invalid profile")
		err = false
	}
	return
}

//GetRegionAndProfile gets region and profile from user.
func GetRegionAndProfile(myLogger logger.LoggerInt) (string, string) {
	profile, err := setProfile(myLogger)
	for !err {
		myLogger.Always("Try again, invalid profile")
		profile, err = setProfile(myLogger)
	}
	showRegions(myLogger)
	region, err1 := setRegions(myLogger)
	for !err1 {
		myLogger.Always("Try again, invalid region")
		region, err1 = setRegions(myLogger)
	}
	return profile, region
}

// Setting directory for temporary files.
func setTemporaryFilesDirectory(context *context.Context) (path string) {
	context.Logger.GetInput("Directory for temporary files", &path)
	context.Logger.Always("Your temporary files directory is: " + path)
	return path
}

// CreateMainYaml creates new configuration file.
func CreateMainYaml(context *context.Context, myProfile string, myRegion string) configuration.Configuration {
	myTemporaryFilesDirectory := setTemporaryFilesDirectory(context)
	myResourceSpecificationURL := ResourceSpecificationURL

	myConfig := configuration.Configuration{
		DefaultProfile:                 myProfile,
		DefaultRegion:                  myRegion,
		SpecificationURL:               myResourceSpecificationURL,
		DefaultDecisionForMFA:          false,
		DefaultDurationForMFA:          3600,
		DefaultVerbosity:               "INFO",
		DefaultTemporaryFilesDirectory: myTemporaryFilesDirectory,
	}

	return myConfig
}

// Array of regions.
func makeArrayRegions() []string {
	var regions = []string{
		"us-east-1",
		"us-east-2",
		"us-west-1",
		"us-west-2",
		"ca-central-1",
		"ca-central-1",
		"eu-west-1",
		"eu-west-2",
		"ap-northeast-1",
		"ap-northeast-2",
		"ap-southeast-1",
		"ap-southeast-2",
		"ap-south-1",
		"sa-east-1",
	}
	return regions
}

// CreateAWSCredentialsFile creates .aws/credentials file based on information from user. The file contains access key and MFA serial.
func CreateAWSCredentialsFile(ctx *context.Context, profile string) {
	if profile != "" {
		ctx.Logger.Always("You haven't got .aws/credentials file for profile " + profile)
		var awsAccessKeyID string
		var awsSecretAccessKey string
		var mfaSerial string

		ctx.Logger.GetInput("awsAccessKeyID", &awsAccessKeyID)
		ctx.Logger.GetInput("awsSecretAccessKey", &awsSecretAccessKey)
		ctx.Logger.GetInput("mfaSerial", &mfaSerial)

		homePath, pathError := myuser.GetUserHomeDir()
		if pathError != nil {
			ctx.Logger.Error(pathError.Error())
		}
		path := homePath + "/.aws/credentials"
		line := "[" + profile + "-long-term" + "]\n"
		appendStringToFile(path, line)
		line = "aws_access_key_id" + " = " + awsAccessKeyID + "\n"
		appendStringToFile(path, line)
		line = "aws_secret_access_key" + " = " + awsSecretAccessKey + "\n"
		appendStringToFile(path, line)
		line = "mfa_serial" + " = " + mfaSerial + "\n"
		appendStringToFile(path, line)
	}
}

// CreateAWSConfigFile creates .aws/config file based on information from user. The file contains profile name, region and type of output.
func CreateAWSConfigFile(ctx *context.Context, profile string, region string) {
	var output string
	ctx.Logger.GetInput("Output", &output)
	homePath, pathError := myuser.GetUserHomeDir()
	if pathError != nil {
		ctx.Logger.Error(pathError.Error())
	}
	path := homePath + "/.aws/config"
	line := "[" + profile + "]\n"
	appendStringToFile(path, line)
	line = "region" + " = " + region + "\n"
	appendStringToFile(path, line)
	line = "output" + " = " + output + "\n"
	appendStringToFile(path, line)
}

func appendStringToFile(path, text string) error {
	f, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.WriteString(text)
	if err != nil {
		return err
	}
	return nil
}
