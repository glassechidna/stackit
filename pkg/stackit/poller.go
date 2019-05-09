package stackit

import (
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/cloudformation"
	"time"
)

type TailStackEvent struct {
	cloudformation.StackEvent
}

func (s *Stackit) PollStackEvents(stackId, token string, callback func(event TailStackEvent)) (*TailStackEvent, error) {
	mostRecentEventTimestamp := time.Now().AddDate(0, 0, -1)
	haveSeenExpectedToken := false

	for {
		time.Sleep(3 * time.Second)

		events := []*cloudformation.StackEvent{}

		err := s.api.DescribeStackEventsPages(&cloudformation.DescribeStackEventsInput{
			StackName: &stackId,
		}, func(page *cloudformation.DescribeStackEventsOutput, lastPage bool) bool {
			for _, event := range page.StackEvents {
				if event.ClientRequestToken != nil && *event.ClientRequestToken == token {
					haveSeenExpectedToken = true
				}

				if haveSeenExpectedToken && event.Timestamp.After(mostRecentEventTimestamp) {
					events = append(events, event)
				}
			}

			earliestEvent := page.StackEvents[len(page.StackEvents)-1]
			shouldPaginate := earliestEvent.Timestamp.After(mostRecentEventTimestamp)
			mostRecentEventTimestamp = *page.StackEvents[0].Timestamp
			return shouldPaginate
		})

		if err != nil {
			if awsErr, ok := err.(awserr.Error); ok {
				code := awsErr.Code()
				if code != "ThrottlingException" {
					return nil, err
				}
			} else {
				return nil, err
			}
		}

		if len(events) == 0 {
			continue
		}

		stack, err := s.Describe(*events[0].StackId)
		if err != nil {
			return nil, err
		}
		for ev_i := len(events) - 1; ev_i >= 0; ev_i-- {
			event := events[ev_i]
			tailEvent := TailStackEvent{*event}
			callback(tailEvent)
		}

		if IsTerminalStatus(*stack.StackStatus) {
			return &TailStackEvent{*events[0]}, nil
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
