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
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewStack(t *testing.T) {
	stackName := "StackName"
	templatePath := "./test_resources/test_template.yaml"
	ctx := stack_mocks.SetupContext(t, []string{"cmd", "create-stack", stackName, templatePath})

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockAWSPI := stack_mocks.NewMockCloudFormationAPI(mockCtrl)
	ctx.CloudFormation = mockAWSPI

	template := stack_mocks.ReadFile(t, templatePath)

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
	ctx := stack_mocks.SetupContext(t, []string{"cmd", "create-stack", stackName, templatePath})
	templateBody := stack_mocks.ReadFile(t, templatePath)
	returnedValue, err := createStackInput(&templateBody, &stackName, ctx)
	if err != nil {
		t.Error(err.Error())
	}

	assert.Equal(t, *returnedValue.StackName, stackName)
	assert.Equal(t, *returnedValue.TemplateBody, templateBody)
}
