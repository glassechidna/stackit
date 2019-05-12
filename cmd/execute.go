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
	"github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/glassechidna/stackit/pkg/stackit"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"io"
)

func executeChangeSet(ctx context.Context, region, profile, stackName, changeSet string, writer io.Writer) error {
	sess := awsSession(profile, region)
	sit := stackit.NewStackit(cloudformation.New(sess), sts.New(sess))
	events := make(chan stackit.TailStackEvent)
	printer := stackit.NewTailPrinterWithOptions(true, true, writer)

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

	var err error
	if len(changeSet) == 0 {
		changeSet, err = findChangeSetIdLocally()
		if err != nil {
			return errors.Wrap(err, "trying to find changeset from filesystem when one wasn't provided to CLI")
		}
	}

	return sit.Execute(ctx, stackName, changeSet, events)
}

func findChangeSetIdLocally() (string, error) {
	prepared, err := loadPreparedOutput()
	if err != nil {
		return "", err
	}

	return *prepared.Output.Id, nil
}

func init() {
	cmd := &cobra.Command{
		Use:   "execute",
		Short: "Execute a pre-created changeset",
		Run: func(cmd *cobra.Command, args []string) {
			region := viper.GetString("region")
			profile := viper.GetString("profile")
			stackName := viper.GetString("stack-name")
			changeSet := viper.GetString("change-set")

			err := executeChangeSet(context.Background(), region, profile, stackName, changeSet, cmd.OutOrStderr())
			if err != nil {
				panic(err)
			}
		},
	}

	RootCmd.AddCommand(cmd)
}
