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

package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"net/url"

	"github.com/ecosystemsoftware/ecosystem/ecosql"
	eco "github.com/ecosystemsoftware/ecosystem/utilities"

	"github.com/lib/pq"
	"github.com/pressly/chi/render"
)

//ApiMagicCode processes a request for a magic code
// func ApiMagicCode(c *gin.Context) {

// 	//Set up the map into which the request body will be read
// 	var r map[string]interface{}
// 	//Attempt to bind the JSON in the request body to the map
// 	if err := c.BindJSON(&r); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{
// 			"code":    http.StatusBadRequest,
// 			"message": err.Error(),
// 		})
// 	} else {
// 		//Try to read 'email'
// 		email, ok := r["email"]
// 		if ok {
// 			//If 'email' is set, request a magic code
// 			err := eco.RequestMagicCode(email.(string), "emailMagicCode.html")
// 			//If sending of the magic code fails (user doesn't exist, email fails etc)
// 			if err != nil {
// 				c.JSON(http.StatusBadRequest, gin.H{
// 					"code":    http.StatusBadRequest,
// 					"message": err.Error(),
// 				})
// 			} else {
// 				//If the magic code goes through OK
// 				c.JSON(http.StatusOK, gin.H{
// 					"code":    http.StatusOK,
// 					"message": "Magic code sent to " + email.(string),
// 				})
// 			}
// 		} else {
// 			//In the case that no email address is provided
// 			c.JSON(http.StatusBadRequest, gin.H{
// 				"code":    http.StatusBadRequest,
// 				"message": "No email address provided",
// 			})
// 		}
// 	}

// }

//ShowList shows a list of records from the database
func ShowList(w http.ResponseWriter, r *http.Request) {

	//Retrieve all the context variables
	//These are assigned by correct routing and middleware, so no need to check existence
	ctx := r.Context()
	role, ok1 := ctx.Value("role").(string)
	userID, ok2 := ctx.Value("userID").(string)
	schema, ok3 := ctx.Value("schema").(string)
	table, ok4 := ctx.Value("table").(string)
	queries, _ := ctx.Value("queries").(url.Values) //Not obligatory

	//In normal operation, routing and middleware will make sure that these variables
	//are always present.  However, to aid in testing of the handler, we include a check
	if !eco.AllOK(ok1, ok2, ok3, ok4) {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, eco.ResponseError{http.StatusBadRequest, pq.ErrorCode(""), "Missing required values on context", "", "", ""})
		return
	}

	//Build the SQL string and execute query
	var json string
	sqlString := eco.QueryBuilder(schema, table, queries).RequestMultipleResultsAsJSONArray().SetQueryRole(role).SetUserID(userID).ToSQLString()
	err := eco.DB.QueryRow(sqlString).Scan(&json) //Only one row is returned as JSON is returned by Postgres

	if err != nil {
		//Check for an sql scan error indicating the json has come back empty
		if strings.Contains(err.Error(), "sql") {
			render.SetContentType(render.ContentTypeJSON)
			w.Write([]byte("[]")) //Send back a blank array
		} else {
			//Work out what error has ocurred, 'translate' to a relevant http error and render an error
			dbError := err.(*pq.Error)
			httpCode := eco.DBErrorCodeToHTTPErrorCode(dbError.Code)
			render.Status(r, httpCode)
			render.JSON(w, r, eco.ResponseError{httpCode, dbError.Code, dbError.Message, schema, table, ""})
		}
	} else {
		render.SetContentType(render.ContentTypeJSON)
		w.Write([]byte(json))
	}

}

//ShowSingle shows a single record from the database
func ShowSingle(w http.ResponseWriter, r *http.Request) {

	//Retrieve all the context variables
	//These are assigned by correct routing and middleware, so no need to check existence
	ctx := r.Context()
	role, ok1 := ctx.Value("role").(string)
	userID, ok2 := ctx.Value("userID").(string)
	schema, ok3 := ctx.Value("schema").(string)
	table, ok4 := ctx.Value("table").(string)
	record, ok5 := ctx.Value("record").(string)

	//In normal operation, routing and middleware will make sure that these variables
	//are always present.  However, to aid in testing of the handler, we include a check
	if !eco.AllOK(ok1, ok2, ok3, ok4, ok5) {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, eco.ResponseError{http.StatusBadRequest, pq.ErrorCode(""), "Missing required values on context", "", "", ""})
		return
	}

	var json string
	sqlString := eco.SqlQuery(fmt.Sprintf(ecosql.ToSelectWhere, schema, table, record)).RequestSingleResultAsJSONObject().SetQueryRole(role).SetUserID(userID).ToSQLString()
	err := eco.DB.QueryRow(sqlString).Scan(&json) //Only one row is returned as JSON is returned by Postgres

	if err != nil {
		//Work out what error has ocurred, 'translate' to a relevant http error and render an error
		dbError := err.(*pq.Error)
		httpCode := eco.DBErrorCodeToHTTPErrorCode(dbError.Code)
		render.Status(r, httpCode)
		render.JSON(w, r, eco.ResponseError{httpCode, dbError.Code, dbError.Message, schema, table, record})

	} else {
		render.SetContentType(render.ContentTypeJSON)
		w.Write([]byte(json))
	}

}

