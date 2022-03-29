package utils

import (
	"html/template"
	"net/http"
)

var templateMap = template.FuncMap{
	"safeURL": func(u string) template.URL {
		return template.URL(u)
	},
}

func CallTemplate(templateName string, data interface{}, w http.ResponseWriter) error {
	templates := template.New("").Funcs(templateMap)
	templates, err := templates.ParseFiles("../client/templates/main.gohtml", "../client/templates/"+templateName+".gohtml")
	if err != nil {
		return err
	}
	templates, err = templates.ParseGlob("../client/templates/components/*.gohtml")
	if err != nil {
		return err
	}

	err = templates.ExecuteTemplate(w, "main", data)
	if err != nil {
		return err
	}
	return nil
}
