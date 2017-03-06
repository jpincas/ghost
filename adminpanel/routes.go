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
	"net/http"
	"path"

	"github.com/ecosystemsoftware/ecosystem/core"
	"github.com/pressly/chi"
	"github.com/spf13/viper"
)

//These are functions that support the generation of custom admin panel(s)
//We do not actually serve the admin panel from here

func init() {

	core.Router.Route("/admin-panel", func(r chi.Router) {
		r.Get("/polymer-imports.html", GetPolymerImports)
		r.Get("/views", ShowConcatenatedJSON) //Serve all bundle view.json files combined
		r.Get("/menu", ShowConcatenatedJSON)  //Serve all bundle menu.json files combined
	})

	//For each bundle installed - add that bundle's public directory contents at TOPLEVEL/admin-panel-config/BUNDLENAME
	//This is where each bindles .json view files will be served from
	bundles := viper.GetStringSlice("bundlesInstalled")
	for _, v := range bundles {
		core.Router.FileServer(fmt.Sprintf("/admin-panel/%s", v), http.Dir(path.Join("bundles", v, "admin-panel")))
	}

}

//Notes

//-  The /polymer-imports.html route is probably too specific to have in this general module, and once we have
//established whether there will be a single, standard admin panel or multiple admin panels, we will know
//whether to break this out into a seperate package (i.e. 'polymer-admin-panel')
