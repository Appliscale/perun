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
	"github.com/Appliscale/perun/stack/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCreateUpdateTerminationProtectionInput(t *testing.T) {
	stackName := "StackName"

	returnedTemplateStruct := createUpdateTerminationProtectionInput(stackName, true)

	assert.Equal(t, stackName, *returnedTemplateStruct.StackName)
	assert.Equal(t, true, *returnedTemplateStruct.EnableTerminationProtection)

	returnedTemplateStruct = createUpdateTerminationProtectionInput(stackName, false)

	assert.Equal(t, stackName, *returnedTemplateStruct.StackName)
	assert.Equal(t, false, *returnedTemplateStruct.EnableTerminationProtection)
}

func TestIsProtectionEnabled(t *testing.T) {
	stackName := "StackName"
	templatePath := "./test_resources/test_template.yaml"

	context := mocks.SetupContext(t, []string{"cmd", "set-stack-policy", stackName, templatePath, "--enable-stack-termination"})
	enabled, err := isProtectionEnable(context)
	assert.False(t, enabled)
	assert.Equal(t, err, nil)

	context = mocks.SetupContext(t, []string{"cmd", "set-stack-policy", stackName, templatePath, "--disable-stack-termination"})
	enabled, err = isProtectionEnable(context)
	assert.Equal(t, err, nil)
	assert.True(t, enabled)

	context = mocks.SetupContext(t, []string{"cmd", "set-stack-policy", stackName, templatePath})
	_, err = isProtectionEnable(context)
	assert.NotEmpty(t, err)
}

func TestSetTerminationProtection(t *testing.T) {
	stackName := "StackName"
	templatePath := "./test_resources/test_template.yaml"
	context := mocks.SetupContext(t, []string{"cmd", "set-stack-policy", stackName, templatePath, "--enable-stack-termination"})

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockAWSPI := mocks.NewMockCloudFormationAPI(mockCtrl)
	context.CloudFormation = mockAWSPI

	templateStruct := createUpdateTerminationProtectionInput(stackName, false)

	mockAWSPI.EXPECT().UpdateTerminationProtection(&templateStruct).Return(nil, nil).Times(1)
	SetTerminationProtection(context)
}
