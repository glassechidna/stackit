package stackit

import (
	"context"
	"errors"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

func TestServiceRoleArnCanBeName(t *testing.T) {
	capi := &mockCfn{}
	capi.On("DescribeStacksWithContext", mock.Anything, mock.Anything, mock.Anything).Return(nil, awserr.New("ValidationError", "", nil))
	capi.On("CreateChangeSetWithContext", mock.Anything, mock.Anything, mock.Anything).Return(nil, errors.New("done")).Run(func(args mock.Arguments) {
		input := args.Get(1).(*cloudformation.CreateChangeSetInput)
		assert.Equal(t, "arn:aws:iam::1234567890:role/MyRoleName", *input.RoleARN)
	})

	sapi := &mockSts{}
	sapi.On("GetCallerIdentityWithContext", mock.Anything, mock.Anything, mock.Anything).Return(&sts.GetCallerIdentityOutput{
		Account: aws.String("1234567890"),
	}, nil)

	s := NewStackit(capi, sapi)

	input := StackitUpInput{
		StackName: "stack-name",
		RoleARN:   "MyRoleName",
	}

	//capi.DescribeStacksF = func(input *cloudformation.DescribeStacksInput) (*cloudformation.DescribeStacksOutput, error) {
	//	capi.DescribeStacksF = func(input *cloudformation.DescribeStacksInput) (*cloudformation.DescribeStacksOutput, error) {
	//		return &cloudformation.DescribeStacksOutput{
	//			Stacks: []*cloudformation.Stack{
	//				{
	//					StackId: aws.String("arn:aws:cloudformation:ap-southeast-2:657110686698:stack/stackset-role/58ed6a10-3e2f-11e9-bc5f-0a9966e9c45e"),
	//				},
	//			},
	//		}, nil
	//	}
	//	return nil, awserr.New("ValidationError", "", nil)
	//}

	ch := make(chan TailStackEvent)
	_, err := s.Prepare(context.Background(), input, ch)
	assert.EqualError(t, err, "creating change set: done")
}

func TestServiceRoleArnDoesntTriggerStsCall(t *testing.T) {
	capi := &mockCfn{}
	capi.On("DescribeStacksWithContext", mock.Anything, mock.Anything, mock.Anything).Return(nil, awserr.New("ValidationError", "", nil))
	capi.On("CreateChangeSetWithContext", mock.Anything, mock.Anything, mock.Anything).Return(nil, errors.New("done")).Run(func(args mock.Arguments) {
		input := args.Get(1).(*cloudformation.CreateChangeSetInput)
		assert.Equal(t, "arn:aws:iam::1234567890:role/MyRoleName", *input.RoleARN)
	})

	sapi := &mockSts{}
	sapi.On("GetCallerIdentityWithContext", mock.Anything, mock.Anything, mock.Anything).Return(nil, nil).Run(func(args mock.Arguments) {
		assert.Fail(t, "shouldn't call sts:GetCallerIdentity if full role arn passed")
	})

	s := NewStackit(capi, sapi)

	input := StackitUpInput{
		StackName: "stack-name",
		RoleARN:   "arn:aws:iam::1234567890:role/MyRoleName",
	}

	ch := make(chan TailStackEvent)
	_, err := s.Prepare(context.Background(), input, ch)
	assert.EqualError(t, err, "creating change set: done")
}

func TestNoOpChangesetTriggersDelete(t *testing.T) {
	t.SkipNow()
}

func TestChangesetErrorIsReported(t *testing.T) {
	t.SkipNow()
}
