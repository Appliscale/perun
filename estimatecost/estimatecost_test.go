package estimatecost

import (
	"github.com/Appliscale/perun/parameters"
	"github.com/Appliscale/perun/stack/stack_mocks"
	"github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"testing"
)

func TestEstimateCosts(t *testing.T) {
	templatePath := "../validator/test_resources/test_template.yaml"
	ctx := stack_mocks.SetupContext(t, []string{"cmd", "estimate-cost", templatePath})

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockAWSPI := stack_mocks.NewMockCloudFormationAPI(mockCtrl)
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

	EstimateCosts(ctx)
}
