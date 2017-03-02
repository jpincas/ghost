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

package images

import (
	"net/http"

	"github.com/ecosystemsoftware/ecosystem/core"
	"github.com/pressly/chi/render"
)

//Image Display Handler
func ShowImage(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	image, _ := ctx.Value("image").(string)
	width, _ := ctx.Value("width").(string)

	imageToServe, err := GetImage(image, width)
	if err != nil {
		render.Status(r, http.StatusNotFound)
		render.JSON(w, r, core.ResponseError{http.StatusNotFound, "", err.Error(), "", "", ""})
	} else {
		http.ServeFile(w, r, imageToServe)
	}
}
