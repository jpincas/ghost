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

package core

import (
	"encoding/json"
	"io/ioutil"

	"github.com/spf13/viper"
)

//Config is the basic structure of the config.json file
type Config struct {
	PgSuperUser      string  `json:"pgSuperUser"`
	PgDBName         string  `json:"pgDBName"`
	PgPort           string  `json:"pgPort"`
	PgServer         string  `json:"pgServer"`
	PgDisableSSL     bool    `json:"pgDisableSSL"`
	ApiPort          string  `json:"apiPort"`
	SmtpHost         string  `json:"smtpHost"`
	SmtpPort         string  `json:"smtpPort"`
	SmtpUserName     string  `json:"smtpUserName"`
	SmtpFrom         string  `json:"smtpFrom"`
	EmailFrom        string  `json:"emailFrom"`
	JWTRealm         string  `json:"jwtRealm"`
	BundlesInstalled Bundles `json:"bundlesInstalled"`
	Host             string  `json:"host"`
	Protocol         string  `json:"protocol"`
}

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

func readConfig() {

	viper.SetConfigName(viper.GetString("configfile"))

	if err := viper.ReadInConfig(); err == nil {
		//Initialise the db config structs for later use
		InitDBConnectionConfigs()
		Log(LogEntry{"CORE.CONFIG", true, "Config file detected and correctly applied:" + viper.ConfigFileUsed()})
	} else {
		LogFatal(LogEntry{"CORE.CONFIG", true, "Config file not found. Aborting"})
	}

}

//createDafaultConfigFile creates the default config.json template with sane defaults
//Will overwrite existing config.json, so ask for confirmation
func createDefaultConfigFile(configFileName string) error {

	config := Config{
		PgSuperUser:      "postgres",
		PgDBName:         "testdb",
		PgPort:           "5432",
		PgServer:         "localhost",
		PgDisableSSL:     true,
		ApiPort:          "3000",
		SmtpHost:         "smtp",
		SmtpPort:         "25",
		SmtpUserName:     "info@yourdomain.com",
		SmtpFrom:         "info@yourdomain.com",
		EmailFrom:        "Your Name",
		JWTRealm:         "Your App Name",
		BundlesInstalled: make([]string, 0, 0),
		Host:             "localhost",
		Protocol:         "http",
	}

	configJSON, _ := json.MarshalIndent(config, "", "\t")
	fileName := configFileName + ".json"
	err := ioutil.WriteFile(fileName, configJSON, 0644)
	if err != nil {
		return err
	}

	return nil

}
