package stackit

import (
	"github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/aws/credentials/stscreds"
	"github.com/aws/aws-sdk-go/aws"
	"log"
	"fmt"
	"time"
	"os"
	"os/signal"
	"encoding/json"
	"github.com/aws/aws-sdk-go/aws/awserr"
)

type StackitUpInput struct {
	StackName *string
	RoleARN *string
	StackPolicyBody *string
	TemplateBody *string
	PreviousTemplate *bool
	Parameters []*cloudformation.Parameter
	Tags []*cloudformation.Tag
	NotificationARNs []*string
	Capabilities []*string
}

func CfnClient(profile, region string) *cloudformation.CloudFormation {
	sessOpts := session.Options{
		SharedConfigState: session.SharedConfigEnable,
		AssumeRoleTokenProvider: stscreds.StdinTokenProvider,
	}

	if len(profile) > 0 {
		sessOpts.Profile = profile
	}

	sess, _ := session.NewSessionWithOptions(sessOpts)
	config := aws.NewConfig()

	if len(region) > 0 {
		config.Region = aws.String(region)
	}

	return cloudformation.New(sess, config)
}

func Up(region, profile string, input StackitUpInput, noDestroy, cancelOnExit bool) {
	cfn := CfnClient(profile, region)

	stackId, mostRecentEventIdSeen, err := doStackUp(input, cfn)
	shouldTail := true

	if err != nil {
		if awsErr, ok := err.(awserr.Error); ok {
			if awsErr.Code() == "ValidationError" && awsErr.Message() == "No updates are to be performed." {
				shouldTail = false
			} else {
				log.Fatal(err.Error())
			}
		} else {
			log.Fatal(err.Error())
		}
	}

	if shouldTail {
		if cancelOnExit {
			sigs := make(chan os.Signal, 1)
			signal.Notify(sigs, os.Interrupt)

			go func() {
				<- sigs

				isNewStackCreation := mostRecentEventIdSeen == nil

				if isNewStackCreation {
					cfn.DeleteStack(&cloudformation.DeleteStackInput{
						StackName: stackId,
					})
				} else {
					cfn.CancelUpdateStack(&cloudformation.CancelUpdateStackInput{
						StackName: stackId,
					})
				}
			}()
		}

		TailStack(stackId, mostRecentEventIdSeen, cfn)
	}

}

func Down(region, profile, stackName string) {
	cfn := CfnClient(profile, region)

	if stackExists(&stackName, cfn) {
		resp, err := cfn.DescribeStackEvents(&cloudformation.DescribeStackEventsInput{
			StackName: &stackName,
		})

		if err != nil {
			log.Fatal(err.Error())
		}

		stackId := resp.StackEvents[0].StackId
		mostRecentEventIdSeen := resp.StackEvents[0].EventId

		cfn.DeleteStack(&cloudformation.DeleteStackInput{
			StackName: &stackName,
		})

		TailStack(stackId, mostRecentEventIdSeen, cfn)
	}
}

func stackExists(stackName *string, cfn *cloudformation.CloudFormation) bool {
	_, err := cfn.DescribeStacks(&cloudformation.DescribeStacksInput{
		StackName: stackName,
	})

	return err == nil
}

func fixedLengthString(length int, str string) string {
	verb := fmt.Sprintf("%%%d.%ds", length, length)
	return fmt.Sprintf(verb, str)
}

func TailStack(stackId, mostRecentEventIdSeen *string, cfn *cloudformation.CloudFormation) {
	if mostRecentEventIdSeen == nil {
		tmp := ""
		mostRecentEventIdSeen = &tmp
	}

	resourceNameLength := 20 // TODO: determine this from template/API

	for {
		time.Sleep(3*time.Second)

		events := []*cloudformation.StackEvent{}

		cfn.DescribeStackEventsPages(&cloudformation.DescribeStackEventsInput{
			StackName: stackId,
		}, func(page *cloudformation.DescribeStackEventsOutput, lastPage bool) bool {
			for _, event := range page.StackEvents {
				if *event.EventId == *mostRecentEventIdSeen {
					return false
				}

				events = append(events, event)
			}
			return true
		})

		if len(events) == 0 {
			continue
		}

		mostRecentEventIdSeen = events[0].EventId

		for ev_i := len(events) - 1; ev_i >= 0; ev_i-- {
			event := events[ev_i]

			formattedTime := event.Timestamp.Format("[03:04:05]")

			reasonPart := ""
			if event.ResourceStatusReason != nil {
				reasonPart = fmt.Sprintf("- %s", *event.ResourceStatusReason)
			}
			fmt.Fprintf(os.Stderr, "%s %s - %s %s\n", formattedTime, fixedLengthString(resourceNameLength, *event.LogicalResourceId), *event.ResourceStatus, reasonPart)
		}

		resp, err := cfn.DescribeStacks(&cloudformation.DescribeStacksInput{StackName: stackId})

		if err != nil {
			log.Fatal(err.Error())
		}

		switch *resp.Stacks[0].StackStatus {
		case
			"CREATE_COMPLETE",
			"CREATE_FAILED",
			"DELETE_COMPLETE",
			"DELETE_FAILED",
			"ROLLBACK_COMPLETE",
			"ROLLBACK_FAILED",
			"UPDATE_COMPLETE",
			"UPDATE_FAILED",
			"UPDATE_ROLLBACK_COMPLETE",
			"UPDATE_ROLLBACK_FAILED":
			return
		default:
			// no-op
		}
	}
}

func doStackUp(input StackitUpInput, cfn *cloudformation.CloudFormation) (*string, *string, error) {
	if stackExists(input.StackName, cfn) {
		describeResp, err := cfn.DescribeStackEvents(&cloudformation.DescribeStackEventsInput{
			StackName: input.StackName,
		})

		if err != nil {
			return nil, nil, err
		}

		resp, err := cfn.UpdateStack(&cloudformation.UpdateStackInput{
			StackName: input.StackName,
			Capabilities: input.Capabilities,
			RoleARN: input.RoleARN,
			StackPolicyDuringUpdateBody: input.StackPolicyBody,
			TemplateBody: input.TemplateBody,
			UsePreviousTemplate: input.PreviousTemplate,
			Parameters: input.Parameters,
			Tags: input.Tags,
			NotificationARNs: input.NotificationARNs,
		})

		return resp.StackId, describeResp.StackEvents[0].EventId, err
	} else {
		resp, err := cfn.CreateStack(&cloudformation.CreateStackInput{
			StackName: input.StackName,
			Capabilities: input.Capabilities,
			RoleARN: input.RoleARN,
			//StackPolicyBody: input.StackPolicyBody,
			TemplateBody: input.TemplateBody,
			Parameters: input.Parameters,
			Tags: input.Tags,
			NotificationARNs: input.NotificationARNs,
		})

		if err != nil {
			return nil, nil, err
		} else {
			return resp.StackId, nil, nil
		}
	}
}
