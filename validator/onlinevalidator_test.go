package validator

import (
	"github.com/Appliscale/perun/parameters"
	"github.com/Appliscale/perun/stack/mocks"
	"github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"testing"
)

func TestIsTemplateValid(t *testing.T) {
	ctx := mocks.SetupContext(t, []string{"cmd", "validate", "templatePath"})

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockAWSPI := mocks.NewMockCloudFormationAPI(mockCtrl)
	ctx.CloudFormation = mockAWSPI

	templateBody := "templateBody"
	mockAWSPI.
		EXPECT().
		ValidateTemplate(&cloudformation.ValidateTemplateInput{TemplateBody: &templateBody}).
		Times(2).
		Return(nil, nil)

	isTemplateValid(ctx, &templateBody)
	awsValidate(ctx, &templateBody)
}

func TestEstimateCost(t *testing.T) {
	templatePath := "./test_resources/test_template.yaml"
	ctx := mocks.SetupContext(t, []string{"cmd", "validate", templatePath})

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockAWSPI := mocks.NewMockCloudFormationAPI(mockCtrl)
	ctx.CloudFormation = mockAWSPI

	templateBodyBytes, err := ioutil.ReadFile(templatePath)
	assert.NoError(t, err)
	templateBody := string(templateBodyBytes)
	templateParameters, err := parameters.ResolveParameters(ctx)
	assert.NoError(t, err)

	url := "url"
	mockAWSPI.
		EXPECT().
		EstimateTemplateCost(&cloudformation.EstimateTemplateCostInput{
			TemplateBody: &templateBody,
			Parameters:   templateParameters,
		}).
		Times(1).
		Return(&cloudformation.EstimateTemplateCostOutput{Url: &url}, nil)

	estimateCosts(ctx, &templateBody)
}
