package linter

import (
	"github.com/Appliscale/perun/context"
	"github.com/Appliscale/perun/helpers"
	"github.com/Appliscale/perun/offlinevalidator/template"
	"io/ioutil"
	"os"
	"path"
	"regexp"
	"strings"
	"strconv"
	"github.com/awslabs/goformation/cloudformation"
)

type Parameter struct {
	Type          string   `json:"Type"`
	Default       string   `json:"Default"`
	AllowedValues []string `json:"AllowedValues"`
	Description   string   `json:"Description"`
}

func CheckStyle(ctx *context.Context) (err error) {

	err, lintConf := GetLinterConfiguration(ctx)

if err != nil {
		os.Exit(1)
	}

	templateExtension := path.Ext(*ctx.CliArguments.TemplatePath)
	templateBytes, err := ioutil.ReadFile(*ctx.CliArguments.TemplatePath)
	if err != nil {
		ctx.Logger.Error(err.Error())
		return
	}
	rawTemplate := string(templateBytes)

	lines := strings.Split(rawTemplate, "\n")

	checkAWSCFSpecificStuff(ctx, rawTemplate, lintConf)
	checkBlankLines(lintConf, rawTemplate, ctx)
	checkLineLengths(lines, lintConf, ctx)

	if templateExtension == ".json" {
		checkJsonIndentation(ctx, lintConf, lines)
		checkJsonSpaces(ctx, lintConf, lines)
	} else if templateExtension == ".yaml" {
		checkYamlIndentation(ctx, lintConf, lines)
		checkYamlQuotes(ctx, lintConf, lines)
		checkYamlLists(ctx, lintConf, rawTemplate)
	}

	return
}

func checkLineLengths(lines []string, lintConf LinterConfiguration, ctx *context.Context) {
	for line := range lines {
		if lintConf.Global.LineLength.Required && len(lines[line]) > int(lintConf.Global.LineLength.Value.(float64)) {
			ctx.Logger.Warning("line " + strconv.Itoa(line) + ": maximum line lenght exceeded")
		}
	}
}

func checkBlankLines(lintConf LinterConfiguration, rawTemplate string, ctx *context.Context) {
	if !lintConf.Global.BlankLinesAllowed && regexp.MustCompile("\n\n").MatchString(rawTemplate) {
		ctx.Logger.Warning("Blank lines are not allowed in current lint configuration")
	}
}
func checkAWSCFSpecificStuff(ctx *context.Context, rawTemplate string, lintConf LinterConfiguration) {
	var perunTemplate template.Template
	parser, err := helpers.GetParser(*ctx.CliArguments.TemplatePath)
	if err != nil {
		ctx.Logger.Error(err.Error())
		return
	}
	var goFormationTemplate cloudformation.Template
	goFormationTemplate, err = parser([]byte(rawTemplate), perunTemplate, ctx.Logger)
	if err != nil {
		ctx.Logger.Error(err.Error())
		return
	}

	if lintConf.Global.RequiredFields.TemplateDescription && goFormationTemplate.Description == "" {
		ctx.Logger.Warning("The template has no description")
	}

	if lintConf.Global.RequiredFields.ParametersDescription {
		for parameterName, parameterValue := range goFormationTemplate.Parameters {
			if parameterValue.(map[string]interface{})["Description"] != nil {
				ctx.Logger.Warning("No description provided for parameter " + parameterName)
			}
		}
	}

	logicalNameRegex := regexp.MustCompile(lintConf.Global.NamingConventions.LogicalNames)
	for resourceName := range goFormationTemplate.Resources {
		if !logicalNameRegex.MatchString(resourceName) {
			ctx.Logger.Warning("Resource '" + resourceName + "' does not meat the given logical Name regex")
		}
	}
}

func checkJsonSpaces(ctx *context.Context, lintConf LinterConfiguration, lines []string) {
	for line := range lines {
		for sign := range lintConf.Json.Spaces.After {
			if strings.Count(lines[line], lintConf.Json.Spaces.After[sign]) != strings.Count(lines[line], lintConf.Json.Spaces.After[sign]+" ") {
				ctx.Logger.Warning("line " + strconv.Itoa(line) + ": no space after '" + string(lintConf.Json.Spaces.After[sign]) + "'")
			}
		}
		for sign := range lintConf.Json.Spaces.Before {
			if strings.Count(lines[line], lintConf.Json.Spaces.Before[sign]) != strings.Count(lines[line], " "+lintConf.Json.Spaces.Before[sign]) {
				ctx.Logger.Warning("line " + strconv.Itoa(line) + ": no space before '" + string(lintConf.Json.Spaces.Before[sign]) + "'")
			}
		}
	}
}

