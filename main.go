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

package main

import (
	_ "github.com/lib/pq"

	"github.com/ecosystemsoftware/ecosystem/cmd"
)

func main() {
	cmd.Execute()
}

// func main() {

// ///////////////////////////////////
// // Static Server for Admin Panel //
// ///////////////////////////////////

// adminServer := gin.Default()

// adminServer.StaticFS("/", http.Dir("ecosystem-admin/build/unbundled/"))

// go adminServer.Run(":8080")

// ////////////////
// // Web Server //
// ////////////////

// webServer := gin.Default()

// //Only attempt to load templates if there is at least one templates file in the directories
// c, _ := ioutil.ReadDir("./templates/custom")
// d, _ := ioutil.ReadDir("./templates/partials")
// p, _ := ioutil.ReadDir("./templates/default")

// if len(c) > 0 || len(d) > 0 || len(p) > 0 {
// 	webServer.LoadHTMLGlob("templates/**/*.html")
// }

// //Resized image route
// //Note format: /img/[IMAGE NAME WITH OPTIONAL PATH]?width=[WIDTH IN PIXELS]
// webServer.GET("/img/*image", eco.ShowImage) //Use star instead fo colons to allow for paths

// //Static file system for 'public' directory
// webServer.StaticFS("/public", http.Dir("public"))

// //Homepage and web categories
// webServer.GET("/", eco.WebShowHomepage)
// webServer.GET("category/:table/:cat", eco.WebShowCategory)

// //Unprotected HTML routes.  Authentiaction middleware is not activated
// //so there is no need for the browser to present a JWT
// //Database will always be queried with role 'web'.  Therefore give priveleges to this role
// //to all tables that are intended to be public
// //This is intended for the main site pages that are public and available to crawlers
// site := webServer.Group(*eco.PublicSiteSlug)

// {
// 	site.GET(":table/:slug", eco.WebShowSingle)
// 	site.GET(":table", eco.WebShowList)
// }

// //Protected HTML routes.
// //Authentication middlware is actiaved so a JWT must be presented by the browser
// // These are used as partials when you want to
// //return formatted HTML specified to the logged in user (e.g. a cart)
// private := webServer.Group(*eco.PrivateSiteSlug)

// {
// 	private.Use(eco.AuthMiddleware.MiddlewareFunc())
// 	private.GET(":table", eco.WebShowList)
// 	private.GET(":table/:slug", eco.WebShowSingle)
// }

// go webServer.Run(":8000")

// ////////////////
// // API Server //
// ////////////////

// apiServer := gin.Default()
// apiServer.Use(eco.AllowCORS)                        //Activate CORS middleware
// apiServer.OPTIONS("/*anything", eco.OptionsHandler) //Must allow unauthorised requests

// //Resized image route
// //Note format: /img/[IMAGE NAME WITH OPTIONAL PATH]?width=[WIDTH IN PIXELS]
// apiServer.GET("/img/*image", eco.ShowImage) //Use star instead fo colons to allow for paths

// //Get JWT
// apiServer.POST("/login", eco.AuthMiddleware.LoginHandler) //for anonymous login post 'anon' for both username and password.  Must post both, otherwise fails
// apiServer.POST("/magiccode", eco.ApiMagicCode)

// api := apiServer.Group("/api")

// {
// 	api.Use(eco.AuthMiddleware.MiddlewareFunc())
// 	api.Use(eco.MakeJSON) //Activate JSON Header middleware
// 	api.GET("/:table", eco.ApiShowList)
// 	api.GET("/:table/:id", eco.ApiShowSingle)
// 	api.POST("/:table", eco.ApiInsertRecord)
// 	api.DELETE("/:table/:id", eco.ApiDeleteRecord)
// 	api.PATCH("/:table/:id", eco.ApiUpdateRecord)
// 	//Experimental: Full Text Search
// 	api.Handle("SEARCH", "/:table/", eco.ReturnBlank) //Useful for when blank searches are sent by client, to avoid errors
// 	api.Handle("SEARCH", "/:table/:searchTerm", eco.SearchList)

// }

// //Start the API
// apiServer.Run(":3000")

// }
