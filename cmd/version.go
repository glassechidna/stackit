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
	"github.com/spf13/cobra"
)

// these are set by goreleaser
var version, commit, date string

func init() {
	cmd := &cobra.Command{
		Use:   "version",
		Short: "Information about this build of stackit",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf(`
Version: %s
Commit: %s
Date: %s
`, version, commit, date)
		},
	}

	RootCmd.AddCommand(cmd)
}
