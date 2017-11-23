// Copyright 2017 Appliscale
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

package configuration

import (
	"testing"
	"os"
	"github.com/Appliscale/perun/cliparser"
	"github.com/Appliscale/perun/logger"
	"github.com/stretchr/testify/assert"
)

var configuration Configuration

func setup(osArgs []string) error {
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()

	os.Args = osArgs
	cliArgs, err := cliparser.ParseCliArguments()
	if err != nil {
		return err
	}

	logger := logger.CreateQuietLogger()

	configuration, err = GetConfiguration(cliArgs, &logger)
	if err != nil {
		return err
	}

	return err
}

func TestSpecificationFileURL(t *testing.T) {
	setup([]string{"cmd", "--mode=validate_offline", "--template=some_path", "--config=test_resources/test_config.yaml"})
	url, _ := configuration.GetSpecificationFileURLForCurrentRegion()
	assert.Equal(t, "https://d1uauaxba7bl26.cloudfront.net/latest/gzip/CloudFormationResourceSpecification.json", url)
}

func TestNoSpecificationForRegion(t *testing.T) {
	setup([]string{"cmd", "--mode=validate_offline", "--template=some_path", "--config=test_resources/test_config.yaml"})

	localConfiguration := Configuration{
		DefaultRegion: "someRegion",
	}

	_, err := localConfiguration.GetSpecificationFileURLForCurrentRegion()
	assert.NotNil(t, err)
}

func TestGettingMFADecisionFromConfigurationFile(t *testing.T) {
	setup([]string{"cmd", "--mode=validate_offline", "--template=some_path", "--config=test_resources/test_config.yaml"})
	assert.Equal(t, false, configuration.DefaultDecisionForMFA)
}

func TestOverrideForMFADecision(t *testing.T) {
	setup([]string{"cmd", "--mode=validate_offline", "--template=some_path", "--config=test_resources/test_config.yaml", "--mfa"})
	assert.Equal(t, true, configuration.DefaultDecisionForMFA)
}

func TestNoMFADecision(t *testing.T) {
	setup([]string{"cmd", "--mode=validate_offline", "--template=some_path", "--sandbox"})
	assert.Equal(t, false, configuration.DefaultDecisionForMFA)
}

func TestGettingDefaultRegionFromConfigurationFile(t *testing.T) {
	setup([]string{"cmd", "--mode=validate_offline", "--template=some_path", "--config=test_resources/test_config.yaml"})
	assert.Equal(t, "us-west-2", configuration.DefaultRegion)
}

func TestCLIOverrideForRegion(t *testing.T) {
	setup([]string{"cmd", "--mode=validate_offline", "--template=some_path", "--config=test_resources/test_config.yaml", "--region=ap-southeast-1"})
	assert.Equal(t, "ap-southeast-1", configuration.DefaultRegion)
}

func TestGettingDurationForMFAFromConfigurationFile(t *testing.T) {
	setup([]string{"cmd", "--mode=validate_offline", "--template=some_path", "--config=test_resources/test_config.yaml"})
	assert.Equal(t, int64(2600), configuration.DefaultDurationForMFA)
}

func TestCLIOverrideForDurationForMFA(t *testing.T) {
	setup([]string{"cmd", "--mode=validate_offline", "--template=some_path", "--config=test_resources/test_config.yaml", "--duration=1600"})
	assert.Equal(t, int64(1600), configuration.DefaultDurationForMFA)
}

func TestNoDurationForMFA(t *testing.T) {
	setup([]string{"cmd", "--mode=validate_offline", "--template=some_path", "--sandbox"})
	assert.Equal(t, int64(0), configuration.DefaultDurationForMFA)
}

func TestTooBigDurationForMFA(t *testing.T) {
	err := setup([]string{"cmd", "--mode=validate_offline", "--template=some_path", "--duration=600000000"})
	assert.NotNil(t, err)
}

func TestTooSmallDurationForMFA(t *testing.T) {
	err := setup([]string{"cmd", "--mode=validate_offline", "--template=some_path", "--duration=-1"})
	assert.NotNil(t, err)
}

func TestZeroDurationForMFA(t *testing.T) {
	setup([]string{"cmd", "--mode=validate_offline", "--template=some_path", "--duration=0", "--sandbox"})
	assert.Equal(t, int64(0), configuration.DefaultDurationForMFA)
}

func TestGettingProfileFromConfigurationFile(t *testing.T) {
	setup([]string{"cmd", "--mode=validate_offline", "--template=some_path", "--config=test_resources/test_config.yaml"})
	assert.Equal(t, "profile", configuration.DefaultProfile)
}

func TestCLIOverrideForProfile(t *testing.T) {
	setup([]string{"cmd", "--mode=validate_offline", "--template=some_path", "--config=test_resources/test_config.yaml", "--profile=cliProfile"})
	assert.Equal(t, "cliProfile", configuration.DefaultProfile)
}

func TestGettingVerbosityFromConfigurationFile(t *testing.T) {
	setup([]string{"cmd", "--mode=validate_offline", "--template=some_path", "--config=test_resources/test_config.yaml"})
	assert.Equal(t, "ERROR", configuration.DefaultVerbosity)
}

func TestCLIOverrideForVerbosity(t *testing.T) {
	setup([]string{"cmd", "--mode=validate_offline", "--template=some_path", "--config=test_resources/test_config.yaml", "--verbosity=INFO"})
	assert.Equal(t, "INFO", configuration.DefaultVerbosity)
}
