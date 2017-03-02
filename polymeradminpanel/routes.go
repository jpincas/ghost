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

package polymeradminpanel

import (
	"fmt"
	"net/http"
	"path"

	"github.com/ecosystemsoftware/ecosystem/core"
	"github.com/pressly/chi"
	"github.com/spf13/afero"
)

func main() {

	//These are functions that support the generation of custom admin panel(s)

	core.Router.Route("/admin", func(r chi.Router) {

		core.Router.Route("/polymer", func(r chi.Router) {
			core.Router.Get("/imports.html", AdminGetImports)
			core.Router.Get("/views", AdminShowConcatenatedJSON)
			core.Router.Get("/menu", AdminShowConcatenatedJSON)

			core.Router.Route("/bundles", func(r chi.Router) {
				//For each bundle present - add that bundle's admin directory contents at TOPLEVEL/custom/BUNDLENAME
				if bundleDirectoryContents, err := afero.ReadDir(core.AppFs, "bundles"); err == nil {
					for _, v := range bundleDirectoryContents {
						if v.IsDir() {
							core.Router.Get(fmt.Sprintf("/%s", v.Name()), http.FileServer(http.Dir(path.Join("bundles", v.Name(), "admin-panel", "polymer"))).(http.HandlerFunc))
						}
					}
				}
			})

		})

	})

}
