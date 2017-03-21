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
)

var cases = []struct {
	in  dbInfo
	out Schema
}{
	{
		//No properties
		dbInfo{},
		Schema{},
	},
	{
		//One property - optional
		dbInfo{
			dbInfoRow{
				Column_name: "id",
				Data_type:   "text",
				Is_nullable: "YES",
			},
		},
		Schema{
			"id": Property{
				DataType: "string",
			},
		},
	},
	{
		//One property - required
		dbInfo{
			dbInfoRow{
				Column_name: "id",
				Data_type:   "text",
				Is_nullable: "NO",
			},
		},
		Schema{
			"id": Property{
				DataType: "string",
				Required: true,
			},
		},
	},
	{
		//Multiple properties - optional
		dbInfo{
			dbInfoRow{
				Column_name: "id",
				Data_type:   "text",
				Is_nullable: "YES",
			},
			dbInfoRow{
				Column_name: "aNumber",
				Data_type:   "numeric",
				Is_nullable: "YES",
			},
			dbInfoRow{
				Column_name: "anArray",
				Data_type:   "ARRAY",
				Is_nullable: "YES",
			},
		},
		Schema{
			"id": Property{
				DataType: "string",
				Required: false,
			},
			"aNumber": Property{
				DataType: "number",
				Required: false,
			},
			"anArray": Property{
				DataType: "array",
				Required: false,
			},
		},
	},
	{
		//Multiple properties - mixed
		dbInfo{
			dbInfoRow{
				Column_name: "id",
				Data_type:   "text",
				Is_nullable: "YES",
			},
			dbInfoRow{
				Column_name: "aNumber",
				Data_type:   "numeric",
				Is_nullable: "NO",
			},
			dbInfoRow{
				Column_name: "anArray",
				Data_type:   "ARRAY",
				Is_nullable: "YES",
			},
			dbInfoRow{
				Column_name: "aBoolean",
				Data_type:   "boolean",
				Is_nullable: "YES",
			},
			dbInfoRow{
				Column_name: "aTimestamp",
				Data_type:   "timestamp without time zone",
				Is_nullable: "YES",
			},
		},
		Schema{
			"id": Property{
				DataType: "string",
				Required: false,
			},
			"aNumber": Property{
				DataType: "number",
				Required: true,
			},
			"anArray": Property{
				DataType: "array",
				Required: false,
			},
			"aBoolean": Property{
				DataType: "boolean",
				Required: false,
			},
			"aTimestamp": Property{
				DataType: "date",
				Required: false,
			},
		},
	},
	{
		//Multiple properties - nullables YES scanned - default text
		dbInfo{
			dbInfoRow{
				Column_name:              "id",
				Data_type:                "text",
				Is_nullable:              "YES",
				Column_default:           sql.NullString{"'this is the default text'::text", true},
				Character_Maximum_Length: sql.NullInt64{10, true},
				Foreign_Table_Name:       sql.NullString{"reftable", true},
				Foreign_Column_Name:      sql.NullString{"refcol", true},
			},
		},
		Schema{
			"id": Property{
				DataType:         "string",
				Required:         false,
				Default:          "this is the default text",
				MaxLength:        10,
				ReferencesTable:  "reftable",
				ReferencesColumn: "refcol",
			},
		},
	},
	{
		//Multiple properties - nullables NOT scanned
		dbInfo{
			dbInfoRow{
				Column_name:              "id",
				Data_type:                "text",
				Is_nullable:              "YES",
				Character_Maximum_Length: sql.NullInt64{10, false},
				Foreign_Table_Name:       sql.NullString{"reftable", false},
				Foreign_Column_Name:      sql.NullString{"refcol", false},
			},
		},
		Schema{
			"id": Property{
				DataType: "string",
				Required: false,
			},
		},
	},
	{
		//Multiple properties - varying defaults
		dbInfo{
			dbInfoRow{
				Column_name:    "id",
				Data_type:      "numeric",
				Is_nullable:    "YES",
				Column_default: sql.NullString{"10", true},
			},
		},
		Schema{
			"id": Property{
				DataType: "number",
				Required: false,
				Default:  10,
			},
		},
	},
	{
		//Multiple properties - varying defaults
		dbInfo{
			dbInfoRow{
				Column_name:    "id",
				Data_type:      "boolean",
				Is_nullable:    "YES",
				Column_default: sql.NullString{"true", true},
			},
		},
		Schema{
			"id": Property{
				DataType: "boolean",
				Required: false,
				Default:  true,
			},
		},
	},
	{
		//Multiple properties - varying defaults
		dbInfo{
			dbInfoRow{
				Column_name:    "id",
				Data_type:      "numeric",
				Is_nullable:    "YES",
				Column_default: sql.NullString{"123.45", true},
			},
		},
		Schema{
			"id": Property{
				DataType: "number",
				Required: false,
				Default:  123.45,
			},
		},
	},
	{
		//Multiple properties - unrecognised default
		dbInfo{
			dbInfoRow{
				Column_name:    "id",
				Data_type:      "text",
				Is_nullable:    "YES",
				Column_default: sql.NullString{"'this is the default text'::nottext", true}, //not recognised as text
			},
		},
		Schema{
			"id": Property{
				DataType: "string",
				Required: false,
			},
		},
	},
}
