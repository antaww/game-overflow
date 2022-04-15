package routes

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"main/sql"
	"main/utils"
	"mime/multipart"
	"net/http"
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
		// max size = 100 Mb
		err := r.ParseMultipartForm(100 << 20)
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

			buf := bytes.NewBuffer(nil)
			if _, err := io.Copy(buf, file); err != nil {
				log.Fatal(err)
			}

			profilePicture += base64.StdEncoding.EncodeToString(buf.Bytes())
			newUser.ProfilePic = profilePicture
		}

		user, err := LoginUser(r)
		if err != nil {
			log.Fatal(err)
		}

		_, err = sql.EditUser(user.Id, newUser)
		if err != nil {
			log.Fatal(err)
		}

		TemplatesData.ConnectedUser = sql.GetUserById(user.Id)

		r.Method = "GET"
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

func IsActiveRoute(w http.ResponseWriter, r *http.Request) {
	//fmt.Println(TemplatesData.ConnectedUser.IsOnline)
	//w.Header().Set("Content-Type", "application/json")
	//json.NewEncoder(w).Encode(TemplatesData.ConnectedUser.IsOnline)

	if r.Method == "POST" {
		// if header is content-type application/json
		if r.Header.Get("Content-Type") != "application/json" {
			return
		}

		var result struct {
			IsOnline bool `json:"isOnline"`
		}
		err := json.NewDecoder(r.Body).Decode(&result)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(result)

		user, err := LoginUser(r)
		if err != nil {
			log.Fatal(err)
		}

		if user == nil {
			return
		}

		err = sql.SetUserOnline(TemplatesData.ConnectedUser.Id, result.IsOnline)
		if err != nil {
			log.Fatal(err)
		}
	}
}
