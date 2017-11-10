// Copyright 2017 Appliscale
//
// Maintainers and Contributors:
//
//   - Piotr Figwer (piotr.figwer@appliscale.io)
//   - Wojciech Gawroński (wojciech.gawronski@appliscale.io)
//   - Kacper Patro (kacper.patro@appliscale.io)
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

// Package cfonlinevalidator privides tools for online cloudformation template
// validation using AWS API.
package cfonlinevalidator

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/Appliscale/cftool/cflogger"
	"github.com/Appliscale/cftool/cfcontext"
	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/go-ini/ini"
	"os/user"
	"time"
	"io/ioutil"
	"errors"
)

const dateFormat = "2006-01-02 15:04:05 MST"

// Validate template and get URL for cost estimation.
func ValidateAndEstimateCosts(context *cfcontext.Context) bool {
	valid := false
	defer printResult(&valid, context.Logger)

	if *context.CliArguments.MFA {
		err := updateSessionToken(context.Config.Profile, context.Config.Region, context.Logger)
		if err != nil {
			context.Logger.Error(err.Error())
			return false
		}
	}

	session, err := createSession(&context.Config.Region, context.Config.Profile, context.Logger)
	if err != nil {
		context.Logger.Error(err.Error())
		return false
	}

	rawTemplate, err := ioutil.ReadFile(*context.CliArguments.TemplatePath)
	if err != nil {
		context.Logger.Error(err.Error())
		return false
	}

	template := string(rawTemplate)
	valid, err = isTemplateValid(session, &template)
	if err != nil {
		context.Logger.Error(err.Error())
		return false
	}

	estimateCosts(session, &template, context.Logger)

	return valid
}

func isTemplateValid(session *session.Session, template *string) (bool, error) {
	cfm := cloudformation.New(session)
	templateStruct := cloudformation.ValidateTemplateInput{
		TemplateBody: template,
	}
	_, error := cfm.ValidateTemplate(&templateStruct)
	if error != nil {
		return false, error
	}

	return true, nil
}

func estimateCosts(session *session.Session, template *string, logger *cflogger.Logger) {
	cfm := cloudformation.New(session)
	templateCostInput := cloudformation.EstimateTemplateCostInput{
		TemplateBody: template,
	}
	output, err := cfm.EstimateTemplateCost(&templateCostInput)

	if err != nil {
		logger.Error(err.Error())
		return
	}

	logger.Info("Costs estimation: " + *output.Url)
}

func createSession(region *string, profile string, logger *cflogger.Logger) (*session.Session, error) {
	logger.Info("Profile: " + profile)
	logger.Info("Region: " + *region)
	session, err := session.NewSessionWithOptions(
		session.Options{
			Config: aws.Config{
				Region: region,
			},
			Profile: profile,
		})
	if err != nil {
		return nil, err
	}

	return session, nil
}

func updateSessionToken(profile string, region string, logger *cflogger.Logger) error {
	user, err := user.Current()
	if err != nil {
		return err
	}

	credentialsFilePath := user.HomeDir + "/.aws/credentials"
	cfg, err := ini.Load(credentialsFilePath)
	if err != nil {
		return err
	}

	section, err := cfg.GetSection(profile)
	if err != nil {
		section, err = cfg.NewSection(profile)
		if err != nil {
			return err
		}
	}

	profileLongTerm := profile + "-long-term"
	sectionLongTerm, err := cfg.GetSection(profileLongTerm)
	if err != nil {
		return err
	}

	sessionToken := section.Key("aws_session_token")
	expiration := section.Key("expiration")

	expirationDate, err := time.Parse(dateFormat, section.Key("expiration").Value())
	if err == nil {
		logger.Info("Session token will expire in " +
								truncate(time.Since(expirationDate)).String() + " (" + expirationDate.Format(dateFormat) + ")")
	}

	mfaDevice := sectionLongTerm.Key("mfa_serial").Value()
	if mfaDevice == "" {
		return errors.New("There is no mfa_serial for the profile " + profileLongTerm)
	}

	if sessionToken.Value() == "" || expiration.Value() == "" || time.Since(expirationDate).Nanoseconds() > 0 {
		session, err := session.NewSessionWithOptions(
			session.Options{
				Config: aws.Config{
					Region: &region,
				},
				Profile: profileLongTerm,
			})
		if err != nil {
			return err
		}

		var tokenCode string
		err = logger.GetInput("MFA token code", &tokenCode)
		if err != nil {
			return err
		}

		var duration int64
		err = logger.GetInput("Duration", &duration)
		if err != nil {
			return err
		}

		stsSession := sts.New(session)
		newToken, err := stsSession.GetSessionToken(&sts.GetSessionTokenInput{
			DurationSeconds: &duration,
			SerialNumber:    aws.String(mfaDevice),
			TokenCode:       &tokenCode,
		})
		if err != nil {
			return err
		}

		section.Key("aws_access_key_id").SetValue(*newToken.Credentials.AccessKeyId)
		section.Key("aws_secret_access_key").SetValue(*newToken.Credentials.SecretAccessKey)
		sessionToken.SetValue(*newToken.Credentials.SessionToken)
		section.Key("expiration").SetValue(newToken.Credentials.Expiration.Format(dateFormat))

		cfg.SaveTo(credentialsFilePath)
	}

	return nil
}

func printResult(valid *bool, logger *cflogger.Logger) {
	if !*valid {
		logger.Error("Template is invalid!")
	} else {
		logger.Info("Template is valid!")
	}
}

func truncate(d time.Duration) time.Duration {
	return -(d - d % (time.Duration(1) * time.Second))
}
