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
	"time"

	"github.com/goware/cors"
	"github.com/pressly/chi"
	"github.com/pressly/chi/middleware"
)

func init() {

	App.Router = chi.NewRouter()

	if App.Config.ActivateCors {

		// for more ideas, see: https://developer.github.com/v3/#cross-origin-resource-sharing
		cors := cors.New(cors.Options{
			AllowedOrigins:   App.Config.CorsAllowedOrigins,
			AllowedMethods:   App.Config.CorsAllowedMethods,
			AllowedHeaders:   App.Config.CorsAllowedHeaders,
			ExposedHeaders:   App.Config.CorsExposedHeaders,
			AllowCredentials: App.Config.CorsAllowCredentials,
			MaxAge:           App.Config.CorsMaxAge, // Maximum value not ignored by any of major browsers
		})

		App.Router.Use(cors.Handler) //Activate CORS middleware

	}

	//Global router middleware setup
	for _, v := range App.Config.GlobalMiddleware {

		switch v {
		case "RequestID":
			App.Router.Use(middleware.RequestID)
		case "RealIP":
			App.Router.Use(middleware.RealIP)
		case "Logger":
			App.Router.Use(middleware.Logger)
		case "Recoverer":
			App.Router.Use(middleware.Recoverer)
		case "CloseNotify":
			// When a client closes their connection midway through a request, the
			// http.CloseNotifier will cancel the request context (ctx).
			App.Router.Use(middleware.CloseNotify)
		case "Timeout":
			// Set a timeout value on the request context (ctx), that will signal
			// through ctx.Done() that the request has timed out and further
			// processing should be stopped.
			App.Router.Use(middleware.Timeout(60 * time.Second))
		}

	}

}
