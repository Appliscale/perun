package stack

import (
	"github.com/Appliscale/perun/cliparser"
	"github.com/Appliscale/perun/configuration"
	"github.com/Appliscale/perun/context"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestcreateSession(t *testing.T) {
	context, _ := context.GetContext(cliparser.ParseCliArguments, configuration.GetConfiguration)
	profile := context.Config.DefaultProfile
	region := context.Config.DefaultRegion
	session, _ := createSession(&context, profile, &region)

	assert.NotEmptyf(t, session, "Incorrect session")
}

func TestdeleteStackInput(t *testing.T) {
	context, _ := context.GetContext(cliparser.ParseCliArguments, configuration.GetConfiguration)
	templateStruct := deleteStackInput(&context)

	assert.NotEmptyf(t, templateStruct.StackName, "StackName could not be empty")
}
