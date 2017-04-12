package ghost

//SQl query strings for application-wide use
const (
	SQLToFindUserByEmail = `SELECT id from users WHERE email = '%s';`
	SQLToGetUsersRole    = `SELECT role from users WHERE id = '%s';`

	//General
	//NO SEMI COLONS AT THE END
	SQLToSelectAllFieldsFrom = `SELECT * FROM %s.%s`
	SQLToSelectWhere         = `SELECT * FROM %s.%s WHERE id = '%v'` //depracated
	SQLToSelectByID          = `SELECT * FROM %s.%s WHERE id = '%v'`
	SQLToSelectWhereXEqualsY = `SELECT * FROM %s.%s WHERE %s = '%v'`

	SQLToInsertReturningJSON            = `INSERT INTO %s.%s(%s) VALUES (%s) returning row_to_json(%s)`
	SQLToInsertAllDefaultsReturningJSON = `INSERT INTO %s.%s DEFAULT VALUES returning row_to_json(%s)`
	SQLToDeleteWhere                    = `DELETE FROM %s.%s WHERE id = '%v'`
	SQLToUpdateWhereReturningJSON       = `UPDATE %s.%s SET (%s) = (%s) WHERE id = '%v' returning row_to_json(%s)`

	//Full text search_path
	SQLToFullTextSearch = `with item as (select to_tsvector(%s::text) @@ to_tsquery('%s') AS found, %s.* FROM %s.%s) select array_to_json(array_agg(row_to_json(item))) FROM item WHERE item.found = TRUE`
)
