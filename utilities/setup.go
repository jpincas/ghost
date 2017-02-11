// +build ignore

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
	"flag"
	"fmt"
	"log"
	"os"
)

//Command line arguments
var pgUser = flag.String("pguser", "postgres", "The Postgres superuser used for initial configuration")
var pgPW = flag.String("pgpw", "", "The Postgres password")
var pgPort = flag.String("pgport", "5432", "The Postgres password")
var pgServer = flag.String("pgserver", "localhost", "The location of the Postgres server")
var pgName = flag.String("pgname", "", "The name of the Postgres DB to connect to")
var pgDisableSSL = flag.Bool("pgdisablessl", false, "Postgres connection SSL mode")
var adminEmail = flag.String("createadminwithemail", "", "The email address of an admin user to create")
var PublicSiteSlug = flag.String("siteslug", "site", "The initial URL slug for the public site")
var PrivateSiteSlug = flag.String("privateslug", "private", "The initial URL slug for the private site")
var secret = flag.String("secret", "", "The secret used to sign JWTs")
var demoMode = flag.Bool("demomode", false, "Run the server in demo mode, allowing easy login")
var demoRole = flag.String("demorole", "admin", "The default role for demo mode")

//email
var smtpHost = flag.String("smtphost", "", "SMTP server address for outgoing email")
var smtpPort = flag.String("smtpport", "25", "SMTP server port for outgoing email")
var smtpUserName = flag.String("smtpuser", "", "SMTP server username for outgoing email")
var smtpPw = flag.String("smtppw", "", "SMTP server password for outgoing email")
var smtpFrom = flag.String("smtpfrom", "", "SMTP server from email address")
var emailFrom = flag.String("emailfrom", "", "Email from header field")

//Database connection
var db *sql.DB

func init() {

	// flag.Parse()
	// dbSetup()
	// folderSetup()
	// emailSetup()

}

