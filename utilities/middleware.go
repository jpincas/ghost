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
	"database/sql"
	"fmt"
	"strings"
	"time"

	"net/http"

	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
	jwt "gopkg.in/appleboy/gin-jwt.v2"
)

//CORS
func AllowCORS(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", "*")
	c.Next()
}

//JSON Header
func MakeJSON(c *gin.Context) {
	c.Header("Content-Type", "application/json")
	c.Next()
}

//the jwt middleware
var AuthMiddleware = &jwt.GinJWTMiddleware{
	Realm:      "ecosystem",
	Key:        []byte(*secret),
	Timeout:    time.Hour * 24 * 365,
	MaxRefresh: time.Hour * 365,
	Authenticator: func(username string, password string, c *gin.Context) (string, bool) {

		// If a new user is being requested, return a JWT with the id as a new GUID
		//Hint: to request a new user, post both username(email) and password as 'anon'
		if strings.ToUpper(username) == "ANON" && strings.ToUpper(password) == "ANON" {
			return fmt.Sprint(uuid.NewV4()), true
		}

		//Otherwise attempt to login with the given email address
		//Lookup the email in the users table
		var id string
		err := db.QueryRow(fmt.Sprintf("SELECT id from users WHERE email = '%s'", username)).Scan(&id)
		//If the user exists in the database, the email is in the magic cache and the password supplied matches the magic code, then
		//return the JWT with the id encoded
		magicCode, emailIsInCache := magicCodeCache.Get(username)

		//For Demo Mode ONLY - bypass the magic code
		//checking and just send back the id
		//To use: just create a user with the role you want (e.g. admin)
		//and tell demo users to log in with that email and password 123456
		if *demoMode && err == nil && password == "123456" {
			return id, true
		}

		if emailIsInCache && err == nil && password == magicCode {
			//delete the email/magic code combo in the cache so it can't be used again
			magicCodeCache.Remove(username)
			return id, true
		}

		//Otherwise, failed login
		return "", false

	},
	Authorizator: func(userID string, c *gin.Context) bool {
		//This is the first level of authorisation:
		//The JWT contains the userId.  We look this up in the users table in the database and if found
		//attach the specified role.  If nothing is found, we default to anon
		//Beyond this, we do not know anything about database privelages - this is handled
		//further down the line

		var role string
		//Search the user table for the user's role
		err := db.QueryRow(fmt.Sprintf("SELECT role from users WHERE id = '%s'", userID)).Scan(&role)

		//If an error comes back
		if err != nil {
			//If its an error due to the id not being found in the user table, then just set the role to anon
			if err == sql.ErrNoRows {
				c.Set("role", "anon")
				return true
			}
			//Else if there is any other error, don't authorise
			return false

		}
		//If a user has been found, return their role.  The database always defaults to anon, so ther will always be a role
		c.Set("role", role)
		//Finally, attach the user id to the context so it can be used further down the line
		c.Set("userID", userID)

		return true

	},
	Unauthorized: func(c *gin.Context, code int, message string) {
		//All authorisation errors at this level are 401.
		//This is to distinguish from permission errors further down the line, which will be 403
		c.JSON(code, gin.H{
			"code":    http.StatusUnauthorized,
			"message": message,
		})
	},
	TokenLookup: "header:Authorization",
}
