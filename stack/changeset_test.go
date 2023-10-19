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
	"github.com/Appliscale/perun/context"
	"github.com/Appliscale/perun/stack/stack_mocks"
	"github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"log"
	"os"
	"testing"
)

func TestNewChangeSet(t *testing.T) {
	stackName := "StackName"
	templatePath := "./test_resources/test_template.yaml"
	changeSetName := "ChangeSetName"
	template := stack_mocks.ReadFile(t, templatePath)
	ctx := stack_mocks.SetupContext(t, []string{"cmd", "create-change-set", stackName, templatePath, changeSetName})

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockAWSPI := stack_mocks.NewMockCloudFormationAPI(mockCtrl)
	ctx.CloudFormation = mockAWSPI

	describeChangeSetInput := cloudformation.DescribeChangeSetInput{
		ChangeSetName: &changeSetName,
		StackName:     &stackName,
	}

	mockAWSPI.EXPECT().WaitUntilChangeSetCreateComplete(&describeChangeSetInput).Return(nil).Times(2)
	output := cloudformation.DescribeChangeSetOutput{
		Changes: []*cloudformation.Change{},
	}
	mockAWSPI.EXPECT().DescribeChangeSet(&describeChangeSetInput).Return(&output, nil).Times(2)
	changeSetInput, err := createChangeSetInput(&template, &stackName, nil, ctx)
	assert.Empty(t, err)
	mockAWSPI.EXPECT().CreateChangeSet(&changeSetInput).Return(nil, nil).Times(2)
	updateStackInput := cloudformation.UpdateStackInput{
		TemplateBody: &template,
		StackName:    &stackName,
	}
	mockAWSPI.EXPECT().UpdateStack(&updateStackInput).Return(nil, nil).Times(1) //This shouldn't be called when user input is no

	testChangeSetCreationWithUserInput("y", NewChangeSet, ctx)
	testChangeSetCreationWithUserInput("n", NewChangeSet, ctx)

}

func TestCreateChangeSetInput(t *testing.T) {
	stackName := "StackName"
	templatePath := "./test_resources/test_template.yaml"
	changeSetName := "ChangeSetName"
	ctx := stack_mocks.SetupContext(t, []string{"cmd", "create-change-set", stackName, templatePath, changeSetName})
	template := stack_mocks.ReadFile(t, templatePath)

	returnedInput, err := createChangeSetInput(&template, &stackName, []*cloudformation.Parameter{}, ctx)
	assert.Empty(t, err)
	assert.Equal(t, *returnedInput.StackName, stackName)
	assert.Equal(t, *returnedInput.TemplateBody, template)
	assert.Equal(t, returnedInput.Parameters, []*cloudformation.Parameter{})
	assert.Equal(t, *returnedInput.ChangeSetName, changeSetName)
}

func TestDescribeChangeSet(t *testing.T) {
	stackName := "StackName"
	templatePath := "./test_resources/test_template.yaml"
	changeSetName := "ChangeSetName"
	ctx := stack_mocks.SetupContext(t, []string{"cmd", "create-change-set", stackName, templatePath, changeSetName})

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockAWSPI := stack_mocks.NewMockCloudFormationAPI(mockCtrl)
	ctx.CloudFormation = mockAWSPI

	describeChangeSetInput := cloudformation.DescribeChangeSetInput{
		ChangeSetName: &changeSetName,
		StackName:     &stackName,
	}

	mockAWSPI.EXPECT().WaitUntilChangeSetCreateComplete(&describeChangeSetInput).Return(nil).Times(1)
	output := cloudformation.DescribeChangeSetOutput{
		Changes: []*cloudformation.Change{},
	}
	mockAWSPI.EXPECT().DescribeChangeSet(&describeChangeSetInput).Return(&output, nil).Times(1)
	describeChangeSet(ctx)

}

func TestShouldExecuteChangeSet(t *testing.T) {
	assert.True(t, testCheckUserInput("Y", shouldExecuteChangeSet))
	assert.False(t, testCheckUserInput("N", shouldExecuteChangeSet))
	assert.True(t, testCheckUserInput("y", shouldExecuteChangeSet))
	assert.False(t, testCheckUserInput("n", shouldExecuteChangeSet))
}

type checkFunction func() bool
type newChangeSetFunction func(*context.Context) error

func testCheckUserInput(userInput string, function checkFunction) bool {
	tmpfile, oldStdin := supportStdInputReplacement(userInput)
	defer os.Remove(tmpfile.Name())        // clean up
	defer func() { os.Stdin = oldStdin }() // Restore original Stdin
	defer tmpfile.Close()

	return function()
}

func testChangeSetCreationWithUserInput(userInput string, function newChangeSetFunction, context *context.Context) error {
	tmpfile, oldStdin := supportStdInputReplacement(userInput)
	defer os.Remove(tmpfile.Name())        // clean up
	defer func() { os.Stdin = oldStdin }() // Restore original Stdin
	defer tmpfile.Close()

	return function(context)
}

func supportStdInputReplacement(userInput string) (*os.File, *os.File) {
	content := []byte(userInput)
	tmpfile, err := ioutil.TempFile("", "example")
	if err != nil {
		log.Fatal(err)
	}
	if _, err := tmpfile.Write(content); err != nil {
		log.Fatal(err)
	}
	if _, err := tmpfile.Seek(0, 0); err != nil {
		log.Fatal(err)
	}
	oldStdin := os.Stdin
	os.Stdin = tmpfile
	return tmpfile, oldStdin
}