func InsertRecord(w http.ResponseWriter, r *http.Request) {

	//Retrieve all the context variables
	//These are assigned by correct routing and middleware, so no need to check existence
	ctx := r.Context()
	role, ok1 := ctx.Value("role").(string)
	userID, ok2 := ctx.Value("userID").(string)
	schema, ok3 := ctx.Value("schema").(string)
	table, ok4 := ctx.Value("table").(string)

	//In normal operation, routing and middleware will make sure that these variables
	//are always present.  However, to aid in testing of the handler, we include a check
	if !eco.AllOK(ok1, ok2, ok3, ok4) {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, eco.ResponseError{http.StatusBadRequest, pq.ErrorCode(""), "Missing required values on context", "", "", ""})
		return
	}

	//To reference the base table from the view (if necessary), only use the portion of the table name before the first hyphen/underscore
	table = strings.Split(table, "_")[0]

	//Holder for database response after insert
	var (
		dbResponse      string
		requestBody     map[string]interface{}
		sqlString       string
		bodyDecodeError error
	)

	//If r.body is not nil (as in, body doesn't even exist), read and decode
	if r.Body != nil {
		d := json.NewDecoder(r.Body)
		bodyDecodeError = d.Decode(&requestBody)
	}

	//For any decode error other than EOF, return a bad request with the decode error
	if bodyDecodeError != nil && bodyDecodeError != io.EOF {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, eco.ResponseError{http.StatusBadRequest, "", bodyDecodeError.Error(), schema, table, ""})
		return
	}

	//Filter for a nil body, blank body or empty JSON and do a DEFAULT insert
	if r.Body == nil || bodyDecodeError == io.EOF || len(requestBody) == 0 {
		//In this special case, the database will default all fields
		sqlString = eco.SqlQuery(fmt.Sprintf(ecosql.ToInsertAllDefaultsReturningJSON, schema, table, table)).RequestSingleResultAsJSONObject().SetQueryRole(role).SetUserID(userID).ToSQLString()
	} else {
		//Otherwise build the insert SQL
		cols, vals := eco.MapToValsAndCols(requestBody)
		sqlString = eco.SqlQuery(fmt.Sprintf(ecosql.ToInsertReturningJSON, schema, table, cols, vals, table)).RequestSingleResultAsJSONObject().SetQueryRole(role).SetUserID(userID).ToSQLString()
	}

	//Hit the database and deal with errors
	if err := eco.DB.QueryRow(sqlString).Scan(&dbResponse); err != nil {
		dbError := err.(*pq.Error)
		httpCode := eco.DBErrorCodeToHTTPErrorCode(dbError.Code)
		render.Status(r, httpCode)
		render.JSON(w, r, eco.ResponseError{httpCode, dbError.Code, dbError.Message, schema, table, ""})
	} else {
		//If there are no database errors
		render.SetContentType(render.ContentTypeJSON)
		w.Write([]byte(dbResponse))
	}

}

func DeleteRecord(w http.ResponseWriter, r *http.Request) {

	//Retrieve all the context variables
	//These are assigned by correct routing and middleware, so no need to check existence
	ctx := r.Context()
	role, ok1 := ctx.Value("role").(string)
	userID, ok2 := ctx.Value("userID").(string)
	schema, ok3 := ctx.Value("schema").(string)
	table, ok4 := ctx.Value("table").(string)
	record, ok5 := ctx.Value("record").(string)

	//In normal operation, routing and middleware will make sure that these variables
	//are always present.  However, to aid in testing of the handler, we include a check
	if !eco.AllOK(ok1, ok2, ok3, ok4, ok5) {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, eco.ResponseError{http.StatusBadRequest, pq.ErrorCode(""), "Missing required values on context", "", "", ""})
		return
	}

	//To reference the base table from the view (if necessary), only use the portion of the table name before the first hyphen/underscore
	table = strings.Split(table, "_")[0]

	sqlString := eco.SqlQuery(fmt.Sprintf(ecosql.ToDeleteWhere, schema, table, record)).SetQueryRole(role).SetUserID(userID).ToSQLString()
	res, err := eco.DB.Exec(sqlString)

	if err != nil {
		//Work out what error has ocurred, 'translate' to a relevant http error and render an error
		dbError := err.(*pq.Error)
		httpCode := eco.DBErrorCodeToHTTPErrorCode(dbError.Code)
		render.Status(r, httpCode)
		render.JSON(w, r, eco.ResponseError{httpCode, dbError.Code, dbError.Message, schema, table, record})
	} else {
		//If 0 rows are affected then nothing has been deleted
		if n, _ := res.RowsAffected(); n == 0 {
			render.Status(r, http.StatusNotFound)
			render.JSON(w, r, eco.ResponseError{http.StatusNotFound, "", "No record with that id", schema, table, record})
		} else {
			render.NoContent(w, r)
		}

	}
}

