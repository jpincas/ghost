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
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"path"

	"database/sql"

	"github.com/jpincas/ghost/ghost"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	sqlToDropSchema                  = `DROP SCHEMA %s CASCADE;`
	sqlToSetSearchPathForBundle      = `SET search_path TO %s, public;`
	sqlToCreateSchema                = `CREATE SCHEMA %s;`
	sqlToGrantBundleAdminPermissions = `GRANT USAGE ON SCHEMA %s TO admin; ALTER DEFAULT PRIVILEGES IN SCHEMA %s GRANT ALL ON TABLES TO admin; ALTER DEFAULT PRIVILEGES IN SCHEMA %s GRANT USAGE ON SEQUENCES TO admin;`
)

var isInstallDemoData, isReinstall, demoDataOnly bool

func init() {
	RootCmd.AddCommand(installCmd)
	RootCmd.AddCommand(unInstallCmd)
	installCmd.Flags().BoolVar(&isInstallDemoData, "demodata", false, "Install bundle demo data if available")
	installCmd.Flags().BoolVar(&demoDataOnly, "demodataonly", false, "Install bundle demo data if available")
	installCmd.Flags().BoolVarP(&isReinstall, "reinstall", "r", false, "Uninstall bundle before installing")
}

// installCmd represents the install command
var installCmd = &cobra.Command{
	Use:   "install [bundle]",
	Short: "Install a ghost bundle",
	Long: `Installs a ghost bundle from the named folder.
	Note: does not download anything, so the bundle folder must
	exist and contain everything.  Previous to installing, either clone
	or download the bundle into the 'bundles' directory`,
	RunE: installBundle,
}

// installCmd represents the install command
var unInstallCmd = &cobra.Command{
	Use:   "uninstall [bundle]",
	Short: "Removes a ghost bundle",
	Long:  `Removes a ghost bundle by deleting the schema`,
	RunE:  unInstallBundle,
}

//uninstallBundle is the removal function for a bundle
func unInstallBundle(cmd *cobra.Command, args []string) error {

	configFile := viper.GetString("configfile")
	ghost.App.Setup(viper.GetString("configfile"))

	//Check for bundle name
	if len(args) < 1 {
		return errors.New("a bundle name must be provided")
	}

	//If user has used -noprompt flag then we don't prompt for confirmation
	var proceedWithInit = false
	if viper.GetBool("noprompt") {
		proceedWithInit = true
	} else {
		proceedWithInit = ghost.AskForConfirmation("This will delete the bundle, causing loss of all data in the schema created by the bundle.  Are you sure you want to do this?")
	}

	if proceedWithInit {

		//Establish a temporary connection as the super user
		db := ghost.SuperUserDBConfig.ReturnDBConnection("")
		defer db.Close()

		//Drop the schema
		//If it doesn't exist, it won't be dropped - no big deal
		db.Exec(fmt.Sprintf(sqlToDropSchema, args[0]))

		//Attempt to updated the bundles installed list
		if err := ghost.App.Config.UnInstallBundle(args[0]); err != nil {
			ghost.Log("INSTALL", false, "Error uninstalling bundle", err)
		}

		configJSON, _ := json.MarshalIndent(ghost.App.Config, "", "\t")
		if err := ioutil.WriteFile(configFile+".json", configJSON, 0644); err != nil {
			ghost.Log("INSTALL", false, "Error updating config file", err)
		}

		ghost.Log("INSTALL", true, "config.json updated", nil)
		ghost.Log("INSTALL", true, "Uninstallation of bundle "+args[0]+" completed", nil)

	}

	return nil

}

//installBundle is the entire installation procedure for an ghost Bundle
func installBundle(cmd *cobra.Command, args []string) error {

	configFile := viper.GetString("configfile")
	ghost.App.Setup(viper.GetString("configfile"))

	//Check for bundle name
	if len(args) < 1 {
		return errors.New("a bundle name must be provided")
	}

	//Establish a temporary connection as the super user
	db := ghost.SuperUserDBConfig.ReturnDBConnection("")
	defer db.Close()

	bundleName := args[0]
	if demoDataOnly {
		installBundleDemoData(bundleName, db)
		return nil
	}

	if isReinstall {
		ghost.Log("INSTALL", true, "Uninstalling bundle "+bundleName+" before reinstalling", nil)
		unInstallBundle(cmd, args)
	}

	installBundleSchema(bundleName, db)

	if isInstallDemoData {
		installBundleDemoData(bundleName, db)
	}

	//Attempt to update the bundles installed list
	if err := ghost.App.Config.InstallBundle(bundleName); err != nil {
		ghost.Log("INSTALL", false, "Error installing bundle", err)
	}

	//Rewrite the config file
	configJSON, _ := json.MarshalIndent(ghost.App.Config, "", "\t")
	if err := ioutil.WriteFile(configFile+".json", configJSON, 0644); err != nil {
		ghost.Log("INSTALL", false, "Error updating config file. Please update manually", err)
	} else {
		ghost.Log("INSTALL", true, "config file updated", err)
	}

	//Bundle installation complete
	ghost.Log("INSTALL", true, "Installation of bundle "+bundleName+" completed", nil)
	return nil

}

