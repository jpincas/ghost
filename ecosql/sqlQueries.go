// Copyright 2017 EcoSystem Software LLP

// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at

// 	http://www.apache.org/licenses/LICENSE-2.0

// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package ecosql

//SQl query strings for application-wide use
const (

	//Install extensions
	ToCreateUUIDExtension = `CREATE EXTENSION IF NOT EXISTS "uuid-ossp";`

	//Create built-in tables
	ToCreateUsersTable         = `CREATE TABLE users (id uuid PRIMARY KEY, email varchar(256) UNIQUE, role varchar(16) NOT NULL default 'anon');`
	ToCreateWebCategoriesTable = `CREATE TABLE web_categories (id text NOT NULL PRIMARY KEY, title text,image text,description text,subtitle text,parent text,priority integer);`

	//Built in user logic and cmd line user creation
	ToCreateFuncToGenerateNewUserID = `CREATE FUNCTION generate_new_user() RETURNS trigger AS $$ BEGIN NEW.id := uuid_generate_v4(); RETURN NEW; END; $$ LANGUAGE plpgsql;`
	ToCreateTriggerOnNewUserInsert  = `CREATE TRIGGER new_user BEFORE INSERT ON users FOR EACH ROW EXECUTE PROCEDURE generate_new_user();`
	ToCreateAdministrator           = `INSERT INTO users(email, role) VALUES ('%s', '%s');`
	ToFindUserByEmail               = `SELECT id from users WHERE email = '%s';`
	ToGetUsersRole                  = `SELECT role from users WHERE id = '%s';`

	//Built in roles
	ToCreateServerRole               = `CREATE ROLE server NOINHERIT LOGIN;`
	ToSetServerRolePassword          = `ALTER ROLE server WITH PASSWORD '%s';`
	ToSetServerPasswordToLastForever = `ALTER ROLE server VALID UNTIL 'infinity';`
	ToCreateAnonRole                 = `CREATE ROLE anon;`
	ToCreateAdminRole                = `CREATE ROLE admin BYPASSRLS;`
	ToCreateWebRole                  = `CREATE ROLE web;`
	ToGrantAdminPermissions          = `GRANT admin To server; GRANT ALL ON ALL TABLES IN SCHEMA public TO admin;GRANT USAGE ON ALL SEQUENCES IN SCHEMA public TO admin; `
	ToGrantBuiltInPermissions        = `GRANT anon, web TO server; GRANT SELECT ON TABLE users TO server;GRANT SELECT ON TABLE web_categories TO web;`

	//Schema manipulation for bundles
	ToCreateSchema           = `CREATE SCHEMA %s;`
	ToDropSchema             = `DROP SCHEMA %s CASCADE;`
	ToSetSearchPathForBundle = `SET search_path TO %s, public;`

	//Web category retrieval and info
	ToSelectWebCategoryWhere   = `SELECT * FROM web_categories WHERE id = '%s';`
	ToGetAllWebCategories      = `SELECT * FROM web_categories ORDER BY priority;`
	ToGetWebCategoriesByParent = `SELECT * FROM web_categories WHERE parent = '%s' ORDER BY priority;`
	ToSelectKeywordedRecords   = `SELECT * FROM %s WHERE keywords @> '{%s}';`

	//Web requests
	ToSelectRecordBySlug = `SELECT * FROM %s WHERE slug = '%s';`

	//General
	ToSelectWhere                    = `SELECT * FROM %s WHERE id = '%s';`
	ToInsertReturningJSON            = `INSERT INTO %s(%s) VALUES (%s) returning row_to_json(%s);`
	ToInsertAllDefaultsReturningJSON = `INSERT INTO %s DEFAULT VALUES returning row_to_json(%s);`
	ToDeleteWhere                    = `DELETE FROM %v WHERE id = '%v';`
	ToUpdateWhereReturningJSON       = `UPDATE %s SET (%s) = (%s) WHERE id = '%v' returning row_to_json(%s);`

	//Full text search_path
	ToFullTextSearch = `with item as (select to_tsvector(%s::text) @@ to_tsquery('%s') AS found, %s.* FROM %s) select array_to_json(array_agg(row_to_json(item))) FROM item WHERE item.found = TRUE OR item.id ILIKE '%s%%'`
)
