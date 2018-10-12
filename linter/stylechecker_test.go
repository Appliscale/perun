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

func setupTestEnv(t *testing.T, filename string, styleConfiguration string) (*context.Context, *gomock.Controller, *mocks.MockLoggerInt, LinterConfiguration) {
	ctx := stack_mocks.SetupContext(t, []string{"cmd", "lint", filename, "--lint-configuration=" + styleConfiguration})
	mockCtrl := gomock.NewController(t)
	mockLogger := mocks.NewMockLoggerInt(mockCtrl)
	ctx.Logger = mockLogger
	err, linterConf := GetLinterConfiguration(ctx)
	assert.Nil(t, err)
	return ctx, mockCtrl, mockLogger, linterConf
}

func TestCheckLineLengths(t *testing.T) {
	ctx, mockCtrl, mockLogger, linterConf := setupTestEnv(t, "./test_resources/blanklines_testtemplate.yaml", "test_resources/test_style.yaml")
	defer mockCtrl.Finish()

	mockLogger.EXPECT().Warning("line " + strconv.Itoa(1) + ": maximum line lenght exceeded").Times(1)

	checkLineLengths([]string{"asdasdasdasdasd"}, linterConf, ctx)
}

func TestCheckBlankLines(t *testing.T) {
	ctx, mockCtrl, mockLogger, linterConf := setupTestEnv(t, "./test_resources/blanklines_testtemplate.yaml", "test_resources/test_style.yaml")
	defer mockCtrl.Finish()

	mockLogger.EXPECT().Warning("Blank lines are not allowed in current lint configuration").Times(1)

	checkBlankLines(linterConf, stack_mocks.ReadFile(t, "./test_resources/blanklines_testtemplate.yaml"), ctx)
	checkBlankLines(linterConf, stack_mocks.ReadFile(t, "./test_resources/noblanklines_testtemplate.yaml"), ctx)
}

func TestCheckAWSSpecificStuff(t *testing.T) {
	ctx, mockCtrl, mockLogger, linterConf := setupTestEnv(t, "./test_resources/nodescription_testtemplate.yaml", "test_resources/test_style.yaml")
	defer mockCtrl.Finish()

	mockLogger.EXPECT().Warning("The template has no description").Times(1)
	mockLogger.EXPECT().Warning("No description provided for parameter TestParameter2")
	mockLogger.EXPECT().Warning("Resource 'S3' does not meet the given logical Name regex: Test.+")

	checkAWSCFSpecificStuff(ctx, stack_mocks.ReadFile(t, "./test_resources/nodescription_testtemplate.yaml"), linterConf)
}

func TestTestCheckAWSSpecificStuffOk(t *testing.T) {
	ctx, mockCtrl, _, linterConf := setupTestEnv(t, "./test_resources/described_testtemplate.yaml", "test_resources/test_style.yaml")
	defer mockCtrl.Finish()

	checkAWSCFSpecificStuff(ctx, stack_mocks.ReadFile(t, "./test_resources/described_testtemplate.yaml"), linterConf)
}

func TestCheckJSONSpaces(t *testing.T) {
	ctx, mockCtrl, mockLogger, linterConf := setupTestEnv(t, "./test_resources/spacesjson_testtemplate.json", "test_resources/test_style.yaml")
	defer mockCtrl.Finish()

	mockLogger.EXPECT().Warning("line 3: no space after ':'")
	mockLogger.EXPECT().Warning("line 2: no space before ':'")

	checkJsonSpaces(ctx, linterConf, strings.Split(stack_mocks.ReadFile(t, "./test_resources/spacesjson_testtemplate.json"), "\n"))
}

func TestCheckYamlDashLists(t *testing.T) {
	ctx, mockCtrl, mockLogger, linterConf := setupTestEnv(t, "./test_resources/nodescription_testtemplate.yaml", "test_resources/test_style.yaml")
	defer mockCtrl.Finish()

	mockLogger.EXPECT().Warning("dash lists are not allowed in current lint configuration").Times(1)
	checkYamlLists(ctx, linterConf, stack_mocks.ReadFile(t, "./test_resources/nodescription_testtemplate.yaml"))
}

func TestCheckYamlInlineLists(t *testing.T) {
	ctx, mockCtrl, mockLogger, linterConf := setupTestEnv(t, "./test_resources/nodescription_testtemplate.yaml", "test_resources/test_styleDash.yaml")
	defer mockCtrl.Finish()

	mockLogger.EXPECT().Warning("inline lists are not allowed in current lint configuration").Times(1)
	checkYamlLists(ctx, linterConf, stack_mocks.ReadFile(t, "./test_resources/inlinelist_testtemplate.yaml"))
}

func TestCheckYamlQuotesNoSingleDouble(t *testing.T) {
	ctx, mockCtrl, mockLogger, linterConf := setupTestEnv(t, "./test_resources/nodescription_testtemplate.yaml", "test_resources/test_styleDash.yaml")
	defer mockCtrl.Finish()

	mockLogger.EXPECT().Warning("line 1: double quotes not allowed").Times(1)
	mockLogger.EXPECT().Warning("line 2: double quotes not allowed").Times(1)

	checkYamlQuotes(ctx, linterConf, []string{"ala: \"makota\"", "asd: \"qwe\""})

	mockLogger.EXPECT().Warning("line 2: single quotes not allowed").Times(1)
	mockLogger.EXPECT().Warning("line 3: single quotes not allowed").Times(1)

	checkYamlQuotes(ctx, linterConf, []string{"asd: asd", "qwe: 'qwe'", "zxc: 'zxc'"})
}

func TestCheckYamlQuotesNoQuotes(t *testing.T) {
	ctx, mockCtrl, mockLogger, linterConf := setupTestEnv(t, "./test_resources/nodescription_testtemplate.yaml", "test_resources/test_style.yaml")
	defer mockCtrl.Finish()

	mockLogger.EXPECT().Warning("line 1: quotes required").Times(1)
	mockLogger.EXPECT().Warning("line 2: quotes required").Times(1)

	checkYamlQuotes(ctx, linterConf, []string{"asd: asd", "qwe: qwe"})
}

func TestCheckYamlIndentation(t *testing.T) {
	ctx, mockCtrl, mockLogger, linterConf := setupTestEnv(t, "./test_resources/nodescription_testtemplate.yaml", "test_resources/test_style.yaml")
	defer mockCtrl.Finish()

	checkYamlIndentation(ctx, linterConf, strings.Split(stack_mocks.ReadFile(t, "./test_resources/blanklines_testtemplate.yaml"), "\n"))

	mockLogger.EXPECT().Error("line 6: indentation error")
	mockLogger.EXPECT().Error("line 8: indentation error")

	checkYamlIndentation(ctx, linterConf, strings.Split(stack_mocks.ReadFile(t, "./test_resources/indenterror_testtemplate.yaml"), "\n"))

}

func TestCheckJSONIndentation(t *testing.T) {
	ctx, mockCtrl, _, linterConf := setupTestEnv(t, "./test_resources/nodescription_testtemplate.yaml", "test_resources/test_style.yaml")
	defer mockCtrl.Finish()

	checkJsonIndentation(ctx, linterConf, strings.Split(stack_mocks.ReadFile(t, "./test_resources/spacesjson_testtemplate.json"), "\n"))
}
