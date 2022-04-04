package routes

import (
	"bytes"
	"encoding/base64"
	"io"
	"log"
	"main/sql"
	"main/utils"
	"mime/multipart"
	"net/http"
	"strconv"
)

func SettingsRoute(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		if TemplatesData.ConnectedUser == nil {
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}
		err := utils.CallTemplate("settings", TemplatesData, w)
		if err != nil {
			log.Fatal(err)
		}
	}

	if r.Method == "POST" {
		err := r.ParseForm()
		if err != nil {
			log.Fatal(err)
		}

		newUser := sql.User{
			Username:    r.FormValue("username"),
			Email:       r.FormValue("email"),
			Description: r.FormValue("description"),
			Locale:      r.FormValue("locale"),
		}

		var profilePicture string
		file, header, err := r.FormFile("profile-picture")
		defer func(file multipart.File) {
			err := file.Close()
			if err != nil {
				log.Fatal(err)
			}
		}(file)
		if err != nil {
			if err != http.ErrMissingFile {
				log.Fatal(err)
			}
		} else {
			profilePicture = "data:" + header.Header.Get("Content-Type") + ";base64,"

			buf := bytes.NewBuffer(nil)
			if _, err := io.Copy(buf, file); err != nil {
				log.Fatal(err)
			}

			profilePicture += base64.StdEncoding.EncodeToString(buf.Bytes())
			newUser.ProfilePic = profilePicture
		}

		cookie, err := r.Cookie("session")
		if err != nil {
			log.Fatal(err)
		}
		user, err := sql.GetUserBySession(cookie.Value)
		if err != nil {
			log.Fatal(err)
		}

		_, err = sql.EditUser(user.Id, newUser)
		if err != nil {
			log.Fatal(err)
		}

		TemplatesData.ConnectedUser = sql.GetUserById(user.Id)

		r.Method = "GET"
		http.Redirect(w, r, "/Settings", http.StatusSeeOther)
	}
}

func LikeRoute(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		cookie, err := r.Cookie("session")
		if err != nil {
			if err == http.ErrNoCookie {
				return
			}

			log.Fatal(err)
		}

		user, err := sql.GetUserBySession(cookie.Value)
		if err != nil {
			log.Fatal(err)
		}

		TemplatesData.ConnectedUser = user

		messageIdArg := r.URL.Query().Get("id")
		messageId, err := strconv.ParseInt(messageIdArg, 10, 64)
		if err != nil {
			log.Fatal(err)
		}

		messageLike, err := sql.MessageGetLikeFrom(messageId, user.Id)
		if err != nil {
			log.Fatal(err)
		}

		if messageLike == nil {
			_, err = sql.LikeMessage(messageId, user.Id)
			if err != nil {
				log.Fatal(err)
			}
		} else if messageLike.Like {
			_, err = sql.DeleteLikeMessage(messageId, user.Id)
			if err != nil {
				log.Fatal(err)
			}
		} else {
			_, err = sql.DeleteDislikeMessage(messageId, user.Id)
			if err != nil {
				log.Fatal(err)
			}

			_, err = sql.LikeMessage(messageId, user.Id)
			if err != nil {
				log.Fatal(err)
			}
		}

		err = TemplatesData.ShownTopic.FetchMessages()
		if err != nil {
			log.Fatal(err)
		}

		err = utils.ReloadActualTemplate(TemplatesData.ShownTopic, w)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func DislikeRoute(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		cookie, err := r.Cookie("session")
		if err != nil {
			if err == http.ErrNoCookie {
				return
			}

			log.Fatal(err)
		}

		user, err := sql.GetUserBySession(cookie.Value)
		if err != nil {
			log.Fatal(err)
		}

		TemplatesData.ConnectedUser = user

		messageIdArg := r.URL.Query().Get("id")
		messageId, _ := strconv.ParseInt(messageIdArg, 10, 64)

		messageLike, err := sql.MessageGetLikeFrom(messageId, user.Id)
		if err != nil {
			log.Fatal(err)
		}
		if messageLike == nil {
			_, err = sql.DislikeMessage(messageId, user.Id)
			if err != nil {
				log.Fatal(err)
			}
		} else if !messageLike.Like {
			_, err = sql.DeleteDislikeMessage(messageId, user.Id)
			if err != nil {
				log.Fatal(err)
			}
		} else {
			_, err = sql.DeleteLikeMessage(messageId, user.Id)
			if err != nil {
				log.Fatal(err)
			}

			_, err = sql.DislikeMessage(messageId, user.Id)
			if err != nil {
				log.Fatal(err)
			}
		}

		_, err = sql.DislikeMessage(TemplatesData.ConnectedUser.Id, messageId)
		if err != nil {
			log.Fatal(err)
		}
	}
}
