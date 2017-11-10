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

package validators

import (
	"testing"
	"github.com/Appliscale/perun/offlinevalidator/template"
	"github.com/stretchr/testify/assert"
	"github.com/Appliscale/perun/logger"
	"os"
)

var logger logger.Logger

func setup() {
	logger = logger.Logger{}
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

func createVpc(cidrBlock string) template.Resource {
	vpc := template.Resource{}
	properties := make(map[string]interface{})
	properties["CidrBlock"] = cidrBlock
	vpc.Properties = properties
	return vpc
}
