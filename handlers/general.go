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
	"net/http"

	eco "github.com/ecosystemsoftware/ecosystem/utilities"

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
