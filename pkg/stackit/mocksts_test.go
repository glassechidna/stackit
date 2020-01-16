package stackit

import (
	"context"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/stretchr/testify/mock"
)

type mockSts struct {
	mock.Mock
}

func (m *mockSts) AssumeRoleWithContext(ctx context.Context, input *sts.AssumeRoleInput, opts ...request.Option) (*sts.AssumeRoleOutput, error) {
	f := m.Called(ctx, input, opts)
	return f.Get(0).(*sts.AssumeRoleOutput), f.Error(1)
}

func (m *mockSts) AssumeRoleWithSAMLWithContext(ctx context.Context, input *sts.AssumeRoleWithSAMLInput, opts ...request.Option) (*sts.AssumeRoleWithSAMLOutput, error) {
	f := m.Called(ctx, input, opts)
	return f.Get(0).(*sts.AssumeRoleWithSAMLOutput), f.Error(1)
}

func (m *mockSts) AssumeRoleWithWebIdentityWithContext(ctx context.Context, input *sts.AssumeRoleWithWebIdentityInput, opts ...request.Option) (*sts.AssumeRoleWithWebIdentityOutput, error) {
	f := m.Called(ctx, input, opts)
	return f.Get(0).(*sts.AssumeRoleWithWebIdentityOutput), f.Error(1)
}

func (m *mockSts) DecodeAuthorizationMessageWithContext(ctx context.Context, input *sts.DecodeAuthorizationMessageInput, opts ...request.Option) (*sts.DecodeAuthorizationMessageOutput, error) {
	f := m.Called(ctx, input, opts)
	return f.Get(0).(*sts.DecodeAuthorizationMessageOutput), f.Error(1)
}

func (m *mockSts) GetAccessKeyInfoWithContext(ctx context.Context, input *sts.GetAccessKeyInfoInput, opts ...request.Option) (*sts.GetAccessKeyInfoOutput, error) {
	f := m.Called(ctx, input, opts)
	return f.Get(0).(*sts.GetAccessKeyInfoOutput), f.Error(1)
}

func (m *mockSts) GetCallerIdentityWithContext(ctx context.Context, input *sts.GetCallerIdentityInput, opts ...request.Option) (*sts.GetCallerIdentityOutput, error) {
	f := m.Called(ctx, input, opts)
	return f.Get(0).(*sts.GetCallerIdentityOutput), f.Error(1)
}

func (m *mockSts) GetFederationTokenWithContext(ctx context.Context, input *sts.GetFederationTokenInput, opts ...request.Option) (*sts.GetFederationTokenOutput, error) {
	f := m.Called(ctx, input, opts)
	return f.Get(0).(*sts.GetFederationTokenOutput), f.Error(1)
}

func (m *mockSts) GetSessionTokenWithContext(ctx context.Context, input *sts.GetSessionTokenInput, opts ...request.Option) (*sts.GetSessionTokenOutput, error) {
	f := m.Called(ctx, input, opts)
	return f.Get(0).(*sts.GetSessionTokenOutput), f.Error(1)
}
