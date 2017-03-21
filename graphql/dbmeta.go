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

package graphql

import (
	"fmt"

	"github.com/ecosystemsoftware/ecosystem/core"
)

type dbTable struct {
	tableName string
	tableType string
}

func getDBTables(dbSchema string) (tables []dbTable, err error) {

	sqlString := core.SqlQuery(fmt.Sprintf(sqlToGetTablesInSchema, dbSchema)).SetQueryRole("admin").ToSQLString()

	rows, err := core.DB.Query(sqlString)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	for rows.Next() {
		var t dbTable
		if err := rows.Scan(&t.tableName, &t.tableType); err != nil {
			return nil, err
		}
		tables = append(tables, t)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return tables, nil
}

func getDBSchemas() (schemas []string, err error) {

	sqlString := core.SqlQuery(sqlToGetSchemasInDB).SetQueryRole("admin").ToSQLString()

	rows, err := core.DB.Query(sqlString)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	for rows.Next() {
		var s string
		if err := rows.Scan(&s); err != nil {
			return nil, err
		}
		schemas = append(schemas, s)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return schemas, nil
}
