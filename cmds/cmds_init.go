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

package ghost

import (
	"os"

	"github.com/jpincas/ghost/ghost"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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
	Long:  `Performs a complete initialisation of the database and folder structure for ghost`,
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
	Short: "Creates ghost folder structure",
	Long: `Performs a complete initialisation of the folder structure for ghost.
	Folders that already exist will not be overwritten.`,
	RunE: initFolders,
}

//initAll
func initAll(cmd *cobra.Command, args []string) error {

	//If user has used -noprompt flag then we don't prompt for confirmation
	var proceedWithInit = false
	if viper.GetBool("noprompt") {
		proceedWithInit = true
	} else {
		proceedWithInit = ghost.AskForConfirmation("This will perform a complete (re)initialisation and may perform overwrites. Do you with to proceed?")
	}

	if proceedWithInit {
		initDB(cmd, args)
		initFolders(cmd, args)
		ghost.LLog("INIT", true, "Successfully completed ghost initialisation", nil)
		return nil
	}

	ghost.LLog("INIT", false, "Aborted by user", nil)

	return nil
}

//initDB initialises the built-in database tables, roles and permissions
func initDB(cmd *cobra.Command, args []string) error {

	ghost.App.Config.Setup(viper.GetString("configfile"))

	//Establish a temporary connection as the super user
	db := ghost.SuperUserDBConfig.ReturnDBConnection("")
	defer db.Close()

	//Run initialisation SQL
	var err error
	_, err = db.Exec(ghost.SQLToCreateAdminRole)
	_, err = db.Exec(ghost.SQLToGrantAdminPermissions) //Do this first so everything created after will have correct admin permissions by default
	_, err = db.Exec(ghost.SQLToCreateUUIDExtension)
	_, err = db.Exec(ghost.SQLToCreateUsersTable)
	_, err = db.Exec(ghost.SQLToCreateFuncToGenerateNewUserID)
	_, err = db.Exec(ghost.SQLToCreateTriggerOnNewUserInsert)
	_, err = db.Exec(ghost.SQLToCreateServerRole)
	_, err = db.Exec(ghost.SQLToCreateAnonRole)
	_, err = db.Exec(ghost.SQLToGrantBuiltInPermissions)

	if err != nil {
		ghost.LLogFatal("INIT", false, "Could not complete database setup", err)
	}

	ghost.LLog("INIT", true, "Successfully completed ghost database initialisation", nil)
	return nil

}

//initFolders initialises the filesystem used by ghost
func initFolders(cmd *cobra.Command, args []string) error {

	ghost.App.Config.Setup(viper.GetString("configfile"))

	var err error
	err = os.Mkdir("./bundles", os.ModePerm)

	if err != nil {
		ghost.LLog("INIT", false, "Could not complete folder setup", err)
	}

	ghost.LLog("INIT", true, "Successfully completed ghost folder initialisation", nil)
	return nil
}
