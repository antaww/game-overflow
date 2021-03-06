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

// CookieRoute is the route to accept or decline the cookies
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

// LegalNoticeRoute is the route with the legal notice
func LegalNoticeRoute(w http.ResponseWriter, r *http.Request) {
	templateData, err := GetTemplateDataFromRoute(w, r)
	if err != nil {
		utils.RouteError(err)
	}

	err = utils.CallTemplate("legal-notice", templateData, w)
	if err != nil {
		utils.RouteError(err)
	}
}

// LoginRoute is the route for handling the login page
func LoginRoute(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		templateData, err := GetTemplateDataFromRoute(w, r)
		if err != nil {
			utils.RouteError(err)
		}

		err = utils.CallTemplate("login", templateData, w)
		if err != nil {
			utils.RouteError(err)
		}
	}

	if r.Method == "PUT" {
		var response struct {
			Success bool   `json:"success"`
			Session string `json:"session"`
			Error   string `json:"error"`
		}

		var data struct {
			Username string `json:"username"`
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

		username := data.Username
		password := data.Password

		userStruct, err := sql.GetUserByUsername(username)
		if err != nil {
			utils.RouteError(err)
		}

		match := utils.CheckPasswordHash(password, userStruct.Password)
		if !match {
			if userStruct.Password == "" {
				response.Error = fmt.Sprintf("User with username <span class=\"login-error-user\">%v</span> not found", username)
			} else {
				response.Error = fmt.Sprintf("Wrong password for user <span class=\"login-error-user\">%v</span>", username)
			}

			err = utils.SendResponse(w, response)
			if err != nil {
				utils.RouteError(err)
			}

			return
		}

		exists, err := sql.LoginByIdentifiants(username, userStruct.Password)
		if err != nil {
			utils.RouteError(err)
		}

		if exists {
			user, err := sql.GetUserByUsername(username)
			if err != nil {
				utils.RouteError(err)
			}

			if user.GetHasCookieEnabled().Bool {
				err = sql.AddSessionCookie(user, w)
				if err != nil {
					utils.RouteError(err)
				}
			}

			session, err := sql.AddSession(user)
			if err != nil {
				utils.RouteError(err)
			}

			response.Success = true
			response.Session = session
		} else {
			response.Success = false
		}

		ip := r.Header.Get("X-Forwarded-For")
		utils.ConnectedUsersWithoutCookies[ip] = response.Session

		err = utils.SendResponse(w, response)
		if err != nil {
			utils.RouteError(err)
		}
	}
}

// LogoutRoute is the route to log out the user
func LogoutRoute(w http.ResponseWriter, r *http.Request) {
	sessionId, err := sql.GetSessionId(r)
	if err == http.ErrNoCookie || sessionId == "" {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	if err != nil {
		utils.RouteError(err)
	}

	connectedUser, err := sql.GetUserBySession(sessionId)
	if err != nil {
		utils.RouteError(err)
	}

	err = sql.DeleteSessionCookie(sessionId, w)
	if err != nil {
		utils.RouteError(err)
	}

	ip := r.Header.Get("X-Forwarded-For")
	delete(utils.ConnectedUsersWithoutCookies, ip)

	err = sql.SetUserOnline(connectedUser.Id, false)
	if err != nil {
		utils.RouteError(err)
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// SignUpRoute is the route for handling the signup page
func SignUpRoute(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		templateData, err := GetTemplateDataFromRoute(w, r)
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
			utils.RouteError(errors.New("password does not match"))
		}

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
