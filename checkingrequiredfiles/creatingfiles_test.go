package checkingrequiredfiles

import (
	"testing"

	"github.com/Appliscale/perun/checkingrequiredfiles/mocks"
	"github.com/Appliscale/perun/stack/stack_mocks"
	"github.com/golang/mock/gomock"
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
