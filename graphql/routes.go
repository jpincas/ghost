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

package graphql

import (
	"github.com/ecosystemsoftware/ecosystem/core"
	"github.com/graphql-go/handler"
	"github.com/pressly/chi"
)

//SetRoutes adds the routes the router
func setRoutes() {

	h := handler.New(&handler.Config{
		Schema: &schema,
		Pretty: true,
	})

	core.Router.Route("/graphiql", func(r chi.Router) {

		r.Get("/", showGraphiqlIndex)
		r.Get("/index.html", showGraphiqlIndex)

	})

	core.Router.Route("/graphql", func(r chi.Router) {

		//JWT Authentication Middlware
		// var jwtMiddleware = jwtmiddleware.New(jwtmiddleware.Options{
		// 	ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
		// 		return []byte(viper.GetString("secret")), nil
		// 	},
		// 	// When set, the middleware verifies that tokens are signed with the specific signing algorithm
		// 	// If the signing method is not constant the ValidationKeyGetter callback can be used to implement additional checks
		// 	// Important to avoid security issues described here: https://auth0.com/blog/2015/03/31/critical-vulnerabilities-in-json-web-token-libraries/
		// 	SigningMethod: jwt.SigningMethodHS256,
		// })

		//Activate JWT middleware right at the base
		//r.Use(jwtMiddleware.Handler)
		//If a valid JWT is present, a user id and role will be assigned
		//r.Use(core.Authorizator)

		r.Handle("/", h)

	})
}
