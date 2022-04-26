package routes

import (
	"encoding/base64"
	"encoding/json"
	"io"
	"log"
	"main/sql"
	"main/utils"
	"mime/multipart"
	"net/http"
	"strconv"
	"strings"
)

// IsActiveRoute is a middleware function that checks if the user is active
func IsActiveRoute(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		var response struct {
			IsOnline  bool   `json:"isOnline"`
			SessionId string `json:"sessionId"`
		}
		err := json.NewDecoder(r.Body).Decode(&response)
		if err != nil {
			return
		}

		var user *sql.User

		if response.SessionId != "" {
			user, err = sql.GetUserBySession(response.SessionId)
			if err != nil {
				log.Fatal(err)
			}
		} else {
			user, err = LoginUser(r)
			if err != nil {
				log.Fatal(err)
			}
		}

		if user == nil {
			return
		}

		err = sql.SetUserOnline(user.Id, response.IsOnline)
		if err != nil {
			log.Fatal(err)
		}
	}
}

// SettingsRoute is a route that handles the settings of the user
func SettingsRoute(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		if TemplatesData.ConnectedUser == nil {
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		TemplatesData.ConnectedUser.CalculateDefaultColor()

		err := utils.CallTemplate("settings", TemplatesData, w)
		if err != nil {
			log.Fatal(err)
		}
	}

	if r.Method == "POST" {
		err := r.ParseMultipartForm(100 << 20)
		if err != nil {
			log.Fatal(err)
		}

		colorValue := r.FormValue("color")
		colorValue = strings.TrimPrefix(colorValue, "#")
		color, err := strconv.ParseInt(colorValue, 16, 32)
		if err != nil {
			log.Fatal(err)
		}

		newUser := sql.User{
			Username:    r.FormValue("username"),
			Email:       r.FormValue("email"),
			Description: r.FormValue("description"),
			Locale:      r.FormValue("locale"),
			Color:       int(color),
		}

		var profilePicture string
		file, header, err := r.FormFile("profile-picture")
		if header != nil {
			defer func(file multipart.File) {
				err := file.Close()
				if err != nil {
					log.Fatal(err)
				}
			}(file)
		}
		if err != nil {
			if err != http.ErrMissingFile {
				log.Fatal(err)
			}
		} else {
			profilePicture = "data:" + header.Header.Get("Content-Type") + ";base64,"

			bytes, err := io.ReadAll(file)
			if err != nil {
				log.Fatal(err)
			}

			profilePicture += base64.StdEncoding.EncodeToString(bytes)
			newUser.ProfilePic = profilePicture
		}

		user, err := sql.GetUserByRequest(r)
		if err != nil {
			log.Fatal(err)
		}

		_, err = sql.EditUser(user.Id, newUser)
		if err != nil {
			log.Fatal(err)
		}

		updatedUser, err := sql.GetUserByRequest(r)
		if err != nil {
			log.Fatal(err)
		}

		TemplatesData.ConnectedUser = updatedUser

		r.Method = "GET"
		http.Redirect(w, r, "/settings", http.StatusSeeOther)
	}
}

// UsersActive is a middleware function that checks if the users sent are active
func UsersActive(w http.ResponseWriter, r *http.Request) {
	if r.Method == "PUT" {
		if r.Header.Get("Content-Type") != "application/json" {
			return
		}

		var response struct {
			Users []string `json:"users"`
		}
		err := json.NewDecoder(r.Body).Decode(&response)
		if err != nil {
			log.Fatal(err)
		}

		usersOnline, err := sql.GetUsersStatus(response.Users)
		if err != nil {
			log.Fatal(err)
		}

		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(usersOnline)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func UserPostsRoute(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		if TemplatesData.ConnectedUser == nil {
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}
		query := r.URL.Query()
		if query.Has("id") {
			queryId := query.Get("id")
			id, err := strconv.ParseInt(queryId, 10, 64)
			if err != nil {
				log.Fatal(err)
			}

			user, err := sql.GetUserById(id)
			if err != nil {
				log.Fatal(err)
			}

			topics, err := sql.GetUserTopics(user.Id)

			if err != nil {
				log.Fatal(err)
			}

			TemplatesData.ShownTopics = topics

			err = utils.CallTemplate("feed", TemplatesData, w)
			if err != nil {
				log.Fatal(err)
			}
		} else {
			http.Redirect(w, r, "/", http.StatusSeeOther)
		}

	}
}
