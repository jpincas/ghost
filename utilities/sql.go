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

	sq "gopkg.in/Masterminds/squirrel.v1"
)

type sqlQuery string

func (s sqlQuery) requestMultipleResultsAsJSONArray() sqlQuery {
	newQuery := sqlQuery(fmt.Sprintf(`WITH results AS (%s) SELECT array_to_json(array_agg(row_to_json(results))) from results`, s))
	return newQuery
}

func (s sqlQuery) requestSingleResultAsJSONObject() sqlQuery {
	newQuery := sqlQuery(fmt.Sprintf(`WITH results AS (%s) SELECT row_to_json(results) from results`, s))
	return newQuery
}

func (s sqlQuery) setQueryRole(role string) sqlQuery {
	newQuery := sqlQuery(fmt.Sprintf(`SET LOCAL ROLE %s; %s `, role, s))
	return newQuery
}

func (s sqlQuery) setUserID(userID string) sqlQuery {
	newQuery := sqlQuery(fmt.Sprintf(`SET my.user_id = '%s'; %s `, userID, s))
	return newQuery
}

func (s sqlQuery) toSQLString() string {
	return fmt.Sprint(s)
}

func queryBuilder(table string, queries url.Values) sqlQuery {

	//Start with all the products from the table
	p := sq.Select("*").From(table)

	//Loop through all the URL qeueries
	for key, value := range queries {

		if key == "orderBy" {
			p = p.OrderBy(value[0])
		} else if key == "limit" {
			l, _ := strconv.ParseUint(value[0], 10, 64)
			p = p.Limit(l)
		} else {
			p = p.Where(fmt.Sprintf(`%s %s`, key, value[0]))
		}

	}

	//Build the basic SQL
	sql, _, _ := p.ToSql()
	//Return as JSON array request
	return sqlQuery(sql)
}
