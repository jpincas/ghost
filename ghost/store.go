package ghost

import (
	"encoding/json"
	"strings"
)

type store struct{}

func (s store) Execute(q *Query) (string, error) {

	if err := q.Build(); err != nil {
		return "", nil
	}

	//Caching case
	//Return the cached result if there is a cache key present
	//AND there is a result from the cache
	if q.cacheKey != "" {
		cacheResult, ok := App.Cache.Get(q.cacheKey)
		if ok {
			LogDebug("STORE", true, "Returning from cache, using key: "+q.cacheKey, nil)
			return cacheResult.(string), nil
		}
	}

	//No caching case
	var JSONResponse string
	if err := App.DB.QueryRow(q.queryString).Scan(&JSONResponse); err != nil {
		//Only one row is returned as JSON is returned by Postgres
		//Empty result
		if strings.Contains(err.Error(), "sql") {
			return "", nil
		}

		//Else its a database error
		return "", err

	}

	//Set the cache if a cache key has been provided
	if q.cacheKey != "" {
		LogDebug("STORE", true, "Caching result with key: "+q.cacheKey, nil)
		App.Cache.Set(q.cacheKey, JSONResponse)
	}
	return JSONResponse, nil

}

//ExecuteAndUnmarshall runs a query against the datastore and returns both
//for lists: []map[string]interfaace{}
//for objects: map[string]interface{}
//The corresponding unused data structure is set to nil
func (s store) ExecuteAndUnmarshall(q *Query) (list []map[string]interface{}, single map[string]interface{}, err error) {

	//Execute to JSON first
	var dbResponse string
	dbResponse, err = s.Execute(q)
	if err != nil {
		return nil, nil, err
	}

	if q.IsList {

		//Check for empty result from DB
		if dbResponse == "" {
			return list, nil, nil
		}

		if err := json.Unmarshal([]byte(dbResponse), &list); err != nil {
			return nil, nil, err
		}

		return list, nil, nil
	}

	//Check for empty result from DB
	if dbResponse == "" {
		return nil, single, nil
	}

	if err := json.Unmarshal([]byte(dbResponse), &single); err != nil {
		return nil, nil, err
	}

	return nil, single, nil

}
