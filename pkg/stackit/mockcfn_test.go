package stackit

import (
	"context"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/stretchr/testify/mock"
)

type mockCfn struct {
	mock.Mock
}

func (m *mockCfn) CancelUpdateStackWithContext(ctx context.Context, input *cloudformation.CancelUpdateStackInput, opts ...request.Option) (*cloudformation.CancelUpdateStackOutput, error) {
	f := m.Called(ctx, input, opts)
	output, _ := f.Get(0).(*cloudformation.CancelUpdateStackOutput)
	return output, f.Error(1)
}

func (m *mockCfn) ContinueUpdateRollbackWithContext(ctx context.Context, input *cloudformation.ContinueUpdateRollbackInput, opts ...request.Option) (*cloudformation.ContinueUpdateRollbackOutput, error) {
	f := m.Called(ctx, input, opts)
	output, _ := f.Get(0).(*cloudformation.ContinueUpdateRollbackOutput)
	return output, f.Error(1)
}

func (m *mockCfn) CreateChangeSetWithContext(ctx context.Context, input *cloudformation.CreateChangeSetInput, opts ...request.Option) (*cloudformation.CreateChangeSetOutput, error) {
	f := m.Called(ctx, input, opts)
	output, _ := f.Get(0).(*cloudformation.CreateChangeSetOutput)
	return output, f.Error(1)
}

func (m *mockCfn) CreateStackWithContext(ctx context.Context, input *cloudformation.CreateStackInput, opts ...request.Option) (*cloudformation.CreateStackOutput, error) {
	f := m.Called(ctx, input, opts)
	output, _ := f.Get(0).(*cloudformation.CreateStackOutput)
	return output, f.Error(1)
}

func (m *mockCfn) CreateStackInstancesWithContext(ctx context.Context, input *cloudformation.CreateStackInstancesInput, opts ...request.Option) (*cloudformation.CreateStackInstancesOutput, error) {
	f := m.Called(ctx, input, opts)
	output, _ := f.Get(0).(*cloudformation.CreateStackInstancesOutput)
	return output, f.Error(1)
}

func (m *mockCfn) CreateStackSetWithContext(ctx context.Context, input *cloudformation.CreateStackSetInput, opts ...request.Option) (*cloudformation.CreateStackSetOutput, error) {
	f := m.Called(ctx, input, opts)
	output, _ := f.Get(0).(*cloudformation.CreateStackSetOutput)
	return output, f.Error(1)
}

func (m *mockCfn) DeleteChangeSetWithContext(ctx context.Context, input *cloudformation.DeleteChangeSetInput, opts ...request.Option) (*cloudformation.DeleteChangeSetOutput, error) {
	f := m.Called(ctx, input, opts)
	output, _ := f.Get(0).(*cloudformation.DeleteChangeSetOutput)
	return output, f.Error(1)
}

func (m *mockCfn) DeleteStackWithContext(ctx context.Context, input *cloudformation.DeleteStackInput, opts ...request.Option) (*cloudformation.DeleteStackOutput, error) {
	f := m.Called(ctx, input, opts)
	output, _ := f.Get(0).(*cloudformation.DeleteStackOutput)
	return output, f.Error(1)
}

func (m *mockCfn) DeleteStackInstancesWithContext(ctx context.Context, input *cloudformation.DeleteStackInstancesInput, opts ...request.Option) (*cloudformation.DeleteStackInstancesOutput, error) {
	f := m.Called(ctx, input, opts)
	output, _ := f.Get(0).(*cloudformation.DeleteStackInstancesOutput)
	return output, f.Error(1)
}

func (m *mockCfn) DeleteStackSetWithContext(ctx context.Context, input *cloudformation.DeleteStackSetInput, opts ...request.Option) (*cloudformation.DeleteStackSetOutput, error) {
	f := m.Called(ctx, input, opts)
	output, _ := f.Get(0).(*cloudformation.DeleteStackSetOutput)
	return output, f.Error(1)
}

