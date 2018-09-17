package mocks

import (
	"github.com/Appliscale/perun/cliparser"
	"github.com/Appliscale/perun/configuration"
	"github.com/Appliscale/perun/configurator"
	"github.com/Appliscale/perun/context"
	"github.com/Appliscale/perun/logger"
	"io/ioutil"
	"testing"
)

func SetupContext(t *testing.T, args []string) *context.Context {
	myLogger := logger.CreateDefaultLogger()
	myLogger.SetVerbosity("ERROR")

	cliArguments, err := cliparser.ParseCliArguments(args)
	if err != nil {
		t.Error(err.Error())
		return &context.Context{}
	}

	config := createDefaultConfiguration()
	iconsistenciesConfig := configuration.ReadInconsistencyConfiguration(&myLogger)

	ctx := context.Context{
		CliArguments:        cliArguments,
		Logger:              &myLogger,
		Config:              config,
		InconsistencyConfig: iconsistenciesConfig,
	}

	return &ctx
}

func ReadFile(t *testing.T, filePath string) string {
	rawTemplate, readFileError := ioutil.ReadFile(filePath)
	if readFileError != nil {
		t.Error(readFileError.Error())
	}
	template := string(rawTemplate)
	return template
}

func createDefaultConfiguration() configuration.Configuration {
	return configuration.Configuration{
		DefaultProfile:        "default",
		DefaultRegion:         "region",
		SpecificationURL:      configurator.ResourceSpecificationURL,
		DefaultDecisionForMFA: false,
		DefaultDurationForMFA: 3600,
		DefaultVerbosity:      "INFO"}
}
