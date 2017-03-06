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

package adminpanel

import (
	"fmt"
	"html/template"
	"net/http"
	"path"
	"strings"

	"github.com/ecosystemsoftware/ecosystem/core"
	"github.com/spf13/afero"
	"github.com/spf13/viper"
)

//ShowConcatenatedJSON 'concatenate's the contents of the [filename].json file in each bundle's admin-panel folder and returns the combined JSON
//where [filename] is the first slug in the URL.
//i.e. a route /views using this handler would concatenate all views.json in bundles
//By 'concatenation' in this context, we mean that it creates a .json file with a top-level object
//with a key for each bundle.
func ShowConcatenatedJSON(w http.ResponseWriter, r *http.Request) {

	//Work out whether this is views or menus etc
	url := r.URL.RequestURI()

	urlParts := strings.Split(url, "/")
	stub := urlParts[2]

	//For each bundle present
	bundles := viper.GetStringSlice("bundlesInstalled")

	var compositeFileContents string

	for _, v := range bundles {

		//Work out the name of the file
		viewsFile := path.Join("bundles", v, "admin-panel", stub+".json")

		//Check it exists
		ok, err := afero.Exists(core.AppFs, viewsFile)
		//If it exists, try to read it
		if ok && err == nil {
			viewsFileContents, err := afero.ReadFile(core.AppFs, viewsFile)
			//If it was read correctly
			if err == nil {

				//Prefix with the bundle name
				viewsFileString := fmt.Sprintf(`"%s":%s`, v, string(viewsFileContents))

				//If this is the first one, don't insert comma, otherwise do
				if compositeFileContents != "" {
					compositeFileContents += ","
				}
				compositeFileContents = fmt.Sprintf(`%s%s`, compositeFileContents, viewsFileString)
			}
		}

	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Write([]byte(fmt.Sprintf(`{%s}`, compositeFileContents)))

}

//GetPolymerImports is specific to the Polymer admin panel implementation
//It dynamically serves an imports.html file constructed from the current config
//and the polymeradmintemplate.go template
func GetPolymerImports(w http.ResponseWriter, r *http.Request) {
	var cf map[string]string
	viper.Unmarshal(&cf)

	bundles := viper.GetStringSlice("bundlesInstalled")

	//TODO: template will try to import actions.html from each bundle, even if it doesn't exist
	html := template.Must(template.New("polymeradminimports.html").Parse(PolymerAdminImportTemplate))
	html.ExecuteTemplate(w, "polymeradminimports.html", map[string]interface{}{
		"config":  cf,
		"bundles": bundles,
	})
}
