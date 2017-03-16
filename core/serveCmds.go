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

//Use the forked version of the go-jwt-middlware, not the auth0 version

package core

import (
	"net/http"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"fmt"
)

var nowebsite, noadminpanel bool
var smtpPW string

func init() {
	RootCmd.AddCommand(serveCmd)

	serveCmd.Flags().String("smtppw", "", "SMTP server password for outgoing mail")
	viper.BindPFlag("smtppw", serveCmd.Flags().Lookup("smtppw"))

	serveCmd.Flags().BoolP("demomode", "d", false, "Run server in demo mode")
	viper.BindPFlag("demomode", serveCmd.Flags().Lookup("demomode"))

	serveCmd.Flags().StringP("secret", "s", "", "Secure secret for signing JWT")
	viper.BindPFlag("secret", serveCmd.Flags().Lookup("secret"))

}

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Starts the EcoSystem server",
	Long:  `Start the EcoSystem API Server`,
	RunE:  serve,
}

func serve(cmd *cobra.Command, args []string) error {

	preServe()
	startServer()
	return nil
}

var ActivatePackages func()

func preServe() {

	//Check to make sure a secret has been provided
	//No default provided as a security measure, server will exit of nothing provided
	if viper.GetString("secret") == "" {
		LogFatal(LogEntry{"CORE.SERVE", false, "No signing secret provided"})
	}

	//Establish a temporary connection as the super user
	dbTemp := SuperUserDBConfig.ReturnDBConnection("")

	//Generate a random server password, set it and get out
	serverPW := RandomString(16)
	_, err := dbTemp.Exec(fmt.Sprintf(SQLToSetServerRolePassword, serverPW))
	if err != nil {
		LogFatal(LogEntry{"CORE.SERVE", false, "Error setting server role password: " + err.Error()})
	}

	dbTemp.Close()

	//Establish a permanent connection
	DB = ServerUserDBConfig.ReturnDBConnection(serverPW)

	ActivatePackages()

}

func startServer() {

	Log(LogEntry{"CORE.SERVE", true, "Server started on port " + viper.GetString("apiPort")})
	http.ListenAndServe(":"+viper.GetString("apiPort"), Router)

}

//  Experimental search features
// 	api.Handle("SEARCH", "/:schema/:table/", ReturnBlank) //Useful for when blank searches are sent by client, to avoid errors
// 	api.Handle("SEARCH", "/:schema/:table/:searchTerm", SearchList)
