package main

import (
	"bufio"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"io"
	"log"
	. "main/sql"
	"main/utils"
	"net/http"
	"os"
	"strings"
)

type TemplatesDataType struct {
	ConnectedUser *User
	Locales       map[string]string
	ShownTopics   []Topic
}

var TemplatesData = TemplatesDataType{}

func main() {
	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatal(err)
	}

	css := http.FileServer(http.Dir("../client/style/"))
	http.Handle("/static/", http.StripPrefix("/static/", css))

	resources := http.FileServer(http.Dir("../backend/resources/"))
	http.Handle("/resources/", http.StripPrefix("/resources/", resources))

	scripts := http.FileServer(http.Dir("../client/scripts/"))
	http.Handle("/scripts/", http.StripPrefix("/scripts/", scripts))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		path := strings.TrimPrefix(r.URL.Path, "/")
		cookie, err := r.Cookie("session")

		if err != nil {
			TemplatesData.ConnectedUser = nil
		} else {
			user, err := GetUserBySession(cookie.Value)
			if err != nil {
				log.Fatal(err)
			}
			TemplatesData.ConnectedUser = user
		}

		if path == "" {
			err := utils.CallTemplate("main", TemplatesData, w)
			if err != nil {
				log.Fatal(err)
			}
		}
	})

	http.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			err := utils.CallTemplate("login", TemplatesData, w)
			if err != nil {
				log.Fatal(err)
			}
		}

		if r.Method == "POST" {
			err := r.ParseForm()
			if err != nil {
				log.Fatal(err)
			}

			username := r.FormValue("username")
			exists, err := LoginByIdentifiants(username, r.FormValue("password"))
			if err != nil {
				log.Fatal(err)
			}

			if exists {
				user := GetUserByUsername(username)
				SessionID(*user, w)
				http.Redirect(w, r, "/", http.StatusSeeOther)
				return
			} else {
				http.Redirect(w, r, "/login", http.StatusSeeOther)
				return
			}
		}
	})

	http.HandleFunc("/sign-up", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			err := utils.CallTemplate("sign-up", TemplatesData, w)
			if err != nil {
				log.Fatal(err)
			}
		}

		if r.Method == "POST" {
			err := r.ParseForm()
			if err != nil {
				log.Fatal(err)
			}

			valid, err := SaveUser(CreateUser(
				r.FormValue("username"),
				r.FormValue("password"),
				r.FormValue("email"),
			))

			if err != nil {
				log.Fatal(err)
			}

			if valid {
				fmt.Println("aa")
				http.Redirect(w, r, "/login", http.StatusSeeOther)
				return
			} else {
				http.Redirect(w, r, "/sign-up", http.StatusSeeOther)
				return
			}
		}
	})

	http.HandleFunc("/confirm-password", func(w http.ResponseWriter, r *http.Request) {
		if TemplatesData.ConnectedUser == nil {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		if r.Method == "PUT" {
			var data struct {
				Password string `json:"password"`
			}

			bytes, err := io.ReadAll(r.Body)
			if err != nil {
				log.Fatal(err)
			}

			err = json.Unmarshal(bytes, &data)
			if err != nil {
				log.Fatal(err)
			}

			cookie, err := r.Cookie("session")
			if err != nil {
				log.Fatal(err)
			}
			user, err := GetUserBySession(cookie.Value)
			if err != nil {
				log.Fatal(err)
			}

			valid := ConfirmPassword(user.Id, data.Password)
			var success = struct {
				Success bool `json:"success"`
			}{
				Success: valid,
			}

			w.Header().Set("Content-Type", "application/json")
			marshal, err := json.Marshal(success)
			if err != nil {
				log.Fatal(err)
			}

			_, err = w.Write(marshal)
			if err != nil {
				log.Fatal(err)
			}
		}
	})

	http.HandleFunc("/admin/edit-username", func(w http.ResponseWriter, r *http.Request) {
		err := utils.CallTemplate("admin-edit-username", TemplatesData, w)
		if err != nil {
			log.Fatal(err)
		}

		if r.Method == "POST" {
			err := r.ParseForm()
			if err != nil {
				log.Fatal(err)
			}

			err = AdminEditUsername(
				r.FormValue("old-username"),
				r.FormValue("new-username"),
			)

			if err != nil {
				log.Fatal(err)
			}
			return
		}
	})

	http.HandleFunc("/edit-username", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			err := utils.CallTemplate("edit-username", TemplatesData, w)
			if err != nil {
				log.Fatal(err)
			}
		}

		if r.Method == "POST" {
			err := r.ParseForm()
			if err != nil {
				log.Fatal(err)
			}

			//get cookie from browser
			cookie, err := r.Cookie("session")
			if err != nil {
				log.Fatal(err)
			}

			//select user from session
			result, err := DB.Query("SELECT id_user FROM sessions WHERE id_session = ?", cookie.Value)
			if err != nil {
				log.Fatal(err)
			}

			//get result from query
			var idUser int64
			if result.Next() {
				err = result.Scan(&idUser)
			}

			//Handle sql errors, close the query to avoid memory leaks
			HandleSQLErrors(result)

			// Get User, save for TemplatesData (to show user logged in in templates)
			userConnected := GetUserById(idUser)
			TemplatesData.ConnectedUser = userConnected

			//edit username of idUser
			exists, err := EditUsername(idUser, r.FormValue("new-username"))
			if err != nil {
				log.Fatal(err)
			}

			if exists {
				http.Redirect(w, r, "/", http.StatusSeeOther)
				return
			} else {
				http.Redirect(w, r, "/edit-username", http.StatusSeeOther)
				return
			}
		}
	})

	http.HandleFunc("/edit-avatar", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("cc")
		if r.Method == "GET" {
			fmt.Println("dd")
			err := utils.CallTemplate("edit-avatar", TemplatesData, w)
			if err != nil {
				log.Fatal(err)
			}
		}

		if r.Method == "POST" {
			err := r.ParseForm()
			if err != nil {
				log.Fatal(err)
			}

			//get cookie from browser
			cookie, err := r.Cookie("session")
			if err != nil {
				log.Fatal(err)
			}

			//select user from session
			result, err := DB.Query("SELECT id_user FROM sessions WHERE id_session = ?", cookie.Value)
			if err != nil {
				log.Fatal(err)
			}

			//get result from query
			var idUser int64
			if result.Next() {
				err = result.Scan(&idUser)
			}

			//Handle sql errors, close the query to avoid memory leaks
			HandleSQLErrors(result)

			// Get User, save for TemplatesData (to show user logged in in templates)
			userConnected := GetUserById(idUser)
			TemplatesData.ConnectedUser = userConnected

			//edit username of idUser
			imgForm, file, err := r.FormFile("avatar")
			if err != nil {
				log.Fatal(err)
			}
			defer imgForm.Close()
			var base64Encoding string
			if strings.HasSuffix(file.Filename, "jpeg") {
				base64Encoding += "data:image/jpeg;base64,"
			} else if strings.HasSuffix(file.Filename, "png") {
				base64Encoding += "data:image/png;base64,"
			}

			reader := bufio.NewReader(imgForm)
			content, _ := io.ReadAll(reader)
			encoded := base64.StdEncoding.EncodeToString(content)
			base64Encoding += encoded

			exists, err := EditAvatar(idUser, base64Encoding)
			if err != nil {
				log.Fatal(err)
			}

			if exists {
				fmt.Println("aa")
				http.Redirect(w, r, "/", http.StatusSeeOther)
				return
			} else {
				fmt.Println("bb")
				http.Redirect(w, r, "/edit-avatar", http.StatusSeeOther)
				return
			}

		}
	})

	http.HandleFunc("/edit-password", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			err := utils.CallTemplate("edit-password", TemplatesData, w)
			if err != nil {
				log.Fatal(err)
			}
		}

		if r.Method == "POST" {
			err := r.ParseForm()
			if err != nil {
				log.Fatal(err)
			}

			//get cookie from browser
			cookie, err := r.Cookie("session")
			if err != nil {
				log.Fatal(err)
			}

			//select user from session
			result, err := DB.Query("SELECT id_user FROM sessions WHERE id_session = ?", cookie.Value)
			if err != nil {
				log.Fatal(err)
			}

			//get result from query
			var idUser int64
			if result.Next() {
				err = result.Scan(&idUser)
			}

			//Handle sql errors, close the query to avoid memory leaks
			HandleSQLErrors(result)

			// Get User, save for TemplatesData (to show user logged in in templates)
			userConnected := GetUserById(idUser)
			TemplatesData.ConnectedUser = userConnected

			//edit username of idUser
			exists, err := EditPassword(idUser, r.FormValue("old-password"), r.FormValue("new-password"))

			if err != nil {
				log.Fatal(err)
			}

			if exists {
				http.Redirect(w, r, "/", http.StatusSeeOther)
				return
			} else {
				http.Redirect(w, r, "/edit-password", http.StatusSeeOther)
				return
			}
		}
	})

	http.HandleFunc("/logout", func(w http.ResponseWriter, r *http.Request) {
		//get cookie from browser
		cookie, err := r.Cookie("session")
		if err != nil {
			log.Fatal(err)
		}

		//logout user
		err = CookieLogout(*cookie, w)

		if err != nil {
			log.Fatal(err)
		}
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	})

	http.HandleFunc("/settings", func(w http.ResponseWriter, r *http.Request) {
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

			newUser := User{
				Username:    r.FormValue("username"),
				Email:       r.FormValue("email"),
				Description: r.FormValue("description"),
				Locale:      r.FormValue("locale"),
			}

			cookie, err := r.Cookie("session")
			if err != nil {
				log.Fatal(err)
			}
			user, err := GetUserBySession(cookie.Value)
			if err != nil {
				log.Fatal(err)
			}

			_, err = EditUser(user.Id, newUser)
			if err != nil {
				log.Fatal(err)
			}

			TemplatesData.ConnectedUser = GetUserById(user.Id)

			r.Method = "GET"
			http.Redirect(w, r, "/settings", http.StatusSeeOther)
		}
	})

	http.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		topic, err := GetPost(1)
		if err != nil {
			log.Fatal(err)
		}

		err = topic.FetchMessages()
		if err != nil {
			log.Fatal(err)
		}

		TemplatesData.ShownTopics = append(TemplatesData.ShownTopics, *topic)

		err = utils.CallTemplate("topic", TemplatesData.ShownTopics[0], w)
		if err != nil {
			log.Fatal(err)
		}
	})

	http.HandleFunc("/feed", func(w http.ResponseWriter, r *http.Request) {
		queries := r.URL.Query()

		if queries.Has("category") {
			category := queries.Get("category")

			topics, err := GetTopicsByCategories(category)
			if err != nil {
				log.Fatal(err)
			}

			TemplatesData.ShownTopics = topics

			err = utils.CallTemplate("feed", TemplatesData, w)
			if err != nil {
				log.Fatal(err)
			}
		}

		return
	})

	// Capture connection properties.
	cfg := mysql.Config{
		User:                 os.Getenv("DB_USER"),
		Passwd:               os.Getenv("DB_PASSWORD"),
		Net:                  "tcp",
		Addr:                 os.Getenv("DB_ADDRESS"),
		DBName:               "forum",
		AllowNativePasswords: true,
		ParseTime:            true,
	}
	// Get a database handle.
	DB, err = sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		log.Fatal(err)
	}

	pingErr := DB.Ping()
	if pingErr != nil {
		log.Fatal(pingErr)
	}
	fmt.Println("Connected!")

	fmt.Println("Server started at localhost:8091")
	err = http.ListenAndServe(":8091", http.HandlerFunc(LogHandler))
	if err != nil {
		log.Fatal(err)
	}
}

func LogHandler(w http.ResponseWriter, r *http.Request) {
	if !strings.HasSuffix(".css", r.URL.String()) && !strings.HasSuffix(".png", r.URL.String()) {
		log.Printf("%v %v", r.Method, r.URL.String())
	}
	http.DefaultServeMux.ServeHTTP(w, r)
}
