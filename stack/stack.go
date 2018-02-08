package stack

import (
	"errors"
	"github.com/Appliscale/perun/context"
	"github.com/Appliscale/perun/utilities"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/go-ini/ini"
	"io/ioutil"
	"os/user"
	"time"
)

const dateFormat = "2006-01-02 15:04:05 MST"

//This function gets template and  name of stack. It creates "CreateStackInput" structure.
func createStackInput(context *context.Context, template *string, stackName *string) cloudformation.CreateStackInput {
	templateStruct := cloudformation.CreateStackInput{
		TemplateBody: template,
		StackName:    stackName,
	}
	return templateStruct
}

// This function reads "StackName" from Stack in CliArguments and file from TemplatePath in CliArguments. It converts these to type string.
func getTemplateFromFile(context *context.Context) (string, string) {

	rawTemplate, readFileError := ioutil.ReadFile(*context.CliArguments.TemplatePath)
	if readFileError != nil {
		context.Logger.Error(readFileError.Error())
	}
	rawStackName := *context.CliArguments.Stack
	template := string(rawTemplate)
	stackName := string(rawStackName)
	return template, stackName
}

// This function uses CreateStackInput variable to create Stack.
func createStack(templateStruct cloudformation.CreateStackInput, session *session.Session) {
	api := cloudformation.New(session)
	api.CreateStack(&templateStruct)
}

//This function uses all functions above and session to create Stack.
func NewStack(context *context.Context) {
	tokenError := updateSessionToken(context.Config.DefaultProfile, context.Config.DefaultRegion, context.Config.DefaultDurationForMFA, context)
	if tokenError != nil {
		context.Logger.Error(tokenError.Error())
	}
	session, createSessionError := createSession(context, context.Config.DefaultProfile, &context.Config.DefaultRegion)
	if createSessionError != nil {
		context.Logger.Error(createSessionError.Error())
	}
	template, stackName := getTemplateFromFile(context)
	templateStruct := createStackInput(context, &template, &stackName)

	createStack(templateStruct, session)
}

//This function bases on "DeleteStackInput" structure and destroys stack. It uses "StackName" to choose which stack will be destroy. Before that it creates session.
func DestroyStack(context *context.Context) {
	delStackInput := deleteStackInput(context)
	session, sessionError := createSession(context, context.Config.DefaultProfile, &context.Config.DefaultRegion)
	if sessionError != nil {
		context.Logger.Error(sessionError.Error())
	}
	api := cloudformation.New(session)
	api.DeleteStack(&delStackInput)
}

//This function gets "StackName" from Stack in CliArguments and creates "DeleteStackInput" structure.
func deleteStackInput(context *context.Context) cloudformation.DeleteStackInput {
	name := *context.CliArguments.Stack
	templateStruct := cloudformation.DeleteStackInput{
		StackName: &name,
	}
	return templateStruct
}

//"createSession" and "updateSessionToken" are from onlinevalidator.go file. They allow to connect with AWS.
func createSession(context *context.Context, profile string, region *string) (*session.Session, error) {
	context.Logger.Info("Profile: " + profile)
	context.Logger.Info("Region: " + *region)

	session, sessionWithOptionError := session.NewSessionWithOptions(
		session.Options{
			Config: aws.Config{
				Region: region,
			},
			Profile: profile,
		})

	if sessionWithOptionError != nil {
		return nil, sessionWithOptionError
	}

	return session, nil
}

func updateSessionToken(profile string, region string, defaultDuration int64, context *context.Context) error {
	user, userError := user.Current()
	if userError != nil {
		return userError
	}

	credentialsFilePath := user.HomeDir + "/.aws/credentials"
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
		return errors.New("There is no mfa_serial for the profile " + profileLongTerm)
	}

	if sessionToken.Value() == "" || expiration.Value() == "" || time.Since(expirationDate).Nanoseconds() > 0 {
		session, sessionError := session.NewSessionWithOptions(
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

		stsSession := sts.New(session)
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

	return nil
}
