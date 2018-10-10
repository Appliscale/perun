package linter

import (
	"github.com/Appliscale/perun/checkingrequiredfiles/mocks"
	"github.com/Appliscale/perun/context"
	"github.com/Appliscale/perun/stack/stack_mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"strconv"
	"strings"
	"testing"
)

func setupTestEnv(t *testing.T, filename string) (*context.Context, *gomock.Controller, *mocks.MockLoggerInt, LinterConfiguration) {
	ctx := stack_mocks.SetupContext(t, []string{"cmd", "lint", filename, "--lint-configuration=test_resources/test_style.yaml"})
	mockCtrl := gomock.NewController(t)
	mockLogger := mocks.NewMockLoggerInt(mockCtrl)
	ctx.Logger = mockLogger
	err, linterConf := GetLinterConfiguration(ctx)
	assert.Nil(t, err)
	return ctx, mockCtrl, mockLogger, linterConf
}

func TestCheckLineLengths(t *testing.T) {
	ctx, mockCtrl, mockLogger, linterConf := setupTestEnv(t, "./test_resources/blanklines_testtemplate.yaml")
	defer mockCtrl.Finish()

	mockLogger.EXPECT().Warning("line " + strconv.Itoa(1) + ": maximum line lenght exceeded").Times(1)

	checkLineLengths([]string{"asdasdasdasdasd"}, linterConf, ctx)
}

func TestCheckBlankLines(t *testing.T) {
	ctx, mockCtrl, mockLogger, linterConf := setupTestEnv(t, "./test_resources/blanklines_testtemplate.yaml")
	defer mockCtrl.Finish()

	mockLogger.EXPECT().Warning("Blank lines are not allowed in current lint configuration").Times(1)

	checkBlankLines(linterConf, stack_mocks.ReadFile(t, "./test_resources/blanklines_testtemplate.yaml"), ctx)
	checkBlankLines(linterConf, stack_mocks.ReadFile(t, "./test_resources/noblanklines_testtemplate.yaml"), ctx)
}

func TestCheckAWSSpecificStuff(t *testing.T) {
	ctx, mockCtrl, mockLogger, linterConf := setupTestEnv(t, "./test_resources/nodescription_testtemplate.yaml")
	defer mockCtrl.Finish()

	mockLogger.EXPECT().Warning("The template has no description").Times(1)
	mockLogger.EXPECT().Warning("No description provided for parameter TestParameter2")
	mockLogger.EXPECT().Warning("Resource 'S3' does not meet the given logical Name regex: test.+")

	checkAWSCFSpecificStuff(ctx, stack_mocks.ReadFile(t, "./test_resources/nodescription_testtemplate.yaml"), linterConf)
}

func TestTestCheckAWSSpecificStuffOk(t *testing.T) {
	ctx, mockCtrl, _, linterConf := setupTestEnv(t, "./test_resources/described_testtemplate.yaml")
	defer mockCtrl.Finish()

	checkAWSCFSpecificStuff(ctx, stack_mocks.ReadFile(t, "./test_resources/described_testtemplate.yaml"), linterConf)
}

func TestCheckJSONSpaces(t *testing.T) {
	ctx, mockCtrl, mockLogger, linterConf := setupTestEnv(t, "./test_resources/spacesjson_testtemplate.json")
	defer mockCtrl.Finish()

	mockLogger.EXPECT().Warning("line 3: no space after ':'")
	mockLogger.EXPECT().Warning("line 2: no space before ':'")

	checkJsonSpaces(ctx, linterConf, strings.Split(stack_mocks.ReadFile(t, "./test_resources/spacesjson_testtemplate.json"), "\n"))
}
