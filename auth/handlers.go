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
	"encoding/json"
	"fmt"
	"net/http"

	ghost "github.com/jpincas/ghost/tools"
	uuid "github.com/satori/go.uuid"
	"github.com/spf13/viper"
)

//ApiMagicCode processes a request for a magic code
func magicCode(w http.ResponseWriter, r *http.Request) {

	//Set content type to JSON
	w.Header().Set("Content-Type", ghost.ContentTypeJSON)

	//Set up the map into which the request body will be read
	var (
		requestBody     map[string]interface{}
		bodyDecodeError error
	)

	//If r.body is not nil (as in, body doesn't even exist), read and decode
	if r.Body != nil {
		d := json.NewDecoder(r.Body)
		bodyDecodeError = d.Decode(&requestBody)
	}

	//Filter for a nil body, blank body or empty JSON - return bad response
	// len(requestBody) also catches decode error
	if r.Body == nil || len(requestBody) == 0 {

		message := "Invalid or absent request body"
		if bodyDecodeError != nil {
			message = bodyDecodeError.Error()
		}

		//Output and return
		w.WriteHeader(http.StatusBadRequest)
		b, _ := json.Marshal(ghost.ResponseError{http.StatusBadRequest, "", message, "", "", ""})
		w.Write([]byte(b))
		return

	}

	//Try to read 'email'
	email, ok := requestBody["email"]
	if ok && email != "" {

		//If 'email' is set, request a magic code
		err := RequestMagicCode(email.(string))
		//If sending of the magic code fails (user doesn't exist, email fails etc)
		if err != nil {

			w.WriteHeader(http.StatusServiceUnavailable)
			b, _ := json.Marshal(ghost.ResponseError{http.StatusServiceUnavailable, "", err.Error(), "", "", ""})
			w.Write([]byte(b))
			return

		}

		//If the magic code goes through OK, just return a blank 200
		w.Write([]byte{})
		return

	}

	//If no email provided
	w.WriteHeader(http.StatusBadRequest)
	b, _ := json.Marshal(ghost.ResponseError{http.StatusBadRequest, "", "No email address provided", "", "", ""})
	w.Write([]byte(b))
	return

}

func requestNewUserToken(w http.ResponseWriter, r *http.Request) {

	tokenString, err := GetUserToken(fmt.Sprint(uuid.NewV4()))

	if err != nil {

		w.WriteHeader(http.StatusServiceUnavailable)
		b, _ := json.Marshal(ghost.ResponseError{http.StatusServiceUnavailable, "", err.Error(), "", "", ""})
		w.Write([]byte(b))
		return

	}

	b, _ := json.Marshal(map[string]string{
		"token": tokenString,
	})
	w.Write([]byte(b))
	return

}

func requestLogin(w http.ResponseWriter, r *http.Request) {

	//Set up the map into which the request body will be read
	var (
		requestBody     map[string]interface{}
		bodyDecodeError error
	)

	//If r.body is not nil (as in, body doesn't even exist), read and decode
	if r.Body != nil {
		d := json.NewDecoder(r.Body)
		bodyDecodeError = d.Decode(&requestBody)
	}

	//Filter for a nil body, blank body or empty JSON - return bad response
	// len(requestBody) also catches decode error
	if r.Body == nil || len(requestBody) == 0 {

		message := "Invalid or absent request body"
		if bodyDecodeError != nil {
			message = bodyDecodeError.Error()
		}

		//Output and return
		w.WriteHeader(http.StatusBadRequest)
		b, _ := json.Marshal(ghost.ResponseError{http.StatusBadRequest, "", message, "", "", ""})
		w.Write([]byte(b))
		return

	}

	//Try to read 'email' and 'code'
	email, ok1 := requestBody["email"]
	code, ok2 := requestBody["code"]
	if !ok1 || email == "" {

		//Output and return
		w.WriteHeader(http.StatusBadRequest)
		b, _ := json.Marshal(ghost.ResponseError{http.StatusBadRequest, "", "No email address provided", "", "", ""})
		w.Write([]byte(b))
		return

	} else if !ok2 || code == "" {

		//Output and return
		w.WriteHeader(http.StatusBadRequest)
		b, _ := json.Marshal(ghost.ResponseError{http.StatusBadRequest, "", "No magic code provided", "", "", ""})
		w.Write([]byte(b))
		return

	}

	//Lookup the email in the users table
	var id string
	err := ghost.DB.QueryRow(fmt.Sprintf(ghost.SQLToFindUserByEmail, email)).Scan(&id)
	cachedCode, emailIsInCache := MagicCodeCache.Get(email.(string))

	//For Demo Mode ONLY - bypass the magic code
	//checking and just send back the id
	//To use: just create a user with the role you want (e.g. admin)
	//and tell demo users to log in with that email and password 123456
	if viper.GetBool("demomode") && err == nil && code == "123456" {

		tokenString, err := GetUserToken(id)
		if err != nil {

			//Output and return
			w.WriteHeader(http.StatusServiceUnavailable)
			b, _ := json.Marshal(ghost.ResponseError{http.StatusServiceUnavailable, "", err.Error(), "", "", ""})
			w.Write([]byte(b))
			return

		}

		b, _ := json.Marshal(map[string]string{
			"token": tokenString,
		})
		w.Write([]byte(b))
		return

	} else if emailIsInCache && err == nil && code.(string) == cachedCode {

		//If the user exists in the database, the email is in the magic cache and the password supplied matches the magic code,
		//delete the email/magic code combo in the cache so it can't be used again
		MagicCodeCache.Remove(email.(string))
		tokenString, err := GetUserToken(id)
		if err != nil {

			//Output and return
			w.WriteHeader(http.StatusServiceUnavailable)
			b, _ := json.Marshal(ghost.ResponseError{http.StatusServiceUnavailable, "", err.Error(), "", "", ""})
			w.Write([]byte(b))
			return

		}

		b, _ := json.Marshal(map[string]string{
			"token": tokenString,
		})
		w.Write([]byte(b))
		return

	}

	//Default to unauthorised
	w.WriteHeader(http.StatusUnauthorized)
	b, _ := json.Marshal(ghost.ResponseError{http.StatusUnauthorized, "", "Could not log in with those credentials", "", "", ""})
	w.Write([]byte(b))
	return

}
