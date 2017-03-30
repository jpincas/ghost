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

package introspect

import (
	"fmt"

	ghost "github.com/jpincas/ghost/tools"
)

//DBTable describes a database table
type DBTable struct {
	TableName string
	TableType string
}

//GetDBTables returns a list of all tables in a schema with table name and table type fields
func GetDBTables(dbSchema string) (tables []DBTable, err error) {

	sqlString := ghost.SqlQuery(fmt.Sprintf(sqlToGetTablesInSchema, dbSchema)).SetQueryRole("admin").ToSQLString()

	rows, err := ghost.DB.Query(sqlString)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	for rows.Next() {
		var t DBTable
		if err := rows.Scan(&t.TableName, &t.TableType); err != nil {
			return nil, err
		}
		tables = append(tables, t)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return tables, nil
}

//GetDBSchemas returns a list of non 'postgres' schemas in the database
func GetDBSchemas() (schemas []string, err error) {

	sqlString := ghost.SqlQuery(sqlToGetSchemasInDB).SetQueryRole("admin").ToSQLString()

	rows, err := ghost.DB.Query(sqlString)
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
