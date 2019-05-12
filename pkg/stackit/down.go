package stackit

import (
	"context"
	"github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/pkg/errors"
)

func (s *Stackit) Down(ctx context.Context, stackName string, events chan<- TailStackEvent) error {
	stack, err := s.Describe(stackName)

	if stack != nil { // stack exists
		token := generateToken()

		input := &cloudformation.DeleteStackInput{
			StackName:          stack.StackId,
			ClientRequestToken: &token,
		}
		_, err = s.api.DeleteStack(input)
		if err != nil {
			return err
		}

		finalEvent, err := s.PollStackEvents(ctx, *stack.StackId, token, func(event TailStackEvent) {
			events <- event
		})
		if err != nil {
			return err
		}

		if *finalEvent.ResourceStatus == cloudformation.ResourceStatusDeleteFailed {
			token = generateToken()
			input.ClientRequestToken = &token
			input.RetainResources, err = s.resourcesToBeRetainedDuringDelete(*stack.StackId, events)
			if err != nil {
				return errors.Wrap(err, "determining resources to be kept")
			}

			_, err = s.api.DeleteStack(input)
			if err != nil {
				return errors.Wrap(err, "deleting stack")
			}

			_, err = s.PollStackEvents(ctx, *stack.StackId, token, func(event TailStackEvent) {
				events <- event
			})
			if err != nil {
				return errors.Wrap(err, "deleting stack")
			}
		}
	}

	return nil
}

func (s *Stackit) resourcesToBeRetainedDuringDelete(stackName string, events chan<- TailStackEvent) ([]*string, error) {
	var names []*string

	err := s.api.ListStackResourcesPages(&cloudformation.ListStackResourcesInput{StackName: &stackName}, func(page *cloudformation.ListStackResourcesOutput, lastPage bool) bool {
		for _, resource := range page.StackResourceSummaries {
			if *resource.ResourceStatus == cloudformation.ResourceStatusDeleteFailed {
				names = append(names, resource.LogicalResourceId)
			}
		}
		return !lastPage
	})

	return names, errors.Wrap(err, "listing stack resources")
}
