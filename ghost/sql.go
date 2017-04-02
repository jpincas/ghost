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

//ToSQLString transforms an SqlQuery to a plain string
//Generally the last step before execution
func (s SqlQuery) ToSQLString() string {

	LogDebug("SQL", true, fmt.Sprint(s), nil)
	return fmt.Sprint(s)
}

//QueryBuilder builds an SqlQuery from multiple URL query paramaters
//TODO: deprecate this
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

//Query
type Query struct {
	SQL     string
	BaseSQL string
	SQLArgs []interface{}
	IsList  bool
	Role    string
	UserID  string
}

func (q Query) ExecuteToJSON() (string, error) {

	//First, test to see if complete SQL has been provided
	//If it hasn't, then create it based on 'BaseSQL' and 'SQLArgs'
	if q.SQL == "" {
		q.SQL = fmt.Sprintf(q.BaseSQL, q.SQLArgs...)
	}

	//Create the sqlQuery
	sqlQuery := SqlQuery(q.SQL)

	//Return JSON array or object
	if q.IsList {
		sqlQuery = sqlQuery.RequestMultipleResultsAsJSONArray()
	} else {
		sqlQuery = sqlQuery.RequestSingleResultAsJSONObject()
	}

	//Add role and user id if required
	if q.Role != "" {
		sqlQuery = sqlQuery.SetQueryRole(q.Role)
	}
	if q.UserID != "" {
		sqlQuery = sqlQuery.SetUserID(q.UserID)
	}

	//Transform to SQL string
	sqlString := sqlQuery.ToSQLString()

	//Execure the query
	var dbResponse string
	if err := App.DB.QueryRow(sqlString).Scan(&dbResponse); err != nil {
		//Only one row is returned as JSON is returned by Postgres
		//Empty result
		if strings.Contains(err.Error(), "sql") {
			return "", nil
		}

		//Else its a database error
		return "", err

	}

	return dbResponse, nil

}

func (q Query) ExecuteToSlice() ([]map[string]interface{}, error) {

	q.IsList = true

	//Execute to JSON first
	var dbResponse string
	dbResponse, err := q.ExecuteToJSON()
	if err != nil {
		return nil, err
	}

	var result []map[string]interface{}

	//Check for empty result from DB
	if dbResponse == "" {
		return result, nil
	}

	if err := json.Unmarshal([]byte(dbResponse), &result); err != nil {
		return nil, err
	}

	return result, nil

}

func (q Query) ExecuteToMap() (map[string]interface{}, error) {

	q.IsList = false

	//Execute to JSON first
	var dbResponse string
	dbResponse, err := q.ExecuteToJSON()
	if err != nil {
		return nil, err
	}

	var result map[string]interface{}

	//Check for empty result from DB
	if dbResponse == "" {
		return result, nil
	}

	if err := json.Unmarshal([]byte(dbResponse), &result); err != nil {
		return nil, err
	}

	return result, nil

}
