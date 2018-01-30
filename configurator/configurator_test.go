package configurator

import (
	"github.com/Appliscale/perun/logger"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCreateConfiguration(t *testing.T) {
	a := CreateConfiguration()

	assert.Emptyf(t, a, "Path could not be empty")
}

func TestConfigurePath(t *testing.T) {
	logger := logger.CreateDefaultLogger()
	a, b := ConfigurePath(logger)
	assert.Emptyf(t, a, "Path could not be empty")
	assert.Emptyf(t, b, "Filename could not be empty")
}

func TestMakeUserPath(t *testing.T) {
	logger := logger.CreateDefaultLogger()
	path := makeUserPath(logger)
	a := "perun"
	assert.Containsf(t, path, a, "Inccorect path")
}

func TestSetProfile(t *testing.T) {
	logger := logger.CreateDefaultLogger()
	profile := setProfile(logger)

	assert.Emptyf(t, profile, "Name could not be empty")

}

func TestCreateConfig(t *testing.T) {
	config := createConfig()

	assert.Emptyf(t, config.DefaultProfile, "Default profile could not be empty")
	assert.Emptyf(t, config.DefaultRegion, "Default Region could not be empty")
	assert.Emptyf(t, config.SpecificationURL, "SpecificationURL could not be empty")

}

func TestMakeArrayRegions(t *testing.T) {
	region := makeArrayRegions()

	for i := 0; i < len(region); i++ {
		assert.Containsf(t, region[i], resourceSpecificationURL[region[i]], "Incorrect region adnd URL")
	}
}
