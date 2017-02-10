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
)

//ConnectToDB connects to the database and returns a connection pool
func ConnectToDB(dbConnectionString string) *sql.DB {
	//Initialise database
	db, _ := sql.Open("postgres", dbConnectionString)
	//Ping database to check connectivity
	if err := db.Ping(); err != nil {
		log.Fatal("Error connecting to Postgres as super user during setup. ", err.Error())
	} else {
		log.Println("Connected successfully to ", dbConnectionString)
	}
	return db
}

//GetDBConnectionString returns a correctly formated Postgres connection string
func GetDBConnectionString(user string, pwString string, server string, port string, dbName string, disableSSL bool) (dbConnectionString string) {
	//Set the password string if a password has been supplied
	if pwString != "" {
		pwString = ":" + pwString
	}
	dbConnectionString = "postgres://" + user + pwString + "@" + server + ":" + port + "/" + dbName
	//If disabled SSL flag specified,
	if disableSSL {
		dbConnectionString += "?sslmode=disable"
	}
	return
}
