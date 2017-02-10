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

//SQl query strings
const (
	ToCreateUUIDExtension            = `CREATE EXTENSION IF NOT EXISTS "uuid-ossp";`
	ToCreateUsersTable               = `CREATE TABLE IF NOT EXISTS users (id uuid PRIMARY KEY, email varchar(256) UNIQUE, role varchar(16) NOT NULL default 'anon');`
	ToCreateFuncToGenerateNewUserID  = `CREATE OR REPLACE FUNCTION generate_new_user() RETURNS trigger AS $$ BEGIN NEW.id := uuid_generate_v4(); RETURN NEW; END; $$ LANGUAGE plpgsql;`
	ToCreateTriggerOnNewUserInsert   = `CREATE TRIGGER new_user BEFORE INSERT ON users FOR EACH ROW EXECUTE PROCEDURE generate_new_user();`
	ToCreateAdministrator            = `INSERT INTO users(email, role) VALUES ('%s', '%s')`
	ToCreateWebCategoriesTable       = `CREATE TABLE IF NOT EXISTS web_categories (id text NOT NULL PRIMARY KEY, title text,image text,description text,subtitle text,parent text,priority integer);`
	ToCreateServerRole               = `CREATE ROLE server NOINHERIT LOGIN;`
	ToSetServerRolePassword          = `ALTER ROLE server WITH PASSWORD '%s';`
	ToSetServerPasswordToLastForever = `ALTER ROLE server VALID UNTIL 'infinity';`
	ToCreateAnonRole                 = `CREATE ROLE anon;`
	ToCreateAdminRole                = `CREATE ROLE admin BYPASSRLS;`
	ToCreateWebRole                  = `CREATE ROLE web;`
	ToGrantAdminPermissions          = `GRANT admin To server; GRANT ALL ON ALL TABLES IN SCHEMA public TO admin;GRANT USAGE ON ALL SEQUENCES IN SCHEMA public TO admin; `
	ToGrantBuiltInPermissions        = `GRANT anon, web TO server; GRANT SELECT ON TABLE users TO server;GRANT SELECT ON TABLE web_categories TO web;`
)
