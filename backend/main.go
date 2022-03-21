package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strings"
)

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
			err := templ.ExecuteTemplate(w, "index.gohtml", nil)
			if err != nil {
				log.Fatal(err)
			}
		}
	})

	http.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		err := templ.ExecuteTemplate(w, "login.gohtml", nil)
		if err != nil {
			log.Fatal(err)
		}
	})

	fmt.Println("Server started at localhost:8091")
	err = http.ListenAndServe(":8091", nil)
	if err != nil {
		log.Fatal(err)
	}

}
