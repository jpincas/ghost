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
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	sqlmock "gopkg.in/DATA-DOG/go-sqlmock.v1"

	"fmt"

	ghost "github.com/jpincas/ghost/tools"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/suite"
)

//HandlerTests is a basic suite test struct for all api handlers
type HandlerTests struct {
	suite.Suite
	Req  *http.Request
	Rr   *httptest.ResponseRecorder
	Mock sqlmock.Sqlmock
	Ctx  context.Context
}

type AuthHandlerTests struct {
	HandlerTests
}

//Run test suits
func TestAuthHandlerTests(t *testing.T) {
	suite.Run(t, new(AuthHandlerTests))
}

//Tests setup
func (suite *AuthHandlerTests) SetupTest() {
	suite.Req, _ = http.NewRequest("GET", "", nil)
	suite.Rr = httptest.NewRecorder()
}

func (suite *AuthHandlerTests) TestMagicCode_nobody() {

	http.HandlerFunc(magicCode).ServeHTTP(suite.Rr, suite.Req)
	suite.Equal(http.StatusBadRequest, suite.Rr.Code)

}

func (suite *AuthHandlerTests) TestMagicCode_emptybody() {

	b := []byte("")
	suite.Req, _ = http.NewRequest("POST", "", bytes.NewBuffer(b))
	http.HandlerFunc(magicCode).ServeHTTP(suite.Rr, suite.Req)
	suite.Equal(http.StatusBadRequest, suite.Rr.Code)

}

func (suite *AuthHandlerTests) TestMagicCode_badbody() {

	b := []byte("{gdf4}")
	suite.Req, _ = http.NewRequest("POST", "", bytes.NewBuffer(b))
	http.HandlerFunc(magicCode).ServeHTTP(suite.Rr, suite.Req)
	suite.Equal(http.StatusBadRequest, suite.Rr.Code)

}

func (suite *AuthHandlerTests) TestMagicCode_blankemail() {

	b := []byte(`{"email":""}`)
	suite.Req, _ = http.NewRequest("POST", "", bytes.NewBuffer(b))
	http.HandlerFunc(magicCode).ServeHTTP(suite.Rr, suite.Req)
	suite.Equal(http.StatusBadRequest, suite.Rr.Code)

}

func (suite *AuthHandlerTests) TestMagicCode_noemail() {

	b := []byte(`{"gmail":""}`)
	suite.Req, _ = http.NewRequest("POST", "", bytes.NewBuffer(b))
	http.HandlerFunc(magicCode).ServeHTTP(suite.Rr, suite.Req)
	suite.Equal(http.StatusBadRequest, suite.Rr.Code)

}

func (suite *AuthHandlerTests) TestMagicCode_emailsystemnotactivated() {

	b := []byte(`{"email":"me@me.com"}`)
	suite.Req, _ = http.NewRequest("POST", "", bytes.NewBuffer(b))
	http.HandlerFunc(magicCode).ServeHTTP(suite.Rr, suite.Req)
	suite.Equal(http.StatusServiceUnavailable, suite.Rr.Code)

}

func (suite *AuthHandlerTests) TestRequestNewUserToken() {

	http.HandlerFunc(requestNewUserToken).ServeHTTP(suite.Rr, suite.Req)
	suite.Equal(http.StatusOK, suite.Rr.Code)
	suite.NotEmpty(suite.Rr.Body, "Response should not be empty")

}

func (suite *AuthHandlerTests) TestRequestLogin_nobody() {

	http.HandlerFunc(requestLogin).ServeHTTP(suite.Rr, suite.Req)
	suite.Equal(http.StatusBadRequest, suite.Rr.Code)

}

func (suite *AuthHandlerTests) TestRequestLogin_emptybody() {

	b := []byte("")
	suite.Req, _ = http.NewRequest("POST", "", bytes.NewBuffer(b))
	http.HandlerFunc(requestLogin).ServeHTTP(suite.Rr, suite.Req)
	suite.Equal(http.StatusBadRequest, suite.Rr.Code)

}

func (suite *AuthHandlerTests) TestRequestLogin_badbody() {

	b := []byte("{gdf4}")
	suite.Req, _ = http.NewRequest("POST", "", bytes.NewBuffer(b))
	http.HandlerFunc(requestLogin).ServeHTTP(suite.Rr, suite.Req)
	suite.Equal(http.StatusBadRequest, suite.Rr.Code)

}

func (suite *AuthHandlerTests) TestRequestLogin_noemail() {

	b := []byte(`{"email": "", "code": "123456"}`)
	suite.Req, _ = http.NewRequest("POST", "", bytes.NewBuffer(b))
	http.HandlerFunc(requestLogin).ServeHTTP(suite.Rr, suite.Req)
	suite.Equal(http.StatusBadRequest, suite.Rr.Code)

}

