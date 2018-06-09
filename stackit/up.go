package stackit

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws"
	"time"
	"github.com/pkg/errors"
	"github.com/davecgh/go-spew/spew"
)

type StackitUpInput struct {
	StackName        string
	RoleARN          string
	StackPolicyBody  string
	TemplateBody     string
	PreviousTemplate bool
	Parameters       []*cloudformation.Parameter
	Tags             map[string]string
	NotificationARNs []string
	PopulateMissing  bool
}

func mapToTags(tagMap map[string]string) []*cloudformation.Tag {
	tags := []*cloudformation.Tag{}

	for key, val := range tagMap {
		tags = append(tags, &cloudformation.Tag{Key: aws.String(key), Value: aws.String(val)})
	}

	return tags
}

func populateMissing(sess *session.Session, input *StackitUpInput, stack *cloudformation.Stack) error {
	maybeAddParam := func(name string) {
		for _, param := range input.Parameters {
			if *param.ParameterKey == name {
				return
			}
		}
		input.Parameters = append(input.Parameters, &cloudformation.Parameter{
			ParameterKey:     &name,
			UsePreviousValue: aws.Bool(true),
		})
	}

	if len(input.TemplateBody) == 0 {
		input.PreviousTemplate = true

		for _, param := range stack.Parameters {
			maybeAddParam(*param.ParameterKey)
		}
	} else {
		api := cloudformation.New(sess)
		resp, err := api.ValidateTemplate(&cloudformation.ValidateTemplateInput{TemplateBody: &input.TemplateBody})
		if err != nil {
			return err
		}

		for _, param := range resp.Parameters {
			maybeAddParam(*param.ParameterKey)
		}
	}

	return nil
}

func CleanStackExists(sess *session.Session, name string) (bool, *cloudformation.Stack) {
	cfn := cloudformation.New(sess)

	resp, err := cfn.DescribeStacks(&cloudformation.DescribeStacksInput{StackName: &name})
	stackExists := err == nil

	if resp != nil && len(resp.Stacks) > 0 {
		stack := resp.Stacks[0]
		if *stack.StackStatus == "CREATE_FAILED" || *stack.StackStatus == "ROLLBACK_COMPLETE" {
			cfn.DeleteStack(&cloudformation.DeleteStackInput{StackName: &name})
			stackExists = false
			time.Sleep(time.Duration(3) * time.Second) // wait for cloudformation to register stack deletion
		}
	}

	if stackExists {
		return true, resp.Stacks[0]
	} else {
		return false, nil
	}
}

func Up(sess *session.Session, input StackitUpInput, events chan<- TailStackEvent) (string, error) {
	stackExists, stack := CleanStackExists(sess, input.StackName)

	if stackExists {
		if input.PopulateMissing {
			err := populateMissing(sess, &input, stack)
			if err != nil {
				return "", errors.Wrap(err, "populating missing parameters")
			}
		}
		return updateStack(sess, input, events)
	} else {
		return createStack(sess, input, events)
	}
}

func updateStack(sess *session.Session, input StackitUpInput, events chan<- TailStackEvent) (string, error) {
	cfn := cloudformation.New(sess)

	describeResp, err := cfn.DescribeStackEvents(&cloudformation.DescribeStackEventsInput{
		StackName: &input.StackName,
	})

	if err != nil {
		return "", err
	}

	updateInput := &cloudformation.UpdateStackInput{
		StackName:           &input.StackName,
		Capabilities:        aws.StringSlice([]string{"CAPABILITY_IAM", "CAPABILITY_NAMED_IAM"}),
		UsePreviousTemplate: &input.PreviousTemplate,
		Parameters:          input.Parameters,
		Tags:                mapToTags(input.Tags),
		NotificationARNs:    aws.StringSlice(input.NotificationARNs),
	}
	if len(input.RoleARN) > 0 {
		updateInput.RoleARN = &input.RoleARN
	}
	if len(input.StackPolicyBody) > 0 {
		updateInput.StackPolicyDuringUpdateBody = &input.StackPolicyBody
	}
	if len(input.TemplateBody) > 0 {
		updateInput.TemplateBody = &input.TemplateBody
	}
	_, err = cfn.UpdateStack(updateInput)

	event := describeResp.StackEvents[0]
	eventIdToTail := event.EventId
	stackId := *event.StackId

	if err != nil {
		if awsErr, ok := err.(awserr.Error); ok {
			if awsErr.Code() == "ValidationError" && awsErr.Message() == "No updates are to be performed." {
				close(events)
				return stackId, nil
			}
		}
		spew.Dump(err)
		return stackId, err
	}

	return stackId, PollStackEvents(sess, stackId, eventIdToTail, events)
}

func createStack(sess *session.Session, input StackitUpInput, events chan<- TailStackEvent) (string, error) {
	cfn := cloudformation.New(sess)

	createInput := &cloudformation.CreateStackInput{
		StackName:        &input.StackName,
		Capabilities:     aws.StringSlice([]string{"CAPABILITY_IAM", "CAPABILITY_NAMED_IAM"}),
		TemplateBody:     &input.TemplateBody,
		Parameters:       input.Parameters,
		Tags:             mapToTags(input.Tags),
		NotificationARNs: aws.StringSlice(input.NotificationARNs),
	}
	if len(input.RoleARN) > 0 {
		createInput.RoleARN = &input.RoleARN
	}
	if len(input.StackPolicyBody) > 0 {
		createInput.StackPolicyBody = &input.StackPolicyBody
	}
	resp, err := cfn.CreateStack(createInput)

	if err != nil {
		spew.Dump(err)
		return "", err
	} else {
		eventId := ""
		stackId := *resp.StackId
		return stackId, PollStackEvents(sess, stackId, &eventId, events)
	}
}
