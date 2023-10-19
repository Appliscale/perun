// Copyright 2018 Appliscale
//
// Maintainers and contributors are listed in README file inside repository.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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
