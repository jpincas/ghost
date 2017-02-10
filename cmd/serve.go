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

package cmd

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
)

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Starts the EcoSystem server",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {

		// Static Server for Admin Panel
		adminServer := gin.Default()
		adminServer.StaticFS("/", http.Dir("ecosystem-admin/build/unbundled/"))
		go adminServer.Run(":8080")

	},
}

func init() {
	RootCmd.AddCommand(serveCmd)

	// viper.SetConfigName("config")
	// viper.AddConfigPath(".")

	// viper.SetDefault("pgSuperUser", "postgres")
	// viper.SetDefault("pgName", "ecosystem")
	// viper.SetDefault("pgPort", "5432")
	// viper.SetDefault("pgServer", "localhost")
	// viper.SetDefault("pgDisableSSL", false)

	// err := viper.ReadInConfig() // Find and read the config file
	// if err != nil {             // Handle errors reading the config file
	// 	log.Fatal("Could not read config.json: ", err.Error())
	// }

}
