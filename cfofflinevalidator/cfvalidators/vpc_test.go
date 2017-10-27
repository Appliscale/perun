package cfvalidators

import (
	"testing"
	"github.com/Appliscale/cftool/cfofflinevalidator/cftemplate"
	"github.com/stretchr/testify/assert"
	"github.com/Appliscale/cftool/cflogger"
	"os"
)

var logger cflogger.Logger

func setup() {
	logger = cflogger.Logger{}
}

func TestMain(m *testing.M) {
	setup()
	retCode := m.Run()
	os.Exit(retCode)
}

func TestValidVpc(t *testing.T) {
	vpc := createVpc("10.0.0.0/16")
	assert.True(t, IsVpcValid("Example", vpc, &logger))
}

func TestInvalidVpc(t *testing.T) {
	vpc := createVpc("10.0.0.0")
	assert.False(t, IsVpcValid("Example", vpc, &logger))
}

func createVpc(cidrBlock string) cftemplate.Resource {
	vpc := cftemplate.Resource{}
	properties := make(map[string]interface{})
	properties["CidrBlock"] = cidrBlock
	vpc.Properties = properties
	return vpc
}