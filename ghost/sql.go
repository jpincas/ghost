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
package ghost

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"

	"strings"

	sq "gopkg.in/Masterminds/squirrel.v1"
)

//Query is the basic building block of an SQL query
type Query struct {
	//SQL is for when you need to provide complete, preformed SQL
	//This option will override any BaseAQL + args
	SQL string
	//BaseSQL is a formatted SQL string with placeholder for SQLArgs
	BaseSQL string
	//SQLArgs are inserted into the BaseSQL in the order they appear
	SQLArgs []interface{}
	//WhereAnyOfValues appends a WHERE clause to match multiple values
	//If 'WherAnyOfKey' is not provided, it will default to 'id'
	WhereAnyOfValues []interface{}
	//The fieldname for the multiple-matching WHERE clause
	WhereAnyOfKey string
	//Indicate whether to rquest JSON array or object
	//and when unmarshalling, whether map or slice of maps
	IsList bool
	//Role to execute the query as
	Role string
	//UserID to set on the query context
	UserID string
	//CacheLevel specifies the level of caching to use: all, role, user
	//Omitting, or using any other value, will bypass caching
	CacheLevel string
	//cacheKey is the key used to store the SQL query in the cache
	cacheKey string
}

//Execute runs a query against the data store and returns JSON
func (q Query) Execute() (string, error) {

	//First, test to see if complete SQL has been provided
	//If it hasn't, then create it based on 'BaseSQL' and 'SQLArgs'
	if q.SQL == "" {

		q.SQL = fmt.Sprintf(q.BaseSQL, q.SQLArgs...)

		//For a multiple whereanyof
		if len(q.WhereAnyOfValues) != 0 {

			//Default to id
			if q.WhereAnyOfKey == "" {
				q.WhereAnyOfKey = "id"
			}

			q.SQL += fmt.Sprintf(SQLToAddWhereAnyOfValues, q.WhereAnyOfKey, commaSeparatedStringify(q.WhereAnyOfValues...))
		}

	}

	//Create the sqlQuery
	sqlQuery := SqlQuery(q.SQL)

	//Return JSON array or object
	if q.IsList {
		sqlQuery = sqlQuery.RequestMultipleResultsAsJSONArray()
	} else {
		sqlQuery = sqlQuery.RequestSingleResultAsJSONObject()
	}

	//Add sql role and user id, caching as you go depending
	//on the cache level specified
	if q.CacheLevel == "all" {
		q.cacheKey = sqlQuery.ToSQLCacheKey()
	}

	//Add role
	if q.Role != "" {
		sqlQuery = sqlQuery.SetQueryRole(q.Role)
	}

	if q.CacheLevel == "role" {
		q.cacheKey = sqlQuery.ToSQLCacheKey()
	}

	//Add user id
	if q.UserID != "" {
		sqlQuery = sqlQuery.SetUserID(q.UserID)
	}

	if q.CacheLevel == "user" {
		q.cacheKey = sqlQuery.ToSQLCacheKey()
	}

	//Transform to SQL string and reset on the struct
	q.SQL = sqlQuery.ToSQLString()

	//Execute the query
	return App.Store.executeQuery(q)

}

//ExecuteAndUnmarshall runs a query against the datastore and returns both
//for lists: []map[string]interfaace{}
//for objects: map[string]interface{}
//The corresponding unused data structure is set to nil
func (q Query) ExecuteAndUnmarshall() (list []map[string]interface{}, single map[string]interface{}, err error) {

	//Execute to JSON first
	var dbResponse string
	dbResponse, err = q.Execute()
	if err != nil {
		return nil, nil, err
	}

	if q.IsList {

		//Check for empty result from DB
		if dbResponse == "" {
			return list, nil, nil
		}

		if err := json.Unmarshal([]byte(dbResponse), &list); err != nil {
			return nil, nil, err
		}

		return list, nil, nil
	}

	//Check for empty result from DB
	if dbResponse == "" {
		return nil, single, nil
	}

	if err := json.Unmarshal([]byte(dbResponse), &single); err != nil {
		return nil, nil, err
	}

	return nil, single, nil

}

//SqlQuery is an SQL query string with various methods available for transformation
type SqlQuery string

//RequestMultipleResultsAsJSONArray transforms the SQL query to return a JSON array of results
//Use when multiple lines are going to be returned
func (s SqlQuery) RequestMultipleResultsAsJSONArray() SqlQuery {

	newQuery := SqlQuery(fmt.Sprintf(SQLToRequestMultipleResultsAsJSONArray, s))
	return newQuery
}

//RequestSingleResultAsJSONObject transforms the SQL query to return a JSON object of the result
//Used when a single line is going to be returned
func (s SqlQuery) RequestSingleResultAsJSONObject() SqlQuery {

	newQuery := SqlQuery(fmt.Sprintf(SQLToRequestSingleResultAsJSONObject, s))
	return newQuery
}

//SetQueryRole prepends the database role with which to execute the query
func (s SqlQuery) SetQueryRole(role string) SqlQuery {

	newQuery := SqlQuery(fmt.Sprintf(SQLToSetLocalRole, role, s))
	return newQuery
}

//SetUserID prepends the user id variable with which to execute the query
func (s SqlQuery) SetUserID(userID string) SqlQuery {

	newQuery := SqlQuery(fmt.Sprintf(SQLToSetUserID, userID, s))
	return newQuery
}

//ToSQLCacheKey transforms the SQL query into a cacheable string key
func (s SqlQuery) ToSQLCacheKey() string {

	return fmt.Sprint(s)

}

//ToSQLString transforms an SqlQuery to a plain string
//Generally the last step before execution
func (s SqlQuery) ToSQLString() string {

	LogDebug("SQL", true, fmt.Sprint(s), nil)
	return fmt.Sprint(s)
}

//TODO: deprecate this
//QueryBuilder builds an SqlQuery from multiple URL query paramaters
func QueryBuilder(schema string, table string, queries url.Values) SqlQuery {

	//Concat - schema.table
	tn := fmt.Sprintf("%s.%s", schema, table)

	//Start with all the products from the table
	p := sq.Select("*").From(tn)

	//Loop through all the URL qeueries
	for key, value := range queries {

		if strings.ToLower(key) == "orderby" {
			p = p.OrderBy(value[0])
		} else if strings.ToLower(key) == "limit" {
			l, _ := strconv.ParseUint(value[0], 10, 64)
			p = p.Limit(l)
		} else {
			p = p.Where(fmt.Sprintf(`%s %s`, key, value[0]))
		}

	}

	//Build the basic SQL
	sql, _, _ := p.ToSql()
	//Return as JSON array request
	return SqlQuery(sql)
}
