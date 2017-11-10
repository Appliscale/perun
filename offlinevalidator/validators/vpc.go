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
	"github.com/asaskevich/govalidator"
	"github.com/mitchellh/mapstructure"
	"github.com/Appliscale/perun/offlinevalidator/template"
	"github.com/Appliscale/perun/logger"
)

type VpcProperties struct {
	CidrBlock string
	EnableDnsSupport bool
	EnableDnsHostnames bool
	InstanceTenancy string
	Tags []Tag
}

func IsVpcValid(name string, vpc template.Resource, logger *logger.Logger) bool {
	valid := true
	var properties VpcProperties
	mapstructure.Decode(vpc.Properties, &properties)

	if properties.CidrBlock != "" && !govalidator.IsCIDR(properties.CidrBlock) {
		logger.ValidationError(name, "Invalid CIDR format - " + properties.CidrBlock)
		valid = false
	}

	return valid
}