func dbSetup() {

	//First, check to make sure a Postgres database has been specified.
	//There is no sensible defaut prodivded for this setting, server will exit of nothing provided
	if *pgName == "" {
		log.Fatal("No Postgres database name specified")
	}

	//Second, check to make sure a secret has been provided
	//No default provided as a security measure, server will exit of nothing provided
	if *secret == "" {
		log.Fatal("No signing secret provided")
	}

	/////////////////////////////////
	// Initial Setup as Super User //
	/////////////////////////////////

	//Establish a temporary connection as the super user
	//Set the password string if a password has been supplied
	var pwString = ""
	if *pgPW != "" {
		pwString = ":" + *pgPW
	}
	dbTempConnection := "postgres://" + *pgUser + pwString + "@" + *pgServer + ":" + *pgPort + "/" + *pgName

	//If disabled SSL flag specified,
	if *pgDisableSSL {
		dbTempConnection += "?sslmode=disable"
	}

	//Initialise database
	dbTemp, _ := sql.Open("postgres", dbTempConnection)
	//Ping database to check connectivity
	if err := dbTemp.Ping(); err != nil {
		log.Fatal("Error connecting to Postgres as super user during setup.  Check database name, server and SSL mode.")
	}

	//Setup built in tables and functions
	//Sets up the uuid extension if not already set up - must use double quote marks on the extension name
	dbTemp.Exec(`CREATE EXTENSION IF NOT EXISTS "uuid-ossp";`)
	//Sets up the users table
	dbTemp.Exec("CREATE TABLE IF NOT EXISTS users (id uuid PRIMARY KEY, email varchar(256) UNIQUE, role varchar(16) NOT NULL default 'anon');")
	//Sets up user uuid functionality
	dbTemp.Exec("CREATE OR REPLACE FUNCTION generate_new_user() RETURNS trigger AS $$ BEGIN NEW.id := uuid_generate_v4(); RETURN NEW; END; $$ LANGUAGE plpgsql;")
	dbTemp.Exec("CREATE TRIGGER new_user BEFORE INSERT ON users FOR EACH ROW EXECUTE PROCEDURE generate_new_user();")

	//Bootstraps with an admin user
	//If a user with the email already exists, nothing will happen
	if *adminEmail != "" {
		dbTemp.Exec(fmt.Sprintf("INSERT INTO users(email, role) VALUES ('%s', 'admin')", *adminEmail))
	}

	//Built in table for web categories
	dbTemp.Exec("CREATE TABLE IF NOT EXISTS web_categories (id text NOT NULL PRIMARY KEY, title text,image text,description text,subtitle text,parent text,priority integer);")

	//Create the server role - the only thing it can do is switch to other roles
	dbTemp.Exec("CREATE ROLE server NOINHERIT LOGIN;")
	dbTemp.Exec(fmt.Sprintf("ALTER ROLE server WITH PASSWORD '%s';", *secret))
	dbTemp.Exec("ALTER ROLE server VALID UNTIL 'infinity';")
	//Create other built in roles: anon, web and admin
	dbTemp.Exec("CREATE ROLE anon;")
	dbTemp.Exec("CREATE ROLE admin BYPASSRLS;")
	dbTemp.Exec("CREATE ROLE web;")
	//'server' can switch to other roles
	dbTemp.Exec("GRANT anon, admin, web TO server;")
	//Admin role has access to everything
	dbTemp.Exec("GRANT ALL ON ALL TABLES IN SCHEMA public TO admin;")
	dbTemp.Exec("GRANT USAGE ON ALL SEQUENCES IN SCHEMA public TO admin;")
	//Server role has access to the users table in order to look roles at login time
	dbTemp.Exec("GRANT SELECT ON TABLE users TO server;")
	//Web role has access to categories
	dbTemp.Exec("GRANT SELECT ON TABLE web_categories TO web;")

	//Close the temporary connection
	dbTemp.Close()

	///////////////////////////////////
	// Permanent Database Connection //
	///////////////////////////////////

	//Establish the permanent connect with built in role Server
	//Defualt looks like: postgres://server@localhost:5432/PGNAME
	dbConnection := "postgres://server:" + *secret + "@" + *pgServer + ":" + *pgPort + "/" + *pgName

	//If disabled SSL flag specified,
	//Defualt looks like: postgres://server@localhost:5432/PGNAME?sslmode=disable
	if *pgDisableSSL {
		dbConnection += "?sslmode=disable"
	}

	log.Println("Attempting to connect to ", dbConnection)

	//Initialise databse
	db, _ = sql.Open("postgres", dbConnection)
	//Ping database to check connectivity
	if err := db.Ping(); err != nil {
		log.Fatal("Could not establish connection to Postgres database.  Check database name, server and SSL mode.")
	}

}

func folderSetup() {

	os.MkdirAll("./public", os.ModePerm)
	os.MkdirAll("./public/images_resized", os.ModePerm)
	os.MkdirAll("./public/images_source", os.ModePerm)
	os.MkdirAll("./templates", os.ModePerm)

}

func emailSetup() {

	//Default the from name in the email header to the smtp from address if not provided
	if *emailFrom == "" {
		*emailFrom = *smtpFrom
	}

	//Setup the smtp config struct, and mark as not working
	mySMTPServer = smtpServer{
		host:     *smtpHost,
		port:     *smtpPort,
		password: *smtpPw,
		userName: *smtpUserName,
		from:     *smtpFrom,
		fromName: *emailFrom,
		working:  false,
	}

	//Only test the connection if complete credentials are provided
	if *smtpFrom != "" && *smtpPort != "" && *smtpPw != "" && *smtpHost != "" && *smtpUserName != "" {

		//Test the SMTP connection
		if err := testSMTPConnection(mySMTPServer); err != nil {
			//If it fails, warn and exit
			log.Fatal("Error with SMTP server:", err)
		} else {
			//If it passes, setup the config
			mySMTPServer.working = true

			//If the system has been bootstrapped with an admin user and the email system is working, lets send them a welcome email!
			// if *adminEmail != "" {

			// 	sendEmail(
			// 		mySMTPServer,
			// 		[]string{*adminEmail},  //Recipient
			// 		"Welcome to EcoSystem", //Subject
			// 		map[string]string{},    //Data to include in the email
			// 		"emailWelcome.html")    //Email template to use

			// }

		}
	} else {
		//If incomplete credentials are provided, warn, but don't fails
		log.Println("Incomplete email credentials provided.  The email system will not work")
	}

}
