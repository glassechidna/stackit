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
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/glassechidna/stackit/pkg/stackit"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"io/ioutil"
	"os"
	"strings"
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
		alwaysSucceed := viper.GetBool("always-succeed")

		parsed := parseCLIInput(
			serviceRole,
			stackPolicy,
			template,
			paramValues,
			previousParamValues,
			tags,
			notificationArns,
			previousTemplate)

		parsed.StackName = stackName

		events := make(chan stackit.TailStackEvent)

		sess := awsSession(profile, region)
		sit := stackit.NewStackit(cloudformation.New(sess), sts.New(sess))

		ctx := context.Background()

		printer := stackit.NewTailPrinter(cmd.OutOrStderr())
		printerCtx, printerCancel := context.WithCancel(ctx)
		defer printerCancel()

		go func() {
			for {
				select {
				case <-printerCtx.Done():
					return
				case tailEvent := <-events:
					printer.PrintTailEvent(tailEvent)
				}
			}
		}()

		prepared, err := sit.Prepare(ctx, parsed, events)
		if err != nil {
			panic(err)
		}

		if prepared == nil {
			return // no-op change set
		}

		err = sit.Execute(ctx, *prepared.Output.StackId, *prepared.Output.Id, events)
		if err != nil {
			panic(err)
		}

		stackId := *prepared.Output.StackId
		if success, _ := sit.IsSuccessfulState(stackId); !success && !alwaysSucceed {
			os.Exit(1)
		}

		sit.PrintOutputs(stackId, cmd.OutOrStdout())
	},
}

func keyvalSliceToMap(slice []string) map[string]string {
	theMap := map[string]string{}

	for _, paramPair := range slice {
		parts := strings.SplitN(paramPair, "=", 2)
		if len(parts) != 2 {
			fmt.Fprintf(os.Stderr, `ignoring unexpected key-value pair "%s"`+"\n", paramPair)
			continue
		}
		name, value := strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1])
		theMap[name] = value
	}

	return theMap
}

func parseCLIInput(
	serviceRole,
	stackPolicy,
	template string,
	cliParamValues,
	previousParamValues,
	tags,
	notificationArns []string,
	previousTemplate bool) stackit.StackitUpInput {
	input := stackit.StackitUpInput{
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

	paramMap := keyvalSliceToMap(viper.GetStringSlice("parameters"))
	for key, val := range keyvalSliceToMap(cliParamValues) {
		paramMap[key] = val
	}

	params := []*cloudformation.Parameter{}
	for name, value := range paramMap {
		params = append(params, &cloudformation.Parameter{
			ParameterKey:   aws.String(name),
			ParameterValue: aws.String(value),
		})
	}

	for _, param := range previousParamValues {
		params = append(params, &cloudformation.Parameter{
			ParameterKey:     aws.String(param),
			UsePreviousValue: aws.Bool(true),
		})
	}

	input.Parameters = params
	input.NotificationARNs = notificationArns

	if len(tags) > 0 {
		input.Tags = keyvalSliceToMap(tags)
	}

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
	upCmd.PersistentFlags().Bool("always-succeed", false, "Typically stackit will return a nonzero exit code on failure. This disables that.")

	viper.BindPFlags(upCmd.PersistentFlags())
}
