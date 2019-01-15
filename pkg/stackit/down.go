package stackit

import (
	"github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/pkg/errors"
)

func (s *Stackit) Down(events chan<- TailStackEvent) {
	stack, err := s.Describe()

	if stack != nil { // stack exists
		token := generateToken()

		input := &cloudformation.DeleteStackInput{
			StackName:          &s.stackId,
			ClientRequestToken: &token,
		}
		_, err = s.api.DeleteStack(input)
		if err != nil {
			s.error(errors.Wrap(err, "deleting stack"), events)
			return
		}

		finalEvent := s.PollStackEvents(token, func(event TailStackEvent) {
			events <- event
		})

		if *finalEvent.ResourceStatus == cloudformation.ResourceStatusDeleteFailed {
			token = generateToken()
			input.ClientRequestToken = &token
			input.RetainResources = s.resourcesToBeRetainedDuringDelete(events)
			_, err = s.api.DeleteStack(input)
			if err != nil {
				s.error(errors.Wrap(err, "deleting stack"), events)
				return
			}

			s.PollStackEvents(token, func(event TailStackEvent) {
				events <- event
			})
		}
	}

	close(events)
}

func (s *Stackit) resourcesToBeRetainedDuringDelete(events chan<- TailStackEvent) []*string {
	names := []*string{}

	err := s.api.ListStackResourcesPages(&cloudformation.ListStackResourcesInput{StackName: &s.stackId}, func(page *cloudformation.ListStackResourcesOutput, lastPage bool) bool {
		for _, resource := range page.StackResourceSummaries {
			if *resource.ResourceStatus == cloudformation.ResourceStatusDeleteFailed {
				names = append(names, resource.LogicalResourceId)
			}
		}
		return !lastPage
	})
	if err != nil {
		s.error(errors.Wrap(err, "listing stack resources"), events)
		return nil
	}

	return names
}
