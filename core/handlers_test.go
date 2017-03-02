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

package core

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"errors"

	"fmt"

	"net/url"

	"github.com/lib/pq"
	"github.com/stretchr/testify/suite"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
)

//TODO: test for JSON content type in all
/////////////////////////////////////////////////////////////////////
// Tests for func ShowList(w http.ResponseWriter, r *http.Request) //
/////////////////////////////////////////////////////////////////////

type ShowListTests struct {
	HandlerTests
}

//Run test suits
func TestShowListTests(t *testing.T) {
	suite.Run(t, new(ShowListTests))
}

//Tests setup
func (suite *ShowListTests) SetupTest() {
	suite.Req, _ = http.NewRequest("GET", "", nil)
	suite.Rr = httptest.NewRecorder()
	//Setup db mocks
	DB, suite.Mock, _ = sqlmock.New()
	suite.Ctx = context.WithValue(suite.Req.Context(), "role", "role")
	suite.Ctx = context.WithValue(suite.Ctx, "userID", "123456789")
	suite.Ctx = context.WithValue(suite.Ctx, "schema", "schema")
	suite.Ctx = context.WithValue(suite.Ctx, "table", "table")
}

//Test handler with context values unset, which in theory should never happen (routing, middleware etc)
func (suite *ShowListTests) TestShowList_missing_context_values() {

	http.HandlerFunc(ShowList).ServeHTTP(suite.Rr, suite.Req)
	suite.Equal(http.StatusBadRequest, suite.Rr.Code, "Missing context values - should be bad request")
	suite.NotEmpty(suite.Rr.Body, "Should have a JSON body with error message")

}

//If the schema specified in the URL doesn't exist then 404
func (suite *ShowListTests) TestShowList_schema_notexists() {

	suite.Mock.ExpectQuery("noschema").WillReturnError(&pq.Error{
		Message: "Schema doesn't exist",
		Code:    "42P01",
	})
	ctx := context.WithValue(suite.Ctx, "schema", "noschema")

	http.HandlerFunc(ShowList).ServeHTTP(suite.Rr, suite.Req.WithContext(ctx))
	suite.NotEmpty(suite.Rr.Body, "Should have a JSON error body")
	suite.Equal(http.StatusNotFound, suite.Rr.Code, "Should be HTTP Status Not Found when schema doesn't exist")

}

//If the table specified in the URL doesn't exist then 404
func (suite *ShowListTests) TestShowList_schema_exists_table_notexists() {

	suite.Mock.ExpectQuery("schema.notable").WillReturnError(&pq.Error{
		Message: "Table doesn't exist",
		Code:    "42P01",
	})
	ctx := context.WithValue(suite.Ctx, "table", "notable")

	http.HandlerFunc(ShowList).ServeHTTP(suite.Rr, suite.Req.WithContext(ctx))
	suite.NotEmpty(suite.Rr.Body, "Should have a JSON error body")
	suite.Equal(http.StatusNotFound, suite.Rr.Code, "Should be HTTP Status Not Found when table doesn't exist")

}

//If the table specified is empty, then 200 and empty array
func (suite *ShowListTests) TestShowList_schema_exists_table_exists_with_norows() {

	suite.Mock.ExpectQuery("schema.tablewithnorows").WillReturnError(errors.New("sql: no rows"))
	ctx := context.WithValue(suite.Ctx, "table", "tablewithnorows")

	http.HandlerFunc(ShowList).ServeHTTP(suite.Rr, suite.Req.WithContext(ctx))
	suite.Equal(http.StatusOK, suite.Rr.Code, "Should be HTTP OK")
	suite.NotEmpty(suite.Rr.Body, "Should have a JSON body with a blank array")
	suite.Equal("[]", fmt.Sprint(suite.Rr.Body), "Should return a blank array")

}

//If the table specified exists and has rows, then 200 and return the JSON array
func (suite *ShowListTests) TestShowList_schema_exists_table_exists_with_rows() {

	rows := sqlmock.NewRows([]string{"json"}).AddRow(`["some" : "value"]`)
	suite.Mock.ExpectQuery("schema.tablewithrows").WillReturnRows(rows)
	ctx := context.WithValue(suite.Ctx, "table", "tablewithrows")

	http.HandlerFunc(ShowList).ServeHTTP(suite.Rr, suite.Req.WithContext(ctx))
	suite.NotEmpty(suite.Rr.Body, "Should have a JSON body with rows")
	suite.Equal(http.StatusOK, suite.Rr.Code, "Should be HTTP OK")
	suite.Equal(`["some" : "value"]`, fmt.Sprint(suite.Rr.Body), "Should return results array")

}

