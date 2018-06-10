package stackit

import (
	"github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/pkg/errors"
	"github.com/aws/aws-sdk-go/service/cloudformation/cloudformationiface"
)

type StackitUpInput struct {
	RoleARN          string
	StackPolicyBody  string
	TemplateBody     string
	PreviousTemplate bool
	Parameters       []*cloudformation.Parameter
	Tags             map[string]string
	NotificationARNs []string
	PopulateMissing  bool
}

func (s *Stackit) populateMissing(input *StackitUpInput) error {
	stack, _ := s.describe()

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
		resp, err := s.api.ValidateTemplate(&cloudformation.ValidateTemplateInput{TemplateBody: &input.TemplateBody})
		if err != nil {
			return err
		}

		for _, param := range resp.Parameters {
			maybeAddParam(*param.ParameterKey)
		}
	}

	return nil
}

func (s *Stackit) waitForCleanStack(events chan<- TailStackEvent) {
	stack, err := s.describe()
	if err != nil {
		s.error(err, events)
	}

	if stack != nil { // stack already exists
		if *stack.StackStatus == "CREATE_FAILED" || *stack.StackStatus == "ROLLBACK_COMPLETE" {
			token := generateToken()
			s.api.DeleteStack(&cloudformation.DeleteStackInput{StackName: &s.stackId, ClientRequestToken: &token})
			s.PollStackEvents(token, events)
		}
	}
}

func (s *Stackit) Up(input StackitUpInput, events chan<- TailStackEvent) {
	s.waitForCleanStack(events)
	stack, err := s.describe()
	if err != nil {
		s.error(err, events)
	}

	if stack != nil { // stack already exists
		if input.PopulateMissing {
			err := s.populateMissing(&input)
			if err != nil {
				s.error(errors.Wrap(err, "populating missing parameters"), events)
				return
			}
		}
		token, err := s.updateStack(input, events)
		if err != nil {
			s.error(errors.Wrap(err, "updating stack"), events)
			return
		} else if token == "" { // update is a no-op, nothing to change
			close(events)
			return
		}
		s.PollStackEvents(token, events)
	} else {
		token, err := s.createStack(input, events)
		if err != nil {
			s.error(errors.Wrap(err, "creating stack"), events)
			return
		}
		s.PollStackEvents(token, events)
	}
}

func (s *Stackit) updateStack(input StackitUpInput, events chan<- TailStackEvent) (string, error) {
	token := generateToken()

	updateInput := &cloudformation.UpdateStackInput{
		StackName:           &s.stackId,
		Capabilities:        aws.StringSlice([]string{"CAPABILITY_IAM", "CAPABILITY_NAMED_IAM"}),
		UsePreviousTemplate: &input.PreviousTemplate,
		Parameters:          input.Parameters,
		Tags:                mapToTags(input.Tags),
		NotificationARNs:    aws.StringSlice(input.NotificationARNs),
		ClientRequestToken:  &token,
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
	_, err := s.api.UpdateStack(updateInput)

	if err != nil {
		if awsErr, ok := err.(awserr.Error); ok {
			if awsErr.Code() == "ValidationError" && awsErr.Message() == "No updates are to be performed." {
				return "", nil
			}
		}
		return "", err
	}

	return token, nil
}

func (s *Stackit) createStack(input StackitUpInput, events chan<- TailStackEvent) (string, error) {
	token := generateToken()

	createInput := &cloudformation.CreateStackInput{
		StackName:          &s.stackName,
		Capabilities:       aws.StringSlice([]string{"CAPABILITY_IAM", "CAPABILITY_NAMED_IAM"}),
		TemplateBody:       &input.TemplateBody,
		Parameters:         input.Parameters,
		Tags:               mapToTags(input.Tags),
		NotificationARNs:   aws.StringSlice(input.NotificationARNs),
		ClientRequestToken: &token,
	}
	if len(input.RoleARN) > 0 {
		createInput.RoleARN = &input.RoleARN
	}
	if len(input.StackPolicyBody) > 0 {
		createInput.StackPolicyBody = &input.StackPolicyBody
	}

	resp, err := s.api.CreateStack(createInput)
	if err != nil {
		return "", err
	} else {
		s.stackId = *resp.StackId
		return token, nil
	}
}

   type mockCloudFormationClient struct {
       cloudformationiface.CloudFormationAPI
   }

