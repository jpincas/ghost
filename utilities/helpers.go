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

package utilities

import (
	"fmt"
	"net/http"
	"os"
	"path"
	"strings"

	"github.com/lib/pq"
)

//dbErrorCodeToHTTPErrorCode is a helper to translate error codes from the database into meaningful HTTP codes
func dbErrorCodeToHTTPErrorCode(dbCode pq.ErrorCode) (httpCode int) {
	switch {
	case dbCode == "42501":
		httpCode = http.StatusForbidden
	case dbCode == "42P01":
		httpCode = http.StatusNotFound
	default:
		//Default to 400
		httpCode = http.StatusBadRequest
	}
	return httpCode
}

//checkTemplate looks for a specific template, if not, falls back to the default
//and if there is not defaul, returns false
func checkTemplate(table string, listOrSingle string) (bool, string) {

	defaultName := listOrSingle + ".html"
	templateName := table + "-" + defaultName

	targetTemplateFullPath := path.Join("templates/custom", templateName)

	//First, we check for the existence of a specific template for this table
	if _, err := os.Stat(targetTemplateFullPath); os.IsNotExist(err) {
		//If there isn't one, revert to default
		targetTemplateFullPath = path.Join("templates/default", defaultName)
		if _, err := os.Stat(targetTemplateFullPath); os.IsNotExist(err) {
			//If there is no default, return false
			return false, ""
		}
		//If there is a default, return true and the default name
		return true, defaultName
	}
	//If there is a specific template
	return true, templateName
}

//hyphenToUnderscore replaces all hyphens in the string with underscores.
//This is so you can use pretty URLs with hyphens (as recommended by Google)
//whilst still using underscores in DB table names - which means they don't have to be quoted all the time
//https://support.google.com/webmasters/answer/76329?hl=en
func hyphensToUnderscores(table string) string {
	return strings.Replace(table, "-", "_", -1)
}

//mapToValsAndCols iterates over the map resulting from binding a JSON request body
//and creates cols and vals strings to be used in SQL query
func mapToValsAndCols(r map[string]interface{}) (cols, vals string) {
	//Iterate over the map and set the keys and values on the context
	for k, v := range r {
		cols = cols + k + ", "
		//The value can be anything, so we Sprintf to the default format %v to avoid any printing errors
		//We also surround with '' to make everything a string in the SQL, this way strings don't fail
		//and the database automatically types numebers.  If you don't do this, strings will fail.
		vals = vals + fmt.Sprintf(`'%v'`, v) + ", "
	}
	//Trim the trailing commas
	cols = strings.TrimSuffix(cols, ", ")
	vals = strings.TrimSuffix(vals, ", ")
	return
}
