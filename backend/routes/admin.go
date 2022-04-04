package routes

import (
	"log"
	"main/sql"
	"main/utils"
	"net/http"
)

func AdminEditUsernameRoute(w http.ResponseWriter, r *http.Request) {
	err := utils.CallTemplate("admin-edit-username", TemplatesData, w)
	if err != nil {
		log.Fatal(err)
	}

	if r.Method == "POST" {
		err := r.ParseForm()
		if err != nil {
			log.Fatal(err)
		}

		err = sql.AdminEditUsername(
			r.FormValue("old-username"),
			r.FormValue("new-username"),
		)

		if err != nil {
			log.Fatal(err)
		}
		return
	}
}
