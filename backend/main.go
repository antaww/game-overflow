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

func main() {
	templ, err := template.New("").ParseGlob("../client/templates/*.gohtml")
	if err != nil {
		log.Fatal(err)
	}

	css := http.FileServer(http.Dir("../client/style/"))              //define css file
	http.Handle("/static/", http.StripPrefix("/static/", css)) //set css file to static

	resources := http.FileServer(http.Dir("../backend/resources/"))                //define css file
	http.Handle("/resources/", http.StripPrefix("/resources/", resources)) //set css file to static

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		path := strings.TrimPrefix(r.URL.Path, "/")
		if path == "favicon.ico" {
			return
		}
		if path == "" {
			fmt.Println("index page loaded")
			err := templ.ExecuteTemplate(w, "index.gohtml", nil)
			if err != nil {
				log.Fatal(err)
			}
		}
	})

	http.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("login page loaded")
		err := templ.ExecuteTemplate(w, "login.gohtml", nil)
		if err != nil {
			log.Fatal(err)
		}
	})

	http.HandleFunc("/sign-up", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("sign-up page loaded")
		err := templ.ExecuteTemplate(w, "sign-up.gohtml", nil)
		if err != nil {
			log.Fatal(err)
		}

		if r.Method == "POST" {
			err := r.ParseForm()
			if err != nil {
				log.Fatal(err)
			}

			err = SaveUser(NewUser(
				r.FormValue("username"),
				r.FormValue("password"),
				r.FormValue("email"),
			))
			if err != nil {
				log.Fatal(err)
			}
			return
		}

	})

	// Capture connection properties.
	cfg := mysql.Config{
		User:                 "xhmyjae",
		Passwd:               "xhmyjae",
		Net:                  "tcp",
		Addr:                 "10.13.33.123:3306",
		DBName:               "forum",
		AllowNativePasswords: true,
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

	err = http.ListenAndServe(":8091", nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Server started at localhost:8091")

}
