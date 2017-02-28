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

package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/suite"
	gin "gopkg.in/gin-gonic/gin.v1"
)

type HandlerTestSuite struct {
	suite.Suite
	router *gin.Engine
}

func (suite *HandlerTestSuite) SetupSuite() {
	gin.SetMode(gin.ReleaseMode)
	suite.router = gin.New()
	suite.router.GET("/returnblank", ReturnBlank)
	suite.router.GET("/returnhelloworld", ReturnHelloWorld)
}

func (suite *HandlerTestSuite) TestReturnBlank() {

	req, _ := http.NewRequest("GET", "/returnblank", nil)
	rr := httptest.NewRecorder()
	suite.router.ServeHTTP(rr, req)

	suite.Equal(http.StatusOK, rr.Code)
	suite.Equal("", rr.Body.String())

}

func (suite *HandlerTestSuite) TestReturnHelloWorld() {

	req, _ := http.NewRequest("GET", "/returnhelloworld", nil)
	rr := httptest.NewRecorder()
	suite.router.ServeHTTP(rr, req)

	suite.Equal(http.StatusOK, rr.Code)
	suite.Equal("{\"hello\":\"world\"}\n", rr.Body.String())

}

func TestHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(HandlerTestSuite))
}
