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
	PgSuperUser              string  `json:"pgSuperUser"`
	PgDBName                 string  `json:"pgDBName"`
	PgPort                   string  `json:"pgPort"`
	PgServer                 string  `json:"pgServer"`
	PgDisableSSL             bool    `json:"pgDisableSSL"`
	ApiPort                  string  `json:"apiPort"`
	WebsitePort              string  `json:"websitePort"`
	AdminPanelPort           string  `json:"adminPanelPort"`
	AdminPanelServeDirectory string  `json:"adminPanelServeDirectory"`
	PublicSiteSlug           string  `json:"publicSiteSlug"`
	PrivateSiteSlug          string  `json:"privateSiteSlug"`
	SmtpHost                 string  `json:"smtpHost"`
	SmtpPort                 string  `json:"smtpPort"`
	SmtpUserName             string  `json:"smtpUserName"`
	SmtpFrom                 string  `json:"smtpFrom"`
	EmailFrom                string  `json:"emailFrom"`
	JWTRealm                 string  `json:"jwtRealm"`
	AdminPrimaryColor        string  `json:"adminPrimaryColor"`
	AdminSecondaryColor      string  `json:"adminSecondaryColor"`
	AdminTextColor           string  `json:"adminTextColor"`
	AdminErrorColor          string  `json:"adminErrorColor"`
	AdminTitle               string  `json:"adminTitle"`
	AdminLogoFile            string  `json:"adminLogoHorizontal"`
	AdminLogoBundle          string  `json:"adminLogoVertical"`
	BundlesInstalled         Bundles `json:"bundlesInstalled"`
	Host                     string  `json:"host"`
	Protocol                 string  `json:"protocol"`
}

// initConfig reads in config file and ENV variables if set.
func InitConfig() {

	viper.SetDefault("pgSuperUser", "postgres")
	viper.SetDefault("pgDBName", "testdb")
	viper.SetDefault("pgPort", "5432")
	viper.SetDefault("pgServer", "localhost")
	viper.SetDefault("pgDisableSSL", true)
	viper.SetDefault("apiPort", "3000")
	viper.SetDefault("websitePort", "3001")
	viper.SetDefault("adminPanelPort", "3002")
	viper.SetDefault("adminPanelServeDirectory", "ecosystem-admin/build/unbundled")
	viper.SetDefault("publicSiteSlug", "site")
	viper.SetDefault("privateSiteSlug", "private")
	viper.SetDefault("jwtRealm", "yourappname")
	viper.SetDefault("host", "localhost")
	viper.SetDefault("protocol", "http")

	//For the admin panel
	viper.SetDefault("adminPrimaryColor", "#00c4a7")
	viper.SetDefault("adminSecondaryColor", "#7EC9A2")
	viper.SetDefault("adminTextColor", "black")
	viper.SetDefault("adminErrorColor", "red")
	viper.SetDefault("adminTitle", "Admin Panel")
	viper.SetDefault("adminLogoFile", "logo.png")
	viper.SetDefault("adminLogoBundle", "master")
	viper.SetDefault("bundlesInstalled", make([]string, 0, 0))

	viper.SetConfigName(viper.GetString("configfile"))
	viper.AddConfigPath(".")

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		//Initialise the db config structs for later use
		InitDBConnectionConfigs()
		Log(LogEntry{"CORE.CONFIG", true, "Config file detected and correctly applied:" + viper.ConfigFileUsed()})
	} else {
		//Otherwise create one
		Log(LogEntry{"CORE.CONFIG", true, "Config file does not exist. Creating now..."})

		if err := createDefaultConfigFile(); err != nil {
			LogFatal(LogEntry{"CORE.CONFIG", false, "Error creating config file: " + err.Error()})
		}
	}
}

//createDafaultConfigFile creates the default config.json template with sane defaults
//Will overwrite existing config.json, so ask for confirmation
func createDefaultConfigFile() error {

	config := Config{
		PgSuperUser:              "postgres",
		PgDBName:                 "testdb",
		PgPort:                   "5432",
		PgServer:                 "localhost",
		PgDisableSSL:             true,
		ApiPort:                  "3000",
		WebsitePort:              "3001",
		AdminPanelPort:           "3002",
		AdminPanelServeDirectory: "ecosystem-admin/build/unbundled",
		PublicSiteSlug:           "site",
		PrivateSiteSlug:          "private",
		SmtpHost:                 "smtp",
		SmtpPort:                 "25",
		SmtpUserName:             "info@yourdomain.com",
		SmtpFrom:                 "info@yourdomain.com",
		EmailFrom:                "Your Name",
		JWTRealm:                 "Your App Name",
		AdminPrimaryColor:        "#00c4a7",
		AdminSecondaryColor:      "#7EC9A2",
		AdminTextColor:           "black",
		AdminErrorColor:          "red",
		AdminTitle:               "Admin Panel",
		AdminLogoFile:            "logo.png",
		AdminLogoBundle:          "master",
		BundlesInstalled:         make([]string, 0, 0),
		Host:                     "localhost",
		Protocol:                 "http",
	}

	configJSON, _ := json.MarshalIndent(config, "", "\t")
	err := ioutil.WriteFile("config.json", configJSON, 0644)
	if err != nil {
		return err
	}

	return nil

}
