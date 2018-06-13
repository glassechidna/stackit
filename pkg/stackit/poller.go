package stackit

import (
	"github.com/aws/aws-sdk-go/service/cloudformation"
	"time"
)

type TailStackEvent struct {
	cloudformation.StackEvent
	StackitError error
}

func (s *Stackit) PollStackEvents(token string, callback func(event TailStackEvent)) {
	lastSentEventId := ""

	for {
		time.Sleep(3 * time.Second)

		events := []*cloudformation.StackEvent{}

		s.api.DescribeStackEventsPages(&cloudformation.DescribeStackEventsInput{
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

		if len(events) == 0 {
			continue
		}

		lastSentEventId = *events[0].EventId
		stack, err := s.Describe()
		if err != nil {
			callback(TailStackEvent{cloudformation.StackEvent{}, err})
		}
		terminal := IsTerminalStatus(*stack.StackStatus)

		for ev_i := len(events) - 1; ev_i >= 0; ev_i-- {
			done := terminal && ev_i == 0
			if done {
				return
			}

			event := events[ev_i]
			tailEvent := TailStackEvent{*event, nil}
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