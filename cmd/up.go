// Copyright © 2017 Aidan Steele <aidan.steele@glassechidna.com.au>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cloudformation"
	"io/ioutil"
	"strings"
	"github.com/glassechidna/stackit/stackit"
	"os"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/fatih/color"
	"github.com/aws/aws-sdk-go/aws/session"
)

// up --stack-name stackit-test --template sample.yml --param-value DockerImage=nginx --param-value Cluster=app-cluster-Cluster-1C2I18JXK9QNM --tag MyTag=Cool

var paramValues []string
var previousParamValues []string
var tags []string
var notificationArns []string

var upCmd = &cobra.Command{
	Use:   "up",
	Short: "Bring stack up to date",
	Run: func(cmd *cobra.Command, args []string) {
		region := viper.GetString("region")
		profile := viper.GetString("profile")
		stackName := viper.GetString("stack-name")

		serviceRole := viper.GetString("service-role")
		stackPolicy := viper.GetString("stack-policy")
		template := viper.GetString("template")
		previousTemplate := viper.GetBool("previous-template")
		//noDestroy := viper.GetBool("no-destroy")
		//cancelOnExit := !viper.GetBool("no-cancel-on-exit")

		showTimestamps := !viper.GetBool("no-timestamps")
		showColor := !viper.GetBool("no-color")
		printer := stackit.NewTailPrinterWithOptions(showTimestamps, showColor)

		parsed := parseCLIInput(
			stackName,
			serviceRole,
			stackPolicy,
			template,
			paramValues,
			previousParamValues,
			tags,
			notificationArns,
			previousTemplate)

		sess := stackit.AwsSession(profile, region)
		output, err := stackit.Up(sess, parsed)

		if err != nil {
			if awsErr, ok := err.(awserr.Error); ok {
				color.New(color.FgRed).Fprintf(os.Stderr, "%s: %s\n", awsErr.Code(), awsErr.Message())
			} else {
				color.New(color.FgRed).Fprintln(os.Stderr, err.Error())
			}
			os.Exit(1)
		}

		if !output.NoOp {
			for {
				tailEvent := <- *output.Channel
				printer.PrintTailEvent(tailEvent)
				if tailEvent.Done {
					break
				}
			}
		}

		//
		//if cancelOnExit {
		//	stackit.CancelOnInterrupt(sess, stackId, isNewStack)
		//}

		exitIfFailedUpdate(sess, output.StackId)
		stackit.PrintOutputs(sess, &output.StackId)
	},
}

func exitIfFailedUpdate(sess *session.Session, stackId string) {
	cfn := cloudformation.New(sess)
	resp, _ := cfn.DescribeStacks(&cloudformation.DescribeStacksInput{StackName: &stackId})
	status := *resp.Stacks[0].StackStatus

	if status != "CREATE_COMPLETE" && status != "UPDATE_COMPLETE" {
		os.Exit(1)
	}
}

func parseCLIInput(
	stackName,
	serviceRole,
	stackPolicy,
	template string,
	cliParamValues,
	previousParamValues,
	tags,
	notificationArns []string,
	previousTemplate bool) stackit.StackitUpInput {
	input := stackit.StackitUpInput{
		StackName: stackName,
		PopulateMissing: true,
	}

	if len(serviceRole) > 0 {
		input.RoleARN = serviceRole
	}

	if len(stackPolicy) > 0 {
		policyBody, err := ioutil.ReadFile(stackPolicy)
		if err != nil {

		} else {
			input.StackPolicyBody = string(policyBody)
		}
	}

	if len(template) > 0 {
		templateBody, err := ioutil.ReadFile(template)
		if err != nil {

		} else {
			input.TemplateBody = string(templateBody)
		}
	}

	input.PreviousTemplate = previousTemplate

	paramMap := make(map[string]string)

	populateParamMap := func(slice []string) {
		for _, paramPair := range slice {
			parts := strings.SplitN(paramPair, "=", 2)
			name, value := parts[0], parts[1]
			paramMap[name] = value
		}
	}

	configFileParameters := viper.GetStringSlice("parameters")
	populateParamMap(configFileParameters)
	populateParamMap(cliParamValues)

	params := []*cloudformation.Parameter{}
	for name, value := range paramMap {
		params = append(params, &cloudformation.Parameter{
			ParameterKey: aws.String(name),
			ParameterValue: aws.String(value),
		})
	}

	for _, param := range previousParamValues {
		params = append(params, &cloudformation.Parameter{
			ParameterKey: aws.String(param),
			UsePreviousValue: aws.Bool(true),
		})
	}

	input.Parameters = params

	if len(tags) > 0 {
		cfnTags := []*cloudformation.Tag{}

		for _, tagPair := range tags {
			parts := strings.SplitN(tagPair, "=", 2)
			name, value := parts[0], parts[1]

			cfnTags = append(cfnTags, &cloudformation.Tag{
				Key: aws.String(name),
				Value: aws.String(value),
			})
		}

		input.Tags = cfnTags
	}

	if len(notificationArns) > 0 {
		cfnNotificationArns := []*string{}

		for _, notificationArn := range notificationArns {
			cfnNotificationArns = append(cfnNotificationArns, aws.String(notificationArn))
		}

		input.NotificationARNs = cfnNotificationArns
	}

	input.Capabilities = aws.StringSlice([]string{"CAPABILITY_IAM", "CAPABILITY_NAMED_IAM"})

	return input
}

func init() {
	RootCmd.AddCommand(upCmd)

	upCmd.PersistentFlags().String("service-role", "", "")
	upCmd.PersistentFlags().String("stack-policy", "", "")
	upCmd.PersistentFlags().String("template", "", "")
	upCmd.PersistentFlags().StringArrayVar(&paramValues, "param-value", []string{}, "")
	upCmd.PersistentFlags().StringArrayVar(&previousParamValues, "previous-param-value", []string{}, "")
	upCmd.PersistentFlags().StringArrayVar(&tags, "tag", []string{}, "")
	upCmd.PersistentFlags().StringArrayVar(&notificationArns, "notification-arn", []string{}, "")
	upCmd.PersistentFlags().Bool("previous-template", false, "")
	upCmd.PersistentFlags().Bool("no-destroy", false, "")
	upCmd.PersistentFlags().Bool("no-cancel-on-exit", false, "")
	upCmd.PersistentFlags().Bool("no-timestamps", false, "")
	upCmd.PersistentFlags().Bool("no-color", false, "")

	viper.BindPFlags(upCmd.PersistentFlags())
}
