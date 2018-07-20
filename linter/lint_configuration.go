package linter

import (
	"encoding/json"
	"github.com/Appliscale/perun/configuration"
	"github.com/Appliscale/perun/context"
	"github.com/ghodss/yaml"
	"path"
)

type LinterConfiguration struct {
	Yaml   YamlLinterConfiguration   `yaml:"yaml"`
	Json   JsonLinterConfiguration   `yaml:"json"`
	Global GlobalLinterConfiguration `yaml:"global"`
}

type GlobalLinterConfiguration struct {
	LineLength        Check             `yaml:"lineLength"`
	Indent            Check             `yaml:"indent"`
	RequiredFields    RequiredFields    `yaml:"requiredFields"`
	NamingConventions NamingConventions `yaml:"namingConventions"`
	BlankLinesAllowed bool              `yaml:"blankLinesAllowed"`
}
type Check struct {
	Required bool        `yaml:"required"`
	Value    interface{} `yaml:"value"`
}
type NamingConventions struct {
	LogicalNames string `yaml:"logicalNames"`
}
type RequiredFields struct {
	TemplateDescription   bool `yaml:"templateDescription"`
	ParametersDescription bool `yaml:"parametersDescription"`
}

type JsonLinterConfiguration struct {
	Spaces SpacesConfiguration `yaml:"spaces"`
}
type SpacesConfiguration struct {
	After  []string `yaml:"after"`
	Before []string `yaml:"before"`
}

type YamlLinterConfiguration struct {
	AllowedQuotes      Quotes       `yaml:"allowedQuotes"`
	AllowedLists       AllowedLists `yaml:"allowedLists"`
	ContinuationIndent Check        `yaml:"continuationIndent"`
}
type Quotes struct {
	Single   bool `yaml:"single"`
	Double   bool `yaml:"double"`
	Noquotes bool `yaml:"noquotes"`
}
type AllowedLists struct {
	Inline bool `yaml:"inline"`
	Dash   bool `yaml:"dash"`
}

//func (this YamlLinterConfiguration) CheckLogicalName(name string) bool {
//	return regexp.MustCompile(this.logicalNameRegex).MatchString(name)
//}
//func (this YamlLinterConfiguration) CheckExternalName(name string) bool {
//	return regexp.MustCompile(this.externalNameRegex).MatchString(name)
//}

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