//If the role does not have privileges to view the table
func (suite *ShowListTests) TestShowList_schema_exists_table_notauthorised() {

	suite.Mock.ExpectQuery("nonauthedrole").WillReturnError(&pq.Error{
		Message: "Role not authorised for that table",
		Code:    "42501",
	})
	ctx := context.WithValue(suite.Ctx, "role", "nonauthedrole")

	http.HandlerFunc(ShowList).ServeHTTP(suite.Rr, suite.Req.WithContext(ctx))
	suite.NotEmpty(suite.Rr.Body, "Should have a JSON error body")
	suite.Equal(http.StatusForbidden, suite.Rr.Code, "Should be HTTP Status Forbidden when not authorised to view a table")

}

//If the URL contains paramaters
func (suite *ShowListTests) TestShowList_schema_exists_table_exists_with_rows_goodparams() {

	rows := sqlmock.NewRows([]string{"json"}).AddRow(`["some" : "values"]`)
	suite.Mock.ExpectQuery("ORDER").WillReturnRows(rows)
	ctx := context.WithValue(suite.Ctx, "table", "tablewithrows")
	q, _ := url.ParseQuery(`orderby=price`)
	ctx = context.WithValue(ctx, "queries", q)

	http.HandlerFunc(ShowList).ServeHTTP(suite.Rr, suite.Req.WithContext(ctx))
	suite.NotEmpty(suite.Rr.Body, "Should have a JSON body with rows")
	suite.Equal(http.StatusOK, suite.Rr.Code, "Should be HTTP OK")

}

///////////////////////////////////////////////////////////////////////
// Tests for func ShowSingle(w http.ResponseWriter, r *http.Request) //
///////////////////////////////////////////////////////////////////////

type ShowSingleTests struct {
	HandlerTests
}

func TestShowSingleTests(t *testing.T) {
	suite.Run(t, new(ShowSingleTests))
}

//Tests setup
func (suite *ShowSingleTests) SetupTest() {
	suite.Req, _ = http.NewRequest("GET", "", nil)
	suite.Rr = httptest.NewRecorder()
	//Setup db mocks
	DB, suite.Mock, _ = sqlmock.New()
	suite.Ctx = context.WithValue(suite.Req.Context(), "role", "role")
	suite.Ctx = context.WithValue(suite.Ctx, "userID", "123456789")
	suite.Ctx = context.WithValue(suite.Ctx, "schema", "schema")
	suite.Ctx = context.WithValue(suite.Ctx, "table", "table")
	suite.Ctx = context.WithValue(suite.Ctx, "record", "record")
}

//Test handler with context values unset, which in theory should never happen (routing, middleware etc)
func (suite *ShowSingleTests) TestShowSingle_missing_context_values() {

	http.HandlerFunc(ShowSingle).ServeHTTP(suite.Rr, suite.Req)
	suite.Equal(http.StatusBadRequest, suite.Rr.Code, "Missing context values - should be bad request")
	suite.NotEmpty(suite.Rr.Body, "Should have a JSON body with error message")

}

//If the schema specified in the URL doesn't exist then 404
func (suite *ShowSingleTests) TestShowSingle_schema_notexists() {

	suite.Mock.ExpectQuery("noschema").WillReturnError(&pq.Error{
		Message: "Schema doesn't exist",
		Code:    "42P01",
	})
	ctx := context.WithValue(suite.Ctx, "schema", "noschema")

	http.HandlerFunc(ShowSingle).ServeHTTP(suite.Rr, suite.Req.WithContext(ctx))
	suite.NotEmpty(suite.Rr.Body, fmt.Sprint(suite.Rr.Body))
	suite.Equal(http.StatusNotFound, suite.Rr.Code, fmt.Sprint(suite.Rr.Body))

}

//If the table specified in the URL doesn't exist then 404
func (suite *ShowSingleTests) TestShowSingle_schema_exists_table_notexists() {

	suite.Mock.ExpectQuery("schema.notable").WillReturnError(&pq.Error{
		Message: "Table doesn't exist",
		Code:    "42P01",
	})
	ctx := context.WithValue(suite.Ctx, "table", "notable")

	http.HandlerFunc(ShowSingle).ServeHTTP(suite.Rr, suite.Req.WithContext(ctx))
	suite.NotEmpty(suite.Rr.Body, fmt.Sprint(suite.Rr.Body))
	suite.Equal(http.StatusNotFound, suite.Rr.Code, fmt.Sprint(suite.Rr.Body))

}

