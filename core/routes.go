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

package core

import (
	jwtmiddleware "github.com/auth0/go-jwt-middleware"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/pressly/chi"
	"github.com/spf13/viper"
)

func main() {

	//JWT Authentication Middlware
	jwtMiddleware := jwtmiddleware.New(jwtmiddleware.Options{
		ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
			return []byte(viper.GetString("secret")), nil
		},
		// When set, the middleware verifies that tokens are signed with the specific signing algorithm
		// If the signing method is not constant the ValidationKeyGetter callback can be used to implement additional checks
		// Important to avoid security issues described here: https://auth0.com/blog/2015/03/31/critical-vulnerabilities-in-json-web-token-libraries/
		SigningMethod: jwt.SigningMethodHS256,
	})

	//The base slug is 'schema'
	Router.Route("/:schema", func(r chi.Router) {
		//Activate JWT middleware right at the base
		Router.Use(jwtMiddleware.Handler)
		//If a valid JWT is present, a user id and role will be assigned
		Router.Use(Authorizator)
		//Next slug is 'table'
		Router.Route("/:table", func(r chi.Router) {
			//Use middleware to add the schema, table and queries to the context
			Router.Use(AddSchemaAndTableToContext)
			Router.Get("/", ShowList)      // GET /schema/table
			Router.Post("/", InsertRecord) // PUT /schema/table
			//Final level is 'record'
			Router.Route("/:record", func(r chi.Router) {
				//Use middleware to add the record to the context
				Router.Use(AddRecordToContext)
				Router.Get("/", ShowSingle)      // GET /schema/table/record
				Router.Patch("/", UpdateRecord)  // PATCH /schema/table/record
				Router.Delete("/", DeleteRecord) // DELETE /schema/table/record

			})

		})
	})

}
