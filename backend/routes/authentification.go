package routes

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"main/sql"
	"main/utils"
	"net/http"
	"strconv"
)

// ConfirmPasswordRoute is the route for the confirm password request
func ConfirmPasswordRoute(w http.ResponseWriter, r *http.Request) {
	user, err := sql.GetUserByRequest(r)
	if user == nil || err != nil {
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
		templateData, err := GetTemplatesDataFromRoute(w, r)
		if err != nil {
			utils.RouteError(err)
		}

		err = utils.CallTemplate("login", templateData, w)
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
		password := r.FormValue("password")
		userStruct, err := sql.GetUserByUsername(username)
		if err != nil {
			utils.RouteError(err)
		}
		match := utils.CheckPasswordHash(password, userStruct.Password)
		if !match {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		exists, err := sql.LoginByIdentifiants(username, userStruct.Password)
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
	connectedUser, err := sql.GetUserByRequest(r)
	if err != nil {
		if err == http.ErrNoCookie {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		utils.RouteError(err)
	}

	sessionId, err := sql.GetSessionId(r)
	if err != nil {
		utils.RouteError(err)
	}

	err = sql.DeleteSessionCookie(sessionId, w)
	if err != nil {
		utils.RouteError(err)
	}

	err = sql.SetUserOnline(connectedUser.Id, false)
	if err != nil {
		utils.RouteError(err)
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// SignUpRoute is the route for handling the signup page
func SignUpRoute(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		templateData, err := GetTemplatesDataFromRoute(w, r)
		if err != nil {
			utils.RouteError(err)
		}

		err = utils.CallTemplate("sign-up", templateData, w)
		if err != nil {
			utils.RouteError(err)
		}
	}

	if r.Method == "POST" {
		err := r.ParseForm()
		if err != nil {
			utils.RouteError(err)
		}

		HashedPassword, err := utils.HashPassword(r.FormValue("password"))
		if err != nil {
			utils.RouteError(err)
		}
		match := utils.CheckPasswordHash(r.FormValue("password"), HashedPassword)
		if !match {
			utils.RouteError(errors.New("Password does not match"))
		}

		fmt.Println(HashedPassword)
		valid, err := sql.SaveUser(sql.CreateUser(
			r.FormValue("username"),
			HashedPassword,
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

func CookieRoute(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		queries := r.URL.Query()
		if queries.Has("accept") {
			accept, err := strconv.ParseBool(queries.Get("accept"))
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			user, err := sql.GetUserByRequest(r)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			err = user.SetCookiesEnabled(accept)
			if err != nil {
				utils.RouteError(err)
			}

			w.WriteHeader(http.StatusOK)
		}
	}
}