//If the record specified in the URL doesn't exist then 404
func (suite *ShowSingleTests) TestShowSingle_schema_exists_table_exists_record_notexists() {

	suite.Mock.ExpectQuery("norecord").WillReturnError(&pq.Error{
		Message: "Record doesn't exist",
		Code:    "42P01",
	})
	ctx := context.WithValue(suite.Ctx, "record", "norecord")

	http.HandlerFunc(ShowSingle).ServeHTTP(suite.Rr, suite.Req.WithContext(ctx))
	suite.NotEmpty(suite.Rr.Body, fmt.Sprint(suite.Rr.Body))
	suite.Equal(http.StatusNotFound, suite.Rr.Code, fmt.Sprint(suite.Rr.Body))

}

//If the role does not have privileges to view the table
func (suite *ShowSingleTests) TestShowList_record_notauthorised() {

	suite.Mock.ExpectQuery("nonauthedrole").WillReturnError(&pq.Error{
		Message: "Role not authorised for that record",
		Code:    "42501",
	})
	ctx := context.WithValue(suite.Ctx, "role", "nonauthedrole")

	http.HandlerFunc(ShowSingle).ServeHTTP(suite.Rr, suite.Req.WithContext(ctx))
	suite.NotEmpty(suite.Rr.Body, fmt.Sprint(suite.Rr.Body))
	suite.Equal(http.StatusForbidden, suite.Rr.Code, fmt.Sprint(suite.Rr.Body))

}

//If the table specified exists and has rows, then 200 and return the JSON array
func (suite *ShowSingleTests) TestShowList_record_exists() {

	rows := sqlmock.NewRows([]string{"json"}).AddRow(`["some" : "value"]`)
	suite.Mock.ExpectQuery("record").WillReturnRows(rows)

	http.HandlerFunc(ShowSingle).ServeHTTP(suite.Rr, suite.Req.WithContext(suite.Ctx))
	suite.NotEmpty(suite.Rr.Body, "Should have a JSON body with rows")
	suite.Equal(http.StatusOK, suite.Rr.Code, "Should be HTTP OK")
	suite.Equal(`["some" : "value"]`, fmt.Sprint(suite.Rr.Body), "Should return record JSON")

}

/////////////////////////////////////////////////////////////////////////
// Tests for func InsertRecord(w http.ResponseWriter, r *http.Request) //
/////////////////////////////////////////////////////////////////////////

type InsertRecordTests struct {
	HandlerTests
}

func TestInsertRecordTests(t *testing.T) {
	suite.Run(t, new(InsertRecordTests))
}

//Tests setup
func (suite *InsertRecordTests) SetupTest() {
	suite.Req, _ = http.NewRequest("POST", "", nil)
	suite.Rr = httptest.NewRecorder()
	//Setup db mocks
	DB, suite.Mock, _ = sqlmock.New()
	suite.Ctx = context.WithValue(suite.Req.Context(), "role", "role")
	suite.Ctx = context.WithValue(suite.Ctx, "userID", "123456789")
	suite.Ctx = context.WithValue(suite.Ctx, "schema", "schema")
	suite.Ctx = context.WithValue(suite.Ctx, "table", "table")
}

//Test handler with context values unset, which in theory should never happen (routing, middleware etc)
func (suite *InsertRecordTests) TestInsertRecord_missing_context_values() {

	http.HandlerFunc(InsertRecord).ServeHTTP(suite.Rr, suite.Req)
	suite.Equal(http.StatusBadRequest, suite.Rr.Code, "Missing context values - should be bad request")
	suite.NotEmpty(suite.Rr.Body, "Should have a JSON body with error message")

}

//If the schema specified in the URL doesn't exist then 404
func (suite *InsertRecordTests) TestInsertRecord_schema_notexists() {

	suite.Mock.ExpectQuery("noschema").WillReturnError(&pq.Error{
		Message: "Schema doesn't exist",
		Code:    "42P01",
	})
	ctx := context.WithValue(suite.Ctx, "schema", "noschema")

	http.HandlerFunc(InsertRecord).ServeHTTP(suite.Rr, suite.Req.WithContext(ctx))
	suite.NotEmpty(suite.Rr.Body, "Should have a JSON error body")
	suite.Equal(http.StatusNotFound, suite.Rr.Code, "Should be HTTP Status Not Found when schema doesn't exist")

}

