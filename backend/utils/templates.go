package utils

import (
	"html/template"
	"net/http"
	"strconv"
	"time"
)

var TemplateMap = template.FuncMap{
	"safeURL": func(u string) template.URL {
		return template.URL(u)
	},
	"safeHTML": func(u string) template.HTML {
		return template.HTML(u)
	},
	"formatDate": func(t time.Time) string {
		return t.Format("2006-01-02 15:04:05")
	},
	"formatDateLong": func(t time.Time) string {
		return t.Format("Monday, January 2, 2006 at 3:04pm")
	},
	"formatRelativeDate": func(t time.Time) string {
		now := time.Now()
		diff := now.Sub(t)

		if diff < time.Minute {
			return "just now"
		}

		if diff < time.Hour {
			return strconv.Itoa(int(diff.Minutes())) + " minutes ago"
		}

		if diff < time.Hour*24 {
			return strconv.Itoa(int(diff.Hours())) + " hours ago"
		}

		if diff < time.Hour*24*365 {
			return strconv.Itoa(int(diff.Hours()/24)) + " days ago"
		}

		return strconv.Itoa(int(diff.Hours()/24/365)) + " years ago"
	},
	"decimalToHex": func(i int) string {
		s := "#" + strconv.FormatInt(int64(i), 16)
		return s
	},
}

func CallTemplate(templateName string, data interface{}, w http.ResponseWriter) error {
	templates := template.New("").Funcs(TemplateMap)
	templates, err := templates.ParseFiles("../client/templates/base.gohtml", "../client/templates/"+templateName+".gohtml")
	if err != nil {
		return err
	}

	templates, err = templates.ParseGlob("../client/templates/components/*.gohtml")
	if err != nil {
		return err
	}

	err = templates.ExecuteTemplate(w, "base", data)
	if err != nil {
		return err
	}

	_ = templateName
	return nil
}
