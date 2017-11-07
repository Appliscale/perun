package cfconfiguration

import (
	"testing"
	"os"
	"github.com/Appliscale/cftool/cfcliparser"
	"github.com/Appliscale/cftool/cflogger"
	"github.com/stretchr/testify/assert"
)

var configuration Configuration

func setup() {
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()

	os.Args = []string{"cmd", "--mode=validate_offline", "--template=some_path", "--config=test_resources/test_config.yaml"}
	cliArgs, err := cfcliparser.ParseCliArguments()
	if err != nil {
		panic(err)
	}

	logger := cflogger.CreateDefaultLogger()

	configuration, err = GetConfiguration(cliArgs, &logger)

	if err != nil {
		panic(err)
	}
}

func TestMain(m *testing.M) {
	setup()
	retCode := m.Run()
	os.Exit(retCode)
}

func TestSpecificationFileURL(t *testing.T) {
	url, _ := configuration.GetSpecificationFileURLForCurrentRegion()
	assert.Equal(t, "https://d1uauaxba7bl26.cloudfront.net/latest/gzip/CloudFormationResourceSpecification.json", url)
}

func TestNoSpecificationForRegion(t *testing.T) {
	localConfiguration := Configuration{
		Region: "someRegion",
	}
	_, err := localConfiguration.GetSpecificationFileURLForCurrentRegion()
	assert.NotNil(t, err)
}