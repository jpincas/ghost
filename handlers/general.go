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

package handlers

import (
	"fmt"
	"net/http"
	"path"

	eco "github.com/ecosystemsoftware/ecosystem/utilities"
	"github.com/spf13/afero"
	gin "gopkg.in/gin-gonic/gin.v1"
)

//Image Display Handler
func ShowImage(c *gin.Context) {
	imageToServe, err := eco.GetImage(c.Param("image"), c.DefaultQuery("width", "500"))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": err.Error(),
		})
	} else {
		c.File(imageToServe)
	}
}

//Important for correct custom CORS request.
func OptionsHandler(c *gin.Context) {
	c.Header("Access-Control-Allow-Methods", "POST, GET, PUT, PATCH, DELETE, OPTIONS, SEARCH")
	c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")
	c.String(http.StatusOK, "")
}

func ReturnBlank(c *gin.Context) {
	c.String(http.StatusOK, "")
}

//ReturnHelloWorld is a test handler that can be used when wiring up a custom server to
//check that the EcoSystem utility package is being correctly imported and built
func ReturnHelloWorld(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"hello": "world",
	})
}

//AdminShowViews concatenates the contents of the views.json file in each bundle and returns the combined JSON
func AdminShowViews(c *gin.Context) {
	//For each bundle present
	if bundleDirectoryContents, err := afero.ReadDir(eco.AppFs, "bundles"); err == nil {

		var compositeFileContents string

		for _, v := range bundleDirectoryContents {
			if v.IsDir() {
				//Work out the name of the views file
				viewsFile := path.Join("bundles", v.Name(), "admin-panel", "views.json")
				//Check it exists
				ok, err := afero.Exists(eco.AppFs, viewsFile)
				//If it exists, try to read it
				if ok && err == nil {
					viewsFileContents, err := afero.ReadFile(eco.AppFs, viewsFile)
					//If it was read correctly
					if err == nil {

						//Remove surrounding brackets
						//viewsFileString := strings.TrimSuffix(strings.TrimPrefix(string(viewsFileContents), "{"), "}")
						//Prefix with the bundle name
						viewsFileString := fmt.Sprintf(`"%s":%s`, v.Name(), string(viewsFileContents))

						//If this is the first one, don't insert comma, otherwise do
						if compositeFileContents != "" {
							compositeFileContents += ","
						}
						compositeFileContents = fmt.Sprintf(`%s%s`, compositeFileContents, viewsFileString)
					}
				}
			}
		}

		c.String(http.StatusOK, fmt.Sprintf(`{%s}`, compositeFileContents))

	}
}
