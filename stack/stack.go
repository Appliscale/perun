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

//prepare stackinput from template
func createStackInput(context *context.Context, template *string, stackName *string) cloudformation.CreateStackInput {
	templateStruct := cloudformation.CreateStackInput{
		TemplateBody: template,
		StackName:    stackName,
	}
	return templateStruct
}

//main function to delete stack
func DestroyStack(context *context.Context) {
	delStackInput := deleteStackInput(context)
	session, err := createSession(context, context.Config.DefaultProfile, &context.Config.DefaultRegion)
	if err != nil {
		context.Logger.Error(err.Error())
	}
	api := cloudformation.New(session)
	api.DeleteStack(&delStackInput)
}

// prepare stackinput from --stack
func deleteStackInput(context *context.Context) cloudformation.DeleteStackInput {
	name := *context.CliArguments.Stack
	templateStruct := cloudformation.DeleteStackInput{
		StackName: &name,
	}
	return templateStruct
}

//template and stackname from file as string
func getTemplateFromFile(context *context.Context) (string, string) {

	rawTemplate, err := ioutil.ReadFile(*context.CliArguments.TemplatePath)
	if err != nil {
		context.Logger.Error(err.Error())
	}
	rawStackName := *context.CliArguments.Stack
	template := string(rawTemplate)
	stackName := string(rawStackName)
	return template, stackName
}

// create stack using stackinput
func createStack(templateStruct cloudformation.CreateStackInput, session *session.Session) {
	api := cloudformation.New(session)
	api.CreateStack(&templateStruct)
}

//main function to create new stack
func NewStack(context *context.Context) {
	err := updateSessionToken(context.Config.DefaultProfile, context.Config.DefaultRegion, context.Config.DefaultDurationForMFA, context)
	if err != nil {
		context.Logger.Error(err.Error())
	}
	session, err1 := createSession(context, context.Config.DefaultProfile, &context.Config.DefaultRegion)
	if err1 != nil {
		context.Logger.Error(err1.Error())
	}
	template, stackName := getTemplateFromFile(context)
	templateStruct := createStackInput(context, &template, &stackName)

	createStack(templateStruct, session)
}

//****online validator
func createSession(context *context.Context, profile string, region *string) (*session.Session, error) {
	context.Logger.Info("Profile: " + profile)
	context.Logger.Info("Region: " + *region)

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

func updateSessionToken(profile string, region string, defaultDuration int64, context *context.Context) error {
	user, err := user.Current()
	if err != nil {
		return err
	}

	credentialsFilePath := user.HomeDir + "/.aws/credentials"
	configuration, err := ini.Load(credentialsFilePath)
	if err != nil {
		return err
	}

	section, err := configuration.GetSection(profile)
	if err != nil {
		section, err = configuration.NewSection(profile)
		if err != nil {
			return err
		}
	}

	profileLongTerm := profile + "-long-term"
	sectionLongTerm, err := configuration.GetSection(profileLongTerm)
	if err != nil {
		return err
	}

	sessionToken := section.Key("aws_session_token")
	expiration := section.Key("expiration")

	expirationDate, err := time.Parse(dateFormat, section.Key("expiration").Value())
	if err == nil {
		context.Logger.Info("Session token will expire in " + utilities.TruncateDuration(time.Since(expirationDate)).String() + " (" + expirationDate.Format(dateFormat) + ")")
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
		err = context.Logger.GetInput("MFA token code", &tokenCode)
		if err != nil {
			return err
		}

		var duration int64
		if defaultDuration == 0 {
			err = context.Logger.GetInput("Duration", &duration)
			if err != nil {
				return err
			}
		} else {
			duration = defaultDuration
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

		configuration.SaveTo(credentialsFilePath)
	}

	return nil
}
