// Copyright 2017 Appliscale
//
// Maintainers and Contributors:
//
//   - Piotr Figwer (piotr.figwer@appliscale.io)
//   - Wojciech Gawroński (wojciech.gawronski@appliscale.io)
//   - Kacper Patro (kacper.patro@appliscale.io)
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

func setup() {
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()

	os.Args = []string{"cmd", "--mode=validate_offline", "--template=some_path", "--config=test_resources/test_config.yaml"}
	cliArgs, err := cliparser.ParseCliArguments()
	if err != nil {
		panic(err)
	}

	logger := logger.CreateDefaultLogger()

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
