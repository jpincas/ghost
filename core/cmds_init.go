// Copyright 2017 Jonathan Pincas

// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at

// 	http://www.apache.org/licenses/LICENSE-2.0

// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package core

import (
	"os"

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

	c := AskForConfirmation("This will perform a complete (re)initialisation and may perform overwrites. Do you with to proceed?")

	if c {
		initDB(cmd, args)
		initFolders(cmd, args)
		Log(LogEntry{"CORE.INIT", true, "Successfully completed EcoSystem initialisation"})
		return nil
	}

	Log(LogEntry{"CORE.INIT", false, "Aborted by user"})

	return nil
}

//initDB initialises the built-in database tables, roles and permissions
func initDB(cmd *cobra.Command, args []string) error {

	readConfig()

	//Establish a temporary connection as the super user
	db := SuperUserDBConfig.ReturnDBConnection("")
	defer db.Close()

	//Run initialisation SQL
	var err error
	_, err = db.Exec(SQLToCreateAdminRole)
	_, err = db.Exec(SQLToGrantAdminPermissions) //Do this first so everything created after will have correct admin permissions by default
	_, err = db.Exec(SQLToCreateUUIDExtension)
	_, err = db.Exec(SQLToCreateUsersTable)
	_, err = db.Exec(SQLToCreateFuncToGenerateNewUserID)
	_, err = db.Exec(SQLToCreateTriggerOnNewUserInsert)
	_, err = db.Exec(SQLToCreateServerRole)
	_, err = db.Exec(SQLToCreateAnonRole)
	_, err = db.Exec(SQLToGrantBuiltInPermissions)

	if err != nil {
		LogFatal(LogEntry{"CORE.INIT", false, "Could not complete database setup: " + err.Error()})
	}

	Log(LogEntry{"CORE.INIT", true, "Successfully completed EcoSystem database initialisation"})
	return nil

}

//initFolders initialises the filesystem used by EcoSystem
func initFolders(cmd *cobra.Command, args []string) error {

	readConfig()

	var err error
	err = os.Mkdir("./bundles", os.ModePerm)

	if err != nil {
		Log(LogEntry{"CORE.INIT", false, "Could not complete folder setup: " + err.Error()})
	}

	Log(LogEntry{"CORE.INIT", true, "Successfully completed EcoSystem folder initialisation"})
	return nil
}
