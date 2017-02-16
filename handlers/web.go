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
	"encoding/json"
	"fmt"
	"net/http"
	"path"
	"strings"

	"github.com/ecosystemsoftware/ecosystem/ecosql"
	eco "github.com/ecosystemsoftware/ecosystem/utilities"
	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
)

//WebShowSingle returns HTML for single record from the database table specified in the URL
//with the slug specified in the URL in the 'slug' database field.
//For UNPROTECTED routes, the role will be set to 'web'
//For PROTECTED routes, the role and userID will be set according to the auth middleware
func WebShowSingle(c *gin.Context) {

	table := eco.HyphensToUnderscores(c.Param("table"))
	schema := eco.HyphensToUnderscores(c.Param("schema"))
	slug := c.Param("slug")

	//In order to avoid having to have seperate handlers for protected and unprotected routes
	//If this request has passed through the auth middleware, then set the role to whatever
	//has been established.  If not, default to web
	var role string
	if r, ok := c.Get("role"); ok {
		role = r.(string)
	} else {
		role = "web"
	}

	//If a specific or default template does exist
	//Build the basic SQL query and convert to JSON request
	//Note the use of 'slug' in the query.  Not slug = not visible on the website
	sql := eco.SqlQuery(fmt.Sprintf(ecosql.ToSelectRecordBySlug, schema, table, slug)).RequestSingleResultAsJSONObject().SetQueryRole(role)

	//In order to avoid having to have seperate handlers for protected and unprotected routes
	//If this request has passed through the auth middleware, then set the userId on the sqlquery
	if u, ok := c.Get("userID"); ok {
		sql = sql.SetUserID(u.(string))
	}

	//Turn the SQL query into a string
	sqlString := sql.ToSQLString()

	//Run the DB query
	var dbResponse string
	if err := eco.DB.QueryRow(sqlString).Scan(&dbResponse); err != nil {

		//Check for an sql scan error indicating the json has come back empty
		if strings.Contains(err.Error(), "sql") {
			//Which is a 404
			c.HTML(http.StatusNotFound, "error.html", gin.H{
				"httpCode": http.StatusNotFound,
				"message":  err.Error(),
			})

		} else {
			//Else report a DB error as usual
			dbError := err.(*pq.Error)
			httpCode := eco.DBErrorCodeToHTTPErrorCode(dbError.Code)
			c.HTML(httpCode, "error.html", gin.H{
				"httpCode": httpCode,
				"message":  dbError.Message,
				"dbCode":   dbError.Code,
			})
		}
	} else {
		//If found
		var record map[string]interface{}
		//Attempt to bind the JSON from the DB response to a map of the record
		//Note: there doesn't seem to be a simple native way to scan db results to a generic map or struct easily
		//so the workaround used here is to get JSON from the databse and unmarhall it to a map map[string]interface{}
		//this actually works pretty well
		if err := json.Unmarshal([]byte(dbResponse), &record); err != nil {
			c.HTML(http.StatusInternalServerError, "error.html", gin.H{
				"httpCode": http.StatusInternalServerError,
				"message":  err.Error(),
			})
		} else {
			var s eco.SiteBuilder
			//Build the template name
			template := path.Join(schema, table+"-single.html") //not the directoy location but the defined template name
			c.HTML(http.StatusOK, template, gin.H{
				"record": record,
				"schema": schema,
				"table":  table,
				"site":   s,
			})
		}
	}

}

//WebShowList returns HTML for a list from the database table specified in the URL
//For UNPROTECTED routes, the role will be set to 'web'
//For PROTECTED routes, the role and userID will be set according to the auth middleware
func WebShowList(c *gin.Context) {

	table := eco.HyphensToUnderscores(c.Param("table"))
	schema := eco.HyphensToUnderscores(c.Param("schema"))

	//In order to avoid having to have seperate handlers for protected and unprotected routes
	//If this request has passed through the auth middleware, then set the role to whatever
	//has been established.  If not, default to web
	var role string
	if r, ok := c.Get("role"); ok {
		role = r.(string)
	} else {
		role = "web"
	}

	//Build out the SQL from the URL Query parameters
	queries := c.Request.URL.Query()
	sql := eco.QueryBuilder(schema, table, queries).RequestMultipleResultsAsJSONArray().SetQueryRole(role)

	//In order to avoid having to have seperate handlers for protected and unprotected routes
	//If this request has passed through the auth middleware, then set the userId on the sqlquery
	if u, ok := c.Get("userID"); ok {
		sql = sql.SetUserID(u.(string))
	}

	//Turn the SQL query into a string
	sqlString := sql.ToSQLString()

	//Define the template name
	template := path.Join(schema, table+"-list.html") //not the directoy location but the defined template name

	//Run the DB query
	var dbResponse string
	if err := eco.DB.QueryRow(sqlString).Scan(&dbResponse); err != nil {

		//Check for an sql scan error indicating the json has come back empty
		if strings.Contains(err.Error(), "sql") {

			//Send back the tempalte without records
			var s eco.SiteBuilder
			var records []map[string]interface{}
			c.HTML(http.StatusOK, template, gin.H{
				"records": records,
				"schema":  schema,
				"table":   table,
				"site":    s,
			})

		} else {
			//Else report a DB error as usual
			dbError := err.(*pq.Error)
			httpCode := eco.DBErrorCodeToHTTPErrorCode(dbError.Code)
			c.HTML(httpCode, "error.html", gin.H{
				"httpCode": httpCode,
				"message":  dbError.Message,
				"dbCode":   dbError.Code,
			})
		}
	} else {
		//If found
		var records []map[string]interface{}
		//Attempt to bind the JSON from the DB response to a map of the record
		//Note: there doesn't seem to be a simple native way to scan db results to a generic map or struct easily
		//so the workaround used here is to get JSON from the databse and unmarhall it to a map map[string]interface{}
		//this actually works pretty well
		if err := json.Unmarshal([]byte(dbResponse), &records); err != nil {
			c.HTML(http.StatusInternalServerError, "error.html", gin.H{
				"httpCode": http.StatusInternalServerError,
				"message":  err.Error(),
			})
		} else {
			//Site Context
			var s eco.SiteBuilder
			c.HTML(http.StatusOK, template, gin.H{
				"records": records,
				"schema":  schema,
				"table":   table,
				"site":    s,
			})
		}
	}

}