func (suite *AuthHandlerTests) TestRequestLogin_noemail2() {

	b := []byte(`{"code": "123456"}`)
	suite.Req, _ = http.NewRequest("POST", "", bytes.NewBuffer(b))
	http.HandlerFunc(requestLogin).ServeHTTP(suite.Rr, suite.Req)
	suite.Equal(http.StatusBadRequest, suite.Rr.Code)

}

func (suite *AuthHandlerTests) TestRequestLogin_nocode() {

	b := []byte(`{"email": "me@me.com", "code": ""}`)
	suite.Req, _ = http.NewRequest("POST", "", bytes.NewBuffer(b))
	http.HandlerFunc(requestLogin).ServeHTTP(suite.Rr, suite.Req)
	suite.Equal(http.StatusBadRequest, suite.Rr.Code)

}

func (suite *AuthHandlerTests) TestRequestLogin_nocode2() {

	b := []byte(`{"email": "me@me.com"}`)
	suite.Req, _ = http.NewRequest("POST", "", bytes.NewBuffer(b))
	http.HandlerFunc(requestLogin).ServeHTTP(suite.Rr, suite.Req)
	suite.Equal(http.StatusBadRequest, suite.Rr.Code)

}

//Test demo mode automatic authorisation with code 123456
func (suite *AuthHandlerTests) TestRequestLogin_demomode() {

	ghost.App.DB, suite.Mock, _ = sqlmock.New()
	rows := sqlmock.NewRows([]string{"id"}).AddRow("130e6150-7098-4f72-8842-0e16629f32de")
	suite.Mock.ExpectQuery("is@registered.com").WillReturnRows(rows)

	MagicCodeCache.Set("is@registered.com", "666")
	viper.Set("demomode", true)

	t, _ := GetUserToken("130e6150-7098-4f72-8842-0e16629f32de")
	expectedToken := fmt.Sprintf("{%q:%q}", "token", t)

	b := []byte(`{"email": "is@registered.com", "code": "123456"}`)
	suite.Req, _ = http.NewRequest("POST", "", bytes.NewBuffer(b))

	http.HandlerFunc(requestLogin).ServeHTTP(suite.Rr, suite.Req)
	suite.Equal(http.StatusOK, suite.Rr.Code, fmt.Sprint(suite.Rr.Body))
	suite.Equal(expectedToken, fmt.Sprint(suite.Rr.Body), fmt.Sprint(suite.Rr.Body))

}

//Test demo mode off with same setting
func (suite *AuthHandlerTests) TestRequestLogin_fail() {

	ghost.App.DB, suite.Mock, _ = sqlmock.New()
	rows := sqlmock.NewRows([]string{"id"}).AddRow("130e6150-7098-4f72-8842-0e16629f32de")
	suite.Mock.ExpectQuery("is@registered.com").WillReturnRows(rows)

	MagicCodeCache.Set("is@registered.com", "666")
	viper.Set("demomode", false)

	b := []byte(`{"email": "is@registered.com", "code": "123456"}`)
	suite.Req, _ = http.NewRequest("POST", "", bytes.NewBuffer(b))

	http.HandlerFunc(requestLogin).ServeHTTP(suite.Rr, suite.Req)
	suite.Equal(http.StatusUnauthorized, suite.Rr.Code, fmt.Sprint(suite.Rr.Body))

}

func (suite *AuthHandlerTests) TestRequestLogin_ok() {

	ghost.App.DB, suite.Mock, _ = sqlmock.New()
	rows := sqlmock.NewRows([]string{"id"}).AddRow("130e6150-7098-4f72-8842-0e16629f32de")
	suite.Mock.ExpectQuery("is@registered.com").WillReturnRows(rows)

	MagicCodeCache.Set("is@registered.com", "666")
	viper.Set("demomode", false)

	t, _ := GetUserToken("130e6150-7098-4f72-8842-0e16629f32de")
	expectedToken := fmt.Sprintf("{%q:%q}", "token", t)

	b := []byte(`{"email": "is@registered.com", "code": "666"}`)
	suite.Req, _ = http.NewRequest("POST", "", bytes.NewBuffer(b))

	http.HandlerFunc(requestLogin).ServeHTTP(suite.Rr, suite.Req)
	suite.Equal(http.StatusOK, suite.Rr.Code, fmt.Sprint(suite.Rr.Body))
	suite.Equal(expectedToken, fmt.Sprint(suite.Rr.Body), fmt.Sprint(suite.Rr.Body))

}
