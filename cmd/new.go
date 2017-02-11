// Copyright © 2017 NAME HERE <EMAIL ADDRESS>
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
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"errors"

	"github.com/ecosystemsoftware/eco/ecosql"
	eco "github.com/ecosystemsoftware/eco/utilities"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
)

var isAdmin bool

func init() {
	RootCmd.AddCommand(newCmd)
	newCmd.AddCommand(newConfigfileCmd)
	newCmd.AddCommand(newUserCmd)
	newCmd.AddCommand(newBundleCmd)

	newUserCmd.Flags().BoolVar(&isAdmin, "admin", false, "Create user with admin role")

}

// newCmd represents the new command
var newCmd = &cobra.Command{
	Use:   "new [object]",
	Short: "Create new files, users and bundles",
}

//newConfigfileCmd creates a template config.json with sane defaults.
var newConfigfileCmd = &cobra.Command{
	Use:   "configfile",
	Short: "Create a default config.json",
	Long: `Creates a template for config.json with all available
	configuration parameters`,
	RunE: createDefaultConfigFile,
}

var newUserCmd = &cobra.Command{
	Use:   "user [email]",
	Short: "Create a new user",
	Long: `Creates an entry in the database for a new user with
	associated email address.  The default role
	for this command is 'anon'. Use the -admin flag to create an admin user`,
	RunE: createNewUser,
}

var newBundleCmd = &cobra.Command{
	Use:   "bundle [name]",
	Short: "Create a new EcoSystem Bundle",
	Long:  `Scaffolds a new bundle including folder structure and required files`,
	RunE:  createNewBundle,
}

//createDafaultConfigFile creates the default config.json template with sane defaults
//Will overwrite existing config.json, so ask for confirmation
func createDefaultConfigFile(cmd *cobra.Command, args []string) error {

	c := eco.AskForConfirmation("This will overwrite any existing config.json. Do you with to proceed?")
	if c {
		config := eco.Config{
			PgSuperUser:         "postgres",
			PgDBName:            "ecosystem",
			PgPort:              "5432",
			PgServer:            "localhost",
			PgDisableSSL:        false,
			ApiPort:             "3000",
			WebsitePort:         "3001",
			AdminPanelPort:      "3002",
			AdminPanelServeType: "unbundled",
			PublicSiteSlug:      "site",
			PrivateSiteSlug:     "private",
			SmtpHost:            "smtp",
			SmtpPort:            "25",
			SmtpUserName:        "info@yourdomain.com",
			SmtpFrom:            "info@yourdomain.com",
			EmailFrom:           "Your Name",
			JWTRealm:            "Your App Name",
		}

		configJSON, _ := json.MarshalIndent(config, "", "\t")
		err := ioutil.WriteFile("config.json", configJSON, 0644)
		if err == nil {
			log.Println("Created default config.json")
		}
		return err
	}

	log.Println("Aborted by user")
	return nil
}

func createNewUser(cmd *cobra.Command, args []string) error {

	if len(args) < 1 {
		return errors.New("user's email must be provided")
	}

	// err := viper.ReadInConfig() // Find and read the config file
	// if err != nil {             // Handle errors reading the config file
	// 	log.Fatal("Could not read config.json: ", err.Error())
	// }

	//Establish a temporary connection as the super user
	db := eco.SuperUserDBConfig.ReturnDBConnection("")
	defer db.Close()

	//Set to the default role
	role := "anon"
	//and overwrite if an admin user is being created
	if isAdmin {
		role = "admin"
	}

	_, err := db.Exec(fmt.Sprintf(ecosql.ToCreateAdministrator, args[0], role))
	if err != nil {
		log.Fatal("Could not create new user:", err.Error())
		return nil
	}

	log.Println("Successfully created new user", args[0], "as", role)
	return nil

}

func createNewBundle(cmd *cobra.Command, args []string) error {

	//Check for bundle name
	if len(args) < 1 {
		return errors.New("a bundle name must be provided")
	}

	//Check that bundle doesn't already exists
	var AppFs = afero.NewOsFs()
	basePath := "./" + args[0]
	exists, _ := afero.IsDir(AppFs, basePath)
	if exists {
		log.Fatal("Bundle", args[0], "already exists. Please provide a different name")
	}

	//Create the folder structure
	err := os.MkdirAll("./"+args[0]+"/templates", os.ModePerm)
	err = os.MkdirAll("./"+args[0]+"/images", os.ModePerm)
	err = os.MkdirAll("./"+args[0]+"/public", os.ModePerm)

	if err != nil {
		log.Fatal("Could not complete folder setup: ", err.Error())
	}

	_, err = os.Create("./" + args[0] + "/install.sql")
	_, err = os.Create("./" + args[0] + "/demodata.sql")

	if err != nil {
		log.Fatal("Could not complete file setup: ", err.Error())
	}

	//Creates the bundles

	log.Println("Successfully created bundle", args[0])
	return nil

}
