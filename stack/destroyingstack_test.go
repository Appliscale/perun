package stack

import (
	"github.com/Appliscale/perun/stack/stack_mocks"
	"github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDestroyStack(t *testing.T) {
	stackName := "StackName"
	context := stack_mocks.SetupContext(t, []string{"cmd", "delete-stack", stackName})

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockAWSPI := stack_mocks.NewMockCloudFormationAPI(mockCtrl)
	context.CloudFormation = mockAWSPI

	templateStruct := cloudformation.DeleteStackInput{
		StackName: &stackName,
	}
	mockAWSPI.EXPECT().DeleteStack(&templateStruct).Return(nil, nil).Times(1)

	DestroyStack(context)
}

func TestDeleteStackInput(t *testing.T) {
	stackName := "StackName"
	context := stack_mocks.SetupContext(t, []string{"cmd", "delete-stack", stackName})
	returnedValue := deleteStackInput(context)
	assert.Equal(t, returnedValue.StackName, &stackName)
}
