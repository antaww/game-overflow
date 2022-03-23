package main

import (
	"database/sql"
	"fmt"
	"github.com/go-sql-driver/mysql"
	"html/template"
	"log"
	. "main/sql"
	"net/http"
	"strings"
)

type TemplatesDataType struct {
	ConnectedUser *User
}

var TemplatesData = TemplatesDataType{}

func main() {
	templates, err := template.New("").ParseGlob("../client/templates/*.gohtml")
	if err != nil {
		log.Fatal(err)
	}

	css := http.FileServer(http.Dir("../client/style/"))       //define css file
	http.Handle("/static/", http.StripPrefix("/static/", css)) //set css file to static

	resources := http.FileServer(http.Dir("../backend/resources/"))        //define css file
	http.Handle("/resources/", http.StripPrefix("/resources/", resources)) //set css file to static

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		path := strings.TrimPrefix(r.URL.Path, "/")
		cookie, err := r.Cookie("session")

		if err != nil {
			TemplatesData.ConnectedUser = nil
		} else {
			user, err := LoginBySession(cookie.Value)
			if err != nil {
				log.Fatal(err)
			}
			TemplatesData.ConnectedUser = user
		}

		if path == "" {
			err := templates.ExecuteTemplate(w, "index.gohtml", TemplatesData)
			if err != nil {
				log.Fatal(err)
			}
		}
	})

	http.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			err := templates.ExecuteTemplate(w, "login", nil)
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
		err := templates.ExecuteTemplate(w, "sign-up", nil)
		if err != nil {
			log.Fatal(err)
		}

		if r.Method == "POST" {
			err := r.ParseForm()
			if err != nil {
				log.Fatal(err)
			}

			err = SaveUser(CreateUser(
				r.FormValue("username"),
				r.FormValue("password"),
				r.FormValue("email"),
			))
			if err != nil {
				log.Fatal(err)
			}
		}
	})

	http.HandleFunc("/admin/edit-username", func(w http.ResponseWriter, r *http.Request) {
		err := templates.ExecuteTemplate(w, "admin-edit-username", nil)
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
		err := templates.ExecuteTemplate(w, "edit-username", nil)
		if err != nil {
			log.Fatal(err)
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
			err = EditUsername(idUser, r.FormValue("new-username"))

			if err != nil {
				log.Fatal(err)
			}
		}
	})

	// Capture connection properties.
	cfg := mysql.Config{
		User:                 "Ayfri",
		Passwd:               "ayfri",
		Net:                  "tcp",
		Addr:                 "10.13.33.123:3306",
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
