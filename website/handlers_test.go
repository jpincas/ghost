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

package website

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ecosystemsoftware/ecosystem/core"
	"github.com/lib/pq"
	"github.com/stretchr/testify/suite"
	sqlmock "gopkg.in/DATA-DOG/go-sqlmock.v1"
)

type ShowListTests struct {
	core.HandlerTests
}

//Run test suits
func TestShowListTests(t *testing.T) {
	suite.Run(t, new(ShowListTests))
}

//Tests setup
func (suite *ShowListTests) SetupSuite() {
	parseTemplates()
}

//Tests setup
func (suite *ShowListTests) SetupTest() {
	suite.Req, _ = http.NewRequest("GET", "", nil)
	suite.Rr = httptest.NewRecorder()
	//Setup db mocks
	core.DB, suite.Mock, _ = sqlmock.New()
	suite.Ctx = context.WithValue(suite.Ctx, "schema", "schema")
	suite.Ctx = context.WithValue(suite.Ctx, "table", "table")
}

//Test handler with context values unset, which in theory should never happen (routing, middleware etc)
// func (suite *ShowListTests) TestShowList_missing_context_values() {

// 	http.HandlerFunc(ShowList).ServeHTTP(suite.Rr, suite.Req)
// 	suite.Equal(http.StatusBadRequest, suite.Rr.Code, "Missing context values - should be bad request")
// 	suite.NotEmpty(suite.Rr.Body, "Should be an HTML error page")
// 	//TODO: test content type

// }

//If the schema specified in the URL doesn't exist then 404
func (suite *ShowListTests) TestShowList_schema_notexists() {

	suite.Mock.ExpectQuery("noschema").WillReturnError(&pq.Error{
		Message: "Schema doesn't exist",
		Code:    "42P01",
	})
	ctx := context.WithValue(suite.Ctx, "schema", "noschema")
	http.HandlerFunc(ShowList).ServeHTTP(suite.Rr, suite.Req.WithContext(ctx))
	suite.NotEmpty(suite.Rr.Body, "Should have a error HTML page")
	suite.Equal(http.StatusNotFound, suite.Rr.Code, "Should be HTTP Status Not Found when schema doesn't exist")

}

// //If the table specified in the URL doesn't exist then 404
// func (suite *ShowListTests) TestShowList_schema_exists_table_notexists() {

// 	suite.Mock.ExpectQuery("schema.notable").WillReturnError(&pq.Error{
// 		Message: "Table doesn't exist",
// 		Code:    "42P01",
// 	})
// 	ctx := context.WithValue(suite.Ctx, "table", "notable")

// 	http.HandlerFunc(ShowList).ServeHTTP(suite.Rr, suite.Req.WithContext(ctx))
// 	suite.NotEmpty(suite.Rr.Body, "Should have HTML error page")
// 	suite.Equal(http.StatusNotFound, suite.Rr.Code, "Should be HTTP Status Not Found when table doesn't exist")

// }

// //If the table specified is empty, then 200 and empty array
// func (suite *ShowListTests) TestShowList_schema_exists_table_exists_with_norows() {

// 	suite.Mock.ExpectQuery("schema.tablewithnorows").WillReturnError(errors.New("sql: no rows"))
// 	ctx := context.WithValue(suite.Ctx, "table", "tablewithnorows")

// 	http.HandlerFunc(ShowList).ServeHTTP(suite.Rr, suite.Req.WithContext(ctx))
// 	suite.Equal(http.StatusOK, suite.Rr.Code, "Should be HTTP OK")
// 	suite.NotEmpty(suite.Rr.Body, "Should have HTML content")

// }

// //If the table specified exists and has rows, then 200 and return the JSON array
// func (suite *ShowListTests) TestShowList_schema_exists_table_exists_with_rows() {

// 	rows := sqlmock.NewRows([]string{"json"}).AddRow(`["some" : "value"]`)
// 	suite.Mock.ExpectQuery("schema.tablewithrows").WillReturnRows(rows)
// 	ctx := context.WithValue(suite.Ctx, "table", "tablewithrows")

// 	http.HandlerFunc(ShowList).ServeHTTP(suite.Rr, suite.Req.WithContext(ctx))
// 	suite.NotEmpty(suite.Rr.Body, "Should have HTML body")
// 	suite.Equal(http.StatusOK, suite.Rr.Code, "Should be HTTP OK")

// }

// //If the role does not have privileges to view the table
// func (suite *ShowListTests) TestShowList_schema_exists_table_notauthorised() {

// 	suite.Mock.ExpectQuery("nonauthedrole").WillReturnError(&pq.Error{
// 		Message: "Role not authorised for that table",
// 		Code:    "42501",
// 	})
// 	ctx := context.WithValue(suite.Ctx, "role", "nonauthedrole")

// 	http.HandlerFunc(ShowList).ServeHTTP(suite.Rr, suite.Req.WithContext(ctx))
// 	suite.NotEmpty(suite.Rr.Body, "Should have HTTP body")
// 	suite.Equal(http.StatusForbidden, suite.Rr.Code, "Should be HTTP Status Forbidden when not authorised to view a table")

// }

//Test handler with context values unset, which in theory should never happen (routing, middleware etc)
// func (suite *ShowListTests) TestShowList_ok() {

// 	http.HandlerFunc(ShowList).ServeHTTP(suite.Rr, suite.Req)
// 	suite.Equal(http.StatusBadRequest, suite.Rr.Code, "Missing context values - should be bad request")
// 	suite.NotEmpty(suite.Rr.Body, "Should be an HTML error page")
// 	//TODO: test content type

// }

//If the schema specified in the URL doesn't exist then 404
// func (suite *ShowListTests) TestShowList_schema_notexists() {

// 	suite.Mock.ExpectQuery("noweb").WillReturnError(&pq.Error{
// 		Message: "",
// 		Code:    "42501",
// 	})
// 	ctx := context.WithValue(suite.Ctx, "schema", "noschema")

// 	http.HandlerFunc(ShowList).ServeHTTP(suite.Rr, suite.Req.WithContext(ctx))
// 	suite.NotEmpty(suite.Rr.Body, "Should have a JSON error body")
// 	suite.Equal(http.StatusNotFound, suite.Rr.Code, "Should be HTTP Status Not Found when schema doesn't exist")

// }
