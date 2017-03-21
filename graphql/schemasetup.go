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
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/ecosystemsoftware/ecosystem/core"
	"github.com/ecosystemsoftware/ecosystem/introspect"
	"github.com/graphql-go/graphql"
)

//tablesAsTypes is a dictionary keyed by table name (prefixed with schema name if not unique)
// contains a graphQL type object with all table fields
var tablesAsTypes = map[string]*graphql.Object{}

//schemaRootSetup performs the initial setup for the schema
//It iterates through schemas and tables creating types
//and root queries
func schemaRootSetup() graphql.Fields {

	//Initialise the list of db schemas
	rootFields := graphql.Fields{}

	//Get the list of db schemas
	dbSchemaList, err := introspect.GetDBSchemas()
	if err != nil {
		core.LogFatal(core.LogEntry{"GRAPHQL", false, err.Error()})
	}

	//For each db schema
	for _, thisSchema := range dbSchemaList {

		//Get the tables in the schema
		dbTableList, err := introspect.GetDBTables(thisSchema)
		if err != nil {
			core.LogFatal(core.LogEntry{"GRAPHQL", false, err.Error()})
		}

		//Loop over the tables
		for _, thisTable := range dbTableList {

			//add it to the map of custom types, with unique name if necessary
			schemaTable := uniqueNamer(thisSchema, thisTable.tableName)
			tablesAsTypes[schemaTable] = graphql.NewObject(
				graphql.ObjectConfig{
					Name:   schemaTable,
					Fields: tableFieldsToGraphQLObjectFields(thisSchema, thisTable.tableName),
				},
			)

			//setup the field as the type
			rootFields[schemaTable] = &graphql.Field{
				Type:    graphql.NewList(tablesAsTypes[schemaTable]),
				Resolve: generateResolver(thisSchema, thisTable.tableName),
			}

		}

	}

	return rootFields

}

//generateResolver outputs a 'Resolve' function dynamically, hardcoding the schema and table variables
//This solves the problem of these values not being available to the resolves, which is what happens
//if the resolve function is directly specified
func generateResolver(schema, table string) func(p graphql.ResolveParams) (interface{}, error) {

	f := func(p graphql.ResolveParams) (interface{}, error) {

		var dbResponse string
		sqlString := core.QueryBuilder(schema, table, nil).RequestMultipleResultsAsJSONArray().SetQueryRole("admin").SetUserID("123456").ToSQLString()
		log.Println(sqlString)

		if err := core.DB.QueryRow(sqlString).Scan(&dbResponse); err != nil {
			//Only one row is returned as JSON is returned by Postgres
			//Empty result
			if strings.Contains(err.Error(), "sql") {
				return nil, nil
			}
		}
		//If found
		var records []map[string]interface{}
		if err := json.Unmarshal([]byte(dbResponse), &records); err != nil {
			return nil, err
		}

		return records, nil
	}

	return f

}

//unique name prefixes a table name with its schema if there is a clash
//bwetween tables names across schemas
func uniqueNamer(schema, table string) string {
	if _, ok := tablesAsTypes[table]; ok {
		return fmt.Sprintf("%s_%s", schema, table)
	}
	return table
}

//tableFieldsToGraphQLObjectFields introspects a table and sets up
//graphQL fields based on the table columns and their types
func tableFieldsToGraphQLObjectFields(schema, table string) graphql.Fields {

	fields := graphql.Fields{}
	tableFields, err := introspect.GetSchema(schema, table)
	if err != nil {
		core.LogFatal(core.LogEntry{"GRAPHQL", false, err.Error()})
	}
	for k, v := range tableFields {
		field := fieldBuilder(v)
		fields[k] = field
	}

	return fields

}

//fieldBuilder decides which type of GraphQl field to assign
func fieldBuilder(pgs pgschema.Property) *graphql.Field {

	switch pgs.DataType {
	case "string":
		return &graphql.Field{
			Type: graphql.String,
		}
	case "number":
		return &graphql.Field{
			Type: graphql.Float,
		}
	default:
		return &graphql.Field{
			Type: graphql.String,
		}
	}

}
