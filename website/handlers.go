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
	"bytes"
	"encoding/json"
	"net/http"
	"net/url"
	"path"
	"strings"

	"github.com/ecosystemsoftware/ecosystem/core"
	"github.com/lib/pq"
	"github.com/pressly/chi/render"
)

//ShowList returns HTML for a list from the database table specified in the URL
//For UNPROTECTED routes, the role will be set to 'web'
//For PROTECTED routes, the role and userID will be set according to the auth middleware
func ShowList(w http.ResponseWriter, r *http.Request) {

	var out bytes.Buffer

	//Retrieve all the context variables
	//These are assigned by correct routing and middleware, so no need to check existence
	ctx := r.Context()
	role, isRole := ctx.Value("role").(string)
	userID, isUserID := ctx.Value("userID").(string)
	schema, ok1 := ctx.Value("schema").(string)
	table, ok2 := ctx.Value("table").(string)
	queries, _ := ctx.Value("queries").(url.Values) //Not obligatory

	//In normal operation, routing and middleware will make sure that these variables
	//are always present.  However, to aid in testing of the handler, we include a check
	if !core.AllOK(ok1, ok2) {

		templates.ExecuteTemplate(&out, "defaulterror.html", page{
			Records:  []map[string]interface{}{},
			Schema:   "",
			Table:    "",
			Site:     *new(SiteBuilder),
			HttpCode: http.StatusBadRequest,
			Message:  "Schema/table not correctly set on context",
		})
		render.Status(r, http.StatusBadRequest)
		render.HTML(w, r, out.String())
		return
	}

	//In order to avoid having to have seperate handlers for protected and unprotected routes
	//If this request has passed through the auth middleware, then set the role to whatever
	//has been established.  If not, default to web
	if !isRole {
		role = "web"
	}

	//Build out the basic SQL from the URL Query parameters
	sql := core.QueryBuilder(schema, table, queries).RequestMultipleResultsAsJSONArray().SetQueryRole(role)

	//and set the userId on the sqlquery
	if isUserID {
		sql = sql.SetUserID(userID)
	}

	//Turn the SQL query into a string
	sqlString := sql.ToSQLString()

	//Define the template name
	templateName := path.Join(schema, table+"-list.html") //not the directoy location but the defined template name

	//Run the DB query
	var dbResponse string
	if err := core.DB.QueryRow(sqlString).Scan(&dbResponse); err != nil {

		//Check for an sql scan error indicating the json has come back empty
		if strings.Contains(err.Error(), "sql") {

			//Send back the tempalte without records
			if err := templates.ExecuteTemplate(w, templateName, page{
				Records: []map[string]interface{}{},
				Schema:  schema,
				Table:   table,
				Site:    *new(SiteBuilder),
			}); err != nil {
				renderTemplateError(&out)
			}

		} else {
			//Else report a DB error as usual
			dbError := err.(*pq.Error)
			httpCode := core.DBErrorCodeToHTTPErrorCode(dbError.Code)

			templates.ExecuteTemplate(&out, "defaulterror.html", page{
				Records:  []map[string]interface{}{},
				Schema:   schema,
				Table:    table,
				Site:     *new(SiteBuilder),
				HttpCode: httpCode,
				Message:  dbError.Message,
				DBCode:   dbError.Code,
			})

			render.Status(r, httpCode)
			render.HTML(w, r, out.String())

		}
	} else {
		//If found
		var records []map[string]interface{}
		//Attempt to bind the JSON from the DB response to a map of the record
		//Note: there doesn't seem to be a simple native way to scan db results to a generic map or struct easily
		//so the workaround used here is to get JSON from the databse and unmarhall it to a map map[string]interface{}
		//this actually works pretty well
		if err := json.Unmarshal([]byte(dbResponse), &records); err != nil {

			render.Status(r, http.StatusInternalServerError)
			templates.ExecuteTemplate(w, "error.html", page{
				Records:  []map[string]interface{}{},
				Schema:   schema,
				Table:    table,
				Site:     *new(SiteBuilder),
				HttpCode: http.StatusInternalServerError,
				Message:  err.Error(),
			})

		} else {

			templates.ExecuteTemplate(w, templateName, page{
				Records: records,
				Schema:  schema,
				Table:   table,
				Site:    *new(SiteBuilder),
			})

		}
	}

}

// //WebShowSingle returns HTML for single record from the database table specified in the URL
// //with the slug specified in the URL in the 'slug' database field.
// //For UNPROTECTED routes, the role will be set to 'web'
// //For PROTECTED routes, the role and userID will be set according to the auth middleware
// func ShowSingle(c *gin.Context) {

// 	table := core.HyphensToUnderscores(c.Param("table"))
// 	schema := core.HyphensToUnderscores(c.Param("schema"))
// 	slug := c.Param("slug")

// 	//In order to avoid having to have seperate handlers for protected and unprotected routes
// 	//If this request has passed through the auth middleware, then set the role to whatever
// 	//has been established.  If not, default to web
// 	var role string
// 	if r, ok := c.Get("role"); ok {
// 		role = r.(string)
// 	} else {
// 		role = "web"
// 	}

