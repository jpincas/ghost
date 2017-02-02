![](https://raw.githubusercontent.com/ecosystemsoftware/ecosystem-website/master/themes/ecosystem/static/images/eco-logo-colour.png)

# EcoSystem Web/API Server

This is the EcoSystem server component, written in Golang.  It is intended to be used alongside the [EcoSystem admin panel Polymer app](https://github.com/ecosystemsoftware/ecosystem-admin).

## Getting Started

These instructions are for building and running the standard version of the server.  The repository is structured so that it can also be pulled in as a Go dependency for a custom build.  If you'd like to do that, check out the [Custom EcoSystem Server](https://github.com/ecosystemsoftware/ecosystem-server-custom) repository for some instructions and a template to work from.

### Prerequisites

- Go installed and environment variables correctly configured
- A Postgres database (local or remote) to which you can connect
- An email server you can use for outgoing emails for which you know the credentials

### Instructions

First, fetch the source code and dependencies and compile into an executable which will go into your $GOPATH/bin directory:

```
$ go get github.com/ecosystem-sofware/ecosystem-server
```

Assuming your [$GOPATH/bin is part of your PATH](https://golang.org/doc/install), you should now be able to run 'ecosystem-server' from anywhere on the command line.

Make a new folder from which to run the server and cd into it.  Run the EcoSystem server with the required command-line flags (see below).  If all flags are set correctly, the server will start and create the necessary folders (if they don't already exist).

```
$ mkdir my-server
$ cd my-server
$ ecosystem-server -pgname=testdatabase -pgdisablessl -secret=mysecret -pguser=eco -createadminwithemail=me@me.com -smtphost=smtp.eco.net -smtpuser=info@ecosystem -smtppw=ecopassword -smtpfrom=info@ecosystem -emailfrom=EcoSystem
```

*Optional* To get started quickly with your EcoSystem installation, copy our default web and email templates, css and starter Javascript into your new server folder.  Presuming you are still in the folder:

```
$ cp -r $GOPATH/src/github.com/ecosystemsoftware/ecosystem-server/public/ ./public
$ cp -r $GOPATH/src/github.com/ecosystemsoftware/ecosystem-server/templates/ ./templates
```

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

**Build freely with EcoSystem**.  The EcoSystem Server and The EcoSystem Admin Panel App are licensed under the [Apache 2.0 License](https://www.apache.org/licenses/LICENSE-2.0).  The content on the [EcoSystem website] (http://www.ecosystem.software) is licensed under a [Creative Commons Attribution-NonCommercial-NoDerivatives 4.0 International License](http://creativecommons.org/licenses/by-nc-nd/4.0/).  Neither licence grants permission to use the trade names, trademarks, service marks, or product names of EcoSystem Software LLP, including the EcoSystem logo and symbol, except as required for reasonable and customary use.