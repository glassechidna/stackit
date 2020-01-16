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
	"fmt"
	"github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/glassechidna/stackit/cmd/honey"
	"github.com/glassechidna/stackit/pkg/stackit"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"io/ioutil"
)

var transformCmd = &cobra.Command{
	Use:   "transform",
	Short: "See processed form of a template with transforms",
	Run: func(cmd *cobra.Command, args []string) {
		region := viper.GetString("region")
		profile := viper.GetString("profile")
		templatePath, _ := cmd.PersistentFlags().GetString("template")

		params := keyvalSliceToMap(args)

		sess := awsSession(profile, region)
		sit := stackit.NewStackit(
			cloudformation.New(sess),
			sts.New(sess),
		)

		original, err := ioutil.ReadFile(templatePath)
		if err != nil {
			panic(err)
		}

		ctx, end := honey.RootContext()
		defer end()

		processed, err := sit.Transform(ctx, string(original), params)
		if err != nil {
			panic(err)
		}

		fmt.Println(*processed)
	},
}

func init() {
	RootCmd.AddCommand(transformCmd)
	transformCmd.PersistentFlags().String("template", "", "")
}
