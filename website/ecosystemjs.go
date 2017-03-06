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

const EcoSystemJS = `
//////////////////////////
// EcoSystem Javascript //
//////////////////////////

var EcoSystem = {
    "apiRoot" : "{{ .config.protocol }}://{{ .config.host }}:{{ .config.apiport }}"
};


//As soon as jq is registered, log in anonymously if not already
$(function () {
    //Check to see if we are already anonymously logged in
    //if not, do it
    if (!localStorage.getItem("token")) {
        $.anonLogOn();
    }
});


//Regular POST request to JSON api with JWT auth
jQuery.extend({
    jsonApiPost: function (params, endpoint) {
        return jQuery.ajax(jQuery.extend(params, {
            type: "POST",
            data: JSON.stringify(params.data),//stringify the input object
            dataType: "json", //expect json back
            contentType: "application/json", //send as json
            processData: false, //don't encode as form
            url: EcoSystem.apiRoot + "/api" + params.endpoint,
            headers: {
                "Authorization": "Bearer " + localStorage.getItem("token")
            }
        }));
    }
});

jQuery.extend({
    jsonApiDelete: function (params, endpoint) {
        return jQuery.ajax(jQuery.extend(params, {
            type: "DELETE",
            dataType: "json", //expect json back
            url: EcoSystem.apiRoot + "/api" + params.endpoint,
            headers: {
                "Authorization": "Bearer " + localStorage.getItem("token")
            }
        }));
    }
});

//Regular GET request to public HTML without auth
jQuery.extend({
    webApiGet: function (params, endpoint) {
        return jQuery.ajax(jQuery.extend(params, {
            type: "GET",
            dataType: "html", //expect html back
            url: "/site" + params.endpoint,
        }));
    }
});



//Regular GET request to private HTML api with JWT auth
jQuery.extend({
    htmlApiGet: function (params, endpoint) {
        return jQuery.ajax(jQuery.extend(params, {
            type: "GET",
            dataType: "html", //expect html back
            url: "/private" + params.endpoint,
            headers: {
                "Authorization": "Bearer " + localStorage.getItem("token")
            }
        }));
    }
});

//Anonymous login request to JSON api
//No 'api' in base of URL (doesn't run through authentication middleware)
//No authorisation header
jQuery.extend({
    anonLogOn: function (params, endpoint) {
        return jQuery.ajax(jQuery.extend(params, {
            type: "POST",
            data: JSON.stringify({ "username": "anon", "password": "anon" }),//stringify the input object
            dataType: "json", //expect json back
            contentType: "application/json", //send as json
            processData: false, //don't encode as form
            url: EcoSystem.apiRoot + "/login",
            success: function (json) {
                if (json.token) {
                    localStorage.setItem("token", json.token)
                }
            },
            error: function (err) {
                alert("Could not log in anonymously");
            }
        }));
    }
});
`
