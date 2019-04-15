package stackit

import (
	"github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/aws/aws-sdk-go/service/cloudformation/cloudformationiface"
	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/aws/aws-sdk-go/service/sts/stsiface"
)

type cfnApi struct {
	cloudformationiface.CloudFormationAPI
	CreateChangeSetF func(input *cloudformation.CreateChangeSetInput) (*cloudformation.CreateChangeSetOutput, error)
	DescribeStackEventsPagesF func(*cloudformation.DescribeStackEventsInput, func(*cloudformation.DescribeStackEventsOutput, bool) bool) error
	DescribeStacksF func (*cloudformation.DescribeStacksInput) (*cloudformation.DescribeStacksOutput, error)

}

func (c *cfnApi) DescribeStacks(input *cloudformation.DescribeStacksInput) (*cloudformation.DescribeStacksOutput, error) {
	return c.DescribeStacksF(input)
}

func (c *cfnApi) CreateChangeSet(input *cloudformation.CreateChangeSetInput) (*cloudformation.CreateChangeSetOutput, error) {
	return c.CreateChangeSetF(input)
}

func (c *cfnApi) DescribeStackEventsPages(input *cloudformation.DescribeStackEventsInput, cb func(*cloudformation.DescribeStackEventsOutput, bool) bool) error {
	return c.DescribeStackEventsPagesF(input, cb)
}

type stsApi struct {
	stsiface.STSAPI
	GetCallerIdentityF func(*sts.GetCallerIdentityInput) (*sts.GetCallerIdentityOutput, error)
}

func (s *stsApi) GetCallerIdentity(input *sts.GetCallerIdentityInput) (*sts.GetCallerIdentityOutput, error) {
	return s.GetCallerIdentityF(input)
}
