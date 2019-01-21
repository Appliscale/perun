package validator

import (
	"github.com/Appliscale/perun/stack/stack_mocks"
	"github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/golang/mock/gomock"
	"testing"
)

func TestIsTemplateValid(t *testing.T) {
	ctx := stack_mocks.SetupContext(t, []string{"cmd", "validate", "templatePath"})

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockAWSPI := stack_mocks.NewMockCloudFormationAPI(mockCtrl)
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
