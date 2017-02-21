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

package utilities

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
	"github.com/spf13/viper"
)

//DB is the main shared database connection pool for the application
var DB *sql.DB

//dbConfig holds all the necessary information for a datbase connection
type dbConfig struct {
	user, pw, server, port, dbName string
	disableSSL                     bool
}

//SuperUserDBConfig is the connection configuration for the super user
var SuperUserDBConfig dbConfig

//ServerUserDBConfig is the connection configuration for the 'server' role user
var ServerUserDBConfig dbConfig

//InitDBConnectionConfigs initialises the config structs once configs have come in from Viper
func InitDBConnectionConfigs() {

	//SuperUserDBConfig is the configuration struct for the connection to the DB as a super user
	//This configuration is generally used only temporarily during setup operations
	SuperUserDBConfig = dbConfig{
		user:       viper.GetString("pgSuperUser"),
		pw:         viper.GetString("pgpw"),
		server:     viper.GetString("pgServer"),
		port:       viper.GetString("pgPort"),
		dbName:     viper.GetString("pgDBName"),
		disableSSL: viper.GetBool("pgDisableSSL"),
	}

	//ServerUserDBConfig is the configuration struct for the connection to the DB as 'server' role
	//This is configuration is used as the permanent shared DB connection pool for the application
	ServerUserDBConfig = dbConfig{
		user:       "server",
		server:     viper.GetString("pgServer"),
		port:       viper.GetString("pgPort"),
		dbName:     viper.GetString("pgDBName"),
		disableSSL: viper.GetBool("pgDisableSSL"),
	}

}

//ReturnDBConnection returns a DB connection pool using the connection parameters in a dbConfig struct
//and an optional server password which can be passed in
func (d dbConfig) ReturnDBConnection(serverPW string) *sql.DB {

	dbConnectionString := d.getDBConnectionString(serverPW)
	return connectToDB(dbConnectionString)

}

//getDBConnectionString returns a correctly formated Postgres connection string from
//the config struct.  If there is no pw in the struct (as is the case for )
func (d dbConfig) getDBConnectionString(serverPW string) (dbConnectionString string) {

	//If this is a connection for a server role, use the password supplied as a parameter
	//Otherwise ignore that parameter
	if d.user == "server" {
		d.pw = serverPW
	}

	//Set the password string if a password has been supplied
	//If not leave it blank - this stops any errors for blank passwords
	pwString := ""
	if d.pw != "" {
		pwString = ":" + d.pw
	}
	dbConnectionString = "postgres://" + d.user + pwString + "@" + d.server + ":" + d.port + "/" + d.dbName
	//If disabled SSL flag specified,
	if d.disableSSL {
		dbConnectionString += "?sslmode=disable"
	}
	return
}

//connectToDB connects to the database and returns a connection pool
func connectToDB(dbConnectionString string) *sql.DB {
	//Initialise database
	log.Println("Attempting to connect to ", dbConnectionString)
	dbConnection, _ := sql.Open("postgres", dbConnectionString)
	//Ping database to check connectivity
	if err := dbConnection.Ping(); err != nil {
		log.Fatal("Error connecting to Postgres as super user during setup. ", err.Error())
	} else {
		log.Println("Connected successfully to ", dbConnectionString)
	}
	return dbConnection
}
