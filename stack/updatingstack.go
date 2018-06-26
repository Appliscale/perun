package stack

import (
	"github.com/Appliscale/perun/context"
	"github.com/Appliscale/perun/mysession"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudformation"
)

func UpdateStack(context *context.Context) (err error) {
	template, stackName, err := getTemplateFromFile(context)
	if err != nil {
		return
	}

	templateStruct := updateStackInput(context, &template, &stackName)
	session := mysession.InitializeSession(context)
	err = updateStack(templateStruct, session)
	if err != nil {
		return
	}
	return
}

func updateStack(updateStackInput cloudformation.UpdateStackInput, session *session.Session) error {
	api := cloudformation.New(session)
	_, err := api.UpdateStack(&updateStackInput)
	return err
}

// This function gets template and  name of stack. It creates "CreateStackInput" structure.
func updateStackInput(context *context.Context, template *string, stackName *string) cloudformation.UpdateStackInput {
	rawCapabilities := *context.CliArguments.Capabilities
	capabilities := make([]*string, len(rawCapabilities))
	for i, capability := range rawCapabilities {
		capabilities[i] = &capability
	}
	templateStruct := cloudformation.UpdateStackInput{
		TemplateBody: template,
		StackName:    stackName,
		Capabilities: capabilities,
	}
	return templateStruct
}
