package stackit

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudformation"
	"time"
	"log"
	"fmt"
	"strings"
)

type TailStackEvent struct {
	cloudformation.StackEvent
}

func PollStackEvents(sess *session.Session, stackId string, startEventId *string, channel chan<- TailStackEvent) error {
	cfn := cloudformation.New(sess)

	go func() {

		for {
			time.Sleep(3*time.Second)

			events := []*cloudformation.StackEvent{}

			cfn.DescribeStackEventsPages(&cloudformation.DescribeStackEventsInput{
				StackName: &stackId,
			}, func(page *cloudformation.DescribeStackEventsOutput, lastPage bool) bool {
				for _, event := range page.StackEvents {
					if *event.EventId == *startEventId {
						return false
					}

					events = append(events, event)
				}
				return true
			})

			if len(events) == 0 {
				continue
			}

			startEventId = events[0].EventId

			resp, err := cfn.DescribeStacks(&cloudformation.DescribeStacksInput{StackName: &stackId})

			if err != nil {
				log.Fatal(err.Error())
			}

			status := *resp.Stacks[0].StackStatus

			for ev_i := len(events) - 1; ev_i >= 0; ev_i-- {
				done := IsTerminalStatus(status) && ev_i == 0
				if done {
					close(channel)
					return
				}

				event := events[ev_i]
				tailEvent := TailStackEvent{*event}
				channel <- tailEvent
			}
		}
	}()

	return nil
}

func fixedLengthString(length int, str string) string {
	verb := fmt.Sprintf("%%%d.%ds", length, length)
	return fmt.Sprintf(verb, str)
}

func isBadStatus(status string) bool {
	return strings.HasSuffix(status, "_FAILED")
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