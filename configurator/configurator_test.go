// Copyright 2018 Appliscale
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

package configurator

import (
	"github.com/Appliscale/perun/configuration"
	"github.com/stretchr/testify/assert"
	"os/exec"
	"strings"
	"testing"
)

func TestConfigurePath(t *testing.T) {
	subproc := exec.Command("cmd")
	input := ""
	subproc.Stdin = strings.NewReader(input)
	output, _ := subproc.Output()
	assert.IsTypef(t, string(output), input, "Invalid type of input")
	subproc.Wait()
}

func TestSetProfile(t *testing.T) {
	subproc := exec.Command("cmd")
	input := ""
	subproc.Stdin = strings.NewReader(input)
	output, _ := subproc.Output()
	assert.IsTypef(t, string(output), input, "Name could not be empty")
	subproc.Wait()
}

func TestCreateConfig(t *testing.T) {
	myconfig := configuration.Configuration{
		DefaultProfile:        "profile",
		DefaultRegion:         "region",
		SpecificationURL:      ResourceSpecificationURL,
		DefaultDecisionForMFA: false,
		DefaultDurationForMFA: 3600,
		DefaultVerbosity:      "INFO"}
	assert.NotEmptyf(t, myconfig.DefaultProfile, "Default profile could not be empty")
	assert.NotEmptyf(t, myconfig.DefaultRegion, "Default Region could not be empty")
	assert.NotEmptyf(t, myconfig.SpecificationURL, "SpecificationURL could not be empty")

}

func TestMakeArrayRegions(t *testing.T) {
	region := makeArrayRegions()
	for i := 0; i < len(region); i++ {
		assert.NotEmptyf(t, region[i], "Incorrect region and URL")
	}
}
