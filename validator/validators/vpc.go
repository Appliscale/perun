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

package validators

import (
	"github.com/Appliscale/perun/logger"
	"github.com/Appliscale/perun/validator/template"
	"github.com/asaskevich/govalidator"
	"github.com/mitchellh/mapstructure"
)

// VpcProperties describes structure of Vpc.
type VpcProperties struct {
	CidrBlock          string
	EnableDnsSupport   bool
	EnableDnsHostnames bool
	InstanceTenancy    string
	Tags               []Tag
}

// IsVpcValid : Checks if CIDR block is valid.
func IsVpcValid(vpc template.Resource, resourceValidation *logger.ResourceValidation) bool {
	valid := true
	var properties VpcProperties
	mapstructure.Decode(vpc.Properties, &properties)

	if properties.CidrBlock != "" && !govalidator.IsCIDR(properties.CidrBlock) {
		resourceValidation.AddValidationError("Invalid CIDR format - " + properties.CidrBlock)
		valid = false
	}

	return valid
}
