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
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"path"

	"database/sql"

	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var isInstallDemoData, isReinstall bool

func init() {
	RootCmd.AddCommand(installCmd)
	RootCmd.AddCommand(unInstallCmd)
	installCmd.Flags().BoolVar(&isInstallDemoData, "demodata", false, "Install bundle demo data if available")
	installCmd.Flags().BoolVarP(&isReinstall, "reinstall", "r", false, "Uninstall bundle before installing")
}

// installCmd represents the install command
var installCmd = &cobra.Command{
	Use:   "install [bundle]",
	Short: "Install an ghost bundle",
	Long: `Installs an ghost bundle from the named folder.
	Note: does not download anything, so the bundle folder must
	exist and contain everything.  Previous to installing, either clone
	or download the bundle into the 'bundles' directory`,
	RunE: installBundle,
}

// installCmd represents the install command
var unInstallCmd = &cobra.Command{
	Use:   "uninstall [bundle]",
	Short: "Removes an ghost bundle",
	Long: `Removes an ghost bundle by deleting the DB schema, deleting template
	files and images,`,
	RunE: unInstallBundle,
}

//uninstallBundle is the removal function for a bundle
func unInstallBundle(cmd *cobra.Command, args []string) error {

	readConfig()

	//Check for bundle name
	if len(args) < 1 {
		return errors.New("a bundle name must be provided")
	}

	//Ask for confirmation
	c := AskForConfirmation("This will delete the bundle, causing loss of all data in the schema created by the bundle.  Are you sure you want to do this?")

	if c {
		//Establish a temporary connection as the super user
		db := SuperUserDBConfig.ReturnDBConnection("")
		defer db.Close()

		//Drop the schema
		//If it doesn't exist, it won't be dropped - no big deal
		db.Exec(fmt.Sprintf(SQLToDropSchema, args[0]))

		//Attempt to updated the bundles installed list
		newBundlesInstalled, err := Bundles(viper.GetStringSlice("bundlesInstalled")).UnInstallBundle(args[0])

		//If there is any error, log it
		if err != nil {
			Log(LogEntry{"ghost.INSTALL", false, "Error updating bundles installed list: " + err.Error()})
		}

		//Otherwise set the viper configuration to the new bundles list and overwrite the config.json
		viper.Set("bundlesInstalled", newBundlesInstalled)
		var config Config
		viper.Unmarshal(&config)
		configJSON, _ := json.MarshalIndent(config, "", "\t")
		err = ioutil.WriteFile(viper.GetString("configfile")+".json", configJSON, 0644)
		if err != nil {
			Log(LogEntry{"ghost.INSTALL", false, "Error updating config file: " + err.Error()})
		}

		Log(LogEntry{"ghost.INSTALL", true, "config.json updated"})
		Log(LogEntry{"ghost.INSTALL", true, "Uninstallation of bundle " + args[0] + " completed"})

	}

	return nil

}

