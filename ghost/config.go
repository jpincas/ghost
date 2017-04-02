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
	"errors"
	"io/ioutil"

	"github.com/spf13/viper"
)

type Bundles []string

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

	//Email Settings
	ActivateEmail bool   `json:"activateEmail"`
	SmtpHost      string `json:"smtpHost"`
	SmtpPort      string `json:"smtpPort"`
	SmtpUserName  string `json:"smtpUserName"`
	SmtpFrom      string `json:"smtpFrom"`
	EmailFrom     string `json:"emailFrom"`

	//Bundles installed
	BundlesInstalled Bundles `json:"bundlesInstalled"`

	//Global middleware activation
	GlobalMiddleware []string `json:"globalMiddleware"`
	Timeout          int      `json:"timout"`

	//CORS Settings
	ActivateCors         bool     `json:"activateCors"`
	CorsAllowedOrigins   []string `json:"corsAllowedOrigins"`
	CorsAllowedMethods   []string `json:"corsAllowedMethods"`
	CorsAllowedHeaders   []string `json:"corsAllowedHeaders"`
	CorsExposedHeaders   []string `json:"corsExposedHeaders"`
	CorsAllowCredentials bool     `json:"corsAllowCredentials"`
	CorsMaxAge           int      `json:"corsMaxAge"`
}

//createDafaultConfigFile creates the default config.json template with sane defaults
//TODO: Will overwrite existing config.json, so ask for confirmation
func CreateDefaultConfigFile(configFileName string) error {

	configJSON, _ := json.MarshalIndent(Defaults, "", "\t")
	fileName := configFileName + ".json"
	err := ioutil.WriteFile(fileName, configJSON, 0644)
	if err != nil {
		return err
	}

	return nil

}

//Setup hydrates the app-wide config object by reading in a specified config file
//Also initialises the database setting config objects
func (c *config) Setup(configFileName string) {

	viper.AddConfigPath(".")
	viper.SetConfigName(configFileName)

	if err := viper.ReadInConfig(); err == nil {

		//Unmarshall the whole config file into a config object
		if err := viper.Unmarshal(c); err != nil {
			LogFatal("CONFIG", true, "Error decoding config file. Aborting", err)
		}

		Log("CONFIG", true, "Config file detected and correctly applied:"+viper.ConfigFileUsed(), nil)

	} else {

		LogFatal("CONFIG", true, "Config file not found. Aborting", nil)

	}

}

func (c *config) InstallBundle(bundleName string) error {

	b := c.BundlesInstalled
	//Check if the bundle is already installed (should only happen if user has messed with config.json)
	//If the name of the bundle being installed coincides with any of the names already in the bundle slice,
	//then just return the original bundle slice
	for _, a := range b {
		if a == bundleName {
			return errors.New("Bundle is already installed")
		}
	}
	//Otherwise append
	b = append(b, bundleName)
	//Reset the bundle list on the config object
	c.BundlesInstalled = b

	return nil

}

func (c *config) UnInstallBundle(bundleName string) error {

	b := c.BundlesInstalled
	//Search for the bundle to be uninstalled
	for index, a := range b {
		if a == bundleName {
			//If found, splice it out
			c.BundlesInstalled = append(b[:index], b[index+1:]...)
			return nil
		}
	}

	return errors.New("Bundle is not installed")

}

func compareBundles(b1, b2 Bundles) bool {
	//If lengths are not equal
	if len(b1) != len(b2) {
		return false
	}

	//If any of the elements are not the same
	for k := range b1 {
		if b1[k] != b2[k] {
			return false
		}
	}

	return true
}
