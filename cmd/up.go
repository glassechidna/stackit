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
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"io"
	"os"
	"strings"
)

func printUntilDone(ctx context.Context, events <-chan stackit.TailStackEvent, w io.Writer) {
	printer := stackit.NewTailPrinter(w)

	for {
		select {
		case tailEvent := <-events:
			printer.PrintTailEvent(tailEvent)
		case <-ctx.Done():
			return
		}
	}
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

func parseCLIInput(cmd *cobra.Command, args []string) stackit.StackitUpInput {
	stackName, _ := RootCmd.PersistentFlags().GetString("stack-name")
	serviceRole, _ := cmd.PersistentFlags().GetString("service-role")
	template, _ := cmd.PersistentFlags().GetString("template")
	tags, _ := cmd.PersistentFlags().GetStringSlice("tag")
	notificationArns, _ := cmd.PersistentFlags().GetStringSlice("notification-arn")

	input := stackit.StackitUpInput{
		StackName:       stackName,
		PopulateMissing: true,
	}

	if len(serviceRole) > 0 {
		input.RoleARN = serviceRole
	}

	if len(template) > 0 {
		var err error
		input.Template, err = pathToTemplate(template)
		if err != nil {
			panic(err)
			// TODO
		}
	} else {
		input.PreviousTemplate = true
	}

	params := []*cloudformation.Parameter{}
	for name, value := range keyvalSliceToMap(args) {
		params = append(params, &cloudformation.Parameter{
			ParameterKey:   aws.String(name),
			ParameterValue: aws.String(value),
		})
	}

	input.Parameters = params
	input.NotificationARNs = notificationArns

	if len(tags) > 0 {
		input.Tags = keyvalSliceToMap(tags)
	}

	return input
}

var errUnsuccessfulStack = errors.New("stack update unsuccessful")

func up(cmd *cobra.Command, args []string) error {
	region := viper.GetString("region")
	profile := viper.GetString("profile")
	input := parseCLIInput(cmd, args)

	sess := awsSession(profile, region)
	sit := stackit.NewStackit(cloudformation.New(sess), sts.New(sess))

	ctx := context.Background()

	printerCtx, printerCancel := context.WithCancel(ctx)
	defer printerCancel()

	if templateFile, ok := input.Template.(*templateReader); ok && templateFile != nil {
		template, err := packageTemplate(ctx, sess, input.StackName, templateFile, cmd.OutOrStderr())
		if err != nil {
			return errors.Wrap(err, "packaging template")
		}
		templateFile.body = *template
	}

	events := make(chan stackit.TailStackEvent)
	go printUntilDone(printerCtx, events, cmd.OutOrStderr())

	prepared, err := sit.Prepare(ctx, input, events)
	if err != nil {
		return err
	}

	if prepared == nil {
		return nil
	}

	err = sit.Execute(ctx, *prepared.Output.StackId, *prepared.Output.Id, events)
	if err != nil {
		return err
	}

	stackId := *prepared.Output.StackId
	if success, _ := sit.IsSuccessfulState(stackId); !success {
		return errUnsuccessfulStack
	}

	sit.PrintOutputs(stackId, cmd.OutOrStdout())
	return nil
}

func init() {
	upCmd := &cobra.Command{
		Use:   "up",
		Short: "Bring stack up to date",
		Run: func(cmd *cobra.Command, args []string) {
			err := up(cmd, args)
			if err == errUnsuccessfulStack {
				defaultExiter(1)
			} else if err != nil {
				panic(err)
			}
		},
	}
	RootCmd.AddCommand(upCmd)

	upCmd.PersistentFlags().String("service-role", "", "")
	upCmd.PersistentFlags().String("template", "", "")
	upCmd.PersistentFlags().StringSliceP("tag", "t", []string{}, "")
	upCmd.PersistentFlags().StringSlice("notification-arn", []string{}, "")
}

var defaultExiter = os.Exit
