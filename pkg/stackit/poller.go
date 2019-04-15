package stackit

import (
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/cloudformation"
	"time"
)

type TailStackEvent struct {
	cloudformation.StackEvent
	StackitError error
}

func (s *Stackit) PollStackEvents(token string, callback func(event TailStackEvent)) TailStackEvent {
	lastSentEventId := ""

	for {
		time.Sleep(3 * time.Second)

		events := []*cloudformation.StackEvent{}

		err := s.api.DescribeStackEventsPages(&cloudformation.DescribeStackEventsInput{
			StackName: &s.stackId,
		}, func(page *cloudformation.DescribeStackEventsOutput, lastPage bool) bool {
			for _, event := range page.StackEvents {
				crt := "nil"
				if event.ClientRequestToken != nil {
					crt = *event.ClientRequestToken
				}

				if token == "" {
					token = crt
				}

				if *event.EventId == lastSentEventId || crt != token {
					return false
				}

				events = append(events, event)
			}
			return true
		})

		if err != nil {
			if awsErr, ok := err.(awserr.Error); ok {
				code := awsErr.Code()
				if code == "ThrottlingException" {
					continue
				}
			}
			event := TailStackEvent{cloudformation.StackEvent{}, err}
			callback(event)
			return event
		}

		if len(events) == 0 {
			continue
		}

		lastSentEventId = *events[0].EventId
		stack, err := s.Describe()
		if err != nil {
			event := TailStackEvent{cloudformation.StackEvent{}, err}
			callback(event)
			return event
		}
		terminal := IsTerminalStatus(*stack.StackStatus)

		for ev_i := len(events) - 1; ev_i >= 0; ev_i-- {
			event := events[ev_i]
			tailEvent := TailStackEvent{*event, nil}

			done := terminal && ev_i == 0
			if done {
				return tailEvent
			}

			callback(tailEvent)
		}
	}
}

func IsTerminalStatus(status string) bool {
	switch status {
	case
		"CREATE_COMPLETE",
		"DELETE_COMPLETE",
		"CREATE_FAILED",
		"DELETE_FAILED",
		"ROLLBACK_COMPLETE",
		"ROLLBACK_FAILED",
		"UPDATE_COMPLETE",
		"UPDATE_FAILED",
		"UPDATE_ROLLBACK_COMPLETE",
		"UPDATE_ROLLBACK_FAILED":
		return true
	default:
		return false
	}
}
