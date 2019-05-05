package stackit

import (
	"errors"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestServiceRoleArnCanBeName(t *testing.T) {
	capi := &cfnApi{}
	sapi := &stsApi{}
	s := NewStackit(capi, sapi, "stack-name")

	input := StackitUpInput{
		RoleARN: "MyRoleName",
	}

	sapi.GetCallerIdentityF = func(input *sts.GetCallerIdentityInput) (*sts.GetCallerIdentityOutput, error) {
		return &sts.GetCallerIdentityOutput{
			Account: aws.String("1234567890"),
		}, nil
	}

	capi.DescribeStacksF = func(input *cloudformation.DescribeStacksInput) (*cloudformation.DescribeStacksOutput, error) {
		return &cloudformation.DescribeStacksOutput{
			Stacks: []*cloudformation.Stack{
				{
					StackId: aws.String("arn:aws:cloudformation:ap-southeast-2:657110686698:stack/stackset-role/58ed6a10-3e2f-11e9-bc5f-0a9966e9c45e"),
				},
			},
		}, nil
	}

	capi.CreateChangeSetF = func(input *cloudformation.CreateChangeSetInput) (*cloudformation.CreateChangeSetOutput, error) {
		assert.Equal(t, "arn:aws:iam::1234567890:role/MyRoleName", *input.RoleARN)
		return nil, errors.New("done")
	}

	ch := make(chan TailStackEvent)
	go s.Up(input, ch)

	ev := <-ch
	assert.EqualError(t, ev.StackitError, "creating change set: done")
}

func TestServiceRoleArnDoesntTriggerStsCall(t *testing.T) {
	capi := &cfnApi{}
	sapi := &stsApi{}
	s := NewStackit(capi, sapi, "stack-name")

	input := StackitUpInput{
		RoleARN: "arn:aws:iam::1234567890:role/MyRoleName",
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
		return &cloudformation.DescribeStacksOutput{
			Stacks: []*cloudformation.Stack{
				{
					StackId: aws.String("arn:aws:cloudformation:ap-southeast-2:657110686698:stack/stack-name/58ed6a10-3e2f-11e9-bc5f-0a9966e9c45e"),
				},
			},
		}, nil
	}

	ch := make(chan TailStackEvent)
	go s.Up(input, ch)

	ev := <-ch
	assert.EqualError(t, ev.StackitError, "creating change set: done")
}
