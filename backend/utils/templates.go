package utils

import (
	"fmt"
	"html/template"
	"net/http"
	"time"
)

var templateMap = template.FuncMap{
	"safeURL": func(u string) template.URL {
		return template.URL(u)
	},
	"formatDate": func(t time.Time) string {
		return t.Format("2006-01-02 15:04:05")
	},
	"MinusOne": func(i int) int {
		return i - 1
	},
	"ConnectedUserMessage": func(i int64) bool {
		fmt.Println(i)
		//if i == routes.TemplatesData.ConnectedUser.Id {
		//	return true
		//}
		return false
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

	_ = templateName
	return nil
}
