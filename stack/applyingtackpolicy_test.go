package stack

import (
	"github.com/Appliscale/perun/stack/stack_mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestApplyStackPolicy(t *testing.T) {
	stackName := "StackName"
	policyPath := "./test_resources/test_stackpolicy.json"
	ctx := stack_mocks.SetupContext(t, []string{"cmd", "set-stack-policy", stackName, policyPath})

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockAWSPI := stack_mocks.NewMockCloudFormationAPI(mockCtrl)
	ctx.CloudFormation = mockAWSPI

	template := stack_mocks.ReadFile(t, policyPath)

	input := createStackPolicyInput(&template, &stackName)
	mockAWSPI.EXPECT().SetStackPolicy(&input).Return(nil, nil).Times(1)

	ApplyStackPolicy(ctx)
}

func TestCreateStackPolicyInput(t *testing.T) {
	stackName := "StackName"
	templateBody := "TestTemplate"
	returnedValue := createStackPolicyInput(&templateBody, &stackName)
	assert.Equal(t, *returnedValue.StackName, stackName)
	assert.Equal(t, *returnedValue.StackPolicyBody, templateBody)
}
