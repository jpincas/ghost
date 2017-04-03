![](https://github.com/jpincas/ghost/raw/master/logo.png)

# Ghost
A command-line application and utility package that helps you build web backend services with Go and PostgreSQL (previously known as EcoSystem). 

## Why?

I tend to write backends with a very thick database layer (PostgreSQL) and a very thin application (Go) layer.  Rather than recreate brittle functionality in application code, I prefer to take advantage of the power of Postgres, with its roles, permissions, grants, views, row-level security, functions and triggers.  I keep my business and data-control logic right in the database and totally avoid application level ORMs and suchlike.  

This leads naturally to mostly generic application code, of the CRUDy, boilerplate type.  In fact, this approach leads to application code that is so generic that you could, in theory, avoid writing any at all.  The excellent projects [PostgREST](https://github.com/begriffs/postgrest) and [PostGraphQL](https://github.com/postgraphql/postgraphql) take the same approach and give you a REST API and/or GraphQL API respectively, without the need to touch any application code.  I strongly recommend taking a look at those projects to see if they fulfill your requirements.

Ghost shares a common philosophy with the above projects - indeed Ghost has a generic 'PostgreSQL-to-REST' subpackage available should you need it.  If you like this database-first approach, but need to write your own, custom server, then Ghost could help you reduce:

-  Repetitive tasks
-  Boilerplate code
-  Reimplementing generic functionality


## How?

Ghost is both a command-line utility and a Go package that you can import into your program.  

The command-line utility helps with mundane tasks like generating configuration files, installing SQL code into the database, creating users and so on.

The Go package (which the command-line utility is built on) gives you quick access to all the low-level utilities necessary for getting your Go/Postgres service off the ground, like configuration, database connection, routing, middleware, email etc.

In the future, Ghost will be extended with handy sub-packages.  At the moment, we have `auth`, which gives you utilites, handlers and even routes, all for dealing with authentication.  You can use the basic utilities only, use the handlers in your own routes, or just take the routes as they come, hook them into the central router and fire up.  All future Ghost pakages will work that way.

## Hello World

You should have Go (> 1.7) already installed and your $GOPATH correctly configured.  You should also have a PostgreSQL server somewhere that you can access - easiest for development would be to have one on *localhost:5432*.

### Install and set up

1) Fetch the source code and dependencies and compile into an executable which will go into your $GOPATH/bin directory.  That's just `go get github.com/jpincas/ghost`.  Assuming your [$GOPATH/bin is part of your PATH](https://golang.org/doc/install), you should now be able to run 'ghost' from anywhere on the command line.

2) Log into your Postgres server and create a new database. Call it *testdb* - that way you won't have to make any changes to the default database configuration.

3) Make a new folder `myghostapp` in your Go path and `cd` into it:

### Use the command-line application to bootstrap your project

1) Just type `ghost` to get going.  This will give you a default `config.json`

2) If you're working locally, the defaults will probably just work out-of-the-box.  Otherwise, open *config.json* and edit the database connection parameters.

3) Ghost needs to create a number of built-in tables, roles, permissions and functions, as well as a few folders, so just type `ghost init` to have it do that for you.

4) Set yourself up as an admin user with full permissions by typing `ghost new user [your@email.com] --admin`.

5) Create a new 'bundle' (more on those later) by entering `ghost new bundle mybundle`.

6) In *mybundle/install/00_install.sql*, paste this SQL to create a table: 

```sql
CREATE TABLE helloworld(
	hello text);
```

7) In *mybundle/demodata/00_demodata.sql*, paste this SQL to add a new row: 

```sql
INSERT INTO helloworld(hello)
VALUES ('world');
```

8) Install your new bundle and demo data with `ghost install mybundle --demodata`

### Create and run a simple custom server

1) Create `main.go` and copy this short program:

```go
package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/jpincas/ghost/ghost"
)

func main() {

	//Hook into Ghost's BeforeServe' callback
	//Add our custom route 'hello' which triggers the 'helloWorld' handler
	ghost.BeforeServe = func() {

		ghost.App.Router.Get("/hello", helloWorld)
	}

	//Run Ghost's 'Serve' command to start the server
	if err := ghost.ServeCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}

}

//our helloWorld handler
func helloWorld(w http.ResponseWriter, r *http.Request) {

	var json string //to hold our db result

	//Use Ghost's convenience functions to build a simple SQL query
	// We ask for all fields from table 'helloworld'
	// executing the query as 'admin' and returning JSON from Postgres
	sql := ghost.SqlQuery(fmt.Sprintf(ghost.SQLToSelectAllFieldsFrom, "mybundle", "helloworld")).RequestSingleResultAsJSONObject().SetQueryRole("admin").ToSQLString()

	ghost.App.DB.QueryRow(sql).Scan(&json)
	w.Write([]byte(json))
	return

}


```

2) Build with `go build` and then run in 'debug' mode with `./myghostapp -s=secret -b`

3) Visit *localhost:3000/hello* with your browser and get the response `{"hello":"world"}`. Notice how Ghost logs the SQL query executed since we ran it in debug mode.