package stack

import (
	"github.com/Appliscale/perun/context"
	"github.com/Appliscale/perun/helpers"
	"github.com/Appliscale/perun/offlinevalidator/template"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudformation"
	gocloudformation "github.com/awslabs/goformation/cloudformation"
	"github.com/olekukonko/tablewriter"
	"io/ioutil"
	"os"
	"reflect"
)

// Get StackResources description.
func getResourceDescription(context *context.Context, session *session.Session) (cloudformation.DescribeStackResourcesOutput, error) {
	stackname := cloudformation.DescribeStackResourcesInput{
		StackName: context.CliArguments.Stack,
	}
	api := cloudformation.New(session)
	describeStackResourcesOutput, describeStackResourcesError := api.DescribeStackResources(&stackname)

	return *describeStackResourcesOutput, describeStackResourcesError

}

// Get information about Stack.
func getStackDescription(context *context.Context, session *session.Session) (cloudformation.DescribeStacksOutput, error) {
	name := cloudformation.DescribeStacksInput{
		StackName: context.CliArguments.Stack,
	}
	api := cloudformation.New(session)
	describeStacksOutput, describeStacksError := api.DescribeStacks(&name)

	return *describeStacksOutput, describeStacksError
}

// Get Stack Policy Body.
func getStackPolicy(context *context.Context, session *session.Session) (cloudformation.GetStackPolicyOutput, error) {
	name := cloudformation.GetStackPolicyInput{
		StackName: context.CliArguments.Stack,
	}
	api := cloudformation.New(session)
	stackPolicyOutput, stackPolicyError := api.GetStackPolicy(&name)

	return *stackPolicyOutput, stackPolicyError
}

//Get user's template.
func getTemplateFromUser(context *context.Context) gocloudformation.Template {
	rawTemplate, err := ioutil.ReadFile(*context.CliArguments.TemplatePath)
	if err != nil {
		context.Logger.Error(err.Error())

	}

	var perunTemplate template.Template

	parser, err := helpers.GetParser(*context.CliArguments.TemplatePath)
	if err != nil {
		context.Logger.Error(err.Error())

	}
	goFormationTemplate, err := parser(rawTemplate, perunTemplate, context.Logger)
	if err != nil {
		context.Logger.Error(err.Error())

	}
	return goFormationTemplate
}

//Making map from template - current state.
func currentTemplateMap(description cloudformation.DescribeStackResourcesOutput) map[int]*cloudformation.StackResource {
	resourceMap := make(map[int]*cloudformation.StackResource)
	for index, resource := range description.StackResources {
		resourceMap[index] = resource
	}

	return resourceMap
}

//Making map with stack's information.
func stackInformationMap(description cloudformation.DescribeStacksOutput) map[int]*cloudformation.Stack {
	resourceMap := make(map[int]*cloudformation.Stack)
	for index, resource := range description.Stacks {
		resourceMap[index] = resource

	}
	return resourceMap
}

//Looking for differences between template and current state.
func Diff(ctx *context.Context) {
	currentSession := context.InitializeSession(ctx)
	resourceDescription, resourcesError := getResourceDescription(ctx, currentSession)
	if resourcesError != nil {
		ctx.Logger.Error(resourcesError.Error())
	}

	stackDescription, stackError := getStackDescription(ctx, currentSession)
	if stackError != nil {
		ctx.Logger.Error(stackError.Error())
	}

	stackPolicy, stackPolicyError := getStackPolicy(ctx, currentSession)
	if stackPolicyError != nil {
		ctx.Logger.Error(stackPolicyError.Error())
	}

	informationAboutStack := stackInformationMap(stackDescription)
	mainTemplate := getTemplateFromUser(ctx)
	resourceInfo := currentTemplateMap(resourceDescription)

	var data [][]string
	var title, before, after string

	title, before, after = comparingStackDescription(mainTemplate, informationAboutStack)
	if title != "" && before != "" && after != "" {
		data = append(data, []string{title, before, after})
	}

	title, before, after = comparingStackTerminationProtection(informationAboutStack)
	data = append(data, []string{title, before, after})

	title, before, after = comparingOutputs(mainTemplate, informationAboutStack)
	if title != "" && before != "" && after != "" {
		data = append(data, []string{title, before, after})
	}

	parameters := comparingParameters(mainTemplate, informationAboutStack)
	for _, parameter := range parameters {
		data = append(data, parameter)
	}

	resources := comparingResources(mainTemplate, resourceInfo)
	for _, resource := range resources {
		data = append(data, resource)
	}

	tags := comparingTags(informationAboutStack)
	for _, tag := range tags {
		data = append(data, tag)
	}

	if stackPolicy.StackPolicyBody != nil {
		title, before, after = comparingStackPolicy(stackPolicy)
		data = append(data, []string{title, before, after})
	}

	showTable(data, ctx)
}