func (m *mockCfn) DeregisterTypeWithContext(ctx context.Context, input *cloudformation.DeregisterTypeInput, opts ...request.Option) (*cloudformation.DeregisterTypeOutput, error) {
	f := m.Called(ctx, input, opts)
	output, _ := f.Get(0).(*cloudformation.DeregisterTypeOutput)
	return output, f.Error(1)
}

func (m *mockCfn) DescribeAccountLimitsWithContext(ctx context.Context, input *cloudformation.DescribeAccountLimitsInput, opts ...request.Option) (*cloudformation.DescribeAccountLimitsOutput, error) {
	f := m.Called(ctx, input, opts)
	output, _ := f.Get(0).(*cloudformation.DescribeAccountLimitsOutput)
	return output, f.Error(1)
}

func (m *mockCfn) DescribeChangeSetWithContext(ctx context.Context, input *cloudformation.DescribeChangeSetInput, opts ...request.Option) (*cloudformation.DescribeChangeSetOutput, error) {
	f := m.Called(ctx, input, opts)
	output, _ := f.Get(0).(*cloudformation.DescribeChangeSetOutput)
	return output, f.Error(1)
}

func (m *mockCfn) DescribeStackDriftDetectionStatusWithContext(ctx context.Context, input *cloudformation.DescribeStackDriftDetectionStatusInput, opts ...request.Option) (*cloudformation.DescribeStackDriftDetectionStatusOutput, error) {
	f := m.Called(ctx, input, opts)
	output, _ := f.Get(0).(*cloudformation.DescribeStackDriftDetectionStatusOutput)
	return output, f.Error(1)
}

func (m *mockCfn) DescribeStackEventsWithContext(ctx context.Context, input *cloudformation.DescribeStackEventsInput, opts ...request.Option) (*cloudformation.DescribeStackEventsOutput, error) {
	f := m.Called(ctx, input, opts)
	output, _ := f.Get(0).(*cloudformation.DescribeStackEventsOutput)
	return output, f.Error(1)
}

func (m *mockCfn) DescribeStackEventsPagesWithContext(ctx context.Context, input *cloudformation.DescribeStackEventsInput, cb func(*cloudformation.DescribeStackEventsOutput, bool) bool, opts ...request.Option) error {
	f := m.Called(ctx, input, opts)
	return f.Error(0)
}

func (m *mockCfn) DescribeStackInstanceWithContext(ctx context.Context, input *cloudformation.DescribeStackInstanceInput, opts ...request.Option) (*cloudformation.DescribeStackInstanceOutput, error) {
	f := m.Called(ctx, input, opts)
	output, _ := f.Get(0).(*cloudformation.DescribeStackInstanceOutput)
	return output, f.Error(1)
}

func (m *mockCfn) DescribeStackResourceWithContext(ctx context.Context, input *cloudformation.DescribeStackResourceInput, opts ...request.Option) (*cloudformation.DescribeStackResourceOutput, error) {
	f := m.Called(ctx, input, opts)
	output, _ := f.Get(0).(*cloudformation.DescribeStackResourceOutput)
	return output, f.Error(1)
}

func (m *mockCfn) DescribeStackResourceDriftsWithContext(ctx context.Context, input *cloudformation.DescribeStackResourceDriftsInput, opts ...request.Option) (*cloudformation.DescribeStackResourceDriftsOutput, error) {
	f := m.Called(ctx, input, opts)
	output, _ := f.Get(0).(*cloudformation.DescribeStackResourceDriftsOutput)
	return output, f.Error(1)
}

func (m *mockCfn) DescribeStackResourceDriftsPagesWithContext(ctx context.Context, input *cloudformation.DescribeStackResourceDriftsInput, cb func(*cloudformation.DescribeStackResourceDriftsOutput, bool) bool, opts ...request.Option) error {
	f := m.Called(ctx, input, opts)
	return f.Error(0)
}

func (m *mockCfn) DescribeStackResourcesWithContext(ctx context.Context, input *cloudformation.DescribeStackResourcesInput, opts ...request.Option) (*cloudformation.DescribeStackResourcesOutput, error) {
	f := m.Called(ctx, input, opts)
	output, _ := f.Get(0).(*cloudformation.DescribeStackResourcesOutput)
	return output, f.Error(1)
}

