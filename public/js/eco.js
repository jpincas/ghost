//////////////////////////
// EcoSystem Javascript //
//////////////////////////

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
            url: "http://localhost:3000/api" + params.endpoint,
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
            url: "http://localhost:3000/api" + params.endpoint,
            headers: {
                "Authorization": "Bearer " + localStorage.getItem("token")
            }
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
            url: "http://localhost:3000/login",
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
