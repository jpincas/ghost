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
	"math/rand"
	"time"

	"github.com/diegobernardes/ttlcache"
)

//EmailPWCache is the cache for storing email/temp pw combinations for passwordless authorisation
var magicCodeCache = initCache(300) //5 minute expiry

func initCache(exp time.Duration) *ttlcache.Cache {
	newCache := ttlcache.NewCache()
	newCache.SetTTL(time.Duration(exp * time.Second))
	return newCache
}

func requestMagicCode(email string) error {

	//If system email is not configured, this can't be done, so exit straight away
	if !mySMTPServer.working {
		return errors.New("System email is not configured, so could not send magic code")
	}

	//First, lookup the email in the users table
	var id string
	err := db.QueryRow(fmt.Sprintf("SELECT id from users WHERE email = '%s'", email)).Scan(&id)

	//If the user doesn't exist in the DB
	if err != nil {
		return errors.New("Email address not in user database")
	}

	//User exists in the DB
	//Create a temporary, one-off password consisting of 6 random characters
	pw := randomString(6)
	//Set it in the cache
	magicCodeCache.Set(email, pw)

	//Set up the data map to go to the email sending function
	data := map[string]string{
		"password": pw,
	}

	//Send it to them by mail
	err = sendEmail(
		mySMTPServer,
		[]string{email},                               //Recipient
		"Your Magic Code from "+mySMTPServer.fromName, //Subject
		data, //Data to include in the email
		"emailMagicCode.html") //Email template to use

	//Return whatever the result of the mail send was, either an error or nil
	return err

}

func randomString(strlen int) string {
	rand.Seed(time.Now().UTC().UnixNano())
	const chars = "abcdefghijklmnopqrstuvwxyz0123456789"
	result := make([]byte, strlen)
	for i := 0; i < strlen; i++ {
		result[i] = chars[rand.Intn(len(chars))]
	}
	return string(result)
}
