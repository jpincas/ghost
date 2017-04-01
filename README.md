#Ghost
A command-line application and utility package that helps you build web backend services with Go and PostgreSQL. 

## Why?

Becuase building a web server project from scratch with Go and PostgreSQL tends to involve a lot of:

-  Repetitive tasks
-  Boilerplate code
-  Reimplementing generic functionality

** Ghost ** aims to tackle these problems whilst leaving you free to write a completely custom backend service.

## How?

Ghost is both a command-line utility and also a Go package that you can import into your program.  

The command-line utility helps with mundane tasks like generating configuration files, installing SQL code into the database, creating users and so on.

The Go package (which the command-line utility is built on) gives you quick access to all the low-level utilities necessary for getting your Go/Postgres service off the ground, like configuration, database connection, routing, middleware, email etc.

In the future, Ghost will be extended with handy sub-packages.  At the moment, we have `auth`, which gives you utilites, handlers and even routes, all for dealing with authentication.  You can use the basic utilities only, use the handlers in your own routes, or just take the routes as they come, hook them into the central router and fire up.  All future Ghost pakages will work that way.

## Quickstart

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

6) In *mybundle/install/00_install.sql* write some (valid) SQL.

7) Install your new bundle with `ghost install mybundle`

### Create your custom server

1) Create a `main.go` and copy this short program:

```
package main

import (
	"fmt"
	"net/http"
	"os"

	cmds "github.com/jpincas/ghost/cmds"
	"github.com/jpincas/ghost/ghost"
)

func main() {

	//Hook into 'BeforeServe' and add simple custom route and handler
	cmds.BeforeServe = func() {
		ghost.App.Router.Get("/hello", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("Hello World"))
		})
	}

	//Run the 'Serve' command
	if err := cmds.ServeCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}

}
```

2) Build with `go build` and then run with `./myghostapp -s=secret`



BeforeServe
