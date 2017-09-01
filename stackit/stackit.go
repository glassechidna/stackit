package stackit

import (
	"github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/aws/credentials/stscreds"
	"github.com/aws/aws-sdk-go/aws"
	"log"
	"fmt"
	"os"
	"os/signal"
	"encoding/json"
)

func AwsSession(profile, region string) *session.Session {
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
		sess.Config = config
	}

	return sess
}

func CancelOnInterrupt(sess *session.Session, stackId *string, isNewStackCreation bool) {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt)

	cfn := cloudformation.New(sess)

	go func() {
		<- sigs

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

func PrintOutputs(sess *session.Session, stackId *string) {
	cfn := cloudformation.New(sess)

	resp, err := cfn.DescribeStacks(&cloudformation.DescribeStacksInput{
		StackName: stackId,
	})

	if err != nil {
		log.Fatal(err.Error())
	}

	outputMap := make(map[string]string)

	for _, output := range resp.Stacks[0].Outputs {
		outputMap[*output.OutputKey] = *output.OutputValue
	}

	bytes, err := json.MarshalIndent(outputMap, "", "  ")
	fmt.Println(string(bytes))
}

func Down(sess *session.Session, stackName string) (*string, *string) {
	cfn := cloudformation.New(sess)

	_, err := cfn.DescribeStacks(&cloudformation.DescribeStacksInput{StackName: &stackName})
	stackExists := err == nil

	if stackExists {
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

		return stackId, mostRecentEventIdSeen

	}

	return nil, nil
}
