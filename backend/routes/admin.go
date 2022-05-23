package routes

import (
	"main/sql"
	"main/utils"
	"net/http"
)

// AdminEditUsernameRoute is a route for editing a user's username
func AdminEditUsernameRoute(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		templateData, err := GetTemplateDataFromRoute(w, r)
		if err != nil {
			utils.RouteError(err)
		}

		err = utils.CallTemplate("admin-edit-username", templateData, w)
		if err != nil {
			utils.RouteError(err)
		}
	}

	if r.Method == "POST" {
		err := r.ParseForm()
		if err != nil {
			utils.RouteError(err)
		}

		err = sql.AdminEditUsername(
			r.FormValue("old-username"),
			r.FormValue("new-username"),
		)

		if err != nil {
			utils.RouteError(err)
		}
	}
}
