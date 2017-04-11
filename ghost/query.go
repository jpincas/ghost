package ghost

import "fmt"

//Query is the basic building block of an SQL query
type Query struct {
	//userQueryString is for when you need to provide complete, preformed SQL
	//This option will override any BaseAQL + args
	OverrideQueryString string
	//BaseSQL is a formatted SQL string with placeholder for SQLArgs
	BaseSQL string
	//SQLArgs are inserted into the BaseSQL in the order they appear
	SQLArgs []interface{}
	//WhereAnyOfValues appends a WHERE clause to match multiple values
	//If 'WherAnyOfKey' is not provided, it will default to 'id'
	WhereAnyOfValues []interface{}
	//The fieldname for the multiple-matching WHERE clause
	WhereAnyOfKey string
	//Indicate whether to rquest JSON array or object
	//and when unmarshalling, whether map or slice of maps
	IsList bool
	//Role to execute the query as
	Role string
	//UserID to set on the query context
	UserID string
	//CacheLevel specifies the level of caching to use: all, role, user
	//Omitting, or using any other value, will bypass caching
	CacheLevel string
	//cacheKey is the key used to store the SQL query in the cache
	cacheKey string
	//queryString is the output sql string ready to be executed
	queryString string
}

//Build runs a query against the data store and returns JSON
func (q *Query) Build() error {

	//If any override SQL is present, set the output query to its value
	//and exit immediately
	if q.OverrideQueryString != "" {
		q.queryString = q.OverrideQueryString
		return nil
	}

	tempQueryString := queryBuilder(fmt.Sprintf(q.BaseSQL, q.SQLArgs...))

	//For a multiple whereanyof
	if len(q.WhereAnyOfValues) != 0 {

		//Default to id
		if q.WhereAnyOfKey == "" {
			q.WhereAnyOfKey = "id"
		}

		tempQueryString = tempQueryString.addWhereAnyOfValues(q.WhereAnyOfKey, q.WhereAnyOfValues)
	}

	//Return JSON array or object
	if q.IsList {
		tempQueryString = tempQueryString.requestMultipleResultsAsJSONArray()
	} else {
		tempQueryString = tempQueryString.requestSingleResultAsJSONObject()
	}

	//Add sql role and user id, caching as you go depending
	//on the cache level specified
	if q.CacheLevel == "all" {
		q.cacheKey = tempQueryString.toSQLCacheKey()
	}

	//Add role
	if q.Role != "" {
		tempQueryString = tempQueryString.setQueryRole(q.Role)
	}

	if q.CacheLevel == "role" {
		q.cacheKey = tempQueryString.toSQLCacheKey()
	}

	//Add user id
	if q.UserID != "" {
		tempQueryString = tempQueryString.setUserID(q.UserID)
	}

	if q.CacheLevel == "user" {
		q.cacheKey = tempQueryString.toSQLCacheKey()
	}

	//Transform to SQL string and reset on the struct
	q.queryString = tempQueryString.toSQLString()

	//Execute the query
	return nil

}

//TODO: deprecate this
//QueryBuilder builds an SqlQuery from multiple URL query paramaters
// func QueryBuilder(schema string, table string, queries url.Values) SqlQuery {

// 	//Concat - schema.table
// 	tn := fmt.Sprintf("%s.%s", schema, table)

// 	//Start with all the products from the table
// 	p := sq.Select("*").From(tn)

// 	//Loop through all the URL qeueries
// 	for key, value := range queries {

// 		if strings.ToLower(key) == "orderby" {
// 			p = p.OrderBy(value[0])
// 		} else if strings.ToLower(key) == "limit" {
// 			l, _ := strconv.ParseUint(value[0], 10, 64)
// 			p = p.Limit(l)
// 		} else {
// 			p = p.Where(fmt.Sprintf(`%s %s`, key, value[0]))
// 		}
// 	}

// 	//Build the basic SQL
// 	sql, _, _ := p.ToSql()
// 	//Return as JSON array request
// 	return SqlQuery(sql)

// }
