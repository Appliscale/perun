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

package context

import (
	"errors"
	"github.com/Appliscale/perun/cliparser"
	"os"
	"os/user"
	"time"

	"github.com/Appliscale/perun/utilities"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/go-ini/ini"
)

const dateFormat = "2006-01-02 15:04:05"

// InitializeSession creates session.
func InitializeSession(context *Context) *session.Session {
	tokenError := UpdateSessionToken(context.Config.DefaultProfile, context.Config.DefaultRegion, context.Config.DefaultDurationForMFA, context)
	if tokenError != nil {
		context.Logger.Error(tokenError.Error())
		os.Exit(1)
	}
	currentSession, sessionError := CreateSession(context, context.Config.DefaultProfile, &context.Config.DefaultRegion)
	if sessionError != nil {
		context.Logger.Error(sessionError.Error())
		os.Exit(1)
	}
	return currentSession
}

// CreateSession creates session based on profile and region.
func CreateSession(context *Context, profile string, region *string) (*session.Session, error) {
	context.Logger.Info("Creating new session. Profile: " + profile + " Region: " + *region)

	currentSession, sessionWithOptionError := session.NewSessionWithOptions(
		session.Options{
			Config: aws.Config{
				Region: region,
			},
			Profile: profile,
		})

	if sessionWithOptionError != nil {
		return nil, sessionWithOptionError
	}

	return currentSession, nil
}

// UpdateSessionToken updates session token in .aws/credentials.
func UpdateSessionToken(profile string, region string, defaultDuration int64, context *Context) error {
	if *context.CliArguments.MFA || *context.CliArguments.Mode == cliparser.MfaMode {
		currentUser, userError := user.Current()
		if userError != nil {
			return userError
		}

		credentialsFilePath := currentUser.HomeDir + "/.aws/credentials"
		configuration, loadCredentialsError := ini.Load(credentialsFilePath)
		if loadCredentialsError != nil {
			return loadCredentialsError
		}

		section, sectionError := configuration.GetSection(profile)
		if sectionError != nil {
			section, sectionError = configuration.NewSection(profile)
			if sectionError != nil {
				return sectionError
			}
		}

		profileLongTerm := profile + "-long-term"
		sectionLongTerm, profileLongTermError := configuration.GetSection(profileLongTerm)
		if profileLongTermError != nil {
			return profileLongTermError
		}

		sessionToken := section.Key("aws_session_token")
		expiration := section.Key("expiration")

		expirationDate, dataError := time.Parse(dateFormat, section.Key("expiration").Value())
		if dataError == nil {
			context.Logger.Info("Session token will expire in " + utilities.TruncateDuration(time.Since(expirationDate)).String() + " (" + expirationDate.Format(dateFormat) + ")")
		}

		mfaDevice := sectionLongTerm.Key("mfa_serial").Value()
		if mfaDevice == "" {
			return errors.New("There is no mfa_serial for the profile " + profileLongTerm + ". If you haven't used --mfa option you can change the default decision for MFA in the configuration file")
		}

		if sessionToken.Value() == "" || expiration.Value() == "" || time.Since(expirationDate).Nanoseconds() > 0 {
			currentSession, sessionError := session.NewSessionWithOptions(
				session.Options{
					Config: aws.Config{
						Region: &region,
					},
					Profile: profileLongTerm,
				})
			if sessionError != nil {
				return sessionError
			}

			var tokenCode string
			sessionError = context.Logger.GetInput("MFA token code", &tokenCode)
			if sessionError != nil {
				return sessionError
			}

			var duration int64
			if defaultDuration == 0 {
				sessionError = context.Logger.GetInput("Duration", &duration)
				if sessionError != nil {
					return sessionError
				}
			} else {
				duration = defaultDuration
			}

			stsSession := sts.New(currentSession)
			newToken, tokenError := stsSession.GetSessionToken(&sts.GetSessionTokenInput{
				DurationSeconds: &duration,
				SerialNumber:    aws.String(mfaDevice),
				TokenCode:       &tokenCode,
			})
			if tokenError != nil {
				return tokenError
			}

			section.Key("aws_access_key_id").SetValue(*newToken.Credentials.AccessKeyId)
			section.Key("aws_secret_access_key").SetValue(*newToken.Credentials.SecretAccessKey)
			sessionToken.SetValue(*newToken.Credentials.SessionToken)
			section.Key("expiration").SetValue(newToken.Credentials.Expiration.Format(dateFormat))

			configuration.SaveTo(credentialsFilePath)
		}
	}
	return nil
}
