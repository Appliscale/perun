package stack

import (
	"github.com/Appliscale/perun/stack/stack_mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUpdateStack(t *testing.T) {
	stackName := "StackName"
	templatePath := "./test_resources/test_template.yaml"
	context := stack_mocks.SetupContext(t, []string{"cmd", "update-stack", stackName, templatePath})

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockAWSPI := stack_mocks.NewMockCloudFormationAPI(mockCtrl)
	context.CloudFormation = mockAWSPI

	template := stack_mocks.ReadFile(t, templatePath)
	templateStruct := updateStackInput(context, &template, &stackName)

	mockAWSPI.EXPECT().UpdateStack(&templateStruct).Return(nil, nil).Times(1)
	UpdateStack(context)
}

func TestDoUpdateStack(t *testing.T) {
	stackName := "StackName"
	templatePath := "./test_resources/test_template.yaml"
	context := stack_mocks.SetupContext(t, []string{"cmd", "update-stack", stackName, templatePath})

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockAWSPI := stack_mocks.NewMockCloudFormationAPI(mockCtrl)
	context.CloudFormation = mockAWSPI

	template := stack_mocks.ReadFile(t, templatePath)
	templateStruct := updateStackInput(context, &template, &stackName)

	mockAWSPI.EXPECT().UpdateStack(&templateStruct).Return(nil, nil).Times(1)
	UpdateStack(context)
}

func TestUpdateStackInput(t *testing.T) {
	stackName := "StackName"
	templatePath := "./test_resources/test_template.yaml"
	context := stack_mocks.SetupContext(t, []string{"cmd", "update-stack", stackName, templatePath})

	template := stack_mocks.ReadFile(t, templatePath)

	returnedTemplateStruct := updateStackInput(context, &template, &stackName)

	assert.Equal(t, *returnedTemplateStruct.TemplateBody, template)
	assert.Equal(t, *returnedTemplateStruct.StackName, stackName)
}
