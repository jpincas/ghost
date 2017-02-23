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
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"log"

	eco "github.com/ecosystemsoftware/ecosystem/utilities"
)

var configFileName string

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

	cobra.OnInitialize(initConfig)
	RootCmd.PersistentFlags().StringP("pgpw", "p", "", "Postgres superuser password")
	viper.BindPFlag("pgpw", RootCmd.PersistentFlags().Lookup("pgpw"))

	RootCmd.PersistentFlags().StringVarP(&configFileName, "configfile", "c", "config", "Name of config file (without extension)")

}

func justCheckForConfigFile(cmd *cobra.Command, args []string) error {
	//This function actually does nothing
	//The only benefit of running it is that it sparks initConfig which checks for a config file
	//And creates one if necessary
	//This just means the whole process can be started by typiing 'ecosystem', which is cool!
	return nil
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {

	viper.SetDefault("pgSuperUser", "postgres")
	viper.SetDefault("pgDBName", "testdb")
	viper.SetDefault("pgPort", "5432")
	viper.SetDefault("pgServer", "localhost")
	viper.SetDefault("pgDisableSSL", true)
	viper.SetDefault("apiPort", "3000")
	viper.SetDefault("websitePort", "3001")
	viper.SetDefault("adminPanelPort", "3002")
	viper.SetDefault("adminPanelServeDirectory", "ecosystem-admin/build/unbundled")
	viper.SetDefault("publicSiteSlug", "site")
	viper.SetDefault("privateSiteSlug", "private")
	viper.SetDefault("jwtRealm", "yourappname")
	viper.SetDefault("host", "localhost")
	viper.SetDefault("protocol", "http")

	//For the admin panel
	viper.SetDefault("adminPrimaryColor", "#00c4a7")
	viper.SetDefault("adminSecondaryColor", "#7EC9A2")
	viper.SetDefault("adminTextColor", "black")
	viper.SetDefault("adminErrorColor", "red")
	viper.SetDefault("adminTitle", "Admin Panel")
	viper.SetDefault("adminLogoFile", "logo.png")
	viper.SetDefault("adminLogoBundle", "master")
	viper.SetDefault("bundlesInstalled", make([]string, 0, 0))

	viper.SetConfigName(configFileName)
	viper.AddConfigPath(".")

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		//Initialise the db config structs for later use
		eco.InitDBConnectionConfigs()
		log.Println("Config file detected and correctly applied:", viper.ConfigFileUsed())
	} else {
		//Otherwise create one
		log.Println("Config file does not exist. Creating now...")
		if err := createDefaultConfigFile(); err != nil {
			log.Println("Error creating config file: ", err.Error())
		}
	}
}

//createDafaultConfigFile creates the default config.json template with sane defaults
//Will overwrite existing config.json, so ask for confirmation
func createDefaultConfigFile() error {

	config := eco.Config{
		PgSuperUser:              "postgres",
		PgDBName:                 "testdb",
		PgPort:                   "5432",
		PgServer:                 "localhost",
		PgDisableSSL:             true,
		ApiPort:                  "3000",
		WebsitePort:              "3001",
		AdminPanelPort:           "3002",
		AdminPanelServeDirectory: "ecosystem-admin/build/unbundled",
		PublicSiteSlug:           "site",
		PrivateSiteSlug:          "private",
		SmtpHost:                 "smtp",
		SmtpPort:                 "25",
		SmtpUserName:             "info@yourdomain.com",
		SmtpFrom:                 "info@yourdomain.com",
		EmailFrom:                "Your Name",
		JWTRealm:                 "Your App Name",
		AdminPrimaryColor:        "#00c4a7",
		AdminSecondaryColor:      "#7EC9A2",
		AdminTextColor:           "black",
		AdminErrorColor:          "red",
		AdminTitle:               "Admin Panel",
		AdminLogoFile:            "logo.png",
		AdminLogoBundle:          "master",
		BundlesInstalled:         make([]string, 0, 0),
		Host:                     "localhost",
		Protocol:                 "http",
	}

	configJSON, _ := json.MarshalIndent(config, "", "\t")
	err := ioutil.WriteFile("config.json", configJSON, 0644)
	if err != nil {
		return err
	}

	return nil

}
