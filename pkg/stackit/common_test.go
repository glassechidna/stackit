package stackit

import (
	"context"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/glassechidna/awsctx/service/cloudformationctx"
	"github.com/glassechidna/awsctx/service/stsctx"
)

type cfnApi struct {
	cloudformationctx.CloudFormation
	CreateChangeSetF          func(input *cloudformation.CreateChangeSetInput) (*cloudformation.CreateChangeSetOutput, error)
	DescribeStackEventsPagesF func(*cloudformation.DescribeStackEventsInput, func(*cloudformation.DescribeStackEventsOutput, bool) bool) error
	DescribeStacksF           func(*cloudformation.DescribeStacksInput) (*cloudformation.DescribeStacksOutput, error)
}

func (c *cfnApi) DescribeStacksWithContext(ctx context.Context, input *cloudformation.DescribeStacksInput, opts ...request.Option) (*cloudformation.DescribeStacksOutput, error) {
	return c.DescribeStacksF(input)
}

func (c *cfnApi) CreateChangeSetWithContext(ctx context.Context, input *cloudformation.CreateChangeSetInput, opts ...request.Option) (*cloudformation.CreateChangeSetOutput, error) {
	return c.CreateChangeSetF(input)
}

func (c *cfnApi) DescribeStackEventsPagesWithContext(ctx context.Context, input *cloudformation.DescribeStackEventsInput, cb func(*cloudformation.DescribeStackEventsOutput, bool) bool, opts ...request.Option) error {
	return c.DescribeStackEventsPagesF(input, cb)
}

type stsApi struct {
	stsctx.STS
	GetCallerIdentityF func(*sts.GetCallerIdentityInput) (*sts.GetCallerIdentityOutput, error)
}

func (s *stsApi) GetCallerIdentityWithContext(ctx context.Context, input *sts.GetCallerIdentityInput, opts ...request.Option) (*sts.GetCallerIdentityOutput, error) {
	return s.GetCallerIdentityF(input)
}
