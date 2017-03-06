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

package auth

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/ecosystemsoftware/ecosystem/core"
	"github.com/pressly/chi/render"
	uuid "github.com/satori/go.uuid"
	"github.com/spf13/viper"
)

//ApiMagicCode processes a request for a magic code
func MagicCode(w http.ResponseWriter, r *http.Request) {

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
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, core.ResponseError{http.StatusBadRequest, "", message, "", "", ""})

	} else {
		//Try to read 'email'
		email, ok := requestBody["email"]
		if ok && email != "" {
			//If 'email' is set, request a magic code
			err := RequestMagicCode(email.(string), "emailMagicCode.html")
			//If sending of the magic code fails (user doesn't exist, email fails etc)
			if err != nil {
				render.Status(r, http.StatusServiceUnavailable)
				render.JSON(w, r, core.ResponseError{http.StatusServiceUnavailable, "", err.Error(), "", "", ""})
			} else {
				//If the magic code goes through OK
				render.NoContent(w, r)
			}
		} else {
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, core.ResponseError{http.StatusBadRequest, "", "No email address provided", "", "", ""})
		}
	}
}

func RequestNewUserToken(w http.ResponseWriter, r *http.Request) {

	tokenString, err := GetUserToken(fmt.Sprint(uuid.NewV4()))

	if err != nil {
		render.Status(r, http.StatusServiceUnavailable)
		render.JSON(w, r, core.ResponseError{http.StatusServiceUnavailable, "", err.Error(), "", "", ""})
	} else {
		render.JSON(w, r,
			map[string]string{
				"token": tokenString,
			})
	}

}

func RequestLogin(w http.ResponseWriter, r *http.Request) {

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
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, core.ResponseError{http.StatusBadRequest, "", message, "", "", ""})

	} else {
		//Try to read 'email' and 'code'
		email, ok1 := requestBody["email"]
		code, ok2 := requestBody["code"]
		if !ok1 || email == "" {
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, core.ResponseError{http.StatusBadRequest, "", "No email address provided", "", "", ""})
		} else if !ok2 || code == "" {
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, core.ResponseError{http.StatusBadRequest, "", "No magic code provided", "", "", ""})
		} else {

			//Lookup the email in the users table
			var id string
			err := core.DB.QueryRow(fmt.Sprintf(core.SQLToFindUserByEmail, email)).Scan(&id)
			cachedCode, emailIsInCache := MagicCodeCache.Get(email.(string))

			//For Demo Mode ONLY - bypass the magic code
			//checking and just send back the id
			//To use: just create a user with the role you want (e.g. admin)
			//and tell demo users to log in with that email and password 123456
			if viper.GetBool("demomode") && err == nil && code == "123456" {

				tokenString, err := GetUserToken(id)
				if err != nil {
					render.Status(r, http.StatusServiceUnavailable)
					render.JSON(w, r, core.ResponseError{http.StatusServiceUnavailable, "", err.Error(), "", "", ""})
				} else {
					render.PlainText(w, r, tokenString)
				}

			} else if emailIsInCache && err == nil && code.(string) == cachedCode {

				//If the user exists in the database, the email is in the magic cache and the password supplied matches the magic code,
				//delete the email/magic code combo in the cache so it can't be used again
				MagicCodeCache.Remove(email.(string))
				tokenString, err := GetUserToken(id)
				if err != nil {
					render.Status(r, http.StatusServiceUnavailable)
					render.JSON(w, r, core.ResponseError{http.StatusServiceUnavailable, "", err.Error(), "", "", ""})
				} else {
					render.PlainText(w, r, tokenString)
				}

			} else {

				//Default to unauthorised
				render.Status(r, http.StatusUnauthorized)
				render.JSON(w, r, core.ResponseError{http.StatusServiceUnavailable, "", "Could not log in with those credentials", "", "", ""})
			}

		}

	}

}

// func ReturnBlank(c *gin.Context) {
// 	c.String(http.StatusOK, "")
// }

// //ReturnHelloWorld is a test handler that can be used when wiring up a custom server to
// //check that the EcoSystem utility package is being correctly imported and built
// func ReturnHelloWorld(c *gin.Context) {
// 	c.JSON(http.StatusOK, gin.H{"hello": "world"})
// }
