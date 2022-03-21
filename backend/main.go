package main

import (
	"database/sql"
	"fmt"
	"github.com/go-sql-driver/mysql"
	"html/template"
	"log"
	"net/http"
	"strings"
)

var db *sql.DB

func main() {
	templ, err := template.New("").ParseGlob("../templates/*.gohtml")
	if err != nil {
		log.Fatal(err)
	}

	css := http.FileServer(http.Dir("../style/"))              //define css file
	http.Handle("/static/", http.StripPrefix("/static/", css)) //set css file to static

	resources := http.FileServer(http.Dir("../resources/"))                //define css file
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
	})

	// Capture connection properties.
	cfg := mysql.Config{
		User:                 "root",
		Passwd:               "",
		Net:                  "tcp",
		Addr:                 "127.0.0.1:3306",
		DBName:               "forum",
		AllowNativePasswords: true,
	}
	// Get a database handle.
	db, err = sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		log.Fatal(err)
	}

	pingErr := db.Ping()
	if pingErr != nil {
		log.Fatal(pingErr)
	}
	fmt.Println("Connected!")

	db.Exec("UPDATE users SET id_user = 4, username = 'test', is_online = 0, password = 'password', email = 'email', locale = 'fr', profile_pic = null, description = null, created_at = '2022-03-21 12:38:58', role_type = 'admin' WHERE id_user = 4")

	err = http.ListenAndServe(":8091", nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Server started at localhost:8091")

}
