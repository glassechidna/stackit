package stackit

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws"
	"time"
)

type StackitUpInput struct {
	StackName *string
	RoleARN *string
	StackPolicyBody *string
	TemplateBody *string
	PreviousTemplate *bool
	Parameters []*cloudformation.Parameter
	Tags []*cloudformation.Tag
	NotificationARNs []*string
	Capabilities []*string
}

type StackUpOutput struct {
	Channel *chan TailStackEvent
	StackId string
	NoOp bool
}

func Up(sess *session.Session, input StackitUpInput) (*StackUpOutput, error) {
	cfn := cloudformation.New(sess)

	resp, err := cfn.DescribeStacks(&cloudformation.DescribeStacksInput{StackName: input.StackName})
	stackExists := err == nil

	if resp != nil && len(resp.Stacks) > 0 {
		stack := resp.Stacks[0]
		if *stack.StackStatus == "CREATE_FAILED" || *stack.StackStatus == "ROLLBACK_COMPLETE" {
			cfn.DeleteStack(&cloudformation.DeleteStackInput{StackName: input.StackName})
			stackExists = false
			time.Sleep(time.Duration(3) * time.Second) // wait for cloudformation to register stack deletion
		}
	}

	if stackExists {
		return updateStack(sess, input)
	} else {
		return createStack(sess, input)
	}
}

func updateStack(sess *session.Session, input StackitUpInput) (*StackUpOutput, error) {
	cfn := cloudformation.New(sess)

	describeResp, err := cfn.DescribeStackEvents(&cloudformation.DescribeStackEventsInput{
		StackName: input.StackName,
	})

	if err != nil {
		return nil, err
	}
	_, err = cfn.UpdateStack(&cloudformation.UpdateStackInput{
		StackName: input.StackName,
		Capabilities: input.Capabilities,
		RoleARN: input.RoleARN,
		StackPolicyDuringUpdateBody: input.StackPolicyBody,
		TemplateBody: input.TemplateBody,
		UsePreviousTemplate: input.PreviousTemplate,
		Parameters: input.Parameters,
		Tags: input.Tags,
		NotificationARNs: input.NotificationARNs,
	})

	event := describeResp.StackEvents[0]
	eventIdToTail := event.EventId
	stackId := *event.StackId

	if err != nil {
		if awsErr, ok := err.(awserr.Error); ok {
			if awsErr.Code() == "ValidationError" && awsErr.Message() == "No updates are to be performed." {
				return &StackUpOutput{
					Channel: nil,
					StackId: stackId,
					NoOp: true,
				}, nil
			}
		}
		return nil, err
	}

	channel := DoTailStack(sess, &stackId, eventIdToTail)

	return &StackUpOutput{
		Channel: &channel,
		StackId: stackId,
		NoOp: false,
	}, nil
}

func createStack(sess *session.Session, input StackitUpInput) (*StackUpOutput, error) {
	cfn := cloudformation.New(sess)

	resp, err := cfn.CreateStack(&cloudformation.CreateStackInput{
		StackName: input.StackName,
		Capabilities: input.Capabilities,
		RoleARN: input.RoleARN,
		//StackPolicyBody: input.StackPolicyBody,
		TemplateBody: input.TemplateBody,
		Parameters: input.Parameters,
		Tags: input.Tags,
		NotificationARNs: input.NotificationARNs,
	})

	if err != nil {
		return nil, err
	} else {
		eventId := ""
		channel := DoTailStack(sess, resp.StackId, aws.String(eventId))

		return &StackUpOutput{
			Channel: &channel,
			StackId: *resp.StackId,
			NoOp: false,
		}, nil
	}
}
