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
	"log"
	"os"

	"github.com/ecosystemsoftware/eco/ecosql"
	eco "github.com/ecosystemsoftware/eco/utilities"
	"github.com/spf13/cobra"
)

func init() {

	RootCmd.AddCommand(initCmd)
	initCmd.AddCommand(initDBCmd)
	initCmd.AddCommand(initFoldersCmd)

}

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Complete initial setup of database and folder structure",
	Long:  `Performs a complete initialisation of the database and folder structure for EcoSystem`,
	RunE:  initAll,
}

// initDBCmd initialiseds the database
var initDBCmd = &cobra.Command{
	Use:   "db",
	Short: "Perform the database initialisation for built in tables, roles and permissions",
	Long: `Executes the initialisation SQL which sets up the built-in tables, as well
	as creating built-in roles anon,admin, web and server and assigning permissions.
	Tables will not be overwritten if they already exist.`,
	RunE: initDB,
}

// initCmd initialises the folder structure
var initFoldersCmd = &cobra.Command{
	Use:   "folders",
	Short: "Creates EcoSystem folder structure",
	Long: `Performs a complete initialisation of the folder structure for EcoSystem.
	Folders that already exist will not be overwritten.`,
	RunE: initFolders,
}

//initAll
func initAll(cmd *cobra.Command, args []string) error {

	c := eco.AskForConfirmation("This will perform a complete (re)initialisation and may perform overwrites. Do you with to proceed?")

	if c {
		initDB(cmd, args)
		initFolders(cmd, args)
		log.Println("Successfully completed EcoSystem initialisation")
		return nil
	}

	log.Println("Aborted by user")
	return nil
}

//initDB initialises the built-in database tables, roles and permissions
func initDB(cmd *cobra.Command, args []string) error {

	// err := viper.ReadInConfig() // Find and read the config file
	// if err != nil {             // Handle errors reading the config file
	// 	log.Fatal("Could not read config.json: ", err.Error())
	// }

	//Establish a temporary connection as the super user
	db := eco.SuperUserDBConfig.ReturnDBConnection("")
	defer db.Close()

	//Run initialisation SQL
	var err error
	_, err = db.Exec(ecosql.ToCreateUUIDExtension)
	_, err = db.Exec(ecosql.ToCreateUsersTable)
	_, err = db.Exec(ecosql.ToCreateFuncToGenerateNewUserID)
	_, err = db.Exec(ecosql.ToCreateTriggerOnNewUserInsert)
	_, err = db.Exec(ecosql.ToCreateWebCategoriesTable)
	_, err = db.Exec(ecosql.ToCreateServerRole)
	_, err = db.Exec(ecosql.ToCreateAnonRole)
	_, err = db.Exec(ecosql.ToCreateAdminRole)
	_, err = db.Exec(ecosql.ToCreateWebRole)
	_, err = db.Exec(ecosql.ToGrantBuiltInPermissions)
	_, err = db.Exec(ecosql.ToGrantAdminPermissions)

	if err != nil {
		log.Fatal("Could not complete database setup: ", err.Error())
	}

	log.Println("Successfully completed EcoSystem database initialisation")
	return nil

}

//initFolders initialises the filesystem used by EcoSystem
func initFolders(cmd *cobra.Command, args []string) error {

	var err error
	err = os.MkdirAll("./public/images_resized", os.ModePerm)
	err = os.Mkdir("./public/images_source", os.ModePerm)
	err = os.Mkdir("./templates", os.ModePerm)
	err = os.Mkdir("./bundles", os.ModePerm)
	err = os.Mkdir("./ecosystem-admin", os.ModePerm)

	if err != nil {
		log.Fatal("Could not complete folder setup: ", err.Error())
	}

	log.Println("Successfully completed EcoSystem folder setup")
	return nil
}
