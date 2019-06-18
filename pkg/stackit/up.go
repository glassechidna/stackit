package stackit

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/davecgh/go-spew/spew"
	"github.com/glassechidna/stackit/pkg/stackit/changeset"
	"github.com/pkg/errors"
	"strings"
	"time"
)

type Template interface {
	fmt.Stringer
	Name() string
}

type StackitUpInput struct {
	StackName        string
	RoleARN          string
	StackPolicyBody  string
	Template         Template
	PreviousTemplate bool
	Parameters       []*cloudformation.Parameter
	Tags             map[string]string
	NotificationARNs []string
	PopulateMissing  bool
}

func (s *Stackit) populateMissing(input *StackitUpInput) error {
	stack, _ := s.Describe(input.StackName)

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

	if input.Template == nil {
		input.PreviousTemplate = true

		for _, param := range stack.Parameters {
			maybeAddParam(param.ParameterKey, nil)
		}
	} else {
		resp, err := s.api.ValidateTemplate(&cloudformation.ValidateTemplateInput{TemplateBody: aws.String(input.Template.String())})
		if err != nil {
			return err
		}

		for _, param := range resp.Parameters {
			maybeAddParam(param.ParameterKey, param.DefaultValue)
		}
	}

	return nil
}

func (s *Stackit) ensureStackReady(ctx context.Context, stackName string, events chan<- TailStackEvent) error {
	stack, err := s.Describe(stackName)
	if err != nil {
		return err
	}

	cleanup := func(stackId string) error {
		token := generateToken()
		_, err := s.api.DeleteStack(&cloudformation.DeleteStackInput{StackName: &stackId, ClientRequestToken: &token})
		if err != nil {
			return err
		}

		_, err = s.PollStackEvents(ctx, stackId, token, func(event TailStackEvent) {
			events <- event
		})
		return err
	}

	if stack != nil { // stack already exists
		if !IsTerminalStatus(*stack.StackStatus) && *stack.StackStatus != "REVIEW_IN_PROGRESS" {
			_, err = s.PollStackEvents(ctx, *stack.StackId, "", func(event TailStackEvent) {
				events <- event
			})
			if err != nil {
				return err
			}
		}

		stack, err = s.Describe(*stack.StackId)
		if err != nil {
			return err
		}

		if *stack.StackStatus == "CREATE_FAILED" || *stack.StackStatus == "ROLLBACK_COMPLETE" {
			return cleanup(*stack.StackId)
		} else if *stack.StackStatus == "REVIEW_IN_PROGRESS" {
			resp, err := s.api.ListStackResources(&cloudformation.ListStackResourcesInput{StackName: stack.StackId})
			if err != nil {
				return err
			}
			if len(resp.StackResourceSummaries) == 0 {
				return cleanup(*stack.StackId)
			}
		}
	}

	return nil
}

func (s *Stackit) awsAccountId() (string, error) {
	resp, err := s.stsApi.GetCallerIdentity(&sts.GetCallerIdentityInput{})
	if err != nil {
		return "", errors.Wrap(err, "getting aws account id")
	}
	return *resp.Account, nil
}

type PrepareOutput struct {
	Input        *cloudformation.CreateChangeSetInput
	Output       *cloudformation.CreateChangeSetOutput
	Changes      []*cloudformation.Change
	TemplateBody string
}

func (s *Stackit) Prepare(ctx context.Context, input StackitUpInput, events chan<- TailStackEvent) (*PrepareOutput, error) {
	err := s.ensureStackReady(ctx, input.StackName, events)
	if err != nil {
		return nil, errors.Wrap(err, "waiting for stack to be in a clean state")
	}

	stack, err := s.Describe(input.StackName)
	if err != nil {
		return nil, errors.Wrap(err, "describing stack")
	}

	if input.PopulateMissing && stack != nil {
		err := s.populateMissing(&input)
		if err != nil {
			return nil, errors.Wrap(err, "populating missing parameters")
		}
	}

	token := generateToken()

	createInput := &cloudformation.CreateChangeSetInput{
		ChangeSetName:       aws.String(fmt.Sprintf("%s-csid-%d", input.StackName, time.Now().Unix())),
		StackName:           &input.StackName,
		Capabilities:        aws.StringSlice([]string{"CAPABILITY_IAM", "CAPABILITY_NAMED_IAM"}),
		Parameters:          input.Parameters,
		Tags:                mapToTags(input.Tags),
		NotificationARNs:    aws.StringSlice(input.NotificationARNs),
		ClientToken:         &token,
		UsePreviousTemplate: &input.PreviousTemplate,
	}

	if input.Template != nil {
		createInput.TemplateBody = aws.String(input.Template.String())
	}

	if roleArn := input.RoleARN; len(roleArn) > 0 {
		if !strings.HasPrefix(roleArn, "arn:aws:iam") {
			accountId, err := s.awsAccountId()
			if err != nil {
				return nil, errors.Wrap(err, "retrieving aws account id from sts")
			}
			roleArn = fmt.Sprintf("arn:aws:iam::%s:role/%s", accountId, roleArn)
		}
		createInput.RoleARN = &roleArn
	}

	if stack != nil { // stack already exists
		createInput.ChangeSetType = aws.String(cloudformation.ChangeSetTypeUpdate)
	} else {
		createInput.ChangeSetType = aws.String(cloudformation.ChangeSetTypeCreate)
	}

	resp, err := s.api.CreateChangeSet(createInput)
	if err != nil {
		return nil, errors.Wrap(err, "creating change set")
	}

	change, err := changeset.Wait(s.api, *resp.Id)
	if _, ok := err.(*changeset.NoOpChangesetError); ok {
		_, err = s.api.DeleteChangeSet(&cloudformation.DeleteChangeSetInput{ChangeSetName: resp.Id})
		return nil, errors.Wrap(err, "waiting for no-op changeset to delete")
	}

	if err != nil {
		spew.Dump(err)
		return nil, errors.Wrap(err, "waiting for changeset to stabilise")
	}

	getResp, err := s.api.GetTemplate(&cloudformation.GetTemplateInput{
		ChangeSetName: resp.Id,
		StackName:     resp.StackId,
		TemplateStage: aws.String(cloudformation.TemplateStageProcessed),
	})
	if err != nil {
		return nil, errors.Wrap(err, "getting processed template body")
	}

	return &PrepareOutput{
		Input:        createInput,
		Output:       resp,
		Changes:      change.Changes,
		TemplateBody: *getResp.TemplateBody,
	}, nil
}

func (s *Stackit) Execute(ctx context.Context, stackId, changeSetId string, events chan<- TailStackEvent) error {
	token := generateToken()

	_, err := s.api.ExecuteChangeSet(&cloudformation.ExecuteChangeSetInput{
		ChangeSetName:      &changeSetId,
		ClientRequestToken: &token,
	})

	if err != nil {
		return errors.Wrap(err, "executing change set")
	}

	_, err = s.PollStackEvents(ctx, stackId, token, func(event TailStackEvent) {
		events <- event
	})

	return nil
}