// 	//If a specific or default template does exist
// 	//Build the basic SQL query and convert to JSON request
// 	//Note the use of 'slug' in the query.  Not slug = not visible on the website
// 	sql := core.SqlQuery(fmt.Sprintf(ecosql.ToSelectRecordBySlug, schema, table, slug)).RequestSingleResultAsJSONObject().SetQueryRole(role)

// 	//In order to avoid having to have seperate handlers for protected and unprotected routes
// 	//If this request has passed through the auth middleware, then set the userId on the sqlquery
// 	if u, ok := c.Get("userID"); ok {
// 		sql = sql.SetUserID(u.(string))
// 	}

// 	//Turn the SQL query into a string
// 	sqlString := sql.ToSQLString()

// 	//Run the DB query
// 	var dbResponse string
// 	if err := core.DB.QueryRow(sqlString).Scan(&dbResponse); err != nil {

// 		//Check for an sql scan error indicating the json has come back empty
// 		if strings.Contains(err.Error(), "sql") {
// 			//Which is a 404
// 			c.HTML(http.StatusNotFound, "error.html", gin.H{
// 				"httpCode": http.StatusNotFound,
// 				"message":  err.Error(),
// 			})

// 		} else {
// 			//Else report a DB error as usual
// 			dbError := err.(*pq.Error)
// 			httpCode := core.DBErrorCodeToHTTPErrorCode(dbError.Code)
// 			c.HTML(httpCode, "error.html", gin.H{
// 				"httpCode": httpCode,
// 				"message":  dbError.Message,
// 				"dbCode":   dbError.Code,
// 			})
// 		}
// 	} else {
// 		//If found
// 		var record map[string]interface{}
// 		//Attempt to bind the JSON from the DB response to a map of the record
// 		//Note: there doesn't seem to be a simple native way to scan db results to a generic map or struct easily
// 		//so the workaround used here is to get JSON from the databse and unmarhall it to a map map[string]interface{}
// 		//this actually works pretty well
// 		if err := json.Unmarshal([]byte(dbResponse), &record); err != nil {
// 			c.HTML(http.StatusInternalServerError, "error.html", gin.H{
// 				"httpCode": http.StatusInternalServerError,
// 				"message":  err.Error(),
// 			})
// 		} else {
// 			var s core.SiteBuilder
// 			//Build the template name
// 			template := path.Join(schema, table+"-single.html") //not the directoy location but the defined template name
// 			c.HTML(http.StatusOK, template, gin.H{
// 				"record": record,
// 				"schema": schema,
// 				"table":  table,
// 				"site":   s,
// 			})
// 		}
// 	}

// }

// //WebShowEntryPage shows either the schema homepage or the site-level homepage
// func WebShowEntryPage(c *gin.Context) {
// 	schema := core.HyphensToUnderscores(c.Param("schema"))
// 	var template string
// 	if schema != "" {
// 		template = schema + "/index.html"
// 	} else {
// 		template = "index.html"
// 	}
// 	var s core.SiteBuilder

// 	c.HTML(http.StatusOK, template, gin.H{
// 		"site": s,
// 	})
// }

// func ShowCategory(c *gin.Context) {

// 	table := core.HyphensToUnderscores(c.Param("table"))
// 	schema := core.HyphensToUnderscores(c.Param("schema"))
// 	cat := c.Param("cat")

// 	var catJs string
// 	sql := core.SqlQuery(fmt.Sprintf(ecosql.ToSelectWebCategoryWhere, cat)).RequestSingleResultAsJSONObject().SetQueryRole("web").ToSQLString()

// 	//Attempt to find the category
// 	if err := core.DB.QueryRow(sql).Scan(&catJs); err != nil {
// 		c.HTML(http.StatusNotFound, "error.html", gin.H{
// 			"message":  err.Error(),
// 			"httpCode": http.StatusNotFound,
// 		})
// 	} else {
// 		//Set up the map into which the request body will be read
// 		var catData map[string]interface{}
// 		//Attempt to bind the JSON in the request body to the map
// 		if err := json.Unmarshal([]byte(catJs), &catData); err != nil {
// 			c.HTML(http.StatusInternalServerError, "error.html", gin.H{
// 				"message":  err.Error(),
// 				"httpCpde": http.StatusInternalServerError,
// 			})
// 		} else {
// 			//In the case that the category is found AND the JSON is unmarshalled correctly
// 			var itemsJs string
// 			itemsSQL := core.SqlQuery(fmt.Sprintf(ecosql.ToSelectKeywordedRecords, schema, table, cat)).RequestMultipleResultsAsJSONArray().SetQueryRole("web").ToSQLString()
// 			core.DB.QueryRow(itemsSQL).Scan(&itemsJs)
// 			var itemsData []map[string]interface{}
// 			json.Unmarshal([]byte(itemsJs), &itemsData)
// 			//Site Context
// 			var s core.SiteBuilder
// 			template := path.Join(schema, table+"-categories.html") //not the directoy location but the defined template name
// 			c.HTML(http.StatusOK, template, gin.H{
// 				"category": catData,
// 				"table":    table,
// 				"schema":   schema,
// 				"records":  itemsData,
// 				"site":     s,
// 			})
// 		}

// 	}

// }

// func GetEcoSystemJS(c *gin.Context) {
// 	var cf map[string]string
// 	viper.Unmarshal(&cf)

// 	c.HTML(http.StatusOK, "ecosystem.js", gin.H{
// 		"config": cf,
// 	})
// }