func installBundleSchema(bundleName string, db *sql.DB) {

	//Check that bundle installation folder exists
	basePath := "./bundles/" + bundleName + "/install"
	exists, err := afero.IsDir(ghost.App.FileSystem, basePath)
	if !exists || err != nil {
		ghost.LogFatal("INSTALL", false, "Bundle '"+bundleName+"' install folder not found or unreadable.", err)
	}

	//Check for error reading directory or zero files
	filesInDirectory, err := afero.ReadDir(ghost.App.FileSystem, basePath)
	if err != nil || len(filesInDirectory) == 0 {
		ghost.LogFatal("INSTALL", false, "No installation files could be read for bundle", err)
	}

	ghost.Log("INSTALL", true, "Installing bundle '"+bundleName+"'", nil)

	//Set up a schema for the bundle
	err = setupDBSchema(db, bundleName)
	if err != nil {
		//IF there is any type of error, drop the schema, log and exit
		db.Exec(fmt.Sprintf(sqlToDropSchema, bundleName))
		ghost.LogFatal("INSTALL", false, "Schema creation failed", err)
	}

	//Set the search path to the bundle schema so that all SQL commands take
	//place within the schema
	_, err = db.Exec(fmt.Sprintf(sqlToSetSearchPathForBundle, bundleName))
	if err != nil {
		db.Exec(fmt.Sprintf(sqlToDropSchema, bundleName))
		ghost.LogFatal("INSTALL", false, "Failed to set schema search path", err)
	}

	//Iterate over the installation files
	for _, file := range filesInDirectory {
		//Ignore directories
		if !file.IsDir() {
			//Attempt to processes the sqlfile
			err := processBundleFile(db, path.Join(basePath, file.Name()))
			if err != nil {
				//IF there is any type of error, drop the schema, log and exit
				db.Exec(fmt.Sprintf(sqlToDropSchema, bundleName))
				ghost.LogFatal("INSTALL", false, "Installation of '"+file.Name()+"' failed", err)
			}
			ghost.Log("INSTALL", true, file.Name()+" installed OK", nil)
		}
	}

}

func installBundleDemoData(bundleName string, db *sql.DB) {

	ghost.Log("INSTALL", true, "Installing demo data", nil)

	basePath := "./bundles/" + bundleName + "/demodata"

	//Check for error reading directory or zero files
	filesInDirectory, err := afero.ReadDir(ghost.App.FileSystem, basePath)
	if err != nil || len(filesInDirectory) == 0 {
		//IF there is any type of error, drop the schema, log and exit
		db.Exec(fmt.Sprintf(sqlToDropSchema, bundleName))
		ghost.LogFatal("INSTALL", false, "No demo data files could be read for bundle", err)
	}

	log.Println("installing demo data...")

	//Set the search path to the bundle schema so that all SQL commands take
	//place within the schema
	_, err = db.Exec(fmt.Sprintf(sqlToSetSearchPathForBundle, bundleName))
	if err != nil {
		db.Exec(fmt.Sprintf(sqlToDropSchema, bundleName))
		ghost.LogFatal("INSTALL", false, "Failed to set schema search path", err)
	}

	//Iterate over the demodata files
	for _, file := range filesInDirectory {
		//Ignore directories
		if !file.IsDir() {
			//Attempt to processes the sqlfile
			err := processBundleFile(db, path.Join(basePath, file.Name()))
			if err != nil {
				//IF there is any type of error, drop the schema, log and exit
				db.Exec(fmt.Sprintf(sqlToDropSchema, bundleName))
				ghost.LogFatal("INSTALL", false, "Installation of '"+file.Name()+"' failed", err)
			}

			ghost.Log("INSTALL", true, file.Name()+" installed OK", nil)

		}
	}

}

func processBundleFile(db *sql.DB, filename string) error {

	//Attempt to read file
	sqlBytes, err := afero.ReadFile(ghost.App.FileSystem, filename)

	if err != nil {
		return err
	}

	//Run the SQL
	if _, err = db.Exec(string(sqlBytes)); err != nil {
		return err
	}

	return nil

}

func setupDBSchema(db *sql.DB, bundleName string) error {

	//Attempt to create a schema matching the bundle's name,
	_, err := db.Exec(fmt.Sprintf(sqlToCreateSchema, bundleName))

	if err != nil {
		return err
	}

	//Set admin privileges for everything in this schema going forwards
	_, err = db.Exec(fmt.Sprintf(sqlToGrantBundleAdminPermissions, bundleName, bundleName, bundleName))

	if err != nil {
		return err
	}

	return nil

}
