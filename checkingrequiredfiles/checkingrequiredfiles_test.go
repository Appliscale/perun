package checkingrequiredfiles

import (
	"github.com/Appliscale/perun/logger"
	"github.com/stretchr/testify/assert"
	"testing"
)

var sink logger.Logger

func TestGetProfilesFromFile(t *testing.T) {
	profiles := getProfilesFromFile("test_resources/.aws/config", &sink)
	assert.NotNilf(t, profiles, "Profiles are nil")
}

func TestIsProfileInCredentials(t *testing.T) {
	answer := isProfileInCredentials("default", "test_resources/.aws/credentials", &sink)
	assert.Truef(t, answer, "This profile isn't in credentials")
}

func TestFindRegionForProfile(t *testing.T) {
	region := findRegionForProfile("default", "test_resources/.aws/config", &sink)
	assert.NotNilf(t, region, "Region is nil")
}

func TestFindNewProfileInCredentials(t *testing.T) {
	credentials := []string{"default", "test"}
	config := []string{"default"}
	profiles := findNewProfileInCredentials(credentials, config)
	assert.NotNilf(t, profiles, "Profiles are nil")
}
