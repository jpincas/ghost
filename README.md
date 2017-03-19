![](https://raw.githubusercontent.com/ecosystemsoftware/ecosystem-website/master/themes/ecosystem/static/images/ecosystem-logo.png)

### Quick Update
**16th March 2017**: I've now completed the first large refactor and updated master.  The Gin router has been switched for Chi, which is is compatible with default Go handlers, making life a little easier. Test coverage is coming along nicely. Functionality has been reorganised into 'packages' (which are just Go packages), allowing for much easier extensibility.  *Core, Auth and Email* packages are included in the standard build - everything else (HTML server, image resizer and admin-panel server) has been stripped out and will be rewritten and uploaded as seperate packages that can be added to a [custom build](https://github.com/ecosystemsoftware/ecosystem-custom-server).  We're rethinking our strategy on the admin panel side of things, as the Polymer app was proving hard to install for newcomers and was causing some build problems when trying to reference external HTML imports.  We might go with Elm.  I'm also considering the possiblity of ditching the HTML server in favour of a static site integrator powered by Hugo.  Please get in touch if you'd like to help out!


# EcoSystem

EcoSystem is a platform that allows you to quickly develop completely custom data-driven websites, web-applications and backend systems.  EcoSystem doesn't assume anything about the data structure or logic of your business because you code it directly and freely at the database layer using the power of PostgreSQL.  EcoSystem then augments your database layer with a generic web and JSON API server written in Go, giving you a *bridge* to your data and logic with minimal boilerplate.

## What's in the box?

**Short version:** EcoSystem is a JSON API server written in Go (this repository).

**Slightly longer version:** this repository contains the Go source code for the EcoSystem application, which is actually a command-line tool that allows you to easily initialise, configure and launch the EcoSystem server, as well as perform operations on **bundles** (more on those later).  If you are familiar with Go, you can also easily build a custom version of the server by including the packages (cmd, ecosql, handlers and utilities) as dependencies in your project. If you'd like to do that, check out the [Custom EcoSystem Server](https://github.com/ecosystemsoftware/ecosystem-server-custom) repository for some instructions and a template to work from.


## Getting Started

If you have Go installed and access to a Postgres server, it should take you less than 60 seconds to get up and running.

### Prerequisites

- Go installed and environment variables correctly configured
- A Postgres database server (local or remote) to which you can connect

### Step 1: Install EcoSystem

To start with, fetch the source code and dependencies and compile into an executable which will go into your $GOPATH/bin directory.  That's just `go get github.com/ecosystemsoftware/ecosystem`.  Assuming your [$GOPATH/bin is part of your PATH](https://golang.org/doc/install), you should now be able to run 'ecosystem' from anywhere on the command line.

### Step 2:  Create a database

Log into your Postgres server and create a new database. Call it *testdb* - that way you won't have to make any changes to the default database configuration.

### Step 3: Create a project folder

Make a new folder from which to run the server and `cd` into it:

### Step 4: Create a configuration file

Just type `ecosystem` to get going.  Since you don't have a defualt *config,json*, EcoSystem will generate one for you.

### Step 5: Configure the database connection

If you're working locally, the defaults will probably just work out-of-the-box.  Otherwise, open *config.json* and edit the database connection parameters.

### Step 6: Initialise EcoSystem

EcoSystem needs to create a number of built-in tables, roles, permissions and functions, as well as a few folders, so just type `ecosystem init` to have it do that for you.

### Step 7: Create an admin user

Set yourself up as an admin user with full permissions by typing `ecosystem new user [your@email.com] --admin`

### Step 8: Download and install a bundle

The easiest way to get started with EcoSystem is by installing an existing bundle.  Bundles are just folders that group together everything EcoSystem needs.  Once you get familiar with EcoSystem, you'll create your own bundles, but for now, let's clone a simple demo bundle and install it with demo data:

```
$ cd bundles
$ git clone git@github.com:ecosystemsoftware/eco_bundle_dogshelter.git
$ cd ..
$ ecosystem install eco_bundle_dogshelter --demodata
```

### Step 9: Run the server

Run the server in 'demo' mode (disabling authorisation) with `ecosystem serve -s=secret -d`


## Congrats!  You now have a backend powered by EcoSystem

The JSON API is at `/api` (believe it or not).  Remember schema and table/view names map directly to API endpoints, so `ecosystem_bundle_dogshelter.dogs_available` is at `localhost:3000/api/eco_bundle_dogshelter/dogs_available`

All requests to the API need to go through authorisation, so if you'd like to send some test requests, you'll need to start by requesting a JWT from **/login** - just include a JSON body in the request like this:

```
{
    "username": [your@email.com] //The email of the admin account you set up earlier
    "password": "123456" //Since you're runnning in demo mode, just use the password 123456 
}
```

grab the JWT that comes back and include it with your requests in an authorization header: `Authorization: Bearer xxxxxx...` - you'll have full admin permissions since you set yourself up as an admin user earlier.


***The above is real whistlestop tour of EcoSystem and you probably have questions.  We're working hard to release better documentation, videos, tutorials etc, but in the meantime, please check out [the developer section of the EcoSystem website](www.ecosystem.software/developers)***

## Configuration

### config.json

| Attribute                | Function                                 | Default                         |
| ------------------------ | ---------------------------------------- | ------------------------------- |
| pgSuperUser              | Username of database superuser with which to connect for initial setup | postgres                        |
| pgDBName                 | Name of the Postgres DB to connect to    | testdb                          |
| pgPort                   | Server port for the Postgres connection  | 5432                            |
| pgServer                 | Database server to connect to            | localhost                       |
| pgDisableSSL             | Disable SSL mode on DB connection        | FALSE                           |
| apiPort                  | Port to serve the EcoSystem API on       | 3000                            |
| smtpHost                 | SMTP server for outgoing emails          |                                 |
| smtpPort                 | SMTP port for outgoing emails            |                                 |
| smtpUserName             | SMTP authentication username             |                                 |
| smtpFrom                 | 'From' address for outgoing emails       |                                 |
| emailFrom                | 'From' name for outgoing emails          |                                 |
| jwtRealm                 | Realm parameter for JWT authentication tokens | yourappname                     |
| bundlesInstalled         | An automatically maintained list of bundles installed.  If you use `ecosystem install` and `ecosystem uninstall` commands, you shouldn't need to touch this. |               
| host                     | The server name or IP address            | localhost                       |
| protocol                 | Serve with http or https                 | http                            |

### Command-Line

For convenience and security, some configurations are specified on the command line when using `ecosystem` commands. In general, type `ecosystem --help` for a list of commands and available flags.

Running `ecosystem ` on its own, with no command or arguments, verifies that the configuration file is present and readable.  If not, a new default *config.json* will be created for you, with all possible attributes.

|             |            | Flags                                    |                                          |
| ----------- | ---------- | ---------------------------------------- | ---------------------------------------- |
| All         |            | -p, —pgpw                                | **Optional:** Postgres super user password, if required |
|             |            | -c, —configfile                          | **Optional:** Specify a different configuration file from the default *config.json* |
| `init`      |            |                                          | Performs full initialisation - DB and folders |
| `init`      | `db`       |                                          | Performs the DB initialisation for built-in tables, roles and permissions |
| `init`      | `folders`  |                                          | Creates the EcoSystem folder structure   |
| `install`   | `[bundle]` | —demodata; -r, —reinstall                | Install named EcoSystem bundle.  Bundle folder must be downloaded/cloned into /bundles first. **Optional:** Include demo data with the bundle install. **Optional:** Uninstall the bundle before installing |
| `uninstall` | `[bundle]` |                                          | Uninstall named EcoSystem bundle.  Will not delete the bundle folder from /bundles |
| `new`       | `bundle`   |                                          | Creates the folder structure for a new EcoSystem bundle |
| `new`       | `user`     | —admin                                   | Creates a new user. **Optional:** make user with admin permissions |
| `serve`     |            | -d, —demomode; -s, —secret; —smtppw | Starts the EcoSystem server. **Optional:** Run the server in demo mode, allowing users to log in with magic code '123456', rather than having to request a code. **Required:** Encryption secret for JWT signing.  Remember to use the same secret every time you start the server, or JWTs previously issued will not be valid. |

### SMTP Configuration

SMTP settings for outgoing email are not strictly required, but the login system will not work without them, since it uses a passwordless process via email.  You can disable passwordless authentication for testing by running the server in demo mode.

If email credentials are provided, the connection will be tested, but startup will not fail if the test fails - the email system will simply be disabled and a warning displayed.

### Demo mode

As a convenience, the server can be run in *demo mode* with the flag `--demomode`.  In this mode, passwordless authentication is disabled, and users can log on with the code '123456'.  Note that the user must still be created in the database and a role assigned.

For example, if you wanted to demo admin panel functionality, you might create a user with the role 'admin' and the email 'test@test.com' - you'd then tell users to log in with that email and password '123456'.

## Licence

**Build freely with EcoSystem**.  The EcoSystem Server and The EcoSystem Admin Panel App are licensed under the [Apache 2.0 License](https://www.apache.org/licenses/LICENSE-2.0).  The content on the [EcoSystem website] (http://www.ecosystem.software) is licensed under a [Creative Commons Attribution-NonCommercial-NoDerivatives 4.0 International License](http://creativecommons.org/licenses/by-nc-nd/4.0/).  Neither licence grants permission to use the trade names, trademarks, service marks, or product names of EcoSystem Software LLP, including the EcoSystem logo and symbol, except as required for reasonable and customary use.