func (m *mockCfn) DescribeStackSetWithContext(ctx context.Context, input *cloudformation.DescribeStackSetInput, opts ...request.Option) (*cloudformation.DescribeStackSetOutput, error) {
	f := m.Called(ctx, input, opts)
	output, _ := f.Get(0).(*cloudformation.DescribeStackSetOutput)
	return output, f.Error(1)
}

func (m *mockCfn) DescribeStackSetOperationWithContext(ctx context.Context, input *cloudformation.DescribeStackSetOperationInput, opts ...request.Option) (*cloudformation.DescribeStackSetOperationOutput, error) {
	f := m.Called(ctx, input, opts)
	output, _ := f.Get(0).(*cloudformation.DescribeStackSetOperationOutput)
	return output, f.Error(1)
}

func (m *mockCfn) DescribeStacksWithContext(ctx context.Context, input *cloudformation.DescribeStacksInput, opts ...request.Option) (*cloudformation.DescribeStacksOutput, error) {
	f := m.Called(ctx, input, opts)
	output, _ := f.Get(0).(*cloudformation.DescribeStacksOutput)
	return output, f.Error(1)
}

func (m *mockCfn) DescribeStacksPagesWithContext(ctx context.Context, input *cloudformation.DescribeStacksInput, cb func(*cloudformation.DescribeStacksOutput, bool) bool, opts ...request.Option) error {
	f := m.Called(ctx, input, opts)
	return f.Error(0)
}

func (m *mockCfn) DescribeTypeWithContext(ctx context.Context, input *cloudformation.DescribeTypeInput, opts ...request.Option) (*cloudformation.DescribeTypeOutput, error) {
	f := m.Called(ctx, input, opts)
	output, _ := f.Get(0).(*cloudformation.DescribeTypeOutput)
	return output, f.Error(1)
}

func (m *mockCfn) DescribeTypeRegistrationWithContext(ctx context.Context, input *cloudformation.DescribeTypeRegistrationInput, opts ...request.Option) (*cloudformation.DescribeTypeRegistrationOutput, error) {
	f := m.Called(ctx, input, opts)
	output, _ := f.Get(0).(*cloudformation.DescribeTypeRegistrationOutput)
	return output, f.Error(1)
}

func (m *mockCfn) DetectStackDriftWithContext(ctx context.Context, input *cloudformation.DetectStackDriftInput, opts ...request.Option) (*cloudformation.DetectStackDriftOutput, error) {
	f := m.Called(ctx, input, opts)
	output, _ := f.Get(0).(*cloudformation.DetectStackDriftOutput)
	return output, f.Error(1)
}

func (m *mockCfn) DetectStackResourceDriftWithContext(ctx context.Context, input *cloudformation.DetectStackResourceDriftInput, opts ...request.Option) (*cloudformation.DetectStackResourceDriftOutput, error) {
	f := m.Called(ctx, input, opts)
	output, _ := f.Get(0).(*cloudformation.DetectStackResourceDriftOutput)
	return output, f.Error(1)
}

func (m *mockCfn) DetectStackSetDriftWithContext(ctx context.Context, input *cloudformation.DetectStackSetDriftInput, opts ...request.Option) (*cloudformation.DetectStackSetDriftOutput, error) {
	f := m.Called(ctx, input, opts)
	output, _ := f.Get(0).(*cloudformation.DetectStackSetDriftOutput)
	return output, f.Error(1)
}

func (m *mockCfn) EstimateTemplateCostWithContext(ctx context.Context, input *cloudformation.EstimateTemplateCostInput, opts ...request.Option) (*cloudformation.EstimateTemplateCostOutput, error) {
	f := m.Called(ctx, input, opts)
	output, _ := f.Get(0).(*cloudformation.EstimateTemplateCostOutput)
	return output, f.Error(1)
}

