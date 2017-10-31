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
	logger := cflogger.CreateDefaultLogger()

	cliArguments, err := cfcliparser.ParseCliArguments()
	if err != nil {
		logger.Error(err.Error())
		return
	}

	logger.Quiet = *cliArguments.Quiet
	logger.Yes = *cliArguments.Yes
	logger.SetVerbosity(*cliArguments.Verbosity)

	config, err := cfconfiguration.GetConfiguration(cliArguments, &logger)
	if err != nil {
		logger.Error(err.Error())
		return
	}

	context = Context{
		CliArguments: cliArguments,
		Logger: &logger,
		Config: config,
	}

	return
}