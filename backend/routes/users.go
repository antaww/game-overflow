package routes

import (
	"encoding/base64"
	"encoding/json"
	"io"
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
				return
			}
		} else {
			user, err = LoginUser(r)
			if err != nil {
				utils.RouteError(err)
			}
		}

		if user == nil {
			return
		}

		err = sql.SetUserOnline(user.Id, response.IsOnline)
		if err != nil {
			utils.RouteError(err)
		}
	}
}

// ProfileRoute is a route that returns the user profile
func ProfileRoute(w http.ResponseWriter, r *http.Request) {
	var userId int64
	var user *sql.User

	query := r.URL.Query()
	userIdString := query.Get("id")
	userId, err := strconv.ParseInt(userIdString, 10, 64)
	if err != nil {
		if user, err = sql.GetUserByRequest(r); err != nil {
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		} else {
			userId = user.Id
			query.Del("id")
		}
	}

	if user == nil {
		user, err = sql.GetUserById(userId)
		if err != nil {
			utils.RouteError(err)
		}
	}

	if user.Username == "" {
		http.Redirect(w, r, "/profile", http.StatusSeeOther)
		return
	}

	if user == nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
	} else {
		TemplatesData.ShownUser = user
		err := utils.CallTemplate("profile", TemplatesData, w)
		if err != nil {
			utils.RouteError(err)
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
			utils.RouteError(err)
		}
	}

	if r.Method == "POST" {
		err := r.ParseMultipartForm(100 << 20)
		if err != nil {
			utils.RouteError(err)
		}

		colorValue := r.FormValue("color")
		colorValue = strings.TrimPrefix(colorValue, "#")
		color, err := strconv.ParseInt(colorValue, 16, 32)
		if err != nil {
			utils.RouteError(err)
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
					utils.RouteError(err)
				}
			}(file)
		}
		if err != nil {
			if err != http.ErrMissingFile {
				utils.RouteError(err)
			}
		} else {
			profilePicture = "data:" + header.Header.Get("Content-Type") + ";base64,"

			bytes, err := io.ReadAll(file)
			if err != nil {
				utils.RouteError(err)
			}

			profilePicture += base64.StdEncoding.EncodeToString(bytes)
			newUser.ProfilePic = profilePicture
		}

		user, err := sql.GetUserByRequest(r)
		if err != nil {
			utils.RouteError(err)
		}

		_, err = sql.EditUser(user.Id, newUser)
		if err != nil {
			utils.RouteError(err)
		}

		updatedUser, err := sql.GetUserByRequest(r)
		if err != nil {
			utils.RouteError(err)
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
			utils.RouteError(err)
		}

		usersOnline, err := sql.GetUsersStatus(response.Users)
		if err != nil {
			utils.RouteError(err)
		}

		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(usersOnline)
		if err != nil {
			utils.RouteError(err)
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
				utils.RouteError(err)
			}

			user, err := sql.GetUserById(id)
			if err != nil {
				utils.RouteError(err)
			}

			topics, err := sql.GetUserTopics(user.Id)

			if err != nil {
				utils.RouteError(err)
			}

			TemplatesData.ShownTopics = topics

			err = utils.CallTemplate("feed", TemplatesData, w)
			if err != nil {
				utils.RouteError(err)
			}
		} else {
			http.Redirect(w, r, "/", http.StatusSeeOther)
		}

	}
}

func UserLikesRoute(w http.ResponseWriter, r *http.Request) {
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
				utils.RouteError(err)
			}

			user, err := sql.GetUserById(id)
			if err != nil {
				utils.RouteError(err)
			}

			if TemplatesData.ConnectedUser.Id != user.Id {
				http.Redirect(w, r, "/", http.StatusSeeOther)
				return
			}

			topics, err := sql.GetUserLikesTopics(user.Id)
			if err != nil {
				utils.RouteError(err)
			}

			messages, err := sql.GetUserLikedMessages(user.Id)
			if err != nil {
				utils.RouteError(err)
			}

			TemplatesData.ShownTopics = topics
			TemplatesData.ShownMessages = messages

			err = utils.CallTemplate("feed", TemplatesData, w)
			if err != nil {
				utils.RouteError(err)
			}
		} else {
			http.Redirect(w, r, "/", http.StatusSeeOther)
		}
	}
}

func FollowUserRoute(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		if TemplatesData.ConnectedUser == nil {
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		body := r.Body
		defer body.Close()

		var response struct {
			Id string `json:"id"`
		}
		err := json.NewDecoder(body).Decode(&response)
		if err != nil {
			utils.RouteError(err)
		}
		idUserFollowed, err := strconv.ParseInt(response.Id, 10, 64)
		if err != nil {
			utils.RouteError(err)
		}

		err = sql.FollowUser(idUserFollowed, TemplatesData.ConnectedUser.Id)
		if err != nil {
			w.WriteHeader(http.StatusForbidden)
		}
	}
	w.WriteHeader(http.StatusOK)
}
