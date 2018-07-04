package awsapi

import (
	"github.com/aws/aws-sdk-go/service/cloudformation"
)

type CloudFormationAPI interface {
	CreateStack(input *cloudformation.CreateStackInput) (*cloudformation.CreateStackOutput, error)
	DeleteStack(input *cloudformation.DeleteStackInput) (*cloudformation.DeleteStackOutput, error)
	UpdateStack(input *cloudformation.UpdateStackInput) (*cloudformation.UpdateStackOutput, error)
	SetStackPolicy(input *cloudformation.SetStackPolicyInput) (*cloudformation.SetStackPolicyOutput, error)
	UpdateTerminationProtection(input *cloudformation.UpdateTerminationProtectionInput) (*cloudformation.UpdateTerminationProtectionOutput, error)

	EstimateTemplateCost(input *cloudformation.EstimateTemplateCostInput) (*cloudformation.EstimateTemplateCostOutput, error)
	ValidateTemplate(input *cloudformation.ValidateTemplateInput) (*cloudformation.ValidateTemplateOutput, error)

	CreateChangeSet(input *cloudformation.CreateChangeSetInput) (*cloudformation.CreateChangeSetOutput, error)
	DescribeChangeSet(input *cloudformation.DescribeChangeSetInput) (*cloudformation.DescribeChangeSetOutput, error)
	WaitUntilChangeSetCreateComplete(input *cloudformation.DescribeChangeSetInput) error
}

type AWSCloudFormationAPI struct {
	api *cloudformation.CloudFormation
}

func NewAWSCloudFormation(api *cloudformation.CloudFormation) CloudFormationAPI {
	awsCloudFormationAPI := AWSCloudFormationAPI{
		api: api,
	}
	return &awsCloudFormationAPI
}

func (cf *AWSCloudFormationAPI) CreateStack(input *cloudformation.CreateStackInput) (*cloudformation.CreateStackOutput, error) {
	return cf.api.CreateStack(input)
}
func (cf *AWSCloudFormationAPI) DeleteStack(input *cloudformation.DeleteStackInput) (*cloudformation.DeleteStackOutput, error) {
	return cf.api.DeleteStack(input)
}
func (cf *AWSCloudFormationAPI) UpdateStack(input *cloudformation.UpdateStackInput) (*cloudformation.UpdateStackOutput, error) {
	return cf.api.UpdateStack(input)
}
func (cf *AWSCloudFormationAPI) SetStackPolicy(input *cloudformation.SetStackPolicyInput) (*cloudformation.SetStackPolicyOutput, error) {
	return cf.api.SetStackPolicy(input)
}
func (cf *AWSCloudFormationAPI) UpdateTerminationProtection(input *cloudformation.UpdateTerminationProtectionInput) (*cloudformation.UpdateTerminationProtectionOutput, error) {
	return cf.api.UpdateTerminationProtection(input)
}

func (cf *AWSCloudFormationAPI) EstimateTemplateCost(input *cloudformation.EstimateTemplateCostInput) (*cloudformation.EstimateTemplateCostOutput, error) {
	return cf.api.EstimateTemplateCost(input)
}
func (cf *AWSCloudFormationAPI) ValidateTemplate(input *cloudformation.ValidateTemplateInput) (*cloudformation.ValidateTemplateOutput, error) {
	return cf.api.ValidateTemplate(input)
}

func (cf *AWSCloudFormationAPI) CreateChangeSet(input *cloudformation.CreateChangeSetInput) (*cloudformation.CreateChangeSetOutput, error) {
	return cf.api.CreateChangeSet(input)
}

func (cf *AWSCloudFormationAPI) DescribeChangeSet(input *cloudformation.DescribeChangeSetInput) (*cloudformation.DescribeChangeSetOutput, error) {
	return cf.api.DescribeChangeSet(input)
}
func (cf *AWSCloudFormationAPI) WaitUntilChangeSetCreateComplete(input *cloudformation.DescribeChangeSetInput) error {
	return cf.api.WaitUntilChangeSetCreateComplete(input)
}
