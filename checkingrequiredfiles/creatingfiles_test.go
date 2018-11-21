package checkingrequiredfiles

import (
	"github.com/Appliscale/perun/logger"
	"testing"

	"github.com/Appliscale/perun/checkingrequiredfiles/mocks"
	"github.com/Appliscale/perun/stack/stack_mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestUseProfileFromConfig(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	mockLogger := mocks.NewMockLoggerInt(mockCtrl)
	defer mockCtrl.Finish()
	profile := "Test"
	profilesInConfig := []string{"Test", "Test1"}
	messages := [3]string{
		"Available profiles from config:",
		profilesInConfig[0],
		profilesInConfig[1],
	}

	for _, mes := range messages {
		mockLogger.EXPECT().Always(mes).Times(1)
	}
	mockLogger.EXPECT().GetInput("Which profile should perun use as a default?", &profile).Return(nil).Times(1)

	useProfileFromConfig(profilesInConfig, profile, mockLogger)
}

func TestAddNewProfileFromCredentialsToConfig(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	mockLogger := mocks.NewMockLoggerInt(mockCtrl)
	defer mockCtrl.Finish()
	profile := "perun"
	homePath := "./test_resources"
	templatePath := "../stack/test_resources/test_template.yaml"
	ctx := stack_mocks.SetupContext(t, []string{"cmd", "create-stack", "teststack", templatePath})

	mockLogger.EXPECT().GetInput("I found profile "+profile+" in credentials, but not in config. \nCreate new profile in config? Y/N", gomock.Any()).Return(nil).Times(1)

	addNewProfileFromCredentialsToConfig("default", homePath, ctx, mockLogger)
}

func TestAddProfileToCredentials(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	mockLogger := mocks.NewMockLoggerInt(mockCtrl)
	defer mockCtrl.Finish()
	profile := "perun"
	homePath := "./test_resources"
	templatePath := "../stack/test_resources/test_template.yaml"
	ctx := stack_mocks.SetupContext(t, []string{"cmd", "create-stack", "teststack", templatePath})

	mockLogger.EXPECT().Always("Profile " + profile + " has already credentials").Times(1)

	addProfileToCredentials(profile, homePath, ctx, mockLogger)
}

func TestConfigIsPresent(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	mockLogger := mocks.NewMockLoggerInt(mockCtrl)
	defer mockCtrl.Finish()
	profile := "default"
	homePath := "./test_resources"
	templatePath := "../stack/test_resources/test_template.yaml"
	ctx := stack_mocks.SetupContext(t, []string{"cmd", "create-stack", "teststack", templatePath})

	mockLogger.EXPECT().GetInput("Default profile exists, do you want to use it *Y* or create your own *N*?", gomock.Any()).Return(nil).Times(1)

	configIsPresent(profile, homePath, ctx, mockLogger)
}

func TestCreateCredentials(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	mockLogger := mocks.NewMockLoggerInt(mockCtrl)
	defer mockCtrl.Finish()
	profile := "test1"
	homePath := "./test_resources"
	templatePath := "../stack/test_resources/test_template.yaml"
	ctx := stack_mocks.SetupContext(t, []string{"cmd", "create-stack", "teststack", templatePath})

	mockLogger.EXPECT().GetInput("I found profile "+profile+" in .aws/config without credentials, add? Y/N", gomock.Any())

	createCredentials(profile, homePath, ctx, mockLogger)

}

func TestGetIamInstanceProfileAssociations(t *testing.T) {
	sink := logger.CreateDefaultLogger()
	output, _ := getIamInstanceProfileAssociations(sink, "us-east-1")

	assert.Emptyf(t, output, "Should be empty")
}

func TestGetRegion(t *testing.T) {
	region, _, _ := getRegion()
	assert.Equalf(t, region, "", "Should be nil")
}

func TestWorkingOnEC2(t *testing.T) {
	sink := logger.CreateDefaultLogger()
	_, _, err := workingOnEC2(sink)
	assert.NotNilf(t, err, "Should be non-nil")
}

func TestCreateEC2context(t *testing.T) {
	templatePath := "../stack/test_resources/test_template.yaml"
	ctx := stack_mocks.SetupContext(t, []string{"cmd", "validate", templatePath})
	mockCtrl := gomock.NewController(t)
	mockLogger := mocks.NewMockLoggerInt(mockCtrl)
	defer mockCtrl.Finish()
	a := ""
	mockLogger.EXPECT().GetInput("Directory for temporary files", gomock.Any())
	mockLogger.EXPECT().Always("Your temporary files directory is: " + a).Times(1)
	mockLogger.EXPECT().Error("stat ./test_resources/.config/perun: no such file or directory").Times(1)
	createEC2context("test", "./test_resources", "test", ctx, mockLogger)
}
