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
		SpecificationURL:      resourceSpecificationURL,
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