func checkYamlLists(ctx *context.Context, lintConf LinterConfiguration, template string) {
	preprocessed := regexp.MustCompile("#.*\n").ReplaceAllString(template, "\n")
	dashListRegex := regexp.MustCompile(".*- .*")
	inlineListRegex := regexp.MustCompile(`.*: \[.*].*`)
	if !lintConf.Yaml.AllowedLists.Dash && dashListRegex.MatchString(preprocessed) {
		ctx.Logger.Warning("dash list are not allowed in current lint configuration")
	}
	if !lintConf.Yaml.AllowedLists.Inline && inlineListRegex.MatchString(preprocessed) {
		ctx.Logger.Warning("inline lists are not allowed in current lint configuration")
	}
}

func checkYamlQuotes(ctx *context.Context, lintConf LinterConfiguration, lines []string) {
	for line := range lines {
		if !lintConf.Yaml.AllowedQuotes.Double && strings.Contains(lines[line], "\"") {
			ctx.Logger.Warning("line " + strconv.Itoa(line) + ": double quotes not allowed")
		}
		if !lintConf.Yaml.AllowedQuotes.Single && strings.Contains(lines[line], "\"") {
			ctx.Logger.Warning("line " + strconv.Itoa(line) + ": single quotes not allowed")
		}
		noQuotesRegex := regexp.MustCompile(".*: [^\"']*")
		if !lintConf.Yaml.AllowedQuotes.Noquotes && noQuotesRegex.MatchString(lines[line]) {
			ctx.Logger.Warning("line " + strconv.Itoa(line) + ": quotes required")
		}
	}
}

func checkYamlIndentation(ctx *context.Context, lintConf LinterConfiguration, lines []string) {
	indent := int(lintConf.Global.Indent.Value.(float64))
	last_spaces := 0
	for line := range lines {
		curr_spaces := countLeadingSpace(lines[line])
		if lintConf.Global.Indent.Required {
			if curr_spaces%indent != 0 {
				ctx.Logger.Error("line " + strconv.Itoa(line) + ": indentation error")
			}
		}

		if last_spaces < curr_spaces {
			if last_spaces + indent != curr_spaces ||
				wrongYAMLContinuationIndent(lintConf, lines, line, last_spaces, curr_spaces) {
				ctx.Logger.Error("line " + strconv.Itoa(line) + ": indentation error")
			}
		}
		last_spaces = curr_spaces
	}
}

func wrongYAMLContinuationIndent(lintConf LinterConfiguration, lines []string, line int, last_spaces int, curr_spaces int) bool {
	return lintConf.Yaml.ContinuationIndent.Required && !strings.Contains(lines[line], ": ") && !strings.Contains(lines[line], "- ") &&
		!strings.HasSuffix(lines[line], ":") && last_spaces + int(lintConf.Yaml.ContinuationIndent.Value.(float64)) != curr_spaces
}

func checkJsonIndentation(ctx *context.Context, lintConf LinterConfiguration, lines []string) {
	var last_spaces = 0
	if lintConf.Global.Indent.Required {
		indent := int(lintConf.Global.Indent.Value.(float64))
		jsonIndentations := map[string]int{
			",": 0,
			"{": indent,
			"}": -indent,
			"[": indent,
			"]": -indent,
		}
		last_spaces = 0
		for line := range lines {
			if line == 0 {
				continue
			}
			prevLine := lines[line-1]
			indentation, found := jsonIndentations[string(prevLine[len(prevLine)-1])]
			if !found {
				indentation = -indentation
			}
			curr_spaces := countLeadingSpace(lines[line])
			if curr_spaces-last_spaces != indentation {
				ctx.Logger.Error("line " + strconv.Itoa(line) + ": indentation error")
			}
			last_spaces = curr_spaces
		}
	}
}

func countLeadingSpace(line string) int {
	i := 0
	for _, runeValue := range line {
		if runeValue == ' ' {
			i++
		} else {
			break
		}
	}
	return i
}
