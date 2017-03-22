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

package graphql

import (
	"net/http"

	"github.com/ecosystemsoftware/ecosystem/core"
	"github.com/spf13/viper"
)

var cf map[string]string

func showGraphiqlIndex(w http.ResponseWriter, r *http.Request) {

	viper.Unmarshal(&cf)

	w.Header().Set("Content-Type", core.ContentTypeHTML)
	templates.ExecuteTemplate(w, "graphiql_index.html", cf)
	return

}

func showGraphiqlJS(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", core.ContentTypeJS)
	templates.ExecuteTemplate(w, "graphiql.js", cf)
	return

}

func showGraphiqlCSS(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", core.ContentTypeCSS)
	templates.ExecuteTemplate(w, "graphiql.css", cf)
	return

}

//ShowList shows a list of records from the database
func serveQuery(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Write([]byte("GraphQL connected"))

}
