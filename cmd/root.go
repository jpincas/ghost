// Copyright 2017 EcoSystem Software LLP

// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at

// 	http://www.apache.org/licenses/LICENSE-2.0

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

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "ecosystem [command] [arguments]",
	Short: "EcoSystem command line tool",
	Long: `Use to initialise or launch the EcoSystem server or create new users or bundles.
	Use the bare command 'ecosystem' to create a new config.json or verify an existing one.`,
	RunE: justCheckForConfigFile,
}

// Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func init() {

	RootCmd.PersistentFlags().StringP("pgpw", "p", "", "Postgres superuser password")
	viper.BindPFlag("pgpw", RootCmd.PersistentFlags().Lookup("pgpw"))
	RootCmd.PersistentFlags().StringP("configfile", "c", "config", "Name of config file (without extension)")
	viper.BindPFlag("configfile", RootCmd.PersistentFlags().Lookup("configfile"))

}

func justCheckForConfigFile(cmd *cobra.Command, args []string) error {
	//This function actually does nothing
	//The only benefit of running it is that it sparks initConfig which checks for a config file
	//And creates one if necessary
	//This just means the whole process can be started by typiing 'ecosystem', which is cool!
	return nil
}
