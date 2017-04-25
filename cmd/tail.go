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
	"github.com/glassechidna/stackit/stackit"
)

var tailCmd = &cobra.Command{
	Use:   "tail",
	Short: "Tail output of stack change in progress",
	Run: func(cmd *cobra.Command, args []string) {
		region := viper.GetString("region")
		profile := viper.GetString("profile")
		stackName := viper.GetString("stack-name")
		showTimestamps := !viper.GetBool("no-timestamps")
		showColor := !viper.GetBool("no-color")

		cfn := stackit.CfnClient(profile, region)
		stackit.TailStack(&stackName, nil, showTimestamps, showColor, cfn)
	},
}

func init() {
	RootCmd.AddCommand(tailCmd)
}