//WebShowHomepage shows the homepage
func WebShowHomepage(c *gin.Context) {
	var s eco.SiteBuilder
	c.HTML(http.StatusOK, "index.html", gin.H{
		"site": s,
	})
}

func WebShowCategory(c *gin.Context) {

	table := eco.HyphensToUnderscores(c.Param("table"))
	schema := eco.HyphensToUnderscores(c.Param("schema"))
	cat := c.Param("cat")

	var catJs string
	sql := eco.SqlQuery(fmt.Sprintf(ecosql.ToSelectWebCategoryWhere, cat)).RequestSingleResultAsJSONObject().SetQueryRole("web").ToSQLString()

	//Attempt to find the category
	if err := eco.DB.QueryRow(sql).Scan(&catJs); err != nil {
		c.HTML(http.StatusNotFound, "error.html", gin.H{
			"message":  err.Error(),
			"httpCode": http.StatusNotFound,
		})
	} else {
		//Set up the map into which the request body will be read
		var catData map[string]interface{}
		//Attempt to bind the JSON in the request body to the map
		if err := json.Unmarshal([]byte(catJs), &catData); err != nil {
			c.HTML(http.StatusInternalServerError, "error.html", gin.H{
				"message":  err.Error(),
				"httpCpde": http.StatusInternalServerError,
			})
		} else {
			//In the case that the category is found AND the JSON is unmarshalled correctly
			var itemsJs string
			itemsSQL := eco.SqlQuery(fmt.Sprintf(ecosql.ToSelectKeywordedRecords, schema, table, cat)).RequestMultipleResultsAsJSONArray().SetQueryRole("web").ToSQLString()
			eco.DB.QueryRow(itemsSQL).Scan(&itemsJs)
			var itemsData []map[string]interface{}
			json.Unmarshal([]byte(itemsJs), &itemsData)
			//Site Context
			var s eco.SiteBuilder
			template := path.Join(schema, table+"-categories.html") //not the directoy location but the defined template name
			c.HTML(http.StatusOK, template, gin.H{
				"category": catData,
				"table":    table,
				"schema":   schema,
				"records":  itemsData,
				"site":     s,
			})
		}

	}

}

// Possible code for a more complate HTML api with all HTTP verbs
// Currently not used
// func WebInsertRecord(c *gin.Context) {

// 	//To reference the base table from the view (if necessary), only use the portion of the table name before the first hyphen/underscore
// 	table := strings.Split(eco.HyphensToUnderscores(c.Param("table")), "_")[0]

// 	role, _ := c.Get("role")
// 	userID, _ := c.Get("userID")

// 	//Check templates
// 	if ok, template := checkTemplate(schema, table, "single"); !ok {

// 		//If the template doesn't exist
// 		c.String(http.StatusBadRequest, "No template found for this record type and no default template.")

// 	} else {

// 		//Set up the map into which the request body will be read
// 		var r map[string]interface{}

// 		//Attempt to bind the JSON body of the request
// 		if err := c.BindJSON(&r); err != nil {
// 			c.String(http.StatusBadRequest, err.Error())
// 		} else {

// 			//Map the JSON body to vals and cols suitable for SQL
// 			cols, vals := mapToValsAndCols(r)

// 			//Build the SQL
// 			sqlString := eco.SqlQuery(fmt.Sprintf(`INSERT INTO %s(%s) VALUES (%s) returning row_to_json(%s)`, table, cols, vals, table)).RequestSingleResultAsJSONObject().SetQueryRole(role.(string)).SetUserID(userID.(string)).ToSQLString()

// 			//Run the DB query
// 			var dbResponse string
// 			if err := eco.DB.QueryRow(sqlString).Scan(&dbResponse); err != nil {
// 				dbError := err.(*pq.Error)
// 				httpCode := eco.DBErrorCodeToHTTPErrorCode(dbError.Code)
// 				c.HTML(httpCode, "error.html", gin.H{
// 					"httpCode": httpCode,
// 					"message":  dbError.Message,
// 					"dbCode":   dbError.Code,
// 				})
// 			} else {
// 				var newRecord map[string]interface{}
// 				//Attempt to bind the JSON in the request body to the map
// 				if err := json.Unmarshal([]byte(dbResponse), &newRecord); err != nil {
// 					c.HTML(http.StatusInternalServerError, "error.html", gin.H{
// 						"httpCode": http.StatusInternalServerError,
// 						"message":  err.Error(),
// 					})
// 				} else {
// 					c.HTML(http.StatusOK, template, gin.H{
// 						"record": newRecord,
// 						"table":  table,
// 					})
// 				}
// 			}
// 		}

// 	}

// }
