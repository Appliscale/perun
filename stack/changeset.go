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
	"fmt"
	"github.com/Appliscale/perun/context"
	"github.com/Appliscale/perun/parameters"
	"github.com/Appliscale/perun/progress"
	"github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/olekukonko/tablewriter"
	"strings"
)

func createChangeSetInput(template *string, stackName *string, params []*cloudformation.Parameter, context *context.Context) (cloudformation.CreateChangeSetInput, error) {

	templateStruct := cloudformation.CreateChangeSetInput{
		Parameters:    params,
		ChangeSetName: context.CliArguments.ChangeSet,
		TemplateBody:  template,
		StackName:     stackName,
	}
	return templateStruct, nil
}

func createDeleteChangeSetInput(ctx *context.Context) cloudformation.DeleteChangeSetInput {
	return cloudformation.DeleteChangeSetInput{
		ChangeSetName: ctx.CliArguments.ChangeSet,
		StackName:     ctx.CliArguments.Stack,
	}
}

// DeleteChangeSet deletes change set.
func DeleteChangeSet(ctx *context.Context) (err error) {
	templateStruct := createDeleteChangeSetInput(ctx)
	_, err = ctx.CloudFormation.DeleteChangeSet(&templateStruct)
	if err != nil {
		ctx.Logger.Error(err.Error())
		return
	}
	ctx.Logger.Info("Deletion of Change Set " + *ctx.CliArguments.ChangeSet + " request successful")
	return
}

// NewChangeSet create change set and gets parameters.
func NewChangeSet(context *context.Context) (err error) {
	template, stackName, err := getTemplateFromFile(context)
	if err != nil {
		return
	}

	params, err := parameters.ResolveParameters(context)
	if err != nil {
		context.Logger.Error(err.Error())
		return
	}

	templateStruct, templateError := createChangeSetInput(&template, &stackName, params, context)
	if templateError != nil {
		return
	}

	_, err = context.CloudFormation.CreateChangeSet(&templateStruct)

	if err != nil {
		context.Logger.Error(err.Error())
		return
	}

	describeChangeSet(context)

	if shouldExecuteChangeSet() {
		templateStruct := cloudformation.UpdateStackInput{
			Parameters:   params,
			TemplateBody: &template,
			StackName:    &stackName,
		}
		doUpdateStack(context, templateStruct)
	}
	return
}

func shouldExecuteChangeSet() bool {
	println("Do You want to execute the change set? (Y/N) ")
	for true {
		var executeChangeSet string
		fmt.Scanf("%s", &executeChangeSet)
		if strings.ToLower(executeChangeSet) == "n" {
			return false
		} else if strings.ToLower(executeChangeSet) == "y" {
			return true
		}
	}
	return false
}

func describeChangeSet(context *context.Context) error {
	context.Logger.Info("Waiting for change set creation...")
	describeChangeSetInput := cloudformation.DescribeChangeSetInput{
		ChangeSetName: context.CliArguments.ChangeSet,
		StackName:     context.CliArguments.Stack,
	}

	err := context.CloudFormation.WaitUntilChangeSetCreateComplete(&describeChangeSetInput)
	if err != nil {
		context.Logger.Error(err.Error())
		return err
	}

	describeChangeSetOutput, err := context.CloudFormation.DescribeChangeSet(&describeChangeSetInput)
	if err != nil {
		context.Logger.Error(err.Error())
		return err
	}

	_, table := initStackTableWriter()
	for rowNum := range describeChangeSetOutput.Changes {
		currRow := describeChangeSetOutput.Changes[rowNum]
		var physicalResourceId string = ""
		var replacement string = ""
		if currRow.ResourceChange.PhysicalResourceId != nil {
			physicalResourceId = *currRow.ResourceChange.PhysicalResourceId
		}
		if currRow.ResourceChange.Replacement != nil {
			replacement = *currRow.ResourceChange.Replacement
		}
		table.Append([]string{
			*currRow.ResourceChange.Action,
			*currRow.ResourceChange.LogicalResourceId,
			physicalResourceId,
			*currRow.ResourceChange.ResourceType,
			replacement,
		})
	}
	table.Render()
	return nil
}

func initStackTableWriter() (*progress.ParseWriter, *tablewriter.Table) {
	pw := progress.NewParseWriter()
	table := tablewriter.NewWriter(pw)
	table.SetHeader([]string{"Action", "Logical ID", "Physical ID", "Resource Type", "Replacement"})
	table.SetBorder(false)
	table.SetColumnColor(
		tablewriter.Colors{tablewriter.Bold},
		tablewriter.Colors{tablewriter.Bold, tablewriter.FgCyanColor},
		tablewriter.Colors{tablewriter.FgWhiteColor},
		tablewriter.Colors{tablewriter.FgWhiteColor},
		tablewriter.Colors{tablewriter.FgWhiteColor, tablewriter.Bold})
	return pw, table
}
