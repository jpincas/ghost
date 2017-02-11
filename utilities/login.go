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

package utilities

import (
	"errors"
	"fmt"
	"time"

	"github.com/diegobernardes/ttlcache"
	"github.com/ecosystemsoftware/eco/ecosql"
)

//MagicCodeCache is the cache for storing email/temp pw combinations for passwordless authorisation
var MagicCodeCache = initCache(300) //5 minute expiry

func initCache(exp time.Duration) *ttlcache.Cache {
	newCache := ttlcache.NewCache()
	newCache.SetTTL(time.Duration(exp * time.Second))
	return newCache
}

//RequestMagicCode generates a magic code, stores it in the cache against the user's email and sends it to them by email
func RequestMagicCode(email string, templateName string) error {

	//If system email is not configured, this can't be done, so exit straight away
	if !MailServer.working {
		return errors.New("System email is not configured, so could not send magic code")
	}

	//First, lookup the email in the users table
	var id string
	err := DB.QueryRow(fmt.Sprintf(ecosql.ToFindUserByEmail, email)).Scan(&id)

	//If the user doesn't exist in the DB
	if err != nil {
		return errors.New("Email address not in user database")
	}

	//User exists in the DB
	//Create a temporary, one-off password consisting of 6 random characters
	pw := RandomString(6)
	//Set it in the cache
	MagicCodeCache.Set(email, pw)

	//Set up the data map to go to the email sending function
	data := map[string]string{
		"password": pw,
	}

	//Send it to them by mail
	err = MailServer.SendEmail(
		[]string{email},                             //Recipient
		"Your Magic Code from "+MailServer.fromName, //Subject
		data,         //Data to include in the email
		templateName) //Email template to use

	//Return whatever the result of the mail send was, either an error or nil
	return err

}
