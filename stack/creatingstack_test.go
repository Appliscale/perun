package stack

import (
	"github.com/Appliscale/perun/stack/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewStack(t *testing.T) {
	stackName := "StackName"
	templatePath := "./test_resources/test_template.yaml"
	ctx := mocks.SetupContext(t, []string{"cmd", "create-stack", stackName, templatePath})

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockAWSPI := mocks.NewMockCloudFormationAPI(mockCtrl)
	ctx.CloudFormation = mockAWSPI

	template := mocks.ReadFile(t, templatePath)

	input, err := createStackInput(&template, &stackName, ctx)
	if err != nil {
		t.Error(err.Error())
	}
	mockAWSPI.EXPECT().CreateStack(&input).Return(nil, nil).Times(1)

	NewStack(ctx)
}

func TestCreateStackInput(t *testing.T) {
	stackName := "StackName"
	templatePath := "./test_resources/test_template.yaml"
	ctx := mocks.SetupContext(t, []string{"cmd", "create-stack", stackName, templatePath})
	templateBody := mocks.ReadFile(t, templatePath)
	returnedValue, err := createStackInput(&templateBody, &stackName, ctx)
	if err != nil {
		t.Error(err.Error())
	}

	assert.Equal(t, *returnedValue.StackName, stackName)
	assert.Equal(t, *returnedValue.TemplateBody, templateBody)
}
