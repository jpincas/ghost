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
	"fmt"
	"os"

	"errors"

	"github.com/spf13/afero"
	"github.com/spf13/cobra"
)

var isAdmin bool

func init() {
	RootCmd.AddCommand(newCmd)
	newCmd.AddCommand(newUserCmd)
	newCmd.AddCommand(newBundleCmd)
	newUserCmd.Flags().BoolVar(&isAdmin, "admin", false, "Create user with admin role")

}

// newCmd represents the new command
var newCmd = &cobra.Command{
	Use:   "new [object]",
	Short: "Create new files, users and bundles",
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
	Short: "Create a new ghost Bundle",
	Long:  `Scaffolds a new bundle including folder structure and required files`,
	RunE:  createNewBundle,
}

func createNewUser(cmd *cobra.Command, args []string) error {

	readConfig()

	if len(args) < 1 {
		return errors.New("user's email must be provided")
	}

	// err := viper.ReadInConfig() // Find and read the config file
	// if err != nil {             // Handle errors reading the config file
	// 	log.Fatal("Could not read config.json: ", err.Error())
	// }

	//Establish a temporary connection as the super user
	db := SuperUserDBConfig.ReturnDBConnection("")
	defer db.Close()

	//Set to the default role
	role := "anon"
	//and overwrite if an admin user is being created
	if isAdmin {
		role = "admin"
	}

	_, err := db.Exec(fmt.Sprintf(SQLToCreateAdministrator, args[0], role))
	if err != nil {
		LogFatal(LogEntry{"ghost.NEW", true, "Could not create new user:" + err.Error()})
		return nil
	}

	Log(LogEntry{"ghost.NEW", true, "Successfully created new user " + args[0] + " as " + role})
	return nil

}

func createNewBundle(cmd *cobra.Command, args []string) error {

	readConfig()

	//Check for bundle name
	if len(args) < 1 {
		return errors.New("a bundle name must be provided")
	}

	//Check that bundle doesn't already exists
	var AppFs = afero.NewOsFs()
	basePath := "./" + args[0]
	exists, _ := afero.IsDir(AppFs, basePath)
	if exists {
		LogFatal(LogEntry{"ghost.NEW", true, "Bundle " + args[0] + " already exists. Please provide a different name"})
	}

	//Create the folder structure
	err := os.MkdirAll("./"+args[0]+"/templates/pages", os.ModePerm)
	err = os.MkdirAll("./"+args[0]+"/templates/email", os.ModePerm)
	err = os.MkdirAll("./"+args[0]+"/templates/partials", os.ModePerm)
	err = os.MkdirAll("./"+args[0]+"/images", os.ModePerm)
	err = os.MkdirAll("./"+args[0]+"/public", os.ModePerm)
	err = os.MkdirAll("./"+args[0]+"/admin-panel", os.ModePerm)

	if err != nil {
		LogFatal(LogEntry{"ghost.NEW", true, "Could not complete folder setup: " + err.Error()})
	}

	_, err = os.Create("./" + args[0] + "/install.sql")
	_, err = os.Create("./" + args[0] + "/demodata.sql")

	if err != nil {
		LogFatal(LogEntry{"ghost.NEW", true, "Could not complete folder setup: " + err.Error()})
	}

	//Creates the bundles
	Log(LogEntry{"ghost.NEW", true, "Successfully created bundle " + args[0]})
	return nil

}
