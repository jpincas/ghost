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

package utilities

//Config is the basic structure of the config.json file
type Config struct {
	PgSuperUser              string `json:"pgSuperUser"`
	PgDBName                 string `json:"pgDBName"`
	PgPort                   string `json:"pgPort"`
	PgServer                 string `json:"pgServer"`
	PgDisableSSL             bool   `json:"pgDisableSSL"`
	ApiPort                  string `json:"apiPort"`
	WebsitePort              string `json:"websitePort"`
	AdminPanelPort           string `json:"adminPanelPort"`
	AdminPanelServeDirectory string `json:"adminPanelServeDirectory"`
	PublicSiteSlug           string `json:"publicSiteSlug"`
	PrivateSiteSlug          string `json:"privateSiteSlug"`
	SmtpHost                 string `json:"smtpHost"`
	SmtpPort                 string `json:"smtpPort"`
	SmtpUserName             string `json:"smtpUserName"`
	SmtpFrom                 string `json:"smtpFrom"`
	EmailFrom                string `json:"emailFrom"`
	JWTRealm                 string `json:"jwtRealm"`
}
