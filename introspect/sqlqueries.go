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

const (
	sqlToGetTablesInSchema = `SELECT table_name, table_type FROM information_schema.tables WHERE table_schema='%s';`
	sqlToGetSchemasInDB    = `SELECT schema_name FROM information_schema.schemata WHERE schema_owner != 'postgres';`
	sqlToGetTableInfo      = `SET LOCAL ROLE %s; WITH cs AS (
	SELECT
	tc.constraint_name, tc.table_name, kcu.column_name,
	ccu.table_name AS foreign_table_name,
	ccu.column_name AS foreign_column_name
	FROM
	information_schema.table_constraints AS tc
	JOIN information_schema.key_column_usage AS kcu
	  ON tc.constraint_name = kcu.constraint_name
	JOIN information_schema.constraint_column_usage AS ccu
	  ON ccu.constraint_name = tc.constraint_name
	WHERE constraint_type = 'FOREIGN KEY' AND tc.table_name='%s'
)
SELECT i.column_name, i.data_type, i.is_nullable, i.column_default, i.character_maximum_length, cs.foreign_table_name, cs.foreign_column_name
FROM INFORMATION_SCHEMA.COLUMNS i
LEFT JOIN cs as cs
ON cs.column_name = i.column_name
WHERE i.table_schema = '%s' AND i.table_name = '%s';`
)
