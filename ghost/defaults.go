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

//Defaults the the app wide default settings
var Defaults = config{

	//PG Settings
	PgSuperUser:  "postgres",
	PgDBName:     "testdb",
	PgPort:       "5432",
	PgServer:     "localhost",
	PgDisableSSL: true,

	//General Settings
	ApiPort:  "3000",
	JWTRealm: "Your App Name",
	Host:     "localhost",
	Protocol: "http",

	//Email Settings
	ActivateEmail: false,
	SmtpHost:      "smtp",
	SmtpPort:      "25",
	SmtpUserName:  "info@yourdomain.com",
	SmtpFrom:      "info@yourdomain.com",
	EmailFrom:     "Your Name",

	//Bundles installed
	BundlesInstalled: make([]string, 0, 0),

	//Global Middleware
	GlobalMiddleware: []string{"RequestID", "RealIP", "Logger", "Recoverer", "CloseNotify", "Timeout"},
	Timeout:          60,

	//CORS Settings
	ActivateCors:         false,
	CorsAllowedOrigins:   []string{"*"},
	CorsAllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH", "SEARCH"},
	CorsAllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
	CorsExposedHeaders:   []string{"Link"},
	CorsAllowCredentials: true,
	CorsMaxAge:           300,
}