func (m *mockCfn) ExecuteChangeSetWithContext(ctx context.Context, input *cloudformation.ExecuteChangeSetInput, opts ...request.Option) (*cloudformation.ExecuteChangeSetOutput, error) {
	f := m.Called(ctx, input, opts)
	output, _ := f.Get(0).(*cloudformation.ExecuteChangeSetOutput)
	return output, f.Error(1)
}

func (m *mockCfn) GetStackPolicyWithContext(ctx context.Context, input *cloudformation.GetStackPolicyInput, opts ...request.Option) (*cloudformation.GetStackPolicyOutput, error) {
	f := m.Called(ctx, input, opts)
	output, _ := f.Get(0).(*cloudformation.GetStackPolicyOutput)
	return output, f.Error(1)
}

func (m *mockCfn) GetTemplateWithContext(ctx context.Context, input *cloudformation.GetTemplateInput, opts ...request.Option) (*cloudformation.GetTemplateOutput, error) {
	f := m.Called(ctx, input, opts)
	output, _ := f.Get(0).(*cloudformation.GetTemplateOutput)
	return output, f.Error(1)
}

func (m *mockCfn) GetTemplateSummaryWithContext(ctx context.Context, input *cloudformation.GetTemplateSummaryInput, opts ...request.Option) (*cloudformation.GetTemplateSummaryOutput, error) {
	f := m.Called(ctx, input, opts)
	output, _ := f.Get(0).(*cloudformation.GetTemplateSummaryOutput)
	return output, f.Error(1)
}

func (m *mockCfn) ListChangeSetsWithContext(ctx context.Context, input *cloudformation.ListChangeSetsInput, opts ...request.Option) (*cloudformation.ListChangeSetsOutput, error) {
	f := m.Called(ctx, input, opts)
	output, _ := f.Get(0).(*cloudformation.ListChangeSetsOutput)
	return output, f.Error(1)
}

func (m *mockCfn) ListExportsWithContext(ctx context.Context, input *cloudformation.ListExportsInput, opts ...request.Option) (*cloudformation.ListExportsOutput, error) {
	f := m.Called(ctx, input, opts)
	output, _ := f.Get(0).(*cloudformation.ListExportsOutput)
	return output, f.Error(1)
}

func (m *mockCfn) ListExportsPagesWithContext(ctx context.Context, input *cloudformation.ListExportsInput, cb func(*cloudformation.ListExportsOutput, bool) bool, opts ...request.Option) error {
	f := m.Called(ctx, input, opts)
	return f.Error(0)
}

func (m *mockCfn) ListImportsWithContext(ctx context.Context, input *cloudformation.ListImportsInput, opts ...request.Option) (*cloudformation.ListImportsOutput, error) {
	f := m.Called(ctx, input, opts)
	output, _ := f.Get(0).(*cloudformation.ListImportsOutput)
	return output, f.Error(1)
}

func (m *mockCfn) ListImportsPagesWithContext(ctx context.Context, input *cloudformation.ListImportsInput, cb func(*cloudformation.ListImportsOutput, bool) bool, opts ...request.Option) error {
	f := m.Called(ctx, input, opts)
	return f.Error(0)
}

func (m *mockCfn) ListStackInstancesWithContext(ctx context.Context, input *cloudformation.ListStackInstancesInput, opts ...request.Option) (*cloudformation.ListStackInstancesOutput, error) {
	f := m.Called(ctx, input, opts)
	output, _ := f.Get(0).(*cloudformation.ListStackInstancesOutput)
	return output, f.Error(1)
}

func (m *mockCfn) ListStackResourcesWithContext(ctx context.Context, input *cloudformation.ListStackResourcesInput, opts ...request.Option) (*cloudformation.ListStackResourcesOutput, error) {
	f := m.Called(ctx, input, opts)
	output, _ := f.Get(0).(*cloudformation.ListStackResourcesOutput)
	return output, f.Error(1)
}

func (m *mockCfn) ListStackResourcesPagesWithContext(ctx context.Context, input *cloudformation.ListStackResourcesInput, cb func(*cloudformation.ListStackResourcesOutput, bool) bool, opts ...request.Option) error {
	f := m.Called(ctx, input, opts)
	return f.Error(0)
}

