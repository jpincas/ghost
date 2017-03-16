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

//SQl query strings for application-wide use
const (

	//Install extensions
	SQLToCreateUUIDExtension = `CREATE EXTENSION IF NOT EXISTS "uuid-ossp";`

	//Create built-in tables
	SQLToCreateUsersTable = `CREATE TABLE users (id uuid PRIMARY KEY, email varchar(256) UNIQUE, role varchar(16) NOT NULL default 'anon');`

	//Built in user logic and cmd line user creation
	SQLToCreateFuncToGenerateNewUserID = `CREATE FUNCTION generate_new_user() RETURNS trigger AS $$ BEGIN NEW.id := uuid_generate_v4(); RETURN NEW; END; $$ LANGUAGE plpgsql;`
	SQLToCreateTriggerOnNewUserInsert  = `CREATE TRIGGER new_user BEFORE INSERT ON users FOR EACH ROW EXECUTE PROCEDURE generate_new_user();`
	SQLToCreateAdministrator           = `INSERT INTO users(email, role) VALUES ('%s', '%s');`
	SQLToFindUserByEmail               = `SELECT id from users WHERE email = '%s';`
	SQLToGetUsersRole                  = `SELECT role from users WHERE id = '%s';`

	//Built in roles
	SQLToCreateServerRole        = `CREATE ROLE server NOINHERIT LOGIN PASSWORD NULL;`
	SQLToSetServerRolePassword   = `ALTER ROLE server NOINHERIT LOGIN PASSWORD '%s' VALID UNTIL 'infinity';`
	SQLToCreateAnonRole          = `CREATE ROLE anon;`
	SQLToCreateAdminRole         = `CREATE ROLE admin BYPASSRLS;`
	SQLToGrantBuiltInPermissions = `GRANT anon, admin TO server; GRANT SELECT ON TABLE users TO server;`

	//Admin permissions
	SQLToGrantAdminPermissions = `ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT ALL ON TABLES TO admin; ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT USAGE ON SEQUENCES TO admin;`

	//Schema manipulation for bundles
	SQLToCreateSchema                = `CREATE SCHEMA %s;`
	SQLToGrantBundleAdminPermissions = `ALTER DEFAULT PRIVILEGES IN SCHEMA %s GRANT ALL ON TABLES TO admin; ALTER DEFAULT PRIVILEGES IN SCHEMA %s GRANT USAGE ON SEQUENCES TO admin;`
	SQLToDropSchema                  = `DROP SCHEMA %s CASCADE;`
	SQLToSetSearchPathForBundle      = `SET search_path TO %s, public;`

	//General
	//NO SEMI COLONS AT THE END
	SQLToSelectWhere                    = `SELECT * FROM %s.%s WHERE id = '%v'`
	SQLToInsertReturningJSON            = `INSERT INTO %s.%s(%s) VALUES (%s) returning row_to_json(%s)`
	SQLToInsertAllDefaultsReturningJSON = `INSERT INTO %s.%s DEFAULT VALUES returning row_to_json(%s)`
	SQLToDeleteWhere                    = `DELETE FROM %s.%s WHERE id = '%v'`
	SQLToUpdateWhereReturningJSON       = `UPDATE %s.%s SET (%s) = (%s) WHERE id = '%v' returning row_to_json(%s)`

	//JSON Conversion
	SQLToRequestMultipleResultsAsJSONArray = `WITH results AS (%s) SELECT array_to_json(array_agg(row_to_json(results))) from results;`
	SQLToRequestSingleResultAsJSONObject   = `WITH results AS (%s) SELECT row_to_json(results) from results;`

	//Setting local role and user id
	SQLToSetLocalRole = `SET LOCAL ROLE %s; %s `
	SQLToSetUserID    = `SET my.user_id = '%s'; %s `

	//Full text search_path
	SQLToFullTextSearch = `with item as (select to_tsvector(%s::text) @@ to_tsquery('%s') AS found, %s.* FROM %s.%s) select array_to_json(array_agg(row_to_json(item))) FROM item WHERE item.found = TRUE`
)
