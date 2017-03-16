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

package core

import "testing"
import "net/url"

//TODO: this is a fairly useless test - just chains together all the methods
//and runs against a few different SqlQueries and logs the results
func TestSqlQueryTransformers(t *testing.T) {

	cases := []SqlQuery{
		*new(SqlQuery),
		SqlQuery(""),
		SqlQuery("SOME SQL"),
	}

	for _, c := range cases {
		c.RequestSingleResultAsJSONObject().RequestMultipleResultsAsJSONArray().SetQueryRole("tester").SetUserID("123456").ToSQLString()
	}

}

func TestQueryBuilder(t *testing.T) {
	cases := []struct {
		url, schema, table, want string
	}{
		{"value=='test'", "public", "table", `SELECT * FROM public.table WHERE value ='test'`},
		{"", "public", "table", `SELECT * FROM public.table`},
		{"value=>5", "public", "table", `SELECT * FROM public.table WHERE value >5`},
		{"orderBy=price", "public", "table", `SELECT * FROM public.table ORDER BY price`},
		{"orderby=price", "public", "table", `SELECT * FROM public.table ORDER BY price`},
		{"limit=10", "public", "table", `SELECT * FROM public.table LIMIT 10`},
		{"orderby=price&limit=10", "public", "table", `SELECT * FROM public.table ORDER BY price LIMIT 10`},
		{"value=='test'&orderby=price&limit=10", "public", "table", `SELECT * FROM public.table WHERE value ='test' ORDER BY price LIMIT 10`},
	}

	for _, c := range cases {
		queries, _ := url.ParseQuery(c.url)
		got := QueryBuilder(c.schema, c.table, queries).ToSQLString()
		if got != c.want {
			t.Errorf("QueryBuilder(%q, %q) == %q, want %q", c.table, c.url, got, c.want)
		}
	}

}
