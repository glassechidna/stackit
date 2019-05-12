package stackit

import (
	"context"
	"errors"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestServiceRoleArnCanBeName(t *testing.T) {
	capi := &cfnApi{}
	sapi := &stsApi{}
	s := NewStackit(capi, sapi)

	input := StackitUpInput{
		StackName: "stack-name",
		RoleARN:   "MyRoleName",
	}

	sapi.GetCallerIdentityF = func(input *sts.GetCallerIdentityInput) (*sts.GetCallerIdentityOutput, error) {
		return &sts.GetCallerIdentityOutput{
			Account: aws.String("1234567890"),
		}, nil
	}

	capi.DescribeStacksF = func(input *cloudformation.DescribeStacksInput) (*cloudformation.DescribeStacksOutput, error) {
		capi.DescribeStacksF = func(input *cloudformation.DescribeStacksInput) (*cloudformation.DescribeStacksOutput, error) {
			return &cloudformation.DescribeStacksOutput{
				Stacks: []*cloudformation.Stack{
					{
						StackId: aws.String("arn:aws:cloudformation:ap-southeast-2:657110686698:stack/stackset-role/58ed6a10-3e2f-11e9-bc5f-0a9966e9c45e"),
					},
				},
			}, nil
		}
		return nil, awserr.New("ValidationError", "", nil)
	}

	capi.CreateChangeSetF = func(input *cloudformation.CreateChangeSetInput) (*cloudformation.CreateChangeSetOutput, error) {
		assert.Equal(t, "arn:aws:iam::1234567890:role/MyRoleName", *input.RoleARN)
		return nil, errors.New("done")
	}

	ch := make(chan TailStackEvent)
	_, err := s.Prepare(context.Background(), input, ch)
	assert.EqualError(t, err, "creating change set: done")
}

func TestServiceRoleArnDoesntTriggerStsCall(t *testing.T) {
	capi := &cfnApi{}
	sapi := &stsApi{}
	s := NewStackit(capi, sapi)

	input := StackitUpInput{
		StackName: "stack-name",
		RoleARN:   "arn:aws:iam::1234567890:role/MyRoleName",
	}

	sapi.GetCallerIdentityF = func(input *sts.GetCallerIdentityInput) (*sts.GetCallerIdentityOutput, error) {
		assert.Fail(t, "shouldn't call sts:GetCallerIdentity if full role arn passed")
		return nil, errors.New("shouldn't come here")
	}

	capi.CreateChangeSetF = func(input *cloudformation.CreateChangeSetInput) (*cloudformation.CreateChangeSetOutput, error) {
		assert.Equal(t, "arn:aws:iam::1234567890:role/MyRoleName", *input.RoleARN)
		return nil, errors.New("done")
	}

	capi.DescribeStacksF = func(input *cloudformation.DescribeStacksInput) (*cloudformation.DescribeStacksOutput, error) {
		capi.DescribeStacksF = func(input *cloudformation.DescribeStacksInput) (*cloudformation.DescribeStacksOutput, error) {
			return &cloudformation.DescribeStacksOutput{
				Stacks: []*cloudformation.Stack{
					{
						StackId: aws.String("arn:aws:cloudformation:ap-southeast-2:657110686698:stack/stackset-role/58ed6a10-3e2f-11e9-bc5f-0a9966e9c45e"),
					},
				},
			}, nil
		}
		return nil, awserr.New("ValidationError", "", nil)
	}

	ch := make(chan TailStackEvent)
	_, err := s.Prepare(context.Background(), input, ch)
	assert.EqualError(t, err, "creating change set: done")
}
