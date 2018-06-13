package stackit

import (
	"github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/pkg/errors"
)

func (s *Stackit) Down(events chan<- TailStackEvent) {
	stack, err := s.Describe()

	if stack != nil { // stack exists
		token := generateToken()

		_, err = s.api.DeleteStack(&cloudformation.DeleteStackInput{
			StackName:          &s.stackId,
			ClientRequestToken: &token,
		})
		if err != nil {
			s.error(errors.Wrap(err, "deleting stack"), events)
			return
		}

		s.PollStackEvents(token, func(event TailStackEvent) {
			events <- event
		})
	}

	close(events)
}