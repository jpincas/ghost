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
	"fmt"
	"net/url"
	"strconv"

	"strings"

	"github.com/ecosystemsoftware/ecosystem/ecosql"
	sq "gopkg.in/Masterminds/squirrel.v1"
)

//SqlQuery is an SQL query string with various methods available for transformation
type SqlQuery string

//RequestMultipleResultsAsJSONArray transforms the SQL query to return a JSON array of results
//Use when multiple lines are going to be returned
func (s SqlQuery) RequestMultipleResultsAsJSONArray() SqlQuery {
	newQuery := SqlQuery(fmt.Sprintf(ecosql.ToRequestMultipleResultsAsJSONArray, s))
	return newQuery
}

//RequestSingleResultAsJSONObject transforms the SQL query to return a JSON object of the result
//Used when a single line is going to be returned
func (s SqlQuery) RequestSingleResultAsJSONObject() SqlQuery {
	newQuery := SqlQuery(fmt.Sprintf(ecosql.ToRequestSingleResultAsJSONObject, s))
	return newQuery
}

//SetQueryRole prepends the database role with which to execute the query
func (s SqlQuery) SetQueryRole(role string) SqlQuery {
	newQuery := SqlQuery(fmt.Sprintf(ecosql.ToSetLocalRole, role, s))
	return newQuery
}

//SetUserID prepends the user id variable with which to execute the query
func (s SqlQuery) SetUserID(userID string) SqlQuery {
	newQuery := SqlQuery(fmt.Sprintf(ecosql.ToSetUserID, userID, s))
	return newQuery
}

//ToSQLString transforms an SqlQuery to a plain string
//Generally the last step before execution
func (s SqlQuery) ToSQLString() string {
	//Uncomment to turn on SQL loggingfor debugging
	//log.Println(fmt.Sprint(s))
	return fmt.Sprint(s)
}

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
