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
	"database/sql"
	"time"

	"github.com/diegobernardes/ttlcache"
	"github.com/pressly/chi"
	"github.com/spf13/afero"
)

//App is the container for the app-wide constructs like database, router and mailserver
var App application

type application struct {
	//Mailserver is the app-wide SMTP server construct for sending email
	MailServer smtpServer
	//Router is the main router - hook into it with custom routes
	Router *chi.Mux
	//DB is the main database connection pool
	DB *sql.DB
	//Config is the main application configuration object
	Config config
	//FileSystem is the main FileSystem
	FileSystem afero.Fs
	//Store is the data abstraction layer
	//Normally your applications would interact with Store rather than DB or Cache
	Store store
	//Cache is the app wide cache for SQL queries
	QueryCache *ttlcache.Cache
}

//Setup bootstraps the whole application
func (a *application) Setup(configFileName string) {

	//Setup the config
	a.Config.Setup(configFileName)

	//Initialise the db config structs for later use
	SuperUserDBConfig.SetupConnection(true)
	ServerUserDBConfig.SetupConnection(false)

	//Initialise the filesysem
	a.FileSystem = afero.NewOsFs()

	//Initialise the cache
	a.QueryCache.SetTTL(time.Duration(60 * time.Second))

}