func UpdateRecord(w http.ResponseWriter, r *http.Request) {

	//Retrieve all the context variables
	//These are assigned by correct routing and middleware, so no need to check existence
	ctx := r.Context()
	role, ok1 := ctx.Value("role").(string)
	userID, ok2 := ctx.Value("userID").(string)
	schema, ok3 := ctx.Value("schema").(string)
	table, ok4 := ctx.Value("table").(string)
	record, ok5 := ctx.Value("record").(string)

	//In normal operation, routing and middleware will make sure that these variables
	//are always present.  However, to aid in testing of the handler, we include a check
	if !eco.AllOK(ok1, ok2, ok3, ok4, ok5) {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, eco.ResponseError{http.StatusBadRequest, pq.ErrorCode(""), "Missing required values on context", "", "", ""})
		return
	}

	//To reference the base table from the view (if necessary), only use the portion of the table name before the first hyphen/underscore
	table = strings.Split(table, "_")[0]

	//Holder for database response after update
	var (
		dbResponse      string
		requestBody     map[string]interface{}
		bodyDecodeError error
	)

	//If r.body is not nil (as in, body doesn't even exist), read and decode
	if r.Body != nil {
		d := json.NewDecoder(r.Body)
		bodyDecodeError = d.Decode(&requestBody)
	}

	//Filter for a nil body, blank body or empty JSON - return bad response
	if r.Body == nil || len(requestBody) == 0 {

		message := "Invalid or absent request body"
		if bodyDecodeError != nil {
			message = bodyDecodeError.Error()
		}
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, eco.ResponseError{http.StatusBadRequest, "", message, schema, table, record})

	} else {

		//Otherwise build the insert SQL
		cols, vals := eco.MapToValsAndCols(requestBody)
		sqlString := eco.SqlQuery(fmt.Sprintf(ecosql.ToUpdateWhereReturningJSON, schema, table, cols, vals, record, table)).SetQueryRole(role).SetUserID(userID).ToSQLString()
		//Hit the database and deal with errors
		err := eco.DB.QueryRow(sqlString).Scan(&dbResponse)
		//Check for an sql scan error indicating the json has come back empty
		//which means that the record was not found, so 404
		if err != nil && strings.Contains(err.Error(), "sql") {
			render.Status(r, http.StatusNotFound)
			render.JSON(w, r, eco.ResponseError{http.StatusNotFound, "", "No record with that id", schema, table, record})
		} else if err != nil {
			dbError := err.(*pq.Error)
			httpCode := eco.DBErrorCodeToHTTPErrorCode(dbError.Code)
			render.Status(r, httpCode)
			render.JSON(w, r, eco.ResponseError{httpCode, dbError.Code, dbError.Message, schema, table, ""})
		} else {
			//If there are no database errors
			render.SetContentType(render.ContentTypeJSON)
			w.Write([]byte(dbResponse))
		}

	}

}

// func SearchList(c *gin.Context) {

// 	var json string
// 	schema := eco.HyphensToUnderscores(c.Param("schema"))
// 	table := c.Param("table")
// 	searchTerm := c.Param("searchTerm")

// 	role, _ := c.Get("role")
// 	userID, _ := c.Get("userID")

// 	sqlString := eco.SqlQuery(fmt.Sprintf(ecosql.ToFullTextSearch, table, searchTerm, table, schema, table, searchTerm)).SetQueryRole(role.(string)).SetUserID(userID.(string)).ToSQLString()

// 	err := eco.DB.QueryRow(sqlString).Scan(&json) //Only one row is returned as JSON is returned by Postgres

// 	if err != nil {
// 		//Check for an sql scan error indicating the json has come back empty
// 		if strings.Contains(err.Error(), "sql") {
// 			//In this case, no rows is OK - it's just an empty list
// 			c.String(http.StatusOK, json)
// 		} else {
// 			dbError := err.(*pq.Error)
// 			httpCode := eco.DBErrorCodeToHTTPErrorCode(dbError.Code)
// 			c.JSON(httpCode, gin.H{
// 				"code":    httpCode,
// 				"message": dbError.Message,
// 				"dbCode":  dbError.Code,
// 				"table":   table,
// 			})
// 		}
// 	} else {
// 		c.String(http.StatusOK, json)
// 	}

// }
