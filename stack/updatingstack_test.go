package stack

import (
	"github.com/Appliscale/perun/stack/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUpdateStack(t *testing.T) {
	stackName := "StackName"
	templatePath := "./test_resources/test_template.yaml"
	context := mocks.SetupContext(t, []string{"cmd", "update-stack", stackName, templatePath})

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockAWSPI := mocks.NewMockCloudFormationAPI(mockCtrl)
	context.CloudFormation = mockAWSPI

	template := mocks.ReadFile(t, templatePath)
	templateStruct := updateStackInput(context, &template, &stackName)

	mockAWSPI.EXPECT().UpdateStack(&templateStruct).Return(nil, nil).Times(1)
	UpdateStack(context)
}

func TestDoUpdateStack(t *testing.T) {
	stackName := "StackName"
	templatePath := "./test_resources/test_template.yaml"
	context := mocks.SetupContext(t, []string{"cmd", "update-stack", stackName, templatePath})

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockAWSPI := mocks.NewMockCloudFormationAPI(mockCtrl)
	context.CloudFormation = mockAWSPI

	template := mocks.ReadFile(t, templatePath)
	templateStruct := updateStackInput(context, &template, &stackName)

	mockAWSPI.EXPECT().UpdateStack(&templateStruct).Return(nil, nil).Times(1)
	UpdateStack(context)
}

func TestUpdateStackInput(t *testing.T) {
	stackName := "StackName"
	templatePath := "./test_resources/test_template.yaml"
	context := mocks.SetupContext(t, []string{"cmd", "update-stack", stackName, templatePath})

	template := mocks.ReadFile(t, templatePath)

	returnedTemplateStruct := updateStackInput(context, &template, &stackName)

	assert.Equal(t, *returnedTemplateStruct.TemplateBody, template)
	assert.Equal(t, *returnedTemplateStruct.StackName, stackName)
}
