package checkingrequiredfiles

import (
	"github.com/Appliscale/perun/checkingrequiredfiles/mocks"
	"github.com/Appliscale/perun/cliparser"
	"github.com/Appliscale/perun/configuration"
	"github.com/Appliscale/perun/context"
	"github.com/golang/mock/gomock"
	"testing"
)

func TestUseProfileFromConfig(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	mockLogger := mocks.NewMockLoggerInt(mockCtrl)
	defer mockCtrl.Finish()

	profile := "Test"
	profilesInConfig := []string{"Test", "Test1"}
	contx, _ := context.GetContext(cliparser.ParseCliArguments, configuration.GetConfiguration, configuration.ReadInconsistencyConfiguration)

	mockLogger.EXPECT().GetInput(gomock.Any(),gomock.Any()).Return("Test").Times(1)
	useProfileFromConfig(profilesInConfig, profile, contx.Logger)
}




func setup(t *testing.T) (mockCtrl *gomock.Controller, mockCreatingFiles *mocks.MockCreatingFiles, profile string, homePath string, contx context.Context) {
	mockCtrl = gomock.NewController(t)
	mockCreatingFiles = mocks.NewMockCreatingFiles(mockCtrl)
	profile = "test"
	homePath = "test_resources"
	contx, _ = context.GetContext(cliparser.ParseCliArguments, configuration.GetConfiguration, configuration.ReadInconsistencyConfiguration)
	return
}
func TestCreateNewMainYaml(t *testing.T) {
	mockCtrl, mockCreatingFiles, profile, homePath, contx := setup(t)
	defer mockCtrl.Finish()
	mockCreatingFiles.EXPECT().CreateNewMainYaml(profile, homePath, &contx, contx.Logger).Return(contx).Times(1)
	createNewMainYaml(profile, homePath, &contx, contx.Logger)
}

/*
func TestUseProfileFromConfig(t *testing.T) {
	mockCtrl, mockCreatingFiles, profile, _, contx := setup(t)
	defer mockCtrl.Finish()
	profilesInConfig := []string{"Test", "Test1"}
	mockCreatingFiles.EXPECT().UseProfileFromConfig(profilesInConfig, profile, contx.Logger)
	mockCreatingFiles.UseProfileFromConfig(profilesInConfig, profile, contx.Logger)
}
*/

func TestAddNewProfileFromCredentialsToConfig(t *testing.T) {
	mockCtrl, mockCreatingFiles, profile, homePath, contx := setup(t)
	defer mockCtrl.Finish()
	mockCreatingFiles.EXPECT().AddNewProfileFromCredentialsToConfig(profile, homePath, contx, contx.Logger)
	mockCreatingFiles.AddNewProfileFromCredentialsToConfig(profile, homePath, &contx, contx.Logger)
}

func TestAddProfileToCredentials(t *testing.T) {
	mockCtrl, mockCreatingFiles, profile, homePath, contx := setup(t)
	defer mockCtrl.Finish()
	mockCreatingFiles.EXPECT().AddProfileToCredentials(profile, homePath, contx, contx.Logger)
	mockCreatingFiles.AddProfileToCredentials(profile, homePath, &contx, contx.Logger)
}

func TestConfigIsPresent(t *testing.T) {
	mockCtrl, mockCreatingFiles, profile, homePath, contx := setup(t)
	defer mockCtrl.Finish()
	mockCreatingFiles.EXPECT().ConfigIsPresent(profile, homePath, contx, contx.Logger)
	mockCreatingFiles.ConfigIsPresent(profile, homePath, &contx, contx.Logger)
}

func TestNewConfig(t *testing.T) {
	mockCtrl, mockCreatingFiles, profile, homePath, contx := setup(t)
	defer mockCtrl.Finish()
	region := "us-east-1"
	mockCreatingFiles.EXPECT().NewConfigFile(profile, region, homePath, contx, contx.Logger)
	mockCreatingFiles.NewConfigFile(profile, region, homePath, &contx, contx.Logger)
}

func TestCreateCredentials(t *testing.T) {
	mockCtrl, mockCreatingFiles, profile, homePath, contx := setup(t)
	defer mockCtrl.Finish()
	mockCreatingFiles.EXPECT().CreateCredentials(profile, homePath, contx, contx.Logger)
	mockCreatingFiles.CreateCredentials(profile, homePath, &contx, contx.Logger)
}
