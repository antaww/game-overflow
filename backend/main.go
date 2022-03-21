package backend

import (
	"html/template"
	"net/http"
)

func main(){
	templ := template.Must(template.ParseFiles("index.gohtml")) //define html file

	fs := http.FileServer(http.Dir("main.css")) //define css file
	http.Handle("/static/", http.StripPrefix("/static/", fs)) //set css file to static

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			templ.ExecuteTemplate(w, "index.gohtml", nil)
	})

	http.ListenAndServe(":8091", nil)

}