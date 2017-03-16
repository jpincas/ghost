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

//go:generate hardcodetemplates -p=auth

import (
	"errors"
	"fmt"
	"html/template"
	"time"

	"log"

	"github.com/dgrijalva/jwt-go"
	"github.com/diegobernardes/ttlcache"
	"github.com/ecosystemsoftware/ecosystem/core"
	ecomail "github.com/ecosystemsoftware/ecosystem/email"
	"github.com/spf13/viper"
)

//Template holder
var templates *template.Template

//Activate is the main package activation function
func Activate() {
	log.Println("[auth package] Activating...")
	parseTemplates()
	//Set the routes for the package
	setRoutes()
}

func parseTemplates() {

	templates = template.Must(template.New("base").Parse(baseTemplate))
	log.Println("[auth package]", templates.DefinedTemplates())

}

//MagicCodeCache is the cache for storing email/temp pw combinations for passwordless authorisation
var MagicCodeCache = initCache(300) //5 minute expiry

func initCache(exp time.Duration) *ttlcache.Cache {

	if exp < 1 {
		log.Fatal("Cache expiry cannot be zero or negative")
	}

	newCache := ttlcache.NewCache()
	newCache.SetTTL(time.Duration(exp * time.Second))
	return newCache
}

//RequestMagicCode generates a magic code, stores it in the cache against the user's email and sends it to them by email
func RequestMagicCode(email string, template *template.Template) error {

	//If system email is not configured, this can't be done, so exit straight away
	if !ecomail.MailServer.Working {
		return errors.New("System email is not configured, so could not send magic code")
	}

	//First, lookup the email in the users table
	var id string
	err := core.DB.QueryRow(fmt.Sprintf(core.SQLToFindUserByEmail, email)).Scan(&id)

	//If the user doesn't exist in the DB
	if err != nil {
		return errors.New("Email address not in user database")
	}

	//User exists in the DB
	//Create a temporary, one-off password consisting of 6 random characters
	pw := core.RandomString(6)
	//Set it in the cache
	MagicCodeCache.Set(email, pw)

	//Set up the data map to go to the email sending function
	data := map[string]string{
		"password": pw,
	}

	//Send it to them by mail
	err = ecomail.MailServer.SendEmail(
		[]string{email},                                     //Recipient
		"Your Magic Code from "+ecomail.MailServer.FromName, //Subject
		data,     //Data to include in the email
		template) //Email template to use

	//Return whatever the result of the mail send was, either an error or nil
	return err

}

//GetUserToken returns a JWT string encoded with a user id
func GetUserToken(userID string) (string, error) {

	//Error for empty user ID
	if userID == "" {
		return "", errors.New("Empty user ID")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userID": userID,
		//TODO: Rest of claims, expiry etc.
	})

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString([]byte(viper.GetString("secret")))

	return tokenString, err

}
