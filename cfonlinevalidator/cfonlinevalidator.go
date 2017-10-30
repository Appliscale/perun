package cfonlinevalidator

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudformation"
	"io/ioutil"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/Appliscale/cftool/cfconfiguration"
	"github.com/Appliscale/cftool/cflogger"
)

func ValidateAndEstimateCosts(filePath *string, configPath *string) {
	valid := false
	logger := cflogger.Logger{}
	defer printResult(&valid, &logger)

	region, err := cfconfiguration.GetRegion(*configPath)
	if err != nil {
		cflogger.LogError(&logger, err.Error())
		return
	}

	session, err := createSession(&region)
	if err != nil {
		cflogger.LogError(&logger, err.Error())
		return
	}

	rawTemplate, err := ioutil.ReadFile(*filePath)
	if err != nil {
		cflogger.LogError(&logger, err.Error())
		return
	}

	template := string(rawTemplate)
	valid, err = isTemplateValid(session, &template)
	if err != nil {
		cflogger.LogError(&logger, err.Error())
		return
	}

	estimateCosts(session, &template, &logger)
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
		cflogger.LogError(logger, err.Error())
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
	cflogger.PrintErrors(logger)
	if !*valid {
		fmt.Println("Template is invalid!")
	} else {
		fmt.Println("Template is valid!")
	}
}