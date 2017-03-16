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

import (
	"bufio"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/lib/pq"
)

const (
	ContentTypeHTML = `text/html; charset=utf-8`
	ContentTypeJSON = `application/json; charset=utf-8`
)

//ResponseError is the struct containing details of a server error
type ResponseError struct {
	HTTPCode     int          `json:"httpCode"`
	DBErrorCode  pq.ErrorCode `json:"dbCode"`
	ErrorMessage string       `json:"message"`
	Schema       string       `json:"schema"`
	Table        string       `json:"table"`
	Record       string       `json:"record"`
}

//AllOK Takes any number of bools and returns true if all are true, or false if ANY are false
func AllOK(oks ...bool) bool {
	for _, ok := range oks {
		if !ok {
			return false
		}
	}
	return true
}

//DBErrorCodeToHTTPErrorCode is a helper to translate error codes from the database into meaningful HTTP codes
func DBErrorCodeToHTTPErrorCode(dbCode pq.ErrorCode) (httpCode int) {
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

//RandomString generates a random string of int length
func RandomString(strlen int) string {
	rand.Seed(time.Now().UTC().UnixNano())
	const chars = "abcdefghijklmnopqrstuvwxyz0123456789"
	result := make([]byte, strlen)
	for i := 0; i < strlen; i++ {
		result[i] = chars[rand.Intn(len(chars))]
	}
	return string(result)
}

//CheckTemplate looks for a template corresponding to the
// func CheckTemplate(t []string, schema string, table string, listOrSingle string) (bool, string) {

// 	defaultName := listOrSingle + ".html"
// 	templateName := table + "-" + defaultName

// 	//Check for templates following a defined order of priority.  Once one is found,
// 	//return true and the name of the template

// 	//1) BUNDLE/pages/TABLE-[list]or[single].HTML	Top priority is a named template in the bundle folder
// 	if matches, _ := filepath.Glob(path.Join("templates", schema, "pages", templateName)); matches != nil {
// 		return true, path.Join(schema, templateName)
// 	}
// 	//2) */pages/TABLE-[list]or[single].HTML	Second priority is a named template in ANY folder at the bundle level
// 	//This allows a user to define a folder say "my-templates" which would not be removed in the case
// 	//of bundle removal
// 	if matches, _ := filepath.Glob(path.Join("templates", "*", "pages", templateName)); matches != nil {
// 		return true, templateName
// 	}
// 	//3) BUNDLE/pages/[list]or[single].HTML Third priority is a generic 'list.html' or 'single.html' in the bundle folder
// 	if matches, _ := filepath.Glob(path.Join("templates", schema, "pages", defaultName)); matches != nil {
// 		return true, path.Join(schema, defaultName)
// 	}
// 	//4) */pages/[list]or[single].HTML Fourth priority is a generic 'list.html' or 'single.html' in ANY folder at the bundle level
// 	//This allows a user to define a folder say "my-templates" which would not be removed in the case
// 	//of bundle removal
// 	if matches, _ := filepath.Glob(path.Join("templates", "*", "pages", defaultName)); matches != nil {
// 		return true, defaultName
// 	}

// 	//If there is no usable template, return false
// 	return false, ""

// }

//HyphensToUnderscores replaces all hyphens in the string with underscores.
//This is so you can use pretty URLs with hyphens (as recommended by Google)
//whilst still using underscores in DB table names - which means they don't have to be quoted all the time
//https://support.google.com/webmasters/answer/76329?hl=en
func HyphensToUnderscores(table string) string {
	return strings.Replace(table, "-", "_", -1)
}

//MapToValsAndCols iterates over the map resulting from binding a JSON request body
//and creates cols and vals strings to be used in SQL query
func MapToValsAndCols(r map[string]interface{}) (cols, vals string) {
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

// AskForConfirmation asks the user for confirmation. A user must type in "yes" or "no" and
// then press enter. It has fuzzy matching, so "y", "Y", "yes", "YES", and "Yes" all count as
// confirmations. If the input is not recognized, it will ask again. The function does not return
// until it gets a valid response from the user.
func AskForConfirmation(s string) bool {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Printf("%s [y/n]: ", s)

		response, err := reader.ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}

		response = strings.ToLower(strings.TrimSpace(response))

		if response == "y" || response == "yes" {
			return true
		} else if response == "n" || response == "no" {
			return false
		}
	}
}
