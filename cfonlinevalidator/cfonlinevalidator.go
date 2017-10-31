package cfonlinevalidator

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudformation"
	"io/ioutil"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/Appliscale/cftool/cflogger"
	"github.com/Appliscale/cftool/cfcontext"
)

func ValidateAndEstimateCosts(context *cfcontext.Context) bool {
	valid := false
	defer printResult(&valid, context.Logger)

	session, err := createSession(&context.Config.Region)
	if err != nil {
		context.Logger.Error(err.Error())
		return false
	}

	rawTemplate, err := ioutil.ReadFile(*context.CliArguments.FilePath)
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
	/*resp, _ := http.Get(*output.Url)
	bytes, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("HTML:\n\n", string(bytes))*/
}

func createSession(endpoint *string) (*session.Session, error) {
	session, error := session.NewSession(&aws.Config{
		Region: endpoint,
	})
	return session, error
}

func printResult(valid *bool, logger *cflogger.Logger) {
	if !*valid {
		logger.Info("Template is invalid!")
	} else {
		logger.Info("Template is valid!")
	}
}