func (m *mockCfn) ListStackSetOperationResultsWithContext(ctx context.Context, input *cloudformation.ListStackSetOperationResultsInput, opts ...request.Option) (*cloudformation.ListStackSetOperationResultsOutput, error) {
	f := m.Called(ctx, input, opts)
	output, _ := f.Get(0).(*cloudformation.ListStackSetOperationResultsOutput)
	return output, f.Error(1)
}

func (m *mockCfn) ListStackSetOperationsWithContext(ctx context.Context, input *cloudformation.ListStackSetOperationsInput, opts ...request.Option) (*cloudformation.ListStackSetOperationsOutput, error) {
	f := m.Called(ctx, input, opts)
	output, _ := f.Get(0).(*cloudformation.ListStackSetOperationsOutput)
	return output, f.Error(1)
}

func (m *mockCfn) ListStackSetsWithContext(ctx context.Context, input *cloudformation.ListStackSetsInput, opts ...request.Option) (*cloudformation.ListStackSetsOutput, error) {
	f := m.Called(ctx, input, opts)
	output, _ := f.Get(0).(*cloudformation.ListStackSetsOutput)
	return output, f.Error(1)
}

func (m *mockCfn) ListStacksWithContext(ctx context.Context, input *cloudformation.ListStacksInput, opts ...request.Option) (*cloudformation.ListStacksOutput, error) {
	f := m.Called(ctx, input, opts)
	output, _ := f.Get(0).(*cloudformation.ListStacksOutput)
	return output, f.Error(1)
}

func (m *mockCfn) ListStacksPagesWithContext(ctx context.Context, input *cloudformation.ListStacksInput, cb func(*cloudformation.ListStacksOutput, bool) bool, opts ...request.Option) error {
	f := m.Called(ctx, input, opts)
	return f.Error(0)
}

func (m *mockCfn) ListTypeRegistrationsWithContext(ctx context.Context, input *cloudformation.ListTypeRegistrationsInput, opts ...request.Option) (*cloudformation.ListTypeRegistrationsOutput, error) {
	f := m.Called(ctx, input, opts)
	output, _ := f.Get(0).(*cloudformation.ListTypeRegistrationsOutput)
	return output, f.Error(1)
}

func (m *mockCfn) ListTypeRegistrationsPagesWithContext(ctx context.Context, input *cloudformation.ListTypeRegistrationsInput, cb func(*cloudformation.ListTypeRegistrationsOutput, bool) bool, opts ...request.Option) error {
	f := m.Called(ctx, input, opts)
	return f.Error(0)
}

func (m *mockCfn) ListTypeVersionsWithContext(ctx context.Context, input *cloudformation.ListTypeVersionsInput, opts ...request.Option) (*cloudformation.ListTypeVersionsOutput, error) {
	f := m.Called(ctx, input, opts)
	output, _ := f.Get(0).(*cloudformation.ListTypeVersionsOutput)
	return output, f.Error(1)
}

func (m *mockCfn) ListTypeVersionsPagesWithContext(ctx context.Context, input *cloudformation.ListTypeVersionsInput, cb func(*cloudformation.ListTypeVersionsOutput, bool) bool, opts ...request.Option) error {
	f := m.Called(ctx, input, opts)
	return f.Error(0)
}

func (m *mockCfn) ListTypesWithContext(ctx context.Context, input *cloudformation.ListTypesInput, opts ...request.Option) (*cloudformation.ListTypesOutput, error) {
	f := m.Called(ctx, input, opts)
	output, _ := f.Get(0).(*cloudformation.ListTypesOutput)
	return output, f.Error(1)
}

func (m *mockCfn) ListTypesPagesWithContext(ctx context.Context, input *cloudformation.ListTypesInput, cb func(*cloudformation.ListTypesOutput, bool) bool, opts ...request.Option) error {
	f := m.Called(ctx, input, opts)
	return f.Error(0)
}

