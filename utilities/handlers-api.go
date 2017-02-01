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
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"database/sql"

	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
)

func ApiMagicCode(c *gin.Context) {

	//Set up the map into which the request body will be read
	var r map[string]interface{}
	//Attempt to bind the JSON in the request body to the map
	if err := c.BindJSON(&r); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": err.Error(),
		})
	} else {
		//Try to read 'email'
		email, ok := r["email"]
		if ok {
			//If 'email' is set, request a magic code
			err := requestMagicCode(email.(string))
			//If sending of the magic code fails (user doesn't exist, email fails etc)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"code":    http.StatusBadRequest,
					"message": err.Error(),
				})
			} else {
				//If the magic code goes through OK
				c.JSON(http.StatusOK, gin.H{
					"code":    http.StatusOK,
					"message": "Magic code sent to " + email.(string),
				})
			}
		} else {
			//In the case that no email address is provided
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    http.StatusBadRequest,
				"message": "No email address provided",
			})
		}
	}

}

func ApiShowList(c *gin.Context) {

	var json string
	table := hyphensToUnderscores(c.Param("table"))

	//Build out the SQL from the URL Query parameters
	queries := c.Request.URL.Query()
	role, _ := c.Get("role")
	userID, _ := c.Get("userID")

	sqlString := queryBuilder(table, queries).requestMultipleResultsAsJSONArray().setQueryRole(role.(string)).setUserID(userID.(string)).toSQLString()

	err := db.QueryRow(sqlString).Scan(&json) //Only one row is returned as JSON is returned by Postgres

	if err != nil {
		//Check for an sql scan error indicating the json has come back empty
		if strings.Contains(err.Error(), "sql") {
			//In this case, no rows is OK - it's just an empty list
			c.String(http.StatusOK, json)
		} else {
			dbError := err.(*pq.Error)
			httpCode := dbErrorCodeToHTTPErrorCode(dbError.Code)
			c.JSON(httpCode, gin.H{
				"code":    httpCode,
				"message": dbError.Message,
				"dbCode":  dbError.Code,
				"table":   table,
			})
		}
	} else {
		c.String(http.StatusOK, json)
	}

}

func ApiShowSingle(c *gin.Context) {

	var json string

	table := hyphensToUnderscores(c.Param("table"))
	id := c.Param("id")

	role, _ := c.Get("role")
	userID, _ := c.Get("userID")
	sqlString := sqlQuery(fmt.Sprintf(`SELECT * FROM %s WHERE id = '%s'`, table, id)).requestSingleResultAsJSONObject().setQueryRole(role.(string)).setUserID(userID.(string)).toSQLString()

	err := db.QueryRow(sqlString).Scan(&json) //Only one row is returned as JSON is returned by Postgres

	if err != nil {
		//Check for an sql scan error indicating the json has come back empty
		if strings.Contains(err.Error(), "sql") {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    http.StatusBadRequest,
				"message": "record not found",
				"table":   table,
				"record":  id,
			})
		} else {
			//Else report a DB error as usual
			dbError := err.(*pq.Error)
			httpCode := dbErrorCodeToHTTPErrorCode(dbError.Code)
			c.JSON(httpCode, gin.H{
				"code":    httpCode,
				"message": dbError.Message,
				"dbCode":  dbError.Code,
				"table":   table,
				"record":  id,
			})
		}
	} else {
		c.String(http.StatusOK, json)
	}

}

func ApiInsertRecord(c *gin.Context) {

	//To reference the base table from the view (if necessary), only use the portion of the table name before the first hyphen/underscore
	table := strings.Split(hyphensToUnderscores(c.Param("table")), "_")[0]
	var dbResponse string

	role, _ := c.Get("role")
	userID, _ := c.Get("userID")

	//Set up the map into which the request body will be read
	var r map[string]interface{}

	//Attempt to bind body to JSON
	//We DON'T use the built in Gin BindJSON, because that automatically writes a 400 if there is any type
	//of error and we need more granular control for EOF errors
	decoder := json.NewDecoder(c.Request.Body)
	if err := decoder.Decode(&r); err != nil {

		//First filter for a nil body which is an immediate EOF
		if err == io.EOF {

			//In this special case, the database will default all fields
			//Not very common, but can happen if you are inserting a record with all defaults
			sqlString := sqlQuery(fmt.Sprintf(`INSERT INTO %s DEFAULT VALUES returning row_to_json(%s)`, table, table)).requestSingleResultAsJSONObject().setQueryRole(role.(string)).setUserID(userID.(string)).toSQLString()

			//Deal with database errors
			if err := db.QueryRow(sqlString).Scan(&dbResponse); err != nil {
				dbError := err.(*pq.Error)
				httpCode := dbErrorCodeToHTTPErrorCode(dbError.Code)
				c.JSON(httpCode, gin.H{
					"code":    httpCode,
					"message": dbError.Message,
					"dbCode":  dbError.Code,
					"table":   table,
				})
			} else {
				//If there are no database errors
				c.String(http.StatusOK, dbResponse)
			}
		} else {
			//Output the JSON decoding error in the case of any error other than EOF
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    http.StatusBadRequest,
				"message": err.Error(),
			})
		}

	} else {

		//In the case that there are no decoding errors
		//Map the JSON body to vals and cols suitable for SQL
		cols, vals := mapToValsAndCols(r)

		//Build the SQL
		sqlString := sqlQuery(fmt.Sprintf(`INSERT INTO %s(%s) VALUES (%s) returning row_to_json(%s)`, table, cols, vals, table)).requestSingleResultAsJSONObject().setQueryRole(role.(string)).setUserID(userID.(string)).toSQLString()

		//In order to get the return value, we use a QueryRow rather than EXEC and return the whole new row in JSON format
		//from the DB.
		err := db.QueryRow(sqlString).Scan(&dbResponse)

		if err != nil {
			dbError := err.(*pq.Error)
			httpCode := dbErrorCodeToHTTPErrorCode(dbError.Code)
			c.JSON(httpCode, gin.H{
				"code":    httpCode,
				"message": dbError.Message,
				"dbCode":  dbError.Code,
				"table":   table,
			})
		} else {
			c.String(http.StatusOK, dbResponse)
		}
	}

}