//If the table specified in the URL doesn't exist then 404
func (suite *InsertRecordTests) TestInsertRecord_schema_exists_table_notexists() {

	suite.Mock.ExpectQuery("schema.notable").WillReturnError(&pq.Error{
		Message: "Table doesn't exist",
		Code:    "42P01",
	})
	ctx := context.WithValue(suite.Ctx, "table", "notable")

	http.HandlerFunc(InsertRecord).ServeHTTP(suite.Rr, suite.Req.WithContext(ctx))
	suite.NotEmpty(suite.Rr.Body, "Should have a JSON error body")
	suite.Equal(http.StatusNotFound, suite.Rr.Code, "Should be HTTP Status Not Found when table doesn't exist")

}

//If the body is blank, do a default insert and return 200
func (suite *InsertRecordTests) TestInsertRecord_no_body() {

	rows := sqlmock.NewRows([]string{"json"}).AddRow(`["some" : "value"]`)
	suite.Mock.ExpectQuery("DEFAULT VALUES").WillReturnRows(rows)

	http.HandlerFunc(InsertRecord).ServeHTTP(suite.Rr, suite.Req.WithContext(suite.Ctx))
	suite.NotEmpty(suite.Rr.Body, "Should have a JSON response body")
	suite.Equal(http.StatusOK, suite.Rr.Code, "Should be HTTP OK")

}

//If the body is empty,  do a default insert and return 200
func (suite *InsertRecordTests) TestInsertRecord_empty_body() {

	rows := sqlmock.NewRows([]string{"json"}).AddRow(`["some" : "value"]`)
	suite.Mock.ExpectQuery("DEFAULT VALUES").WillReturnRows(rows)

	b := []byte("")
	suite.Req, _ = http.NewRequest("POST", "", bytes.NewBuffer(b))

	http.HandlerFunc(InsertRecord).ServeHTTP(suite.Rr, suite.Req.WithContext(suite.Ctx))
	suite.NotEmpty(suite.Rr.Body, fmt.Sprintf("Error JSON: %s", suite.Rr.Body))
	suite.Equal(http.StatusOK, suite.Rr.Code, fmt.Sprintf("Error JSON: %s", suite.Rr.Body))

}

//If the body is empty json,  do a default insert and return 200
func (suite *InsertRecordTests) TestInsertRecord_empty_json() {

	rows := sqlmock.NewRows([]string{"json"}).AddRow(`["some" : "value"]`)
	suite.Mock.ExpectQuery("DEFAULT VALUES").WillReturnRows(rows)

	b := []byte("{}")
	suite.Req, _ = http.NewRequest("POST", "", bytes.NewBuffer(b))

	http.HandlerFunc(InsertRecord).ServeHTTP(suite.Rr, suite.Req.WithContext(suite.Ctx))
	suite.NotEmpty(suite.Rr.Body, fmt.Sprintf("Error JSON: %s", suite.Rr.Body))
	suite.Equal(http.StatusOK, suite.Rr.Code, fmt.Sprintf("Error JSON: %s", suite.Rr.Body))

}

//If the body is malformed JSON,  return 400
func (suite *InsertRecordTests) TestInsertRecord_malformed_json() {

	b := []byte("{jjkh")
	suite.Req, _ = http.NewRequest("POST", "", bytes.NewBuffer(b))

	http.HandlerFunc(InsertRecord).ServeHTTP(suite.Rr, suite.Req.WithContext(suite.Ctx))
	suite.NotEmpty(suite.Rr.Body, fmt.Sprintf("Error JSON: %s", suite.Rr.Body))
	suite.Equal(http.StatusBadRequest, suite.Rr.Code, fmt.Sprintf("Error JSON: %s", suite.Rr.Body))

}

//If the role does not have privileges to insert into the table
func (suite *InsertRecordTests) TestShowList_record_notauthorised() {

	suite.Mock.ExpectQuery("nonauthedrole").WillReturnError(&pq.Error{
		Message: "Role not authorised for that table",
		Code:    "42501",
	})
	ctx := context.WithValue(suite.Ctx, "role", "nonauthedrole")

	http.HandlerFunc(InsertRecord).ServeHTTP(suite.Rr, suite.Req.WithContext(ctx))
	suite.NotEmpty(suite.Rr.Body, fmt.Sprint(suite.Rr.Body))
	suite.Equal(http.StatusForbidden, suite.Rr.Code, fmt.Sprint(suite.Rr.Body))

}

