package main

import (
	"html/template"
	"log"
	"net/http"
)

func main(){
	templ := template.Must(template.ParseFiles("index.gohtml")) //define html file

	fs := http.FileServer(http.Dir("main.css")) //define css file
	http.Handle("/static/", http.StripPrefix("/static/", fs)) //set css file to static

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		err := templ.ExecuteTemplate(w, "index.gohtml", nil)
		if err != nil {
			log.Fatal(err)
		}
	})

	err := http.ListenAndServe(":8091", nil)
	if err != nil {
		log.Fatal(err)
	}

}