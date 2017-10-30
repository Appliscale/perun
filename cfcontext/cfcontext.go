package cfcontext

import (
	"github.com/Appliscale/cftool/cfcliparser"
	"github.com/Appliscale/cftool/cflogger"
)

type Context struct {
	CliArguments cfcliparser.CliArguments
	Logger* cflogger.Logger
}