//Showing differences in the table.
func showTable(data [][]string, context *context.Context) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{*context.CliArguments.Stack, "Template", "Current"})
	table.SetHeaderColor(tablewriter.Colors{tablewriter.Bold, tablewriter.BgBlueColor},
		tablewriter.Colors{tablewriter.FgHiBlackColor, tablewriter.Bold, tablewriter.BgYellowColor},
		tablewriter.Colors{tablewriter.FgHiBlackColor, tablewriter.Bold, tablewriter.BgGreenColor})

	table.SetColumnColor(tablewriter.Colors{tablewriter.Bold, tablewriter.FgHiBlueColor},
		tablewriter.Colors{tablewriter.Bold, tablewriter.FgHiYellowColor},
		tablewriter.Colors{tablewriter.Bold, tablewriter.FgHiGreenColor})

	table.AppendBulk(data)
	table.Render()
}

//Functions which compare chosen template's elements.
func comparingStackDescription(template gocloudformation.Template, informationAboutStack map[int]*cloudformation.Stack) (string, string, string) {

	if template.Description != "" && informationAboutStack[0].Description != nil {
		stack := informationAboutStack[0]
		descr := stack.Description
		valStack := *descr

		if template.Description != valStack {
			return "Description", template.Description, valStack
		}
	} else if template.Description == "" && informationAboutStack[0].Description != nil {
		stack := informationAboutStack[0]
		descr := stack.Description
		valStack := *descr
		return "Description", "<empty>", valStack

	} else if template.Description != "" && informationAboutStack[0].Description == nil {
		return "Description", template.Description, "<empty>"
	}
	return "", "", ""
}

func comparingStackTerminationProtection(informationAboutStack map[int]*cloudformation.Stack) (string, string, string) {
	stack := informationAboutStack[0]
	prot := stack.EnableTerminationProtection
	valprot := *prot
	var protection string
	if valprot == true {
		protection = "Enable"
	} else {
		protection = "Disable"
	}
	return "Stack Termination Protection", "<no information>", protection
}

func comparingTags(informationAboutStack map[int]*cloudformation.Stack) (data [][]string) {
	stack := informationAboutStack[0]
	tags := stack.Tags
	for _, tag := range tags {
		key := *tag.Key
		value := *tag.Value
		element := key + " " + value
		data = append(data, []string{"Tags", element, "<empty>"})
	}
	return data
}

func comparingOutputs(template gocloudformation.Template, informationAboutStack map[int]*cloudformation.Stack) (string, string, string) {
	if template.Outputs != nil && informationAboutStack[0] != nil {
		stack := informationAboutStack[0]
		output := stack.Outputs
		for _, element := range output {
			index := *element.OutputKey
			if element != template.Outputs[index] {
				return "Output", template.Outputs[index].(string), element.String()
			}
		}
	} else if template.Outputs == nil && informationAboutStack[0] != nil {
		stack := informationAboutStack[0]
		output := stack.Outputs
		for _, element := range output {
			return "Output", "<empty>", element.String()
		}
	} else if template.Outputs != nil && informationAboutStack[0] == nil {
		for a := range template.Outputs {
			return "Output", template.Outputs[a].(string), "<empty>"
		}
	}
	return "", "", ""
}

func comparingParameters(template gocloudformation.Template, informationAboutStack map[int]*cloudformation.Stack) (data [][]string) {

	if template.Parameters != nil && informationAboutStack[0] != nil {
		stack := informationAboutStack[0]
		parameters := stack.Parameters
		for _, element := range parameters {
			index := *element.ParameterKey
			if element != template.Parameters[index] {
				data = append(data, []string{"Parameters", template.Parameters[index].(string), element.String()})
			}
		}
	}
	return data
}

func comparingResources(template gocloudformation.Template, resources map[int]*cloudformation.StackResource) (data [][]string) {

	fromTemplate := template.Resources
	fromInfo := resources

	for resourcename, value := range fromTemplate {
		valueof := reflect.ValueOf(value)
		for i := 0; i < len(fromInfo); i++ {
			logicalResourceID := *fromInfo[i].LogicalResourceId
			if logicalResourceID == resourcename {
				for _, key := range valueof.MapKeys() {
					keyString := key.Interface().(string)
					keyValue := valueof.MapIndex(key)
					keyValueString := keyValue.Interface().(string)

					if keyString == "Type" {
						resourceType := *fromInfo[i].ResourceType

						if resourceType != keyValueString {
							data = append(data, []string{"Resource type", keyValueString, resourceType})
						}
						if keyString == "Description" {
							description := *fromInfo[i].Description
							if description != keyValueString {
								data = append(data, []string{"Resource description", keyValueString, description})
							}

						}
					}

				}
			}
		}
	}
	return data
}

func comparingStackPolicy(stackPolicy cloudformation.GetStackPolicyOutput) (string, string, string) {
	policy := *stackPolicy.StackPolicyBody
	return "StackPolicy", "<no information>", policy
}
