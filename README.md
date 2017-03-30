**Warning: This project is under active heavy development, and the api is changing constantly.  Furthermore, the documentation lags behind the code. Please contact me if you're interested in using it now and I'll help you get set up.**

# Ghost

ghost is a platform that allows you to quickly develop completely custom data-driven websites, web-applications and backend systems.  ghost doesn't assume anything about the data structure or logic of your business because you code it directly and freely at the database layer using the power of PostgreSQL.  ghost then augments your database layer with a generic web and JSON API server written in Go, giving you a *bridge* to your data and logic with minimal boilerplate.

## Features

-  GraphQL (currently very basic)
-  REST API
-  Authorisation/authentication with JWTs
-  Email messaging
-  Server-rendered HTML websites (currently being redeveloped - avaiable soon)
-  Image resizing (currently being redeveloped - avaiable soon)
-  Admin Dashboard (currently being redeveloped - avaiable soon)

## What's in the box?

**Short version:** ghost is a JSON API server written in Go (this repository).

**Slightly longer version:** this repository contains the Go source code for the ghost application, which is actually a command-line tool that allows you to easily initialise, configure and launch the ghost server, as well as perform operations on **bundles** (more on those later).  If you are familiar with Go, you can also easily build a custom version of the server by including the packages (cmd, ecosql, handlers and utilities) as dependencies in your project. If you'd like to do that, check out the [Custom ghost Server](https://github.com/jpincas/ghost-server-custom) repository for some instructions and a template to work from.


## Getting Started

If you have Go installed and access to a Postgres server, it should take you less than 60 seconds to get up and running.

### Prerequisites

- Go installed and environment variables correctly configured
- A Postgres database server (local or remote) to which you can connect

### Step 1: Install ghost

To start with, fetch the source code and dependencies and compile into an executable which will go into your $GOPATH/bin directory.  That's just `go get github.com/jpincas/ghost`.  Assuming your [$GOPATH/bin is part of your PATH](https://golang.org/doc/install), you should now be able to run 'ghost' from anywhere on the command line.

### Step 2:  Create a database

Log into your Postgres server and create a new database. Call it *testdb* - that way you won't have to make any changes to the default database configuration.

### Step 3: Create a project folder

Make a new folder from which to run the server and `cd` into it:

### Step 4: Create a configuration file

Just type `ghost` to get going.  Since you don't have a defualt *config,json*, ghost will generate one for you.

### Step 5: Configure the database connection

If you're working locally, the defaults will probably just work out-of-the-box.  Otherwise, open *config.json* and edit the database connection parameters.

### Step 6: Initialise ghost

ghost needs to create a number of built-in tables, roles, permissions and functions, as well as a few folders, so just type `ghost init` to have it do that for you.

### Step 7: Create an admin user

Set yourself up as an admin user with full permissions by typing `ghost new user [your@email.com] --admin`

### Step 8: Download and install a bundle

The easiest way to get started with ghost is by installing an existing bundle.  Bundles are just folders that group together everything ghost needs.  Once you get familiar with ghost, you'll create your own bundles, but for now, let's clone a simple demo bundle and install it with demo data:

```
$ cd bundles
$ git clone git@github.com:jpincas/eco_bundle_dogshelter.git
$ cd ..
$ ghost install eco_bundle_dogshelter --demodata
```

### Step 9: Run the server

Run the server in 'demo' mode (disabling authorisation) with `ghost serve -s=secret -d`


## Congrats!  You now have a backend powered by ghost

The JSON API is at `/api` (believe it or not).  Remember schema and table/view names map directly to API endpoints, so `ghost_bundle_dogshelter.dogs_available` is at `localhost:3000/api/eco_bundle_dogshelter/dogs_available`

All requests to the API need to go through authorisation, so if you'd like to send some test requests, you'll need to start by requesting a JWT from **/login** - just include a JSON body in the request like this:

```
{
    "username": [your@email.com] //The email of the admin account you set up earlier
    "password": "123456" //Since you're runnning in demo mode, just use the password 123456 
}
```

grab the JWT that comes back and include it with your requests in an authorization header: `Authorization: Bearer xxxxxx...` - you'll have full admin permissions since you set yourself up as an admin user earlier.


***The above is real whistlestop tour of ghost and you probably have questions.  We're working hard to release better documentation, videos, tutorials etc, but in the meantime, please check out [the developer section of the ghost website](www.ghost.software/developers)***

