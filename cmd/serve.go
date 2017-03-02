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

//Use the forked version of the go-jwt-middlware, not the auth0 version

package cmd

import (
	"net/http"
	"time"

	"github.com/goware/cors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"fmt"

	"log"

	"github.com/ecosystemsoftware/ecosystem/core"
	"github.com/pressly/chi/middleware"
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

func preServe() {

	//Check to make sure a secret has been provided
	//No default provided as a security measure, server will exit of nothing provided
	if viper.GetString("secret") == "" {
		log.Fatal("No signing secret provided")
	}

	//Set up the email server and test
	err := core.EmailSetup()
	if err != nil {
		log.Println("Error setting up email system: ", err.Error())
		log.Println("Email system will not function")
	}

	//Establish a temporary connection as the super user
	dbTemp := core.SuperUserDBConfig.ReturnDBConnection("")

	//Generate a random server password, set it and get out
	serverPW := core.RandomString(16)
	_, err = dbTemp.Exec(fmt.Sprintf(core.SQLToSetServerRolePassword, serverPW))
	if err != nil {
		log.Fatal("Error setting server role password: ", err.Error())
	}

	dbTemp.Close()

	//Establish a permanent connection
	core.DB = core.ServerUserDBConfig.ReturnDBConnection(serverPW)

}

func startServer() {

	//////////////////////
	// Middleware Stack //
	//////////////////////

	// Basic CORS
	// for more ideas, see: https://developer.github.com/v3/#cross-origin-resource-sharing
	cors := cors.New(cors.Options{
		// AllowedOrigins: []string{"https://foo.com"}, // Use this to allow specific origin hosts
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH", "SEARCH"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	})

	core.Router.Use(middleware.RequestID)
	core.Router.Use(middleware.RealIP)
	core.Router.Use(middleware.Logger)
	core.Router.Use(middleware.Recoverer)
	core.Router.Use(cors.Handler) //Activate CORS middleware

	// When a client closes their connection midway through a request, the
	// http.CloseNotifier will cancel the request context (ctx).
	core.Router.Use(middleware.CloseNotify)

	// Set a timeout value on the request context (ctx), that will signal
	// through ctx.Done() that the request has timed out and further
	// processing should be stopped.
	core.Router.Use(middleware.Timeout(60 * time.Second))

	http.ListenAndServe(":"+viper.GetString("apiPort"), core.Router)

}

//  Experimental search features
// 	api.Handle("SEARCH", "/:schema/:table/", core.ReturnBlank) //Useful for when blank searches are sent by client, to avoid errors
// 	api.Handle("SEARCH", "/:schema/:table/:searchTerm", core.SearchList)
