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

//Very important: use the forked version of go-jwt-middlware
//jwtmiddleware "github.com/jonbonazza/go-jwt-middleware"
//The auth0 version doesn't properly support getting the claims

package core

import (
	jwt "github.com/dgrijalva/jwt-go"
	jwtmiddleware "github.com/jonbonazza/go-jwt-middleware"
	"github.com/pressly/chi"
	"github.com/spf13/viper"
)

func init() {

	//JWT Authentication Middlware
	var jwtMiddleware = jwtmiddleware.New(jwtmiddleware.Options{
		ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
			return []byte(viper.GetString("secret")), nil
		},
		// When set, the middleware verifies that tokens are signed with the specific signing algorithm
		// If the signing method is not constant the ValidationKeyGetter callback can be used to implement additional checks
		// Important to avoid security issues described here: https://auth0.com/blog/2015/03/31/critical-vulnerabilities-in-json-web-token-libraries/
		SigningMethod: jwt.SigningMethodHS256,
	})

	Router.Route("/api", func(r chi.Router) {
		//The base slug is 'schema'
		r.Route("/:schema", func(r chi.Router) {
			//Activate JWT middleware right at the base
			r.Use(jwtMiddleware.Handler)
			//If a valid JWT is present, a user id and role will be assigned
			r.Use(Authorizator)
			//Next slug is 'table'
			r.Route("/:table", func(r chi.Router) {
				//Use middleware to add the schema, table and queries to the context
				r.Use(AddSchemaAndTableToContext)
				r.Get("/", ShowList)      // GET /schema/table
				r.Post("/", InsertRecord) // PUT /schema/table
				//Final level is 'record'
				r.Route("/:record", func(r chi.Router) {
					//Use middleware to add the record to the context
					r.Use(AddRecordToContext)
					r.Get("/", ShowSingle)      // GET /schema/table/record
					r.Patch("/", UpdateRecord)  // PATCH /schema/table/record
					r.Delete("/", DeleteRecord) // DELETE /schema/table/record
				})
			})
		})
	})

}
