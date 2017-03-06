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

package adminpanel

//This html import is dynamically created by substituting values from
//the main config.json.
//It is imported automatically into the admin panel by most of the components
//For some reason, Viper flattens the case of map keys when unmarshalling, so refer to key all in lower case

//Admin is the main html import template for the admin panel
const PolymerAdminImportTemplate = `
<!--BUNDLE ACTION IMPORTS-->
{{ range .bundles }}<link rel="import" href="/admin-panel/{{ . }}/actions.html">
{{ end }}

<!-- ADMIN PANEL STYLES -->
<dom-module id="custom-styles">
    <template>
        <style>
            :root {
                --secondary-color: {{ .config.adminsecondarycolor }};
                --primary-color: {{ .config.adminprimarycolor }};
                --primary-text-color: {{ .config.admintextcolor }};
                --error-color:{{ .config.adminerrorcolor }};
                /*--icon-button-color: var(--primary-text-color);*/
                /*--sidebar-color: var(--secondary-color);*/
            }
        </style>
    </template>
</dom-module>

<!--ADMIN PANEL PROPERTIES-->
<script>
    // Custom settings are available on all elements by including the behavior EcoCustom.Settings
    // and then referencing the name of the setting directly

    //Create the EcoCustom namespace just in case no bundle has already created it
    var EcoCustom = EcoCustom || {};
    EcoCustom.Properties = EcoCustom.Properties || {};

    //Setting the .properties attribute will mean these properties are merged into the properties
    //of whichever component these are imported into (which is every eco- component in the admin interface).
    //Only import one of these files, otherwise they will overwrite each other

    EcoCustom.Properties.properties = {

        adminPanelTitle: {
            type: String,
            value: "{{ .config.admintitle }}"
        },
        logos: {
            type: Object,
            value: {
                "horizontal": "{{ .config.protocol }}://{{ .config.host }}:{{ .config.apiport }}/admin-panel/{{ .config.adminlogobundle }}/{{ .config.adminlogofile }}",
                "vertical": "{{ .config.protocol }}://{{ .config.host }}:{{ .config.apiport }}/admin-panel/{{ .config.adminlogobundle }}/{{ .config.adminlogofile }}"
            }
        },
        apiRoot: {
            type: String,
            value: "{{ .config.protocol }}://{{ .config.host }}:{{ .config.apiport }}"
        },
        bundlesInstalled: {
            type: Array,
            value: {{ .bundles }}
        }
    }

</script>`