func ApiDeleteRecord(c *gin.Context) {

	//To reference the base table from the view (if necessary), only use the portion of the table name before the first hyphen/underscore
	table := strings.Split(hyphensToUnderscores(c.Param("table")), "_")[0]
	id := c.Param("id")

	role, _ := c.Get("role")
	userID, _ := c.Get("userID")

	sqlString := sqlQuery(fmt.Sprintf(`DELETE FROM %v WHERE id = '%v'`, table, id)).setQueryRole(role.(string)).setUserID(userID.(string)).toSQLString()

	res, err := db.Exec(sqlString)
	if err != nil {
		dbError := err.(*pq.Error)
		httpCode := dbErrorCodeToHTTPErrorCode(dbError.Code)
		c.JSON(httpCode, gin.H{
			"code":    httpCode,
			"message": dbError.Message,
			"dbCode":  dbError.Code,
			"table":   table,
			"record":  id,
		})
	} else {
		//If 0 rows are affected then nothing has been deleted
		if r, _ := res.RowsAffected(); r == 0 {
			c.JSON(http.StatusNotFound, gin.H{
				"code":    http.StatusNotFound,
				"message": "No records with that ID were found, and none were deleted",
				"table":   table,
				"record":  id,
			})
		} else {
			//If successful
			c.JSON(http.StatusOK, gin.H{
				"deleted": r,
			})
		}

	}
}

func ApiUpdateRecord(c *gin.Context) {

	//To reference the base table from the view (if necessary), only use the portion of the table name before the first hyphen/underscore
	table := strings.Split(hyphensToUnderscores(c.Param("table")), "_")[0]
	id := c.Param("id")

	role, _ := c.Get("role")
	userID, _ := c.Get("userID")

	var json string

	//Set up the map into which the request body will be read
	var r map[string]interface{}
	//Attempt to bind the JSON in the request body to the map
	if err := c.BindJSON(&r); err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
	} else {

		//Map the JSON body to vals and cols suitable for SQL
		cols, vals := mapToValsAndCols(r)

		//Build the SQL
		//again, surround the id with '' in case of non-numeric ids
		sqlString := sqlQuery(fmt.Sprintf(`UPDATE %s SET (%s) = (%s) WHERE id = '%v' returning row_to_json(%s)`, table, cols, vals, id, table)).setQueryRole(role.(string)).setUserID(userID.(string)).toSQLString()

		//In order to get the return value, we use a QueryRow rather than EXEC and return the whole new row in JSON format
		//from the DB.
		err := db.QueryRow(sqlString).Scan(&json)

		if err != nil {
			//In this case, if the record is not found, a db error will NOT be returned, just an empty row
			if err == sql.ErrNoRows {
				c.JSON(http.StatusBadRequest, gin.H{
					"code":    http.StatusBadRequest,
					"message": "record not found",
					"table":   table,
					"record":  id,
				})
			} else {
				//Else report a DB error (auth, table doesn't exist etc.) as usual
				dbError := err.(*pq.Error)
				httpCode := dbErrorCodeToHTTPErrorCode(dbError.Code)
				c.JSON(httpCode, gin.H{
					"code":    httpCode,
					"message": dbError.Message,
					"dbCode":  dbError.Code,
					"table":   table,
					"record":  id,
				})
			}
		} else {
			c.String(http.StatusOK, json)
		}

	}

}

func SearchList(c *gin.Context) {

	var json string
	table := c.Param("table")
	searchTerm := c.Param("searchTerm")

	role, _ := c.Get("role")
	userID, _ := c.Get("userID")

	sqlString := sqlQuery(fmt.Sprintf(`with item as (select to_tsvector(%s::text) @@ to_tsquery('%s') AS found, %s.* FROM %s) select array_to_json(array_agg(row_to_json(item))) FROM item WHERE item.found = TRUE OR item.id ILIKE '%s%%'`, table, searchTerm, table, table, searchTerm)).setQueryRole(role.(string)).setUserID(userID.(string)).toSQLString()

	err := db.QueryRow(sqlString).Scan(&json) //Only one row is returned as JSON is returned by Postgres

	if err != nil {
		//Check for an sql scan error indicating the json has come back empty
		if strings.Contains(err.Error(), "sql") {
			//In this case, no rows is OK - it's just an empty list
			c.String(http.StatusOK, json)
		} else {
			dbError := err.(*pq.Error)
			httpCode := dbErrorCodeToHTTPErrorCode(dbError.Code)
			c.JSON(httpCode, gin.H{
				"code":    httpCode,
				"message": dbError.Message,
				"dbCode":  dbError.Code,
				"table":   table,
			})
		}
	} else {
		c.String(http.StatusOK, json)
	}

}
