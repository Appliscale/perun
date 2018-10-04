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

// Package linter provides revision of templates.
package linter

import (
	"encoding/json"
	"github.com/Appliscale/perun/configuration"
	"github.com/Appliscale/perun/context"
	"github.com/ghodss/yaml"
	"path"
	"regexp"
)

// LinterConfiguration contains configuration for two types: Yaml and JSON, and global.
type LinterConfiguration struct {
	Yaml   YamlLinterConfiguration   `yaml:"yaml"`
	Json   JsonLinterConfiguration   `yaml:"json"`
	Global GlobalLinterConfiguration `yaml:"global"`
}

// GlobalLinterConfiguration describes global configuration.
type GlobalLinterConfiguration struct {
	LineLength        Check             `yaml:"lineLength"`
	Indent            Check             `yaml:"indent"`
	RequiredFields    RequiredFields    `yaml:"requiredFields"`
	NamingConventions NamingConventions `yaml:"namingConventions"`
	BlankLinesAllowed bool              `yaml:"blankLinesAllowed"`
}

// Check stores "Required" and its value.
type Check struct {
	Required bool        `yaml:"required"`
	Value    interface{} `yaml:"value"`
}

// NamingConventions describes names.
type NamingConventions struct {
	LogicalNames string `yaml:"logicalNames"`
}

// RequiredFields stores Description.
type RequiredFields struct {
	TemplateDescription   bool `yaml:"templateDescription"`
	ParametersDescription bool `yaml:"parametersDescription"`
}

// JsonLinterConfiguration describes Spaces.
type JsonLinterConfiguration struct {
	Spaces SpacesConfiguration `yaml:"spaces"`
}

// SpacesConfiguration stores information about spaces.
type SpacesConfiguration struct {
	After  []string `yaml:"after"`
	Before []string `yaml:"before"`
}

// YamlLinterConfiguration describes configuration for Yaml.
type YamlLinterConfiguration struct {
	AllowedQuotes      Quotes       `yaml:"allowedQuotes"`
	AllowedLists       AllowedLists `yaml:"allowedLists"`
	ContinuationIndent Check        `yaml:"continuationIndent"`
}

// Quotes describes types of quotes.
type Quotes struct {
	Single   bool `yaml:"single"`
	Double   bool `yaml:"double"`
	Noquotes bool `yaml:"noquotes"`
}

// AllowedLists describes which types of lists are correct.
type AllowedLists struct {
	Inline bool `yaml:"inline"`
	Dash   bool `yaml:"dash"`
}

// CheckLogicalName checks name.
func (this LinterConfiguration) CheckLogicalName(name string) bool {
	return regexp.MustCompile(this.Global.NamingConventions.LogicalNames).MatchString(name)
}

// GetLinterConfiguration gets configuration from file.
func GetLinterConfiguration(ctx *context.Context) (err error, lintConf LinterConfiguration) {

	linterConfigurationFilename := ctx.CliArguments.LinterConfiguration
	rawLintConfiguration := configuration.GetLinterConfigurationFile(linterConfigurationFilename, ctx.Logger)

	if path.Ext(*linterConfigurationFilename) == ".json" {
		err = json.Unmarshal([]byte(rawLintConfiguration), &lintConf)
		if err != nil {
			ctx.Logger.Error(err.Error())
		}
	} else {

		err = yaml.Unmarshal([]byte(rawLintConfiguration), &lintConf)
		if err != nil {
			ctx.Logger.Error(err.Error())
		}
	}
	return
}
