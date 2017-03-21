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

package main

import (
	"github.com/ecosystemsoftware/ecosystem/auth"
	"github.com/ecosystemsoftware/ecosystem/core"
	"github.com/ecosystemsoftware/ecosystem/email"
	"github.com/ecosystemsoftware/ecosystem/graphql"
	"github.com/ecosystemsoftware/ecosystem/rest"
)

func main() {

	//Tell EcoSystem which packages to activate
	core.ActivatePackages = activatePackages

	//Bootstrap the application
	core.RootCmd.Execute()
}

func activatePackages() {
	//Standard packages
	rest.Activate()
	auth.Activate()
	email.Activate()
	graphql.Activate()
}
