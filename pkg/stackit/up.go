package stackit

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/pkg/errors"
	"fmt"
	"time"
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
	stack, _ := s.Describe()

	maybeAddParam := func(name, defaultValue *string) {
		if defaultValue != nil {
			return
		}

		for _, param := range input.Parameters {
			if *param.ParameterKey == *name {
				return
			}
		}
		input.Parameters = append(input.Parameters, &cloudformation.Parameter{
			ParameterKey:     name,
			UsePreviousValue: aws.Bool(true),
		})
	}

	if len(input.TemplateBody) == 0 {
		input.PreviousTemplate = true

		for _, param := range stack.Parameters {
			maybeAddParam(param.ParameterKey, nil)
		}
	} else {
		resp, err := s.api.ValidateTemplate(&cloudformation.ValidateTemplateInput{TemplateBody: &input.TemplateBody})
		if err != nil {
			return err
		}

		for _, param := range resp.Parameters {
			maybeAddParam(param.ParameterKey, param.DefaultValue)
		}
	}

	return nil
}

func (s *Stackit) EnsureStackReady(events chan<- TailStackEvent) error {
	stack, err := s.Describe()
	if err != nil {
		s.error(err, events)
		return err
	}

	cleanup := func() {
		token := generateToken()
		s.api.DeleteStack(&cloudformation.DeleteStackInput{StackName: &s.stackId, ClientRequestToken: &token})
		s.PollStackEvents(token, func(event TailStackEvent) {
			events <- event
		})
	}

	if stack != nil { // stack already exists
		if !IsTerminalStatus(*stack.StackStatus) && *stack.StackStatus != "REVIEW_IN_PROGRESS" {
			s.PollStackEvents("", func(event TailStackEvent) {
				events <- event
			})
		}

		stack, err = s.Describe()
		if err != nil {
			s.error(err, events)
			return err
		}

		if *stack.StackStatus == "CREATE_FAILED" || *stack.StackStatus == "ROLLBACK_COMPLETE" {
			cleanup()
		} else if *stack.StackStatus == "REVIEW_IN_PROGRESS" {
			resp, err := s.api.ListStackResources(&cloudformation.ListStackResourcesInput{StackName: &s.stackId})
			if err != nil {
				s.error(err, events)
				return err
			}
			if len(resp.StackResourceSummaries) == 0 {
				cleanup()
			}
		}
	}

	return nil
}

func (s *Stackit) Up(input StackitUpInput, events chan<- TailStackEvent) {
	s.stackId = ""
	stack, err := s.Describe()
	if err != nil {
		s.error(err, events)
	}

	if input.PopulateMissing && stack != nil {
		err := s.populateMissing(&input)
		if err != nil {
			s.error(errors.Wrap(err, "populating missing parameters"), events)
			return
		}
	}

	token := generateToken()

	createInput := &cloudformation.CreateChangeSetInput{
		ChangeSetName:       aws.String(fmt.Sprintf("csid-%d", time.Now().Unix())),
		StackName:           &s.stackName,
		Capabilities:        aws.StringSlice([]string{"CAPABILITY_IAM", "CAPABILITY_NAMED_IAM"}),
		Parameters:          input.Parameters,
		Tags:                mapToTags(input.Tags),
		NotificationARNs:    aws.StringSlice(input.NotificationARNs),
		ClientToken:         &token,
		UsePreviousTemplate: &input.PreviousTemplate,
	}

	if len(input.TemplateBody) > 0 {
		createInput.TemplateBody = &input.TemplateBody
	}

	if len(input.RoleARN) > 0 {
		createInput.RoleARN = &input.RoleARN
	}

	if stack != nil { // stack already exists
		createInput.ChangeSetType = aws.String(cloudformation.ChangeSetTypeUpdate)
	} else {
		createInput.ChangeSetType = aws.String(cloudformation.ChangeSetTypeCreate)
	}

	resp, err := s.api.CreateChangeSet(createInput)
	if err != nil {
		s.error(errors.Wrap(err, "creating change set"), events)
		return
	}

	s.stackId = *resp.StackId

	change, err := s.waitForChangeset(resp.Id)
	isNoop := change != nil && len(change.Changes) == 0

	if isNoop { // update is a no-op, nothing to change
		s.api.DeleteChangeSet(&cloudformation.DeleteChangeSetInput{ChangeSetName: resp.Id})
	}

	if err != nil {
		s.error(err, events)
		return
	} else if isNoop {
		close(events)
		return
	}

	_, err = s.api.ExecuteChangeSet(&cloudformation.ExecuteChangeSetInput{
		ChangeSetName:      resp.Id,
		ClientRequestToken: &token,
	})
	if err != nil {
		s.error(errors.Wrap(err, "executing change set"), events)
		return
	}

	s.PollStackEvents(token, func(event TailStackEvent) {
		events <- event
	})
	close(events)
}
