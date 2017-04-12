package ghost

import "fmt"

//WhereConfig describes one or more where clauses
type WhereConfig struct {
	Key        string
	Operator   string
	Value      interface{}
	AnyValue   []interface{}
	JoinWithOr bool
}

//Query is the basic building block of an SQL query
type Query struct {
	//userQueryString is for when you need to provide complete, preformed SQL
	//This option will override any BaseAQL + args
	OverrideQueryString string
	//BaseSQL is a formatted SQL string with placeholder for SQLArgs
	BaseSQL string
	//SQLArgs are inserted into the BaseSQL in the order they appear
	SQLArgs []interface{}
	//SELECT fields
	Select []string
	//From Schema.Table
	Schema, Table string
	//Where
	Where []WhereConfig
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
	//CacheExpiry
	CacheExpiry int
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

	tempQueryString := queryBuilder("")

	//If base sql + args have been supplied, use them
	//otherwise build from parameters
	if q.BaseSQL != "" && len(q.SQLArgs) != 0 {
		tempQueryString = queryBuilder(fmt.Sprintf(q.BaseSQL, q.SQLArgs...))
	} else {
		tempQueryString = tempQueryString.basicSelect(q.Schema, q.Table, q.Select)
	}

	//For WHERE clauses
	if len(q.Where) != 0 {
		tempQueryString = tempQueryString.addWhereClauses(q.Where)
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
