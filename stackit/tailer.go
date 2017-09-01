package stackit

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cloudformation"
	"time"
	"log"
	"fmt"
	"strings"
)

type TailStackEvent struct {
	Event *cloudformation.StackEvent
	Done bool
}

func DoTailStack(sess *session.Session, stackId, startEventId *string) chan TailStackEvent {
	channel := make(chan TailStackEvent)

	if stackId == nil {
		stackId = aws.String("")
	}

	go func() {
		cfn := cloudformation.New(sess)

		for {
			time.Sleep(3*time.Second)

			events := []*cloudformation.StackEvent{}

			cfn.DescribeStackEventsPages(&cloudformation.DescribeStackEventsInput{
				StackName: stackId,
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

			resp, err := cfn.DescribeStacks(&cloudformation.DescribeStacksInput{StackName: stackId})

			if err != nil {
				log.Fatal(err.Error())
			}

			status := *resp.Stacks[0].StackStatus

			for ev_i := len(events) - 1; ev_i >= 0; ev_i-- {
				event := events[ev_i]
				done := isTerminalStatus(status) && ev_i == 0

				tailEvent := TailStackEvent{
					Event: event,
					Done:  done,
				}

				channel <- tailEvent
			}
		}
	}()

	return channel
}

func fixedLengthString(length int, str string) string {
	verb := fmt.Sprintf("%%%d.%ds", length, length)
	return fmt.Sprintf(verb, str)
}

func isBadStatus(status string) bool {
	return strings.HasSuffix(status, "_FAILED")
}

func isTerminalStatus(status string) bool {
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