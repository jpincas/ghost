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

package cmds

import (
	"fmt"
	"os"
	"path"

	"errors"

	"github.com/jpincas/ghost/ghost"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const sqlToCreateAdministrator = `INSERT INTO users(email, role) VALUES ('%s', '%s');`

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

	ghost.App.Setup(viper.GetString("configfile"))

	if len(args) < 1 {
		return errors.New("user's email must be provided")
	}

	//Establish a temporary connection as the super user
	db := ghost.SuperUserDBConfig.ReturnDBConnection("")
	defer db.Close()

	//Set to the default role
	role := "anon"
	//and overwrite if an admin user is being created
	if isAdmin {
		role = "admin"
	}

	_, err := db.Exec(fmt.Sprintf(sqlToCreateAdministrator, args[0], role))
	if err != nil {
		ghost.LogFatal("NEW", true, "Could not create new user", err)
		return nil
	}

	ghost.Log("NEW", true, "Successfully created new user"+args[0]+" as "+role, nil)
	return nil

}

func createNewBundle(cmd *cobra.Command, args []string) error {

	fs := afero.NewOsFs()

	//Check for bundle name
	if len(args) < 1 {
		return errors.New("a bundle name must be provided")
	}

	//Check that bundle doesn't already exists
	basePath := path.Join("bundles", args[0])
	exists, _ := afero.IsDir(fs, basePath)
	if exists {
		ghost.LogFatal("NEW", true, "Bundle "+args[0]+" already exists. Please provide a different name", nil)
	}

	//Create the folder structure
	err := os.MkdirAll(path.Join(basePath, "install"), os.ModePerm)
	err = os.MkdirAll(path.Join(basePath, "demodata"), os.ModePerm)

	if err != nil {
		ghost.LogFatal("NEW", true, "Could not complete folder setup", err)
	}

	_, err = os.Create(path.Join(basePath, "install", "00_install.sql"))
	_, err = os.Create(path.Join(basePath, "demodata", "00_demodata.sql"))

	if err != nil {
		ghost.LogFatal("NEW", true, "Could not complete folder setup", err)
	}

	//Creates the bundles
	ghost.Log("NEW", true, "Successfully created bundle "+args[0], err)
	return nil

}
