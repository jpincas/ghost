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
	"context"
	"database/sql"
	"fmt"

	"net/http"

	jwt "github.com/dgrijalva/jwt-go"
	ghost "github.com/jpincas/ghost/tools"
	"github.com/pressly/chi/render"
)

//This is the first level of authorisation:
//The JWT contains the userId.  We look this up in the users table in the database and if found
//attach the specified role.  If nothing is found, we default to anon
//Beyond this, we do not know anything about database privelages - this is handled
//further down the line
func Authorizator(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		ctx := r.Context()
		//This is a bit hairy. The jwt middleware puts the whole parsed Token
		//on the 'user' context value.  From there you have to reference 'Claims'
		//and then cast that to jwt.MapClaims to be able to reference the individual claims
		//Also: use the forked version of the go-jwt-middlware, not the auth0 version
		claims := ctx.Value("user").(*jwt.Token).Claims.(jwt.MapClaims)
		userID := claims["userID"]

		var role string
		//Search the user table for the user's role
		err := ghost.DB.QueryRow(fmt.Sprintf(ghost.SQLToGetUsersRole, userID)).Scan(&role)

		//If an error comes back
		if err != nil {
			//If its an error due to the id not being found in the user table, then just set the role to anon
			if err == sql.ErrNoRows {
				ctx = context.WithValue(ctx, "role", "anon")
			} else {
				//Else if there is any other error, don't authorise
				render.Status(r, http.StatusUnauthorized)
				render.JSON(w, r, ghost.ResponseError{http.StatusUnauthorized, "", err.Error(), "", "", ""})
				return
			}
		} else {
			//If a user has been found, return their role.  The database always defaults to anon, so there will always be a role
			ctx = context.WithValue(ctx, "role", role)
		}

		ctx = context.WithValue(ctx, "userID", userID)

		next.ServeHTTP(w, r.WithContext(ctx))
	})

}
