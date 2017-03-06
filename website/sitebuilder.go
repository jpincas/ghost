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

package website

import (
	"encoding/json"
	"fmt"

	"github.com/ecosystemsoftware/ecosystem/core"
)

type category struct {
	ID, Title, Image, Subtitle, Parent string
	Priority                           int
}

//menuItem is the basic building block for a menu
type menuItem struct {
	ID, Title, Image, Subtitle string
	SubItems                   menu
}

//menu is a list of menuItems
type menu []menuItem

//SiteBuilder provides sitewide data and logic.
//Intended use is to create an instance in handlers and pass to templates as 'site' key on the input struct
type SiteBuilder struct{}

//BuildMenu can be used by templates to retrieve a simple menu, listing web categories
func (s SiteBuilder) BuildMenu(startNode string, schema string) menu {
	var js, sql string

	if startNode == "" {
		sql = core.SqlQuery(fmt.Sprintf(core.SQLToGetAllWebCategories, schema)).RequestMultipleResultsAsJSONArray().SetQueryRole("web").ToSQLString()
	} else {
		sql = core.SqlQuery(fmt.Sprintf(core.SQLToGetWebCategoriesByParent, startNode, schema)).RequestMultipleResultsAsJSONArray().SetQueryRole("web").ToSQLString()
	}

	core.DB.QueryRow(sql).Scan(&js)
	var c []category
	json.Unmarshal([]byte(js), &c)

	var m menu

	//Iteration for LEVEL 1
	for _, v := range c {

		//If there is no parent, then it is a level 1 item
		if v.Parent == startNode {
			level1Item := menuItem{
				ID:       v.ID,
				Title:    v.Title,
				Image:    v.Image,
				Subtitle: v.Subtitle,
			}

			//Iteration for LEVEL 2
			for _, v2 := range c {

				//Identify categories with this category as a parent
				if level1Item.ID == v2.Parent {
					level2Item := menuItem{
						ID:       v2.ID,
						Title:    v2.Title,
						Image:    v2.Image,
						Subtitle: v2.Subtitle,
					}

					//Iteration for LEVEL 3
					for _, v3 := range c {

						//Identify categories with this category as a parent
						if level2Item.ID == v3.Parent {
							level3Item := menuItem{
								ID:       v3.ID,
								Title:    v3.Title,
								Image:    v3.Image,
								Subtitle: v3.Subtitle,
							}

							//Append the level 3 item to the level 2 item
							level2Item.SubItems = append(level2Item.SubItems, level3Item)

						} //END the level 3 IF

					} //END the level 3 iteration

					//Append the level 2 item to the level 1 item
					level1Item.SubItems = append(level1Item.SubItems, level2Item)

				} //END the level 2 IF

			} //END the level 2 iteration

			//Append the level 1 item to the main menu
			m = append(m, level1Item)

		} //END the level 1 IF

	} //END the level 1 iteration

	return m
}
