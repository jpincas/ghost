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

const (
	sqlToGetTablesInSchema = `SELECT table_name, table_type FROM information_schema.tables WHERE table_schema='%s';`
	sqlToGetSchemasInDB    = `SELECT schema_name FROM information_schema.schemata WHERE schema_owner != 'postgres';`
)
