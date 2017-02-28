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

	eco "github.com/ecosystemsoftware/ecosystem/utilities"
	"github.com/stretchr/testify/suite"
	gin "gopkg.in/gin-gonic/gin.v1"
)

type APIHandlerTestSuite struct {
	suite.Suite
	router *gin.Engine
}

func (suite *APIHandlerTestSuite) SetupSuite() {
	suite.router = gin.Default()
	suite.router.GET("/:schema/:table", ApiShowList)

	eco.DB = eco.TestDBConfig.ReturnDBConnection("")
}

func (suite *APIHandlerTestSuite) TestApiShowList() {

	req, _ := http.NewRequest("GET", "/schema/table", nil)
	rr := httptest.NewRecorder()

	suite.router.ServeHTTP(rr, req)
	suite.Equal(http.StatusBadRequest, rr.Code, "No role, no id should be a bad request")

}

func TestAPIHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(APIHandlerTestSuite))
}

//Teardown - discconent from the database
