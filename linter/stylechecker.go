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

package linter

import (
	"github.com/Appliscale/perun/context"
	"github.com/Appliscale/perun/helpers"
	"github.com/Appliscale/perun/validator/template"
	"github.com/awslabs/goformation/cloudformation"
	"io/ioutil"
	"path"
	"regexp"
	"strconv"
	"strings"
)

// Parameter describes all element which parameter in template should have.
type Parameter struct {
	Type          string   `json:"Type"`
	Default       string   `json:"Default"`
	AllowedValues []string `json:"AllowedValues"`
	Description   string   `json:"Description"`
}

// CheckStyle gets linter configuration and run checking.
func CheckStyle(ctx *context.Context) (err error) {

	err, lintConf := GetLinterConfiguration(ctx)
	if err != nil {
		return
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
			ctx.Logger.Warning("line " + strconv.Itoa(line+1) + ": maximum line lenght exceeded")
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
			if parameterValue.(map[string]interface{})["Description"] == nil {
				ctx.Logger.Warning("No description provided for parameter " + parameterName)
			}
		}
	}

	for resourceName := range goFormationTemplate.Resources {
		if !lintConf.CheckLogicalName(resourceName) {
			ctx.Logger.Warning("Resource '" + resourceName + "' does not meet the given logical Name regex: " + lintConf.Global.NamingConventions.LogicalNames)
		}
	}
}

func checkJsonSpaces(ctx *context.Context, lintConf LinterConfiguration, lines []string) {
	reg := regexp.MustCompile(`"([^"]*)"`)
	for line := range lines {
		for sign := range lintConf.Json.Spaces.After {
			if strings.Count(reg.ReplaceAllString(lines[line], "\"*\""), lintConf.Json.Spaces.After[sign]) != strings.Count(reg.ReplaceAllString(lines[line], "\"*\""), lintConf.Json.Spaces.After[sign]+" ") {
				ctx.Logger.Warning("line " + strconv.Itoa(line+1) + ": no space after '" + string(lintConf.Json.Spaces.After[sign]) + "'")
			}
		}
		for sign := range lintConf.Json.Spaces.Before {
			if strings.Count(reg.ReplaceAllString(lines[line], "\"*\""), lintConf.Json.Spaces.Before[sign]) != strings.Count(reg.ReplaceAllString(lines[line], "\"*\""), " "+lintConf.Json.Spaces.Before[sign]) {
				ctx.Logger.Warning("line " + strconv.Itoa(line+1) + ": no space before '" + string(lintConf.Json.Spaces.Before[sign]) + "'")
			}
		}
	}
}

func checkYamlLists(ctx *context.Context, lintConf LinterConfiguration, template string) {
	preprocessed := regexp.MustCompile("#.*\n").ReplaceAllString(template, "\n")
	dashListRegex := regexp.MustCompile(".*- .*")
	inlineListRegex := regexp.MustCompile(`.*: \[.*].*`)
	if !lintConf.Yaml.AllowedLists.Dash && dashListRegex.MatchString(preprocessed) {
		ctx.Logger.Warning("dash lists are not allowed in current lint configuration")
	}
	if !lintConf.Yaml.AllowedLists.Inline && inlineListRegex.MatchString(preprocessed) {
		ctx.Logger.Warning("inline lists are not allowed in current lint configuration")
	}
}

func checkYamlQuotes(ctx *context.Context, lintConf LinterConfiguration, lines []string) {
	for line := range lines {
		if !lintConf.Yaml.AllowedQuotes.Double && strings.Contains(lines[line], "\"") {
			ctx.Logger.Warning("line " + strconv.Itoa(line+1) + ": double quotes not allowed")
		}
		if !lintConf.Yaml.AllowedQuotes.Single && strings.Contains(lines[line], "'") {
			ctx.Logger.Warning("line " + strconv.Itoa(line+1) + ": single quotes not allowed")
		}
		noQuotesRegex := regexp.MustCompile(".*: [^\"']*")
		if !lintConf.Yaml.AllowedQuotes.Noquotes && noQuotesRegex.MatchString(lines[line]) {
			ctx.Logger.Warning("line " + strconv.Itoa(line+1) + ": quotes required")
		}
	}
}

func checkYamlIndentation(ctx *context.Context, lintConf LinterConfiguration, lines []string) {
	indent := int(lintConf.Global.Indent.Value.(float64))
	last_spaces := 0
	for line := range lines {
		if strings.HasPrefix(strings.TrimSpace(lines[line]), "#") {
			continue
		}
		curr_spaces := helpers.CountLeadingSpaces(lines[line])
		if lintConf.Global.Indent.Required {
			if curr_spaces%indent != 0 || (last_spaces < curr_spaces && last_spaces+indent != curr_spaces) {
				ctx.Logger.Error("line " + strconv.Itoa(line+1) + ": indentation error")
			}
		}

		if last_spaces < curr_spaces {
			if wrongYAMLContinuationIndent(lintConf, lines, line, last_spaces, curr_spaces) {
				ctx.Logger.Error("line " + strconv.Itoa(line+1) + ": continuation indent error")
			}
		}
		last_spaces = curr_spaces
	}
}

func wrongYAMLContinuationIndent(lintConf LinterConfiguration, lines []string, line int, last_spaces int, curr_spaces int) bool {
	return lintConf.Yaml.ContinuationIndent.Required && !strings.Contains(lines[line], ": ") && !strings.Contains(lines[line], "- ") &&
		!strings.HasSuffix(lines[line], ":") && last_spaces+int(lintConf.Yaml.ContinuationIndent.Value.(float64)) != curr_spaces
}

func checkJsonIndentation(ctx *context.Context, lintConf LinterConfiguration, lines []string) {
	var last_spaces = 0
	if lintConf.Global.Indent.Required {
		indent := int(lintConf.Global.Indent.Value.(float64))
		jsonIndentations := map[string]int{
			",":  0,
			"{":  indent,
			"}":  -indent,
			"[":  indent,
			"]":  -indent,
			"\"": -indent,
		}
		last_spaces = 0
		for line := range lines {
			if line == 0 || len(lines[line]) == 0 {
				continue
			}
			prevLine := lines[line-1]
			indentation, found := jsonIndentations[string(prevLine[len(prevLine)-1])]
			if !found {
				indentation = -indentation
			}
			curr_spaces := helpers.CountLeadingSpaces(lines[line])
			if curr_spaces-last_spaces != indentation {
				ctx.Logger.Error("line " + strconv.Itoa(line+1) + ": indentation error")
			}
			last_spaces = curr_spaces
		}
	}
}
