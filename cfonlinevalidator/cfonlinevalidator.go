package cfonlinevalidator

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudformation"
	"io/ioutil"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/Appliscale/cftool/cflogger"
	"github.com/Appliscale/cftool/cfcontext"
)

func ValidateAndEstimateCosts(context *cfcontext.Context) {
	valid := false
	defer printResult(&valid, context.Logger)

	session, err := createSession(&context.Config.Region)
	if err != nil {
		context.Logger.LogError(err.Error())
		return
	}

	rawTemplate, err := ioutil.ReadFile(*context.CliArguments.FilePath)
	if err != nil {
		context.Logger.LogError(err.Error())
		return
	}

	template := string(rawTemplate)
	valid, err = isTemplateValid(session, &template)
	if err != nil {
		context.Logger.LogError(err.Error())
		return
	}

	estimateCosts(session, &template, context.Logger)
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
		logger.LogError(err.Error())
		return
	}

	fmt.Println("Costs estimation: " + *output.Url)
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
		fmt.Println("Template is invalid!")
	} else {
		fmt.Println("Template is valid!")
	}
}