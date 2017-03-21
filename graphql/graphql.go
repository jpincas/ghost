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

//go:generate hardcodetemplates -p=graphql

import (
	"html/template"

	"github.com/ecosystemsoftware/ecosystem/core"
	"github.com/graphql-go/graphql"
)

//Shared templates holder
var templates *template.Template

//Main handler
var schema graphql.Schema

//Activate is the main package activation function
func Activate() error {

	//Log activation message
	core.Log(core.LogEntry{"GRAPHQL", true, "Initialising GraphQL package..."})

	//Parse the templates
	parseTemplates()

	//Initialise the schema
	initSchema()

	//Set the routes for the package
	setRoutes()
	return nil
}

func parseTemplates() {

	templates = template.Must(template.New("base").Parse(baseTemplate))
	core.Log(core.LogEntry{"GRAPHQL", true, templates.DefinedTemplates()})

}

//initSchema sets up the root query and main schema
func initSchema() {

	//Database schemas form the base level root query
	rootQuery := graphql.ObjectConfig{Name: "RootQuery", Fields: schemaRootSetup()}
	schemaConfig := graphql.SchemaConfig{Query: graphql.NewObject(rootQuery)}

	//Set the schema
	var err error
	schema, err = graphql.NewSchema(schemaConfig)
	if err != nil {
		core.LogFatal(core.LogEntry{"GRAPHQL", false, err.Error()})
	}

	return

}
