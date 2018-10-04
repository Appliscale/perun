// Copyright 2018 Appliscale
//
// Maintainers and contributors are listed in README file inside repository.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package awsapi contains interface with all functions which use AWS CloudFormation API.
package awsapi

import (
	"github.com/aws/aws-sdk-go/service/cloudformation"
)

// CloudFormationAPI interface.
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
	DeleteChangeSet(input *cloudformation.DeleteChangeSetInput) (*cloudformation.DeleteChangeSetOutput, error)
}

// AWSCloudFormationAPI implements CLoudFormationAPI.
type AWSCloudFormationAPI struct {
	api *cloudformation.CloudFormation
}

// NewAWSCloudFormation creates CloudFormationAPI.
func NewAWSCloudFormation(api *cloudformation.CloudFormation) CloudFormationAPI {
	awsCloudFormationAPI := AWSCloudFormationAPI{
		api: api,
	}
	return &awsCloudFormationAPI
}

// CreateStack creates stack based on user template.
func (cf *AWSCloudFormationAPI) CreateStack(input *cloudformation.CreateStackInput) (*cloudformation.CreateStackOutput, error) {
	return cf.api.CreateStack(input)
}

// DeleteStack destroys template based on stack name.
func (cf *AWSCloudFormationAPI) DeleteStack(input *cloudformation.DeleteStackInput) (*cloudformation.DeleteStackOutput, error) {
	return cf.api.DeleteStack(input)
}

// UpdateStack updates stack template.
func (cf *AWSCloudFormationAPI) UpdateStack(input *cloudformation.UpdateStackInput) (*cloudformation.UpdateStackOutput, error) {
	return cf.api.UpdateStack(input)
}

// SetStackPolicy sets policy based on template or flag.
func (cf *AWSCloudFormationAPI) SetStackPolicy(input *cloudformation.SetStackPolicyInput) (*cloudformation.SetStackPolicyOutput, error) {
	return cf.api.SetStackPolicy(input)
}

// UpdateTerminationProtection allows to change stack protection from termination.
func (cf *AWSCloudFormationAPI) UpdateTerminationProtection(input *cloudformation.UpdateTerminationProtectionInput) (*cloudformation.UpdateTerminationProtectionOutput, error) {
	return cf.api.UpdateTerminationProtection(input)
}

// EstimateTemplateCost shows stack cost.
func (cf *AWSCloudFormationAPI) EstimateTemplateCost(input *cloudformation.EstimateTemplateCostInput) (*cloudformation.EstimateTemplateCostOutput, error) {
	return cf.api.EstimateTemplateCost(input)
}

// ValidateTemplate checks template correctness.
func (cf *AWSCloudFormationAPI) ValidateTemplate(input *cloudformation.ValidateTemplateInput) (*cloudformation.ValidateTemplateOutput, error) {
	return cf.api.ValidateTemplate(input)
}

// CreateChangeSet creates ChangeSet.
func (cf *AWSCloudFormationAPI) CreateChangeSet(input *cloudformation.CreateChangeSetInput) (*cloudformation.CreateChangeSetOutput, error) {
	return cf.api.CreateChangeSet(input)
}

// DescribeChangeSet returns the inputs and a list of changes.
func (cf *AWSCloudFormationAPI) DescribeChangeSet(input *cloudformation.DescribeChangeSetInput) (*cloudformation.DescribeChangeSetOutput, error) {
	return cf.api.DescribeChangeSet(input)
}

// WaitUntilChangeSetCreateComplete uses the AWS CloudFormation API operation
// DescribeChangeSet to wait for a condition to be met before returning.
// If the condition is not met within the max attempt window, an error will
// be returned.
func (cf *AWSCloudFormationAPI) WaitUntilChangeSetCreateComplete(input *cloudformation.DescribeChangeSetInput) error {
	return cf.api.WaitUntilChangeSetCreateComplete(input)
}

// DeleteChangeSet removes ChangeSet.
func (cf *AWSCloudFormationAPI) DeleteChangeSet(input *cloudformation.DeleteChangeSetInput) (*cloudformation.DeleteChangeSetOutput, error) {
	return cf.api.DeleteChangeSet(input)
}
