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
	"encoding/json"
	"fmt"

	"github.com/ecosystemsoftware/ecosystem/core"
)

//GetSchema returns a Schema struct corresponding to dbSchema.dbTable
func GetSchema(dbSchema, dbTable string) (Schema, error) {

	dbInfo, err := readDBInfo(dbSchema, dbTable)
	if err != nil {
		return nil, err
	}
	return dbInfo.dbInfoToSchema(), nil

}

//GetSchemaJSON returns the JSON for the schema corresponding to dbSchema.dbTable
func GetSchemaJSON(dbSchema, dbTable string) (string, error) {

	s, err := GetSchema(dbSchema, dbTable)
	if err != nil {
		return "", err
	}
	return s.outputSchema(), nil

}

//Property is the column from the database
type Property struct {
	DataType         string      `json:"type"`
	Required         bool        `json:"required, omitempty"`
	Default          interface{} `json:"value, omitempty"`
	MaxLength        int64       `json:"maxlength, omitempty"`
	ReferencesTable  string      `json:"reftable, omitempty"`
	ReferencesColumn string      `json:"refcolumn, omitempty"`
}

//Schema is a struct representing the schematised structure of a database table or view
//Basically just a map of column names with their properties
type Schema map[string]Property

//output returns a JSON string of a Schema struct
func (s Schema) outputSchema() string {

	output, _ := json.MarshalIndent(s, "", "\t")
	return string(output)

}

//dbInfoRow is a row from the informational results returned by the DB
//Each row corresponds to the information about a single database column
//For ease of ref, struct properties are named after db columns
type dbInfoRow struct {
	Column_name string
	Data_type   string
	Is_nullable string
	//Although the default can be any data type, the db has to choose to store the default as a fixed type, and it's a string
	Column_default           sql.NullString
	Character_Maximum_Length sql.NullInt64
	Foreign_Table_Name       sql.NullString
	Foreign_Column_Name      sql.NullString
}

//dbInfo is the list of database columns for the table being queried
type dbInfo []dbInfoRow

//dbInfoToSchema translates a db information table to a schema struct ready to output to a JSON schema
func (i dbInfo) dbInfoToSchema() Schema {

	//Create a default schema
	s := Schema{}
	//Iterate over the db fields
	for _, dbRow := range i {
		//Translate the data type to the JSON schema format (string, number, array etc)
		//and append to the properties map
		s[dbRow.Column_name] = Property{
			DataType:         dbTypeToSchemaType(dbRow.Data_type),
			Required:         dbIsNullableToRequired(dbRow.Is_nullable),
			Default:          dbDefaultToDefault(dbRow.Column_default),
			MaxLength:        dbMaxLengthToMaxLength(dbRow.Character_Maximum_Length),
			ReferencesTable:  dbForeignKeyTableToRefTable(dbRow.Foreign_Table_Name),
			ReferencesColumn: dbForeignKeyColumnToRefColumn(dbRow.Foreign_Column_Name),
		}
	}
	return s
}

//readDBInfo returns a structure containing db meta information for table.schema
func readDBInfo(schema string, table string) (DBInfo dbInfo, err error) {
	sqlString := fmt.Sprintf(sqlToGetTableInfo, "admin", table, schema, table)
	rows, err := core.DB.Query(sqlString)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	for rows.Next() {
		var r dbInfoRow
		if err := rows.Scan(&r.Column_name, &r.Data_type, &r.Is_nullable, &r.Column_default, &r.Character_Maximum_Length, &r.Foreign_Table_Name, &r.Foreign_Column_Name); err != nil {
			return nil, err
		}
		DBInfo = append(DBInfo, r)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return DBInfo, nil
}
