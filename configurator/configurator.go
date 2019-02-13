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
	"eu-west-3":      "https://d2d0mfegowb3wk.cloudfront.net",
	"sa-east-1":      "https://d3c9jyj3w509b0.cloudfront.net",
}

// List of available regions.
var Regions = makeArrayRegions()

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
	createConfigurationFile(yourPath+"/"+yourName, ctx.Logger, myProfile, myRegion)
	*ctx, _ = context.GetContext(cliparser.ParseCliArguments, configuration.GetConfiguration, configuration.ReadInconsistencyConfiguration)
	var answer string
	ctx.Logger.GetInput("Do you want to create .aws/credentials for this profile? Y/N", &answer)
	if strings.ToUpper(answer) == "Y" {
		CreateAWSCredentialsFile(ctx, myProfile)
	}
}

// Creating main.yaml in user's path.
func createConfigurationFile(path string, myLogger logger.LoggerInt, myProfile string, myRegion string) {
	myLogger.Always("File will be created in " + path)
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		con := CreateMainYaml(myLogger, myProfile, myRegion)
		configuration.SaveToFile(con, path, myLogger)
	} else {
		var answer string
		myLogger.GetInput("File already exists in this path. Do you want to overwrite this file? Y/N", &answer)
		if strings.ToUpper(answer) == "Y" {
			con := CreateMainYaml(myLogger, myProfile, myRegion)
			configuration.SaveToFile(con, path, myLogger)

		}
	}
}

//List of all available regions.
func showRegions(myLogger logger.LoggerInt) {
	myLogger.Always("Regions:")
	for i := 0; i < len(Regions); i++ {
		pom := strconv.Itoa(i)
		myLogger.Always("Number " + pom + " region " + Regions[i])
	}
}

// Choosing one region.
func setRegions(myLogger logger.LoggerInt) (region string, err bool) {
	var numberRegion int
	myLogger.GetInput("Choose region", &numberRegion)
	if numberRegion >= 0 && numberRegion < len(ResourceSpecificationURL) {
		region = Regions[numberRegion]
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
func setTemporaryFilesDirectory(myLogger logger.LoggerInt) (path string) {
	myLogger.GetInput("Directory for temporary files", &path)
	myLogger.Always("Your temporary files directory is: " + path)
	return path
}

// CreateMainYaml creates new configuration file.
func CreateMainYaml(myLogger logger.LoggerInt, myProfile string, myRegion string) configuration.Configuration {
	myTemporaryFilesDirectory := setTemporaryFilesDirectory(myLogger)
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
	var regions = []string{}
	for region := range ResourceSpecificationURL {
		regions = append(regions, region)
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
		AppendStringToFile(path, line)
		line = "aws_access_key_id" + " = " + awsAccessKeyID + "\n"
		AppendStringToFile(path, line)
		line = "aws_secret_access_key" + " = " + awsSecretAccessKey + "\n"
		AppendStringToFile(path, line)
		line = "mfa_serial" + " = " + mfaSerial + "\n"
		AppendStringToFile(path, line)
	}
}

// CreateAWSConfigFile creates .aws/config file based on information from user. The file contains profile name, region and type of output.
func CreateAWSConfigFile(myLogger logger.LoggerInt, profile string, region string) {
	var output string
	myLogger.GetInput("Output", &output)
	homePath, pathError := myuser.GetUserHomeDir()
	if pathError != nil {
		myLogger.Error(pathError.Error())
	}
	path := homePath + "/.aws/config"
	line := "[" + profile + "]\n"
	AppendStringToFile(path, line)
	line = "region" + " = " + region + "\n"
	AppendStringToFile(path, line)
	line = "output" + " = " + output + "\n"
	AppendStringToFile(path, line)
}

func AppendStringToFile(path, text string) error {
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
