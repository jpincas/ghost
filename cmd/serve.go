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
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"fmt"

	"log"

	"path"

	"github.com/ecosystemsoftware/eco/ecosql"
	"github.com/ecosystemsoftware/eco/handlers"
	eco "github.com/ecosystemsoftware/eco/utilities"
)

var nowebsite, noadminpanel bool
var smtpPW string

func init() {
	RootCmd.AddCommand(serveCmd)
	serveCmd.Flags().BoolVarP(&nowebsite, "nowebsite", "w", false, "Disable website/HTML server")
	serveCmd.Flags().BoolVarP(&noadminpanel, "noadminpanel", "a", false, "Disable admin panel server")

	serveCmd.Flags().String("smtppw", "", "SMTP server password for outgoing mail")
	viper.BindPFlag("smtppw", serveCmd.Flags().Lookup("smtppw"))

	serveCmd.Flags().BoolP("demomode", "d", false, "Run server in demo mode")
	viper.BindPFlag("demomode", serveCmd.Flags().Lookup("demmode"))

	serveCmd.Flags().StringP("secret", "s", "", "Secure secret for signing JWT")
	viper.BindPFlag("secret", serveCmd.Flags().Lookup("secret"))

}

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Starts the EcoSystem server",
	Long: `Start EcoSystem with all 3 servers: api, web and admin panel.
	Use flags to disable web and admin panel serving if you plan to host them elsewhere or ony use the API server`,
	RunE: serve,
}

func serve(cmd *cobra.Command, args []string) error {

	preServe()

	if !noadminpanel {
		serveAdminPanel()
	}
	if !nowebsite {
		serveWebsite()
	}
	serveAPI()

	return nil
}

func preServe() {

	//Check to make sure a secret has been provided
	//No default provided as a security measure, server will exit of nothing provided
	if viper.GetString("secret") == "" {
		log.Fatal("No signing secret provided")
	}

	//Set up the email server and test
	err := eco.EmailSetup()
	if err != nil {
		log.Println("Error setting up email system: ", err.Error())
		log.Println("Email system will not function")
	}

	//Establish a temporary connection as the super user
	dbTemp := eco.SuperUserDBConfig.ReturnDBConnection("")

	//Generate a random server password, set it and get out
	serverPW := eco.RandomString(16)
	_, err = dbTemp.Exec(fmt.Sprintf(ecosql.ToSetServerRolePassword, serverPW))
	_, err = dbTemp.Exec(ecosql.ToSetServerPasswordToLastForever)
	if err != nil {
		log.Fatal("Error setting server role password: ", err.Error())
	}

	dbTemp.Close()

	//Establish a permanent connection
	eco.DB = eco.ServerUserDBConfig.ReturnDBConnection(serverPW)

}

func serveAPI() {

	apiServer := gin.Default()
	apiServer.Use(eco.AllowCORS)                             //Activate CORS middleware
	apiServer.OPTIONS("/*anything", handlers.OptionsHandler) //Must allow unauthorised requests

	//Resized image route
	//Note format: /images/[IMAGE NAME WITH OPTIONAL PATH]?width=[WIDTH IN PIXELS]
	apiServer.GET("/images/*image", handlers.ShowImage) //Use star instead fo colons to allow for paths

	//Get JWT
	apiServer.POST("/login", eco.AuthMiddleware.LoginHandler) //for anonymous login post 'anon' for both username and password.  Must post both, otherwise fails
	apiServer.POST("/magiccode", handlers.ApiMagicCode)

	api := apiServer.Group("/api")

	{
		api.Use(eco.AuthMiddleware.MiddlewareFunc())
		api.Use(eco.MakeJSON) //Activate JSON Header middleware
		api.GET("/:schema/:table", handlers.ApiShowList)
		api.GET("/:schema/:table/:id", handlers.ApiShowSingle)
		api.POST("/:schema/:table", handlers.ApiInsertRecord)
		api.DELETE("/:schema/:table/:id", handlers.ApiDeleteRecord)
		api.PATCH("/:schema/:table/:id", handlers.ApiUpdateRecord)
		//Experimental: Full Text Search
		api.Handle("SEARCH", "/:schema/:table/", handlers.ReturnBlank) //Useful for when blank searches are sent by client, to avoid errors
		api.Handle("SEARCH", "/:schema/:table/:searchTerm", handlers.SearchList)
	}

	//Start the API
	apiServer.Run(":" + viper.GetString("apiPort"))

}

func serveWebsite() {

	webServer := gin.Default()

	//Check for templates and load if any found
	//Must check first otherwise crashes if no templates present
	// templates/BUNDLE_NAME/PAGES or EMAIL or PARTIALS
	if t, err := filepath.Glob("bundles/**/templates/**/*.html"); t != nil && err == nil {
		webServer.LoadHTMLGlob("bundles/**/templates/**/*.html")
	}

	//Resized image route
	//Note format: /images/[IMAGE NAME WITH OPTIONAL PATH]?width=[WIDTH IN PIXELS]
	webServer.GET("/images/*image", handlers.ShowImage) //Use star instead fo colons to allow for paths

	//Bundle public directories
	public := webServer.Group("/public")
	{
		//For each bundle present - add that bundle's public directory contents at TOPLEVEL/public/BUNDLENAME
		if bundleDirectoryContents, err := afero.ReadDir(AppFs, "bundles"); err == nil {
			for _, v := range bundleDirectoryContents {
				if v.IsDir() {
					public.StaticFS(v.Name(), http.Dir(path.Join("bundles", v.Name(), "public")))
				}
			}
		}

	}

	//Homepage and web categories
	webServer.GET("/", handlers.WebShowHomepage)
	webServer.GET("category/:schema/:table/:cat", handlers.WebShowCategory)

	//Unprotected HTML routes.  Authentiaction middleware is not activated
	//so there is no need for the browser to present a JWT
	//Database will always be queried with role 'web'.  Therefore give priveleges to this role
	//to all tables that are intended to be public
	//This is intended for the main site pages that are public and available to crawlers
	site := webServer.Group(viper.GetString("publicSiteSlug"))

	{
		site.GET(":schema/:table/:slug", handlers.WebShowSingle)
		site.GET(":schema/:table", handlers.WebShowList)
	}

	//Protected HTML routes.
	//Authentication middlware is actiaved so a JWT must be presented by the browser
	// These are used as partials when you want to
	//return formatted HTML specified to the logged in user (e.g. a cart)
	private := webServer.Group(viper.GetString("privateSiteSlug"))

	{
		private.Use(eco.AuthMiddleware.MiddlewareFunc())
		private.GET(":schema/:table", handlers.WebShowList)
		private.GET(":schema/:table/:slug", handlers.WebShowSingle)
	}

	go webServer.Run(":" + viper.GetString("websitePort"))

}

func serveAdminPanel() {

	adminServer := gin.Default()

	//Serve the Polymer app at /admin
	adminServer.StaticFS("/admin", http.Dir("ecosystem-admin/build/"+viper.GetString("adminPanelServeType")+"/"))

	//Serve bundle customisation files at /custom/[BUNDLENAME]
	custom := adminServer.Group("/custom")

	//For each bundle present - add that bundle's admin directory contents at TOPLEVEL/custom/BUNDLENAME
	if bundleDirectoryContents, err := afero.ReadDir(AppFs, "bundles"); err == nil {
		for _, v := range bundleDirectoryContents {
			if v.IsDir() {
				custom.StaticFS(v.Name(), http.Dir(path.Join("bundles", v.Name(), "admin-panel")))
			}
		}
	}

	go adminServer.Run(":" + viper.GetString("adminPanelPort"))

}
