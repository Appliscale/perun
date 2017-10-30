package cfcontext

import (
	"github.com/Appliscale/cftool/cfcliparser"
	"github.com/Appliscale/cftool/cflogger"
	"github.com/Appliscale/cftool/cfconfiguration"
)

type Context struct {
	CliArguments cfcliparser.CliArguments
	Logger* cflogger.Logger
	Config cfconfiguration.Configuration
}

func GetContext() (context Context, err error) {
	logger := cflogger.Logger{}
	defer logger.PrintErrors()

	cliArguments, err := cfcliparser.ParseCliArguments()
	if err != nil {
		return
	}

	config, err := cfconfiguration.GetConfiguration(cliArguments)
	if err != nil {
		return
	}

	context = Context{
		CliArguments: cliArguments,
		Logger: &logger,
		Config: config,
	}

	return
}