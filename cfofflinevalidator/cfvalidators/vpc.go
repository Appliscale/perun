package cfvalidators

import (
	"github.com/asaskevich/govalidator"
	"github.com/mitchellh/mapstructure"
	"github.com/Appliscale/cftool/cfofflinevalidator/cftemplate"
	"github.com/Appliscale/cftool/cflogger"
)

type VpcProperties struct {
	CidrBlock string
	EnableDnsSupport bool
	EnableDnsHostnames bool
	InstanceTenancy string
	Tags []Tag
}

func IsVpcValid(name string, vpc cftemplate.Resource, logger *cflogger.Logger) bool {
	valid := true
	var properties VpcProperties
	mapstructure.Decode(vpc.Properties, &properties)

	if properties.CidrBlock != "" && !govalidator.IsCIDR(properties.CidrBlock) {
		logger.ValidationError(name, "Invalid CIDR format - " + properties.CidrBlock)
		valid = false
	}

	return valid
}