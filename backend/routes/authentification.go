package routes

import (
	"encoding/json"
	"io"
	"main/sql"
	"main/utils"
	"net/http"
)

// ConfirmPasswordRoute is the route for the confirm password request
func ConfirmPasswordRoute(w http.ResponseWriter, r *http.Request) {
	if TemplatesData.ConnectedUser == nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	if r.Method == "PUT" {
		var data struct {
			Password string `json:"password"`
		}

		bytesRead, err := io.ReadAll(r.Body)
		if err != nil {
			utils.RouteError(err)
		}

		err = json.Unmarshal(bytesRead, &data)
		if err != nil {
			utils.RouteError(err)
		}

		user, err := sql.GetUserByRequest(r)
		if err != nil {
			utils.RouteError(err)
		}

		valid := sql.ConfirmPassword(user.Id, data.Password)
		var success = struct {
			Success bool `json:"success"`
		}{
			Success: valid,
		}

		w.Header().Set("Content-Type", "application/json")
		marshal, err := json.Marshal(success)
		if err != nil {
			utils.RouteError(err)
		}

		_, err = w.Write(marshal)
		if err != nil {
			utils.RouteError(err)
		}
	}
}

// LoginRoute is the route for handling the login page
func LoginRoute(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		err := utils.CallTemplate("login", TemplatesData, w)
		if err != nil {
			utils.RouteError(err)
		}
	}

	if r.Method == "POST" {
		err := r.ParseForm()
		if err != nil {
			utils.RouteError(err)
		}

		username := r.FormValue("username")
		exists, err := sql.LoginByIdentifiants(username, r.FormValue("password"))
		if err != nil {
			utils.RouteError(err)
		}

		if exists {
			// return when error, it will just cancel the request and invite the user to retry
			user, err := sql.GetUserByUsername(username)
			if err != nil {
				return
			}
			err = sql.AddSessionCookie(*user, w)
			if err != nil {
				return
			}
			http.Redirect(w, r, "/", http.StatusSeeOther)
		} else {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
		}
	}
}

// LogoutRoute is the route to log out the user
func LogoutRoute(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session")
	if err != nil {
		if err == http.ErrNoCookie {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		utils.RouteError(err)
	}

	err = sql.CookieLogout(*cookie, w)

	err = sql.SetUserOnline(TemplatesData.ConnectedUser.Id, false)
	if err != nil {
		utils.RouteError(err)
	}

	if err != nil {
		utils.RouteError(err)
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// SignUpRoute is the route for handling the signup page
func SignUpRoute(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		err := utils.CallTemplate("sign-up", TemplatesData, w)
		if err != nil {
			utils.RouteError(err)
		}
	}

	if r.Method == "POST" {
		err := r.ParseForm()
		if err != nil {
			utils.RouteError(err)
		}

		valid, err := sql.SaveUser(sql.CreateUser(
			r.FormValue("username"),
			r.FormValue("password"),
			r.FormValue("email"),
		))

		if err != nil {
			utils.RouteError(err)
		}

		if valid {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
		} else {
			http.Redirect(w, r, "/sign-up", http.StatusSeeOther)
		}
	}
}
