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
	"io/ioutil"

	"github.com/spf13/viper"
)

//Set sensible viper defaults
func init() {

	viper.SetDefault("pgSuperUser", "postgres")
	viper.SetDefault("pgDBName", "testdb")
	viper.SetDefault("pgPort", "5432")
	viper.SetDefault("pgServer", "localhost")
	viper.SetDefault("pgDisableSSL", true)
	viper.SetDefault("apiPort", "3000")
	viper.SetDefault("jwtRealm", "yourappname")
	viper.SetDefault("host", "localhost")
	viper.SetDefault("protocol", "http")
	viper.AddConfigPath(".")

}

//Configuration is the app-wide configuration object
var Config config

//Config is the basic structure of the config.json file
type config struct {
	//PG Settings
	PgSuperUser  string `json:"pgSuperUser"`
	PgDBName     string `json:"pgDBName"`
	PgPort       string `json:"pgPort"`
	PgServer     string `json:"pgServer"`
	PgDisableSSL bool   `json:"pgDisableSSL"`
	//General Settings
	ApiPort  string `json:"apiPort"`
	JWTRealm string `json:"jwtRealm"`
	Host     string `json:"host"`
	Protocol string `json:"protocol"`
	//CORS Settings
	Cors                 bool     `json:"cors"`
	CorsAllowedOrigins   []string `json:"corsAllowedOrigins"`
	CorsAllowedMethods   []string `json:"corsAllowedMethods"`
	CorsAllowedHeaders   []string `json:"corsAllowedHeaders"`
	CorsExposedHeaders   []string `json:"corsExposedHeaders"`
	CorsAllowCredentials bool     `json:"corsAllowCredentials"`
	CorsMaxAge           int      `json:"corsMaxAge"`

	//Email Settings
	SmtpHost     string `json:"smtpHost"`
	SmtpPort     string `json:"smtpPort"`
	SmtpUserName string `json:"smtpUserName"`
	SmtpFrom     string `json:"smtpFrom"`
	EmailFrom    string `json:"emailFrom"`
	//Bundles installed
	BundlesInstalled Bundles `json:"bundlesInstalled"`
}

//Setup hydrates the app-wide config object by reading in a specified config file
//Also initialises the database setting config objects
func (c *config) Setup(configFileName string) {

	viper.SetConfigName(configFileName)

	if err := viper.ReadInConfig(); err == nil {

		//Initialise the db config structs for later use
		InitDBConnectionConfigs()

		//Unmarshall the whole config file into a config object
		if err := viper.Unmarshal(c); err != nil {

			LogFatal(LogEntry{"ghost.CONFIG", true, "Error decoding config file. Aborting"})

		}

		Log(LogEntry{"ghost.CONFIG", true, "Config file detected and correctly applied:" + viper.ConfigFileUsed()})

	} else {

		LogFatal(LogEntry{"ghost.CONFIG", true, "Config file not found. Aborting"})

	}

}

//createDafaultConfigFile creates the default config.json template with sane defaults
//Will overwrite existing config.json, so ask for confirmation
func createDefaultConfigFile(configFileName string) error {

	c := config{
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
		//CORS Settings
		Cors:                 false,
		CorsAllowedOrigins:   []string{"*"},
		CorsAllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH", "SEARCH"},
		CorsAllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		CorsExposedHeaders:   []string{"Link"},
		CorsAllowCredentials: true,
		CorsMaxAge:           300,
		//Email Settings
		SmtpHost:     "smtp",
		SmtpPort:     "25",
		SmtpUserName: "info@yourdomain.com",
		SmtpFrom:     "info@yourdomain.com",
		EmailFrom:    "Your Name",
		//Bundls installed
		BundlesInstalled: make([]string, 0, 0),
	}

	configJSON, _ := json.MarshalIndent(c, "", "\t")
	fileName := configFileName + ".json"
	err := ioutil.WriteFile(fileName, configJSON, 0644)
	if err != nil {
		return err
	}

	return nil

}
