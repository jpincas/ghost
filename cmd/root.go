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

	"log"

	eco "github.com/ecosystemsoftware/ecosystem/utilities"
)

var configFileName string

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "ecosystem",
	Short: "EcoSystem command line tool",
	Long: `Use to initialise or launch the EcoSystem server or create new users, bundles
	or a config file`,
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

	cobra.OnInitialize(initConfig)
	RootCmd.PersistentFlags().StringP("pgpw", "p", "", "Postgres superuser password")
	viper.BindPFlag("pgpw", RootCmd.PersistentFlags().Lookup("pgpw"))

	RootCmd.PersistentFlags().StringVarP(&configFileName, "configfile", "c", "config", "Name of config file (without extension)")

}

// initConfig reads in config file and ENV variables if set.
func initConfig() {

	viper.SetDefault("pgSuperUser", "postgres")
	viper.SetDefault("pgDBName", "ecosystem")
	viper.SetDefault("pgPort", "5432")
	viper.SetDefault("pgServer", "localhost")
	viper.SetDefault("pgDisableSSL", false)
	viper.SetDefault("apiPort", "3000")
	viper.SetDefault("websitePort", "3001")
	viper.SetDefault("adminPanelPort", "3002")
	viper.SetDefault("adminPanelServeDirectory", "ecosystem-admin")
	viper.SetDefault("publicSiteSlug", "site")
	viper.SetDefault("privateSiteSlug", "private")
	viper.SetDefault("jwtRealm", "yourappname")

	viper.SetConfigName(configFileName)
	viper.AddConfigPath(".")

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		//Initialise the db config structs for later use
		eco.InitDBConnectionConfigs()
		fmt.Println("Config file detected and correctly applied:", viper.ConfigFileUsed())
	} else {
		log.Fatal(err.Error())
	}
}
