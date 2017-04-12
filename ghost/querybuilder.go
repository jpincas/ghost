package ghost

import (
	"fmt"
	"reflect"
	"strings"
)

const (

	//JSON Conversion
	sqlToRequestMultipleResultsAsJSONArray = `WITH results AS (%s) SELECT array_to_json(array_agg(row_to_json(results))) from results;`
	sqlToRequestSingleResultAsJSONObject   = `WITH results AS (%s) SELECT row_to_json(results) from results;`

	//Setting local role and user id
	sqlToSetLocalRole = `SET LOCAL ROLE %s; %s`
	sqlToSetUserID    = `SET my.user_id = '%s'; %s`

	//Basics
	sqlToSelectFieldsFromTableSchema = `SELECT %s FROM %s.%s`

	//Where clauses
	sqlToAddFirstWhereClause          = `%s WHERE %s %s '%s'` //safe to escape here
	sqlToAddFirstWhereAnyClause       = `%s WHERE %s = ANY(ARRAY%s)`
	sqlToAddSubsequentWhereClauses    = `%s %s %s %s '%s'` //safe to escape here
	sqlToAddSubsequentWhereAnyClauses = `%s %s %s = ANY(ARRAY%s)`
)

//queryBuilder is an SQL query string with various methods available for transformation
type queryBuilder string

//basicSelect is the simple type of base query
func (s queryBuilder) basicSelect(schema string, table string, selectFields []string) queryBuilder {

	return queryBuilder(fmt.Sprintf(sqlToSelectFieldsFromTableSchema, toListString(selectFields), schema, table))

}

//addWhere clauses appends multiple where clauses conjoined with AND or OR
func (s queryBuilder) addWhereClauses(whereClauses []WhereConfig) queryBuilder {

	for k, v := range whereClauses {

		//Default the key to id
		if v.Key == "" {
			v.Key = "id"
		}

		//For the first where clause
		if k == 0 {

			if len(v.AnyValue) != 0 {
				//For strings, must surround with ''
				//But for numbers, doing so causes an error
				//This is unlike regular behaviour (not in arrays),
				//where postgres CAN deal with numbers in ''
				valueType := reflect.TypeOf(v.AnyValue[0]).Name()
				s = queryBuilder(fmt.Sprintf(sqlToAddFirstWhereAnyClause, s, v.Key, toCsvSqlArrayString(v.AnyValue, valueType)))
			} else {
				s = queryBuilder(fmt.Sprintf(sqlToAddFirstWhereClause, s, v.Key, v.Operator, v.Value))
			}

		} else {

			conjunction := "AND"
			if v.JoinWithOr {
				conjunction = "OR"
			}

			if len(v.AnyValue) != 0 {
				valueType := reflect.TypeOf(v.AnyValue[0]).Name()
				s = queryBuilder(fmt.Sprintf(sqlToAddSubsequentWhereAnyClauses, s, conjunction, v.Key, toCsvSqlArrayString(v.AnyValue, valueType)))
			} else {
				s = queryBuilder(fmt.Sprintf(sqlToAddSubsequentWhereClauses, s, conjunction, v.Key, v.Operator, v.Value))
			}
		}

	}

	return s

}

//RequestMultipleResultsAsJSONArray transforms the SQL query to return a JSON array of results
//Use when multiple lines are going to be returned
func (s queryBuilder) requestMultipleResultsAsJSONArray() queryBuilder {

	return queryBuilder(fmt.Sprintf(sqlToRequestMultipleResultsAsJSONArray, s))

}

//RequestSingleResultAsJSONObject transforms the SQL query to return a JSON object of the result
//Used when a single line is going to be returned
func (s queryBuilder) requestSingleResultAsJSONObject() queryBuilder {

	return queryBuilder(fmt.Sprintf(sqlToRequestSingleResultAsJSONObject, s))

}

//SetQueryRole prepends the database role with which to execute the query
func (s queryBuilder) setQueryRole(role string) queryBuilder {

	return queryBuilder(fmt.Sprintf(sqlToSetLocalRole, role, s))

}

//SetUserID prepends the user id variable with which to execute the query
func (s queryBuilder) setUserID(userID string) queryBuilder {

	return queryBuilder(fmt.Sprintf(sqlToSetUserID, userID, s))

}

//ToSQLCacheKey transforms the SQL query into a cacheable string key
func (s queryBuilder) toSQLCacheKey() string {

	return fmt.Sprint(s)

}

//ToSQLString transforms an SqlQuery to a plain string
//Generally the last step before execution
func (s queryBuilder) toSQLString() string {

	LogDebug("SQL", true, fmt.Sprint(s), nil)
	return fmt.Sprint(s)
}

//Helpers
//toCsvSqlArrayString
func toCsvSqlArrayString(i []interface{}, valueType string) string {

	tempArrayString := "["

	//For strings, wrap in ''
	if valueType == "string" {

		for k, v := range i {
			if k == 0 {
				tempArrayString += fmt.Sprintf(`'%s'`, v)
			} else {
				tempArrayString += fmt.Sprintf(`, '%s'`, v)
			}
		}

		// For anything else other than string, dont wrap
	} else {

		for k, v := range i {
			if k == 0 {
				tempArrayString += fmt.Sprintf(`%v`, v)
			} else {
				tempArrayString += fmt.Sprintf(`, %v`, v)
			}
		}

	}

	tempArrayString += "]"

	return tempArrayString

}

//toListString removes the [] from a slice and returns the comma separated string
func toListString(l []string) string {

	return strings.Replace(strings.TrimPrefix(strings.TrimSuffix(fmt.Sprintf(`%s`, l), "]"), "["), " ", ",", -1)

}