//If the record is inserted, 200
func (suite *InsertRecordTests) TestInsertRecord_ok() {

	rows := sqlmock.NewRows([]string{"json"}).AddRow(`{"some" : "value"}`)
	suite.Mock.ExpectQuery("table").WillReturnRows(rows)

	b := []byte(`{"some" : "value"}`)
	suite.Req, _ = http.NewRequest("POST", "", bytes.NewBuffer(b))

	http.HandlerFunc(InsertRecord).ServeHTTP(suite.Rr, suite.Req.WithContext(suite.Ctx))
	suite.NotEmpty(suite.Rr.Body, fmt.Sprint(suite.Rr.Body))
	suite.Equal(http.StatusOK, suite.Rr.Code, fmt.Sprint(suite.Rr.Body))
	suite.Equal(`{"some" : "value"}`, fmt.Sprint(suite.Rr.Body), "Should return inserted record")

}

/////////////////////////////////////////////////////////////////////////
// Tests for func DeleteRecord(w http.ResponseWriter, r *http.Request) //
/////////////////////////////////////////////////////////////////////////

type DeleteRecordTests struct {
	HandlerTests
}

func TestDeleteRecordTests(t *testing.T) {
	suite.Run(t, new(DeleteRecordTests))
}

//Tests setup
func (suite *DeleteRecordTests) SetupTest() {
	suite.Req, _ = http.NewRequest("DELETE", "", nil)
	suite.Rr = httptest.NewRecorder()
	//Setup db mocks
	DB, suite.Mock, _ = sqlmock.New()
	suite.Ctx = context.WithValue(suite.Req.Context(), "role", "role")
	suite.Ctx = context.WithValue(suite.Ctx, "userID", "123456789")
	suite.Ctx = context.WithValue(suite.Ctx, "schema", "schema")
	suite.Ctx = context.WithValue(suite.Ctx, "table", "table")
	suite.Ctx = context.WithValue(suite.Ctx, "record", "record")
}

//Test handler with context values unset, which in theory should never happen (routing, middleware etc)
func (suite *DeleteRecordTests) TestDeleteRecord_missing_context_values() {

	http.HandlerFunc(DeleteRecord).ServeHTTP(suite.Rr, suite.Req)
	suite.Equal(http.StatusBadRequest, suite.Rr.Code, "Missing context values - should be bad request")
	suite.NotEmpty(suite.Rr.Body, "Should have a JSON body with error message")

}

//If the schema specified in the URL doesn't exist then 404
func (suite *DeleteRecordTests) TestDelete_schema_notexists() {

	suite.Mock.ExpectExec("noschema").WillReturnError(&pq.Error{
		Message: "Schema doesn't exist",
		Code:    "42P01",
	})
	ctx := context.WithValue(suite.Ctx, "schema", "noschema")

	http.HandlerFunc(DeleteRecord).ServeHTTP(suite.Rr, suite.Req.WithContext(ctx))
	suite.NotEmpty(suite.Rr.Body, fmt.Sprint(suite.Rr.Body))
	suite.Equal(http.StatusNotFound, suite.Rr.Code, fmt.Sprint(suite.Rr.Body))

}

//If the table specified in the URL doesn't exist then 404
func (suite *DeleteRecordTests) TestDelete_schema_exists_table_notexists() {

	suite.Mock.ExpectExec("schema.notable").WillReturnError(&pq.Error{
		Message: "Table doesn't exist",
		Code:    "42P01",
	})
	ctx := context.WithValue(suite.Ctx, "table", "notable")

	http.HandlerFunc(DeleteRecord).ServeHTTP(suite.Rr, suite.Req.WithContext(ctx))
	suite.NotEmpty(suite.Rr.Body, fmt.Sprint(suite.Rr.Body))
	suite.Equal(http.StatusNotFound, suite.Rr.Code, fmt.Sprint(suite.Rr.Body))

}