## Configuration

### config.json

| Attribute                | Function                                 | Default                         |
| ------------------------ | ---------------------------------------- | ------------------------------- |
| pgSuperUser              | Username of database superuser with which to connect for initial setup | postgres                        |
| pgDBName                 | Name of the Postgres DB to connect to    | testdb                          |
| pgPort                   | Server port for the Postgres connection  | 5432                            |
| pgServer                 | Database server to connect to            | localhost                       |
| pgDisableSSL             | Disable SSL mode on DB connection        | FALSE                           |
| apiPort                  | Port to serve the ghost API on       | 3000                            |
| smtpHost                 | SMTP server for outgoing emails          |                                 |
| smtpPort                 | SMTP port for outgoing emails            |                                 |
| smtpUserName             | SMTP authentication username             |                                 |
| smtpFrom                 | 'From' address for outgoing emails       |                                 |
| emailFrom                | 'From' name for outgoing emails          |                                 |
| jwtRealm                 | Realm parameter for JWT authentication tokens | yourappname                     |
| bundlesInstalled         | An automatically maintained list of bundles installed.  If you use `ghost install` and `ghost uninstall` commands, you shouldn't need to touch this. |               
| host                     | The server name or IP address            | localhost                       |
| protocol                 | Serve with http or https                 | http                            |

### Command-Line

For convenience and security, some configurations are specified on the command line when using `ghost` commands. In general, type `ghost --help` for a list of commands and available flags.

Running `ghost ` on its own, with no command or arguments, verifies that the configuration file is present and readable.  If not, a new default *config.json* will be created for you, with all possible attributes.

|             |            | Flags                                    |                                          |
| ----------- | ---------- | ---------------------------------------- | ---------------------------------------- |
| All         |            | -p, —pgpw                                | **Optional:** Postgres super user password, if required |
|             |            | -c, —configfile                          | **Optional:** Specify a different configuration file from the default *config.json* |
| `init`      |            |                                          | Performs full initialisation - DB and folders |
| `init`      | `db`       |                                          | Performs the DB initialisation for built-in tables, roles and permissions |
| `init`      | `folders`  |                                          | Creates the ghost folder structure   |
| `install`   | `[bundle]` | —demodata; -r, —reinstall                | Install named ghost bundle.  Bundle folder must be downloaded/cloned into /bundles first. **Optional:** Include demo data with the bundle install. **Optional:** Uninstall the bundle before installing |
| `uninstall` | `[bundle]` |                                          | Uninstall named ghost bundle.  Will not delete the bundle folder from /bundles |
| `new`       | `bundle`   |                                          | Creates the folder structure for a new ghost bundle |
| `new`       | `user`     | —admin                                   | Creates a new user. **Optional:** make user with admin permissions |
| `serve`     |            | -d, —demomode; -s, —secret; —smtppw | Starts the ghost server. **Optional:** Run the server in demo mode, allowing users to log in with magic code '123456', rather than having to request a code. **Required:** Encryption secret for JWT signing.  Remember to use the same secret every time you start the server, or JWTs previously issued will not be valid. |

### SMTP Configuration

SMTP settings for outgoing email are not strictly required, but the login system will not work without them, since it uses a passwordless process via email.  You can disable passwordless authentication for testing by running the server in demo mode.

If email credentials are provided, the connection will be tested, but startup will not fail if the test fails - the email system will simply be disabled and a warning displayed.

### Demo mode

As a convenience, the server can be run in *demo mode* with the flag `--demomode`.  In this mode, passwordless authentication is disabled, and users can log on with the code '123456'.  Note that the user must still be created in the database and a role assigned.

For example, if you wanted to demo admin panel functionality, you might create a user with the role 'admin' and the email 'test@test.com' - you'd then tell users to log in with that email and password '123456'.

## Licence

**Build freely with ghost**.  The ghost Server and The ghost Admin Panel App are licensed under the [Apache 2.0 License](https://www.apache.org/licenses/LICENSE-2.0).  The content on the [ghost website] (http://www.ghost.software) is licensed under a [Creative Commons Attribution-NonCommercial-NoDerivatives 4.0 International License](http://creativecommons.org/licenses/by-nc-nd/4.0/).  Neither licence grants permission to use the trade names, trademarks, service marks, or product names of ghost Software LLP, including the ghost logo and symbol, except as required for reasonable and customary use.

