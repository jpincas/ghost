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

// func TestCheckTemplate(t *testing.T) {

// 	//Use afero for test file system
// 	appFS := afero.NewMemMapFs()
// 	appFS.MkdirAll("templates/testbundle/pages", 0755)
// 	appFS.MkdirAll("templates/mytemplates/pages", 0755)

// 	var wantExists, gotExists bool
// 	var wantName, gotName string

// 	//DEFAULT CASE: no template
// 	wantExists = false
// 	gotExists, gotName = CheckTemplate("testbundle", "mytable", "single")

// 	if gotExists != wantExists {
// 		t.Errorf("CheckTemplate(%q, %q, %q) == %v, %q, want %v", "testbundle", "mytable", "single", gotExists, gotName, wantExists)
// 	}

// 	//4th PRIORITY: default template in non-bundle folder
// 	afero.WriteFile(appFS, "templates/mytemplates", []byte("default template"), 0644)
// 	t.Log("Template file single.html written to 'mytemplates' folder")
// 	wantExists = true
// 	wantName = "single.html"
// 	gotExists, gotName = CheckTemplate("testbundle", "mytable", "single")

// 	if gotExists != wantExists || gotName != wantName {
// 		t.Errorf("CheckTemplate(%q, %q, %q) == %v, %q, want %v, %q", "testbundle", "mytable", "single", gotExists, gotName, wantExists, wantName)
// 	}

// }
