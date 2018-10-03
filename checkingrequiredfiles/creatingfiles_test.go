package checkingrequiredfiles

import (
	"testing"

	"github.com/Appliscale/perun/checkingrequiredfiles/mocks"
	mockContext "github.com/Appliscale/perun/stack/mocks"
	"github.com/golang/mock/gomock"
)

func TestUseProfileFromConfig(t *testing.T) { //It's ok
	mockCtrl := gomock.NewController(t)
	mockLogger := mocks.NewMockLoggerInt(mockCtrl)
	defer mockCtrl.Finish()
	profilesInConfig := []string{"Test", "Test1"}
	messages := [3]string{
		"Available profiles from config:",
		profilesInConfig[0],
		profilesInConfig[1],
	}
	profile := "Test"
	for _, mes := range messages {
		mockLogger.EXPECT().Always(mes).Times(1)
	}
	mockLogger.EXPECT().GetInput("Which profile should perun use as a default?", &profile).Return(nil).Times(1)
	useProfileFromConfig(profilesInConfig, profile, mockLogger)
}

func TestAddNewProfileFromCredentialsToConfig(t *testing.T) { // "I found profile perun..." want "Region, Output" x3
	mockCtrl := gomock.NewController(t)
	mockLogger := mocks.NewMockLoggerInt(mockCtrl)
	defer mockCtrl.Finish()

	profile := "perun"
	region := "us-east-1"
	output := "json"
	homePath := "./test_resources"
	templatePath := "../stack/test_resources/test_template.yaml"
	ctx := mockContext.SetupContext(t, []string{"cmd", "create-stack", "teststack", templatePath})

	ans := "Y"

	data := [3]string{
		ans,
		region,
		output,
	}
	messages := [3]string{
		"I found profile " + profile + " in credentials, but not in config. \nCreate new profile in config? Y/N",
		"Region",
		"Output",
	}
	for i, mes := range messages {
		mockLogger.EXPECT().GetInput(mes, &data[i]).Return(nil).Times(1)
	}

	addNewProfileFromCredentialsToConfig("default", homePath, ctx, mockLogger)
}

func TestAddProfileToCredentials(t *testing.T) { //missing calls
	mockCtrl := gomock.NewController(t)
	mockLogger := mocks.NewMockLoggerInt(mockCtrl)
	defer mockCtrl.Finish()
	profile := "test1"
	homePath := "./test_resources"
	templatePath := "../stack/test_resources/test_template.yaml"
	awsAccessKeyID := "TESTTEST"
	awsSecretAccessKey := "testtest"
	mfaSerial := "arn:aws:iam:"

	messages := [3]string{
		"awsAccessKeyID",
		"awsSecretAccessKey",
		"mfaSerial",
	}
	data := [3]string{
		awsAccessKeyID,
		awsSecretAccessKey,
		mfaSerial,
	}
	ctx := mockContext.SetupContext(t, []string{"cmd", "create-stack", "teststack", templatePath})
	mockLogger.EXPECT().Always("You haven't got .aws/credentials file for profile " + profile).Times(1)

	for i, mes := range messages {
		mockLogger.EXPECT().GetInput(mes, &data[i]).Return(nil).Times(1)
	}

	addProfileToCredentials(profile, homePath, ctx, mockLogger)
}

func TestConfigIsPresent(t *testing.T) { // Available profiles from config, default profiles doesn't exist ( but it does)
	mockCtrl := gomock.NewController(t)
	mockLogger := mocks.NewMockLoggerInt(mockCtrl)
	defer mockCtrl.Finish()
	profile := "prof"
	homePath := "./test_resources"
	templatePath := "../stack/test_resources/test_template.yaml"
	answer := "Y"
	ctx := mockContext.SetupContext(t, []string{"cmd", "create-stack", "teststack", templatePath})
	mockLogger.EXPECT().GetInput("Default profile exists, do you want to use it *Y* or create your own *N*?", &answer).Return(nil).Times(1)
	configIsPresent(profile, homePath, ctx, mockLogger)
}

func TestCreateCredentials(t *testing.T) { //unexpected call, arguments doesn't match
	mockCtrl := gomock.NewController(t)
	mockLogger := mocks.NewMockLoggerInt(mockCtrl)
	defer mockCtrl.Finish()
	profile := "test1"
	homePath := "./test_resources"
	templatePath := "../stack/test_resources/test_template.yaml"
	answer := "Y"
	ctx := mockContext.SetupContext(t, []string{"cmd", "create-stack", "teststack", templatePath})
	mockLogger.EXPECT().GetInput("I found profile "+profile+" in .aws/config without credentials, add? Y/N", &answer)
	createCredentials(profile, homePath, ctx, mockLogger)

}
