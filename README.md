![](https://raw.githubusercontent.com/ecosystemsoftware/ecosystem-website/master/themes/ecosystem/static/images/eco-logo-colour.png)

# EcoSystem

EcoSystem is a platform that allows you to quickly develop completely custom data-driven websites, web-applications and backend systems.  EcoSystem doesn't assume anything about the data structure or logic of your business because you code it directly and freely at the database layer using the power of PostgreSQL.  EcoSystem then augments your database layer with a generic web and JSON API server written in Go and a generic backend administration panel created with Polymer, giving you a *bridge* to your data and logic with minial boilerplate.

## What's in the box?

**Short version:** the EcoSystem platform consists of two parts - an HTML/JSON server written in Go (this repository) and a [configurable admin panel created with Polymer](https://github.com/ecosystemsoftware/ecosystem-admin).

**Slightly longer version:** this repository contains the Go source code for the EcoSystem application, which is actually a command-line tool that allows you to easily initialise, configure and launch the EcoSystem server, as well as perform operations on **bundles** (more on those later).  If you are familiar with Go, you can also easily build a custom version of the server by including the packages (cmd, ecosql, handlers and utilities) as dependencies in your project. If you'd like to do that, check out the [Custom EcoSystem Server](https://github.com/ecosystemsoftware/ecosystem-server-custom) repository for some instructions and a template to work from.

**For extra points:** the EcoSystem server also serves the admin panel application for you, but you need to clone/download into your project directory first (see Step 9 below).

-bundles
-cli

## Getting Started

If you have Go installed and access to a Postgres server, it should take you less than 60 seconds to get up and running.

### Prerequisites

- Go installed and environment variables correctly configured
- A Postgres database server (local or remote) to which you can connect

### Step 1: Install EcoSystem

Fetch the source code and dependencies and compile into an executable which will go into your $GOPATH/bin directory:

```
$ go get github.com/ecosystemsofware/ecosystem
```

Assuming your [$GOPATH/bin is part of your PATH](https://golang.org/doc/install), you should now be able to run 'ecosystem' from anywhere on the command line.

### Step 2:  Create a database

Log into your Postgres server and create a new database.  I'll assume you'll call it `test`, but you don't have to.

### Step 3: Create a project folder

Make a new folder from which to run the server and cd into it:

```
$ mkdir test
$ cd test
```

### Step 4: Create a configuration file

Generate a default config.json file which we'll use to set up the database connection:

```
$ ecosystem new configfile
```

### Step 5: Configure the database connection

Open `config.json` and edit the database connection parameters so EcoSystem can connect to the `test` database you created.  If you're working locally, you'll probably only have to change the database name and disable SSL mode.

### Step 6: Initialise EcoSystem

EcoSystem needs to create a number of built in tables, roles, permissions and functions, as well as a few folders:

```
$ ecosystem init
```
If all goes well, you'll see messages confirming successful initialisation.

### Step 7: Create an admin user

Set yourself up as an admin user with full permissions:

```
$ ecosystem new user [your@email.com] --admin
```

### Step 8: Download and install a bundle

The easiest way to get started with EcoSystem is by installing an existing bundle.  Bundles are just folders that group together everything EcoSystem needs: database code, sample data, website templates, admin panel configurations etc.  Once you get familiar with EcoSystem, you'll create your own bundles, but for now, let's clone a simple demo bundle and install it with demo data:

```
$ cd bundles
$ git clone git@github.com:ecosystemsoftware/eco_bundle_dogshelter.git
$ cd ..
$ ecosystem install eco_bundle_dogshelter --demodata
```

There's nothing magic going on under the hood here, we simply created a new folder in bundles called 'eco_bundle_dogshelter' which will be the name of the bundle, and copied in some files.  We then ran the database code in `install.sql` to set up the bundle's data and logic.  Finally we ran `demodata.sql` to install the bundle's demo data.  Everything else, like website template code and admin panel config code, lives in the bundle folder and the EcoSystem server knows to find it there.

### Step 9: Download the EcoSystem admin panel app

From your main project folder, clone the EcoSystem admin panel:

```
$ git clone git@github.com:ecosystemsoftware/ecosystem-admin.git
```

### Step 10: Set up the bundle imports

*Note:  We working on a way of automating this step - until now, it's manual I'm afraid.*

The bundle contains various setting files for the admin panel, for them to take effect, you'll need the admin panel to import them.  The procedure is as follows:

1.  Open the file `bundles/eco_bundle_dogshelter/admin-panel/bundle-imports.html`
2.  Copy all the lines that start `<link rel="import"...`
3.  Open the file `ecosystem-admin/src/eco-bundle-imports.html`
4.  Paste the lines previously copied.

### Step 10: Run the server

Run the server in 'demo' mode (disabling authorisation).  You need to supply a 'secret' which is used by the server for encryption.  Here I use 'secret', but you can use anything.

```
$ ecosystem serve -s=secret -d
```

## Enjoy!

Here's a quick rundown of what you now have in place:

### localhost:3000 - JSON API##

This is the basis of every EcoSystem implementation.  The built-in admin panel uses it, and you can use it to build any kind of API-backed application.  All requests to the API need to go through authorisation, so if you'd like to send some test requests, you'll need to start by requesting a JWT from **/login** - just include a JSON body in the request like this:

```
{
    "username": [your@email.com] //The email of the admin account you set up earlier
    "password": "123456" //Since you're runnning in demo mode, just use the password 123456 
}
```

grab the JWT that comes back and include it with your requests in an authorization header: `Authorization: Bearer xxxxxx...` - you'll have full admin permissions since you set yourself up as an admin user earlier.

### localhost:3001 - Website ###

This is optional, but if you only need a traditional server-generated HTML website, then EcoSystem gives you that out-of-the-box.  Head over to *localhost:3001* to the see the website that was generated with the simple templates in the bundle.  See how the URL path is basically just */SCHEMA/TABLE/RECORD*? That's the beauty of EcoSystem - it's simple!  Golang teplates are incredibly easy to use, so putting together a functional website is quick.

There's more in the box not covered by this short tutorial.  As well as full HTML pages, the web server can also deliver HTML partials and use the authorisation system, so you can easily use AJAX to swap in templated data specific to your current browser - allowing you to build fully dynamic server-generated pages, like carts and wishlist.

Of course, there's nothing to stop you using *both* the HTML API and JSON API to create hybrid sites using combination of server and client side rendering...

And finally, EcoSystem includes a growing collection of utility functions for easier HTML generation, like an automatic menu builder.

### localhost:3002 - Admin Panel ###

Every business needs an administration panel to control its data.  We're not talking about a visual website builder here, we're talking about a proper database visualisation tool to allows administrators to interact with custom data and logic.  Since EcoSystem knows nothing about your application, the admin panel is fully configurable with simple JSON - tell it what you want to see and do, and the admin panel is generated for you.  The easy way to understand how it works is by looking at the configuration files in the `bundle/admin-panel` folder - it's fairly self-explanatory.

Head over to *localhost:3002/admin*, log in with your email address and the password '123456' to get a feel for the admin panel.

***Hey there! I'm fully aware that the above is real whistlestop tour of EcoSystem and you probably have a hundred questions.  We're working like crazy to release better documentation, videos, tutorials etc, but in the meantime, please check out [the developer section of the EcoSystem website](www.ecosystem.software/developers)***

<!--## Commands

## Command Line Flags

The following is a list of available command line flags when starting the server.

| Flag                  | Function                                 | Required | Default        |
| --------------------- | ---------------------------------------- | -------- | -------------- |
| -pgname               | The name of the Postgres database to connect to | YES      |                |
| -pgserver             | The server address for the Postgres connection |          | localhost |
| -pgport               | The server port for the Postgres connection |          | 5432 |
| -pgdisablessl         | Disbles SSL mode in the Postgres connection (for development) |          | FALSE          |
| -pguser               | Username of database superuser with which to connect to the database for initial setup |       | postgres               |
| -pgpw                 | Postgres connection password |          | localhost |
| -secret               | The secret used to sign JWTs             | YES      |                |
| -createadminwithemail | Bootstrap the installation with an admin user with this email address |          |                |
| -siteslug             | The URL slug for the public facing website |          | "site"         |
| -privateslug          | The URL slug for the authenticated HTML api |          | "private"      |
| -smtphost             | SMTP server for outgoing emails          |          |                |
| -smtpuser             | SMTP authentication username             |          |                |
| -smtppw               | SMTP authentication password             |          |                |
| -smtpfrom             | 'From' address for outgoing emails       |          |                |
| -emailFrom            | 'From' name for outgoing emails          |          | =-smtpfrom     |
| -demomode           | Run server in demo mode or not. Demo mode basically bypasses login. See below.      |          | FALSE  |
| -demorole           | Role to use when running in demo mode     |          | "admin"    |

## SMTP Configuration

SMTP settings for outgoing email are not strictly required, but the login system will not work without them.  If incomplete credentials are provided, a warning will be displayed, but startup will not fail.

If complete credentials are provided, the system will ping the server to test the connection.  If the ping fails, startup will exit.

### Demo mode

As a convenience, the server can be run in *demo mode*.  In this mode, any requests to log in will bypass email and magic code checking (i.e. anything can be entered) and return a JWT encoded with a random UUID.  Subsequent calls to the API with this JWT will be authorised with the role specified with the flag `-demorole` which defaults to "admin".

### Licence

**Build freely with EcoSystem**.  The EcoSystem Server and The EcoSystem Admin Panel App are licensed under the [Apache 2.0 License](https://www.apache.org/licenses/LICENSE-2.0).  The content on the [EcoSystem website] (http://www.ecosystem.software) is licensed under a [Creative Commons Attribution-NonCommercial-NoDerivatives 4.0 International License](http://creativecommons.org/licenses/by-nc-nd/4.0/).  Neither licence grants permission to use the trade names, trademarks, service marks, or product names of EcoSystem Software LLP, including the EcoSystem logo and symbol, except as required for reasonable and customary use.-->