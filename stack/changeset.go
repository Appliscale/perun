package stack

import (
	"fmt"
	"github.com/Appliscale/perun/context"
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

func NewChangeSet(context *context.Context) (err error) {
	template, stackName, err := getTemplateFromFile(context)
	if err != nil {
		return
	}

	params, err := getParameters(context)
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

func describeChangeSet(context *context.Context) (err error) {
	context.Logger.Info("Waiting for change set creation...")
	describeChangeSetInput := cloudformation.DescribeChangeSetInput{
		ChangeSetName: context.CliArguments.ChangeSet,
		StackName:     context.CliArguments.Stack,
	}

	err = context.CloudFormation.WaitUntilChangeSetCreateComplete(&describeChangeSetInput)
	if err != nil {
		context.Logger.Error(err.Error())
		return
	}

	describeChangeSetOutput, err := context.CloudFormation.DescribeChangeSet(&describeChangeSetInput)
	if err != nil {
		context.Logger.Error(err.Error())
		return
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
	return
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
