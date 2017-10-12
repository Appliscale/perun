package cfonlinevalidator

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudformation"
	"io/ioutil"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
)

func ValidateAndEstimateCosts(filePath *string, region *string) {
	session, _ := createSession(region)
	rawTemplate, error := ioutil.ReadFile(*filePath)
	if error != nil {
		fmt.Println(error)
		return
	}
	template := string(rawTemplate)
	valid, error := isTemplateValid(session, &template)
	if valid != true {
		fmt.Println(error)
	} else {
		fmt.Println("Template is valid.")
	}

	estimateCosts(session, &template)
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

func estimateCosts(session *session.Session, template *string) {
	cfm := cloudformation.New(session)
	templateCostInput := cloudformation.EstimateTemplateCostInput{
		TemplateBody: template,
	}
	output, error := cfm.EstimateTemplateCost(&templateCostInput)

	if error != nil {
		fmt.Println(error)
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