//installBundle is the entire installation procedure for an ghost Bundle
func installBundle(cmd *cobra.Command, args []string) error {

	readConfig()

	//Check for bundle name
	if len(args) < 1 {
		return errors.New("a bundle name must be provided")
	}

	//Check that bundle installation folder exists
	basePath := "./bundles/" + args[0] + "/install"

	exists, err := afero.IsDir(AppFs, basePath)

	if !exists || err != nil {
		//Exit if doesn't exist
		LogFatal(LogEntry{"ghost.INSTALL", false, "Bundle '" + args[0] + "' install folder not found or unreadable."})
	}

	//Uninstall first if requested
	if isReinstall {
		Log(LogEntry{"ghost.INSTALL", true, "Uninstalling bundle " + args[0] + " before reinstalling"})
		unInstallBundle(cmd, args)
	}

	//Check for error reading directory or zero files
	filesInDirectory, err := afero.ReadDir(AppFs, basePath)
	if err != nil || len(filesInDirectory) == 0 {
		LogFatal(LogEntry{"ghost.INSTALL", false, "No installation files could be read for bundle"})
		return nil
	}

	Log(LogEntry{"ghost.INSTALL", true, "Installing bundle '" + args[0] + "'"})

	//Establish a temporary connection as the super user
	db := SuperUserDBConfig.ReturnDBConnection("")
	defer db.Close()

	//Set up a schema for the bundle
	err = setupDBSchema(db, args[0])
	if err != nil {
		//IF there is any type of error, drop the schema, log and exit
		db.Exec(fmt.Sprintf(SQLToDropSchema, args[0]))
		LogFatal(LogEntry{"ghost.INSTALL", false, "Schema creation failed with error " + err.Error()})
		return nil
	}

	//Iterate over the installation files
	for _, file := range filesInDirectory {
		//Ignore directories
		if !file.IsDir() {
			//Attempt to processes the sqlfile
			err := processBundleFile(db, path.Join(basePath, file.Name()))
			if err != nil {
				//IF there is any type of error, drop the schema, log and exit
				db.Exec(fmt.Sprintf(SQLToDropSchema, args[0]))
				LogFatal(LogEntry{"ghost.INSTALL", false, "Installation of '" + file.Name() + "' failed with error: " + err.Error()})
				return nil
			}
			Log(LogEntry{"ghost.INSTALL", true, file.Name() + " installed OK"})
		}
	}

	//If the user has asked for demo data
	if isInstallDemoData {

		Log(LogEntry{"ghost.INSTALL", true, "Installing demo data"})

		basePath := "./bundles/" + args[0] + "/demodata"

		//Check for error reading directory or zero files
		filesInDirectory, err := afero.ReadDir(AppFs, basePath)
		if err != nil || len(filesInDirectory) == 0 {
			//IF there is any type of error, drop the schema, log and exit
			db.Exec(fmt.Sprintf(SQLToDropSchema, args[0]))
			LogFatal(LogEntry{"ghost.INSTALL", false, "No demo data files could be read for bundle"})
			return nil
		}

		//Iterate over the demodata files
		for _, file := range filesInDirectory {
			//Ignore directories
			if !file.IsDir() {
				//Attempt to processes the sqlfile
				err := processBundleFile(db, path.Join(basePath, file.Name()))
				if err != nil {
					//IF there is any type of error, drop the schema, log and exit
					db.Exec(fmt.Sprintf(SQLToDropSchema, args[0]))
					LogFatal(LogEntry{"ghost.INSTALL", false, "Installation of '" + file.Name() + "' failed with error: " + err.Error()})
					return nil
				}

				Log(LogEntry{"ghost.INSTALL", true, file.Name() + " installed OK"})

			}
		}

	}

	//Attempt to updated the bundles installed list
	newBundlesInstalled, err := Bundles(viper.GetStringSlice("bundlesInstalled")).InstallBundle(args[0])

	//If there is any error, return it
	if err != nil {
		Log(LogEntry{"ghost.INSTALL", false, "Error updating bundles installed list: " + err.Error()})
	}

	//Otherwise set the viper configuration to the new bundles list and overwrite the config.json
	viper.Set("bundlesInstalled", newBundlesInstalled)
	var config Config
	viper.Unmarshal(&config)
	configJSON, _ := json.MarshalIndent(config, "", "\t")
	err = ioutil.WriteFile(viper.GetString("configfile")+".json", configJSON, 0644)
	if err != nil {
		Log(LogEntry{"ghost.INSTALL", false, "Error updating config file: " + err.Error()})
	}

	Log(LogEntry{"ghost.INSTALL", true, "config file updated"})

	//Bundle installation complete
	Log(LogEntry{"ghost.INSTALL", true, "Installation of bundle " + args[0] + " completed"})
	return nil

}

func processBundleFile(db *sql.DB, filename string) error {

	//Attempt to read file
	sqlBytes, err := afero.ReadFile(AppFs, filename)

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
	_, err := db.Exec(fmt.Sprintf(SQLToCreateSchema, bundleName))

	if err != nil {
		return err
	}

	//Set admin privileges for everything in this schema going forwards
	_, err = db.Exec(fmt.Sprintf(SQLToGrantBundleAdminPermissions, bundleName, bundleName, bundleName))

	if err != nil {
		return err
	}

	//Set the search path to the bundle schema so that all SQL commands take
	//place within the schema
	_, err = db.Exec(fmt.Sprintf(SQLToSetSearchPathForBundle, bundleName))

	if err != nil {
		return err
	}

	return nil

}
