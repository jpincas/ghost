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

func (s *smtpServer) Setup() {

	Log("EMAIL", true, "Initialising email system...", nil)

	//Setup the smtp config struct, and mark as not working
	//Read in the configuration parameters from Viper
	s.host = App.Config.SmtpHost
	s.port = App.Config.SmtpPort
	s.password = viper.GetString("smtpPW")
	s.userName = App.Config.SmtpUserName
	s.from = App.Config.SmtpFrom
	s.FromName = App.Config.SmtpFrom
	s.Working = false

	//Test the SMTP connection
	if err := s.TestConnection(); err != nil {
		LogFatal("EMAIL", false, "Error initialising email server", err)
	}

	//If it passes, setup the config
	s.Working = true
	Log("EMAIL", true, "Email system correctly initialised", nil)

}

//SendEmail is used internally by ghost modules to send transactional emails
func (s smtpServer) SendEmail(to []string, subject string, data map[string]string, templates *template.Template, templateToUse string) (err error) {

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

	//This is ridiculously sensitive - even a blank line at the beginning of the file will
	//cause the email send to fail

	buffer := new(bytes.Buffer)
	err = templates.ExecuteTemplate(buffer, templateToUse, &parameters)

	if err != nil {
		return err
	}

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
