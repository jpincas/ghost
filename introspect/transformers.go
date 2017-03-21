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
	"database/sql"
	"strconv"
	"strings"
)

//dbTypeToSchemaType maps the database types to the schema types, defaulting to text
func dbTypeToSchemaType(dbType string) string {

	switch strings.ToLower(dbType) {
	case "numeric", "bigint", "integer":
		return "number"
	case "array":
		return "array"
	case "timestamp without time zone":
		return "date"
	case "boolean":
		return "boolean"
	default:
		return "string"
	}

}

//dbIsNullableToRequired transforms the db's bools to actual bools
func dbIsNullableToRequired(dbIsNullable string) bool {

	if dbIsNullable == "NO" {
		return true
	}

	return false
}

func dbDefaultToDefault(x sql.NullString) interface{} {

	if !x.Valid {
		return nil
	}

	xs := x.String

	//Integers
	if i, err := strconv.Atoi(xs); err == nil {
		return i
	}
	//Floats
	if f, err := strconv.ParseFloat(xs, 64); err == nil {
		return f
	}
	//Bools
	if b, err := strconv.ParseBool(xs); err == nil {
		return b
	}
	//Text
	//Get the actual default text from the db string
	if strings.Contains(xs, "::text") {
		a := strings.Split(xs, "'")
		return a[1]
	}

	return nil

}

func dbMaxLengthToMaxLength(x sql.NullInt64) int64 {
	if !x.Valid {
		return 0
	}
	return x.Int64
}

func dbForeignKeyTableToRefTable(x sql.NullString) string {

	if !x.Valid {
		return ""
	}
	return x.String

}

func dbForeignKeyColumnToRefColumn(x sql.NullString) string {

	if !x.Valid {
		return ""
	}
	return x.String

}
