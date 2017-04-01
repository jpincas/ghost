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

package auth

import (
	"fmt"
	"testing"

	sqlmock "gopkg.in/DATA-DOG/go-sqlmock.v1"

	"log"

	ghost "github.com/jpincas/ghost/tools"
	ecomail "github.com/jpincas/ghost/email"
	"github.com/spf13/viper"
)

func TestInitCache(t *testing.T) {
	initCache(1000)
	//TODO: test conditions that cause fatal error (exp < 1)
}

func TestRequestMagicCodeEmailDisabled(t *testing.T) {
	setup()
	err := RequestMagicCode("user@notindb")
	if err == nil {
		t.Error("Email system is not working.  Magic code request should fail")
	}
	teardown()
}

func TestRequestMagicCodeUserNotInDB(t *testing.T) {
	setup()
	//Flag the email server as enabled
	ecomail.MailServer.Working = true
	err := RequestMagicCode("user@notindb")
	if err.Error() != "Email address not in user database" {
		t.Error("User is not in App.DB, should return an error")
	}
	teardown()
}

//This will return an error in actually sending the email
//because the email server won't be setup, but we're not testing that here
func TestRequestMagicCodeUserInDB(t *testing.T) {
	setup()
	//Flag the email server as enabled
	ecomail.MailServer.Working = true
	err := RequestMagicCode("user@isindb")
	if err.Error() == "Email address not in user database" {
		t.Error(err.Error())
	}
	teardown()
}

//Get token
func TestGetToken(t *testing.T) {

	//Set the secret
	viper.Set("secret", "secret")

	//With a proper user Id
	s, err := GetUserToken("692e8a64-7676-4790-b3f8-a86a5083d5bb")
	if err != nil {
		t.Error(err.Error(), "Token: ", s)
	}

	//With empty string
	s, err = GetUserToken("")
	if err == nil {
		t.Error("Empty user ID should be an error")
	}

}

func setup() {

	//Set up a mock App.DB
	var (
		mock sqlmock.Sqlmock
		err  error
	)
	ghost.App.DB, mock, err = sqlmock.New()

	if err != nil {
		log.Fatalf("An error '%s' was not expected when opening a stub database connection", err)
	}

	//Give it some data
	var columns = []string{"id"}
	dataRows := sqlmock.NewRows(columns).FromCSVString("692e8a64-7676-4790-b3f8-a86a5083d5bb")
	mock.ExpectQuery(fmt.Sprintf(ghost.SQLToFindUserByEmail, "user@isindb")).WillReturnRows(dataRows)

	//Parse the package templates
	parseTemplates()

}

func teardown() {
	ghost.App.DB.Close()
}
