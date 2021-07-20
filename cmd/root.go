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
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "stackit",
	Short: "Synchronous CLI tool for CloudFormation stack operations",
	Long: `stackit is a CLI tool to synchronously and idempotently operate on AWS
CloudFormation stacks - a perfect complement for continuous integration systems
and developers who prefer the comfort of the command line.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	//	Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	fmt.Fprintf(os.Stderr, "stackit version %s\n", version)

	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	RootCmd.PersistentFlags().String("region", "", "")
	RootCmd.PersistentFlags().String("profile", "", "")
	RootCmd.PersistentFlags().String("stack-name", "", "")
	RootCmd.PersistentFlags().String("stackit-s3-suffix", "", "")
	RootCmd.PersistentFlags().String("stackit-tags", "", "")

	// Here you will define your flags and configuration settings.
	// Cobra supports Persistent Flags, which, if defined here,
	// will be global for your application.

	RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.stackit.yaml)")
	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	RootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	viper.BindPFlags(RootCmd.PersistentFlags())
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	viper.SetConfigName(".stackit") // name of config file (without extension)
	viper.AddConfigPath(".")        // adding home directory as first search path
	viper.AddConfigPath("$HOME")    // adding home directory as first search path
	viper.AutomaticEnv()            // read in environment variables that match

	if cfgFile != "" { // enable ability to specify config file via flag
		viper.SetConfigFile(cfgFile)
	}

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}
