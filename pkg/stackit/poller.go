package stackit

import (
	"context"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/aws/aws-sdk-go/service/cloudformation/cloudformationiface"
	"time"
)

type TailStackEvent struct {
	cloudformation.StackEvent
}

func eventsWhile(api cloudformationiface.CloudFormationAPI, stackId string, include func(event *cloudformation.StackEvent) bool) ([]*cloudformation.StackEvent, error) {
	var events []*cloudformation.StackEvent

	err := api.DescribeStackEventsPages(&cloudformation.DescribeStackEventsInput{
		StackName: &stackId,
	}, func(page *cloudformation.DescribeStackEventsOutput, lastPage bool) bool {
		for _, event := range page.StackEvents {
			if !include(event) {
				return false
			}

			events = append(events, event)
		}

		return true
	})

	return events, err
}

func (s *Stackit) PollStackEvents(ctx context.Context, stackId, token string, callback func(event TailStackEvent)) (*TailStackEvent, error) {
	var mostRecent *time.Time
	tick := time.NewTicker(3 * time.Second)

	for {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-tick.C:
			var events []*cloudformation.StackEvent
			var err error
			if mostRecent == nil {
				events, err = eventsWhile(s.api, stackId, func(event *cloudformation.StackEvent) bool {
					return event.ClientRequestToken != nil && *event.ClientRequestToken == token
				})
			} else {
				events, err = eventsWhile(s.api, stackId, func(event *cloudformation.StackEvent) bool {
					return event.Timestamp.After(*mostRecent)
				})
			}

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

			mostRecent = events[0].Timestamp

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