//If the record specified in the URL doesn't exist then 404
func (suite *DeleteRecordTests) TestDelete_record_notexists() {

	res := sqlmock.NewResult(0, 0)
	suite.Mock.ExpectExec("norecord").WillReturnResult(res)

	ctx := context.WithValue(suite.Ctx, "table", "norecord")

	http.HandlerFunc(DeleteRecord).ServeHTTP(suite.Rr, suite.Req.WithContext(ctx))
	suite.NotEmpty(suite.Rr.Body, fmt.Sprint(suite.Rr.Body))
	suite.Equal(http.StatusNotFound, suite.Rr.Code, fmt.Sprint(suite.Rr.Body))

}

//If the user does not have privileges to delete that record
func (suite *DeleteRecordTests) TestDeleteRecord_record_notauthorised() {

	suite.Mock.ExpectExec("nonautheduser").WillReturnError(&pq.Error{
		Message: "User not authorised to delete that record",
		Code:    "42501",
	})
	ctx := context.WithValue(suite.Ctx, "userID", "nonautheduser")

	http.HandlerFunc(DeleteRecord).ServeHTTP(suite.Rr, suite.Req.WithContext(ctx))
	suite.NotEmpty(suite.Rr.Body, fmt.Sprint(suite.Rr.Body))
	suite.Equal(http.StatusForbidden, suite.Rr.Code, fmt.Sprint(suite.Rr.Body))

}

//If the record is deleted correctly return 204 and no body
func (suite *DeleteRecordTests) TestDelete_ok() {

	res := sqlmock.NewResult(1, 1)
	suite.Mock.ExpectExec("record").WillReturnResult(res)

	http.HandlerFunc(DeleteRecord).ServeHTTP(suite.Rr, suite.Req.WithContext(suite.Ctx))
	//TODO: test empty body - suite.Empty doesn't seem to work
	suite.Equal(http.StatusNoContent, suite.Rr.Code, fmt.Sprint(suite.Rr.Body))

}

//TODO - test hypehn underscore test

/////////////////////////////////////////////////////////////////////////
// Tests for func UpdateRecord(w http.ResponseWriter, r *http.Request) //
/////////////////////////////////////////////////////////////////////////

type UpdateRecordTests struct {
	HandlerTests
}

func TestUpdateRecordTests(t *testing.T) {
	suite.Run(t, new(UpdateRecordTests))
}

//Tests setup
func (suite *UpdateRecordTests) SetupTest() {
	suite.Req, _ = http.NewRequest("PATCH", "", nil)
	suite.Rr = httptest.NewRecorder()
	//Setup db mocks
	DB, suite.Mock, _ = sqlmock.New()
	suite.Ctx = context.WithValue(suite.Req.Context(), "role", "role")
	suite.Ctx = context.WithValue(suite.Ctx, "userID", "123456789")
	suite.Ctx = context.WithValue(suite.Ctx, "schema", "schema")
	suite.Ctx = context.WithValue(suite.Ctx, "table", "table")
	suite.Ctx = context.WithValue(suite.Ctx, "record", "record")
}

//Test handler with context values unset, which in theory should never happen (routing, middleware etc)
func (suite *UpdateRecordTests) TestDeleteRecord_missing_context_values() {

	http.HandlerFunc(UpdateRecord).ServeHTTP(suite.Rr, suite.Req)
	suite.Equal(http.StatusBadRequest, suite.Rr.Code, "Missing context values - should be bad request")
	suite.NotEmpty(suite.Rr.Body, "Should have a JSON body with error message")

}

//If the body is blank, 400
func (suite *UpdateRecordTests) TestUpdateRecord_no_body() {

	http.HandlerFunc(UpdateRecord).ServeHTTP(suite.Rr, suite.Req.WithContext(suite.Ctx))
	suite.NotEmpty(suite.Rr.Body, "Should have a JSON response body")
	suite.Equal(http.StatusBadRequest, suite.Rr.Code, "Should be HTTP Bad Request")

}

//If the body is empty, 400
func (suite *UpdateRecordTests) TestUpdateRecord_empty_body() {

	b := []byte("")
	suite.Req, _ = http.NewRequest("PATCH", "", bytes.NewBuffer(b))

	http.HandlerFunc(UpdateRecord).ServeHTTP(suite.Rr, suite.Req.WithContext(suite.Ctx))
	suite.NotEmpty(suite.Rr.Body, "Should have a JSON response body")
	suite.Equal(http.StatusBadRequest, suite.Rr.Code, "Should be HTTP Bad Request")

}