func (m *mockCfn) RecordHandlerProgressWithContext(ctx context.Context, input *cloudformation.RecordHandlerProgressInput, opts ...request.Option) (*cloudformation.RecordHandlerProgressOutput, error) {
	f := m.Called(ctx, input, opts)
	output, _ := f.Get(0).(*cloudformation.RecordHandlerProgressOutput)
	return output, f.Error(1)
}

func (m *mockCfn) RegisterTypeWithContext(ctx context.Context, input *cloudformation.RegisterTypeInput, opts ...request.Option) (*cloudformation.RegisterTypeOutput, error) {
	f := m.Called(ctx, input, opts)
	output, _ := f.Get(0).(*cloudformation.RegisterTypeOutput)
	return output, f.Error(1)
}

func (m *mockCfn) SetStackPolicyWithContext(ctx context.Context, input *cloudformation.SetStackPolicyInput, opts ...request.Option) (*cloudformation.SetStackPolicyOutput, error) {
	f := m.Called(ctx, input, opts)
	output, _ := f.Get(0).(*cloudformation.SetStackPolicyOutput)
	return output, f.Error(1)
}

func (m *mockCfn) SetTypeDefaultVersionWithContext(ctx context.Context, input *cloudformation.SetTypeDefaultVersionInput, opts ...request.Option) (*cloudformation.SetTypeDefaultVersionOutput, error) {
	f := m.Called(ctx, input, opts)
	output, _ := f.Get(0).(*cloudformation.SetTypeDefaultVersionOutput)
	return output, f.Error(1)
}

func (m *mockCfn) SignalResourceWithContext(ctx context.Context, input *cloudformation.SignalResourceInput, opts ...request.Option) (*cloudformation.SignalResourceOutput, error) {
	f := m.Called(ctx, input, opts)
	output, _ := f.Get(0).(*cloudformation.SignalResourceOutput)
	return output, f.Error(1)
}

func (m *mockCfn) StopStackSetOperationWithContext(ctx context.Context, input *cloudformation.StopStackSetOperationInput, opts ...request.Option) (*cloudformation.StopStackSetOperationOutput, error) {
	f := m.Called(ctx, input, opts)
	output, _ := f.Get(0).(*cloudformation.StopStackSetOperationOutput)
	return output, f.Error(1)
}

func (m *mockCfn) UpdateStackWithContext(ctx context.Context, input *cloudformation.UpdateStackInput, opts ...request.Option) (*cloudformation.UpdateStackOutput, error) {
	f := m.Called(ctx, input, opts)
	output, _ := f.Get(0).(*cloudformation.UpdateStackOutput)
	return output, f.Error(1)
}

func (m *mockCfn) UpdateStackInstancesWithContext(ctx context.Context, input *cloudformation.UpdateStackInstancesInput, opts ...request.Option) (*cloudformation.UpdateStackInstancesOutput, error) {
	f := m.Called(ctx, input, opts)
	output, _ := f.Get(0).(*cloudformation.UpdateStackInstancesOutput)
	return output, f.Error(1)
}

func (m *mockCfn) UpdateStackSetWithContext(ctx context.Context, input *cloudformation.UpdateStackSetInput, opts ...request.Option) (*cloudformation.UpdateStackSetOutput, error) {
	f := m.Called(ctx, input, opts)
	output, _ := f.Get(0).(*cloudformation.UpdateStackSetOutput)
	return output, f.Error(1)
}

func (m *mockCfn) UpdateTerminationProtectionWithContext(ctx context.Context, input *cloudformation.UpdateTerminationProtectionInput, opts ...request.Option) (*cloudformation.UpdateTerminationProtectionOutput, error) {
	f := m.Called(ctx, input, opts)
	output, _ := f.Get(0).(*cloudformation.UpdateTerminationProtectionOutput)
	return output, f.Error(1)
}

func (m *mockCfn) ValidateTemplateWithContext(ctx context.Context, input *cloudformation.ValidateTemplateInput, opts ...request.Option) (*cloudformation.ValidateTemplateOutput, error) {
	f := m.Called(ctx, input, opts)
	output, _ := f.Get(0).(*cloudformation.ValidateTemplateOutput)
	return output, f.Error(1)
}

