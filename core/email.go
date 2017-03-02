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

package core

import (
	"bytes"
	"fmt"
	"html/template"
	"net/smtp"
	"strings"

	"github.com/spf13/viper"
)

//SMTPServer holds all the necessary connection configuration for an SMTP server
type smtpServer struct {
	host, port, userName, password, from, FromName string
	Working                                        bool
}

//MailServer is the shared SMTP server for the application
var MailServer smtpServer

//EmailSetup sets up the shared SMTP connection, tests it and marks whether it is working or not
func EmailSetup() error {

	//Setup the smtp config struct, and mark as not working
	//Read in the configuration parameters from Viper
	MailServer = smtpServer{
		host:     viper.GetString("smtpHost"),
		port:     viper.GetString("smtpPort"),
		password: viper.GetString("smtpPW"),
		userName: viper.GetString("smtpUserName"),
		from:     viper.GetString("smtpFrom"),
		FromName: viper.GetString("emailFrom"),
		Working:  false,
	}

	//Test the SMTP connection
	if err := MailServer.TestConnection(); err != nil {
		return err
	}

	//If it passes, setup the config
	MailServer.Working = true
	return nil

}

//SendEmail is used internally by ECOSystem modules to send transactional emails
func (s smtpServer) SendEmail(to []string, subject string, data map[string]string, templateToUse string) (err error) {

	//Prepare the date for the email template
	parameters := struct {
		From    string
		To      string
		Subject string
		Data    map[string]string
	}{
		s.FromName,
		strings.Join([]string(to), ","),
		subject,
		data,
	}

	//Email templating.  Note that the strcuture of the header part of the email template is extremely important if fields
	//aren't correct and line breaking is wrong.  It should look like this, with the exact same line breaks:
	// To: {{ .To }}
	// From: {{ .From }}
	// Subject: {{ .Subject }}
	// MIME-version: 1.0
	// Content-Type: text/html; charset="UTF-8"

	buffer := new(bytes.Buffer)
	t, err := template.New(templateToUse).ParseGlob("templates/**/email/*")
	err = t.Execute(buffer, &parameters)

	auth := smtp.PlainAuth("", s.userName, s.password, s.host)

	err = smtp.SendMail(
		fmt.Sprintf("%s:%s", s.host, s.port),
		auth,
		s.from,
		to,
		buffer.Bytes())

	return err
}

//testConnection tests an SMTP connection
func (s smtpServer) TestConnection() error {
	//First try connecting
	c, err := smtp.Dial(fmt.Sprintf("%s:%s", s.host, s.port))
	if err != nil {
		return err
	}
	//If that works, try authenticating
	auth := smtp.PlainAuth("", s.userName, s.password, s.host)
	if err := c.Auth(auth); err != nil {
		return err
	}
	//If that all worked, return no error
	return nil
}