//If the body is empty json, 400
func (suite *UpdateRecordTests) TestUpdateRecord_empty_json() {

	b := []byte("{}")
	suite.Req, _ = http.NewRequest("PATCH", "", bytes.NewBuffer(b))

	http.HandlerFunc(UpdateRecord).ServeHTTP(suite.Rr, suite.Req.WithContext(suite.Ctx))
	suite.NotEmpty(suite.Rr.Body, "Should have a JSON response body")
	suite.Equal(http.StatusBadRequest, suite.Rr.Code, "Should be HTTP Bad Request")

}

//If the body is malformed JSON, return 400
func (suite *UpdateRecordTests) TestUpdateRecord_malformed_json() {

	b := []byte("{jjkh")
	suite.Req, _ = http.NewRequest("PATCH", "", bytes.NewBuffer(b))

	http.HandlerFunc(UpdateRecord).ServeHTTP(suite.Rr, suite.Req.WithContext(suite.Ctx))
	suite.NotEmpty(suite.Rr.Body, "Should have a JSON response body")
	suite.Equal(http.StatusBadRequest, suite.Rr.Code, "Should be HTTP Bad Request")

}

//If update successful then 200
func (suite *UpdateRecordTests) TestUpdate_ok() {

	rows := sqlmock.NewRows([]string{"json"}).AddRow(`{"some" : "value"}`)
	suite.Mock.ExpectQuery("record").WillReturnRows(rows)

	b := []byte(`{"some" : "value"}`)
	suite.Req, _ = http.NewRequest("PATCH", "", bytes.NewBuffer(b))

	http.HandlerFunc(UpdateRecord).ServeHTTP(suite.Rr, suite.Req.WithContext(suite.Ctx))
	suite.NotEmpty(suite.Rr.Body, fmt.Sprint(suite.Rr.Body))
	suite.Equal(`{"some" : "value"}`, fmt.Sprint(suite.Rr.Body), "Should return patched record")
	suite.Equal(http.StatusOK, suite.Rr.Code, fmt.Sprint(suite.Rr.Body))

}

//If record not found then 404
func (suite *UpdateRecordTests) TestUpdate_record_not_found() {

	rows := sqlmock.NewRows([]string{"json"})
	suite.Mock.ExpectQuery("norecord").WillReturnRows(rows) //Empty result, which is what happens when no such record

	b := []byte(`{"some" : "value"}`)
	suite.Req, _ = http.NewRequest("PATCH", "", bytes.NewBuffer(b))

	ctx := context.WithValue(suite.Ctx, "record", "norecord")

	http.HandlerFunc(UpdateRecord).ServeHTTP(suite.Rr, suite.Req.WithContext(ctx))
	suite.Equal(http.StatusNotFound, suite.Rr.Code, fmt.Sprint(suite.Rr.Body))

}

//If the schema specified in the URL doesn't exist then 404
func (suite *UpdateRecordTests) TestUpdate_schema_notexists() {

	suite.Mock.ExpectQuery("noschema").WillReturnError(&pq.Error{
		Message: "Schema doesn't exist",
		Code:    "42P01",
	})

	b := []byte(`{"some" : "value"}`)
	suite.Req, _ = http.NewRequest("PATCH", "", bytes.NewBuffer(b))

	ctx := context.WithValue(suite.Ctx, "schema", "noschema")

	http.HandlerFunc(UpdateRecord).ServeHTTP(suite.Rr, suite.Req.WithContext(ctx))
	suite.NotEmpty(suite.Rr.Body, fmt.Sprint(suite.Rr.Body))
	suite.Equal(http.StatusNotFound, suite.Rr.Code, fmt.Sprint(suite.Rr.Body))

}

//If the table specified in the URL doesn't exist then 404
func (suite *UpdateRecordTests) TestUpdate_schema_exists_table_notexists() {

	suite.Mock.ExpectQuery("schema.notable").WillReturnError(&pq.Error{
		Message: "Table doesn't exist",
		Code:    "42P01",
	})

	b := []byte(`{"some" : "value"}`)
	suite.Req, _ = http.NewRequest("PATCH", "", bytes.NewBuffer(b))

	ctx := context.WithValue(suite.Ctx, "table", "notable")

	http.HandlerFunc(UpdateRecord).ServeHTTP(suite.Rr, suite.Req.WithContext(ctx))
	suite.NotEmpty(suite.Rr.Body, fmt.Sprint(suite.Rr.Body))
	suite.Equal(http.StatusNotFound, suite.Rr.Code, fmt.Sprint(suite.Rr.Body))

}
