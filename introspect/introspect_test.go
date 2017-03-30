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
	"fmt"
	"testing"

	sqlmock "gopkg.in/DATA-DOG/go-sqlmock.v1"

	ghost "github.com/jpincas/ghost/tools"
	"github.com/stretchr/testify/assert"
)

var columns = []string{"column_name", "data_type", "is_nullable", "column_default", "character_maximum_length", "foreign_table_name", "foreign_column_name"}

func TestGetSchema(t *testing.T) {
	var err error
	var mock sqlmock.Sqlmock
	ghost.DB, mock, err = sqlmock.New()
	if err != nil {
		t.Errorf("An error '%s' was not expected when opening a stub database connection", err)
	}
	defer ghost.DB.Close()

	rows := []string{"id,text,NO,'adefault'::text,10,reftable,refcol", "otherfield,numeric,YES,20,5,reftable,refcol"}

	dataRows := sqlmock.NewRows(columns)
	for _, row := range rows {
		dataRows.FromCSVString(row)
	}

	mock.ExpectQuery(fmt.Sprintf(`i.table_schema = '%s' AND i.table_name = '%s'`, "schema", "table")).WillReturnRows(dataRows)

	s, _ := GetSchema("schema", "table")
	expectedSchema := Schema{
		"id": Property{
			DataType:         "string",
			Required:         true,
			Default:          "adefault",
			MaxLength:        10,
			ReferencesTable:  "reftable",
			ReferencesColumn: "refcol",
		},
		"otherfield": Property{
			DataType:         "number",
			Required:         false,
			Default:          20,
			MaxLength:        5,
			ReferencesTable:  "reftable",
			ReferencesColumn: "refcol",
		},
	}

	assert.Equal(t, expectedSchema, s, "Should be equal")

	//Close the DB connection and test again
	ghost.DB.Close()
	_, err = GetSchema("schema", "table")
	assert.Error(t, err)

}

func TestGetSchemaJSON(t *testing.T) {
	var err error
	var mock sqlmock.Sqlmock
	ghost.DB, mock, err = sqlmock.New()
	if err != nil {
		t.Errorf("An error '%s' was not expected when opening a stub database connection", err)
	}
	defer ghost.DB.Close()

	rows := []string{"id,text,NO,'adefault'::text,10,reftable,refcol", "otherfield,numeric,YES,20,5,reftable,refcol"}

	dataRows := sqlmock.NewRows(columns)
	for _, row := range rows {
		dataRows.FromCSVString(row)
	}

	mock.ExpectQuery(fmt.Sprintf(`i.table_schema = '%s' AND i.table_name = '%s'`, "schema", "table")).WillReturnRows(dataRows)

	j, _ := GetSchemaJSON("schema", "table")

	//TODO: Find a way to compare the output JSON to a fixed JSON file whilst ignoring spaces, tabs, newlines etc.
	assert.NotEmpty(t, j, "Should have content")

	//Close the DB connection and test again
	ghost.DB.Close()
	_, err = GetSchemaJSON("schema", "table")
	assert.Error(t, err)

}

func TestDBInfoToSchema(t *testing.T) {

	assert := assert.New(t)

	//Run the tests
	for _, c := range cases {
		assert.Equal(c.out, c.in.dbInfoToSchema(), "Structs should match")
	}

}

func TestReadDBInfo(t *testing.T) {

	var err error
	var mock sqlmock.Sqlmock
	ghost.DB, mock, err = sqlmock.New()
	if err != nil {
		t.Errorf("An error '%s' was not expected when opening a stub database connection", err)
	}
	defer ghost.DB.Close()

	assert := assert.New(t)

	cases := []struct {
		schema, table string
		rows          []string
		want          dbInfo
	}{
		//One line
		{"testschema", "1", []string{"id,text,NO,'adefault'::text,10,reftable,refcol"},
			dbInfo{
				dbInfoRow{"id", "text", "NO", sql.NullString{"'adefault'::text", true}, sql.NullInt64{10, true}, sql.NullString{"reftable", true}, sql.NullString{"refcol", true}},
			},
		},
		//Two lines
		{"testschema", "2", []string{"id,text,NO,'adefault'::text,10,reftable,refcol", "otherfield,numeric,YES,20,5,reftable,refcol"},
			dbInfo{
				dbInfoRow{"id", "text", "NO", sql.NullString{"'adefault'::text", true}, sql.NullInt64{10, true}, sql.NullString{"reftable", true}, sql.NullString{"refcol", true}},
				dbInfoRow{"otherfield", "numeric", "YES", sql.NullString{"20", true}, sql.NullInt64{5, true}, sql.NullString{"reftable", true}, sql.NullString{"refcol", true}},
			},
		},
		//TODO
		//I would like to run this test but I think there is an issue with the mock driver
		//For database nulls it seems to send "" instead of a real db null
		//This causes the sql package to try to scan the "" into the target type ,rather than setting the valid property to false in the sql.Nullx struct

		// {"testschema", "3", []string{"id,text,NO,null,null,null,null"},
		// 	dbInfo{
		// 		dbInfoRow{"id", "text", "NO", sql.NullString{"'adefault'::text", false}, sql.NullInt64{10, false}, sql.NullString{"reftable", false}, sql.NullString{"refcol", false}},
		// 	},
		// },
	}

	for _, c := range cases {

		dataRows := sqlmock.NewRows(columns)
		for _, row := range c.rows {
			dataRows.FromCSVString(row)
		}
		//Implicitly tests the sqlQuery constant
		mock.ExpectQuery(fmt.Sprintf(`i.table_schema = '%s' AND i.table_name = '%s'`, c.schema, c.table)).WillReturnRows(dataRows)
		dbInfo, err := readDBInfo(c.schema, c.table)
		if err != nil {
			t.Errorf(err.Error())
		} else {
			assert.Equal(c.want, dbInfo, "Should be equal")
		}

	}

}
