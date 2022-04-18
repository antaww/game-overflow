package routes

import (
	"log"
	"main/sql"
	"main/utils"
	"net/http"
	"regexp"
	"strings"
	"time"
)

type TemplatesDataType struct {
	ConnectedUser *sql.User
	Locales       map[string]string
	ShownTopics   []sql.Topic
	ShownTopic    sql.Topic
}

// GetCategories returns all categories
func (t TemplatesDataType) GetCategories() []sql.Category {
	categories, err := sql.GetCategories()
	if err != nil {
		log.Println(err)
	}
	return categories
}

var TemplatesData = TemplatesDataType{
	Locales: map[string]string{"en": "English", "fr": "Français"},
}

var PageLoadedTime time.Time

// IndexRoute is the route for the home page
func IndexRoute(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/")
	_, err := LoginUser(r)
	if err != nil {
		log.Fatal(err)
	}

	if path == "" {
		err := utils.CallTemplate("main", TemplatesData, w)
		if err != nil {
			log.Fatal(err)
		}
	}
}

// LogHandler is a middleware that logs the request and connects the user using the session cookie
func LogHandler(w http.ResponseWriter, r *http.Request) {
	matches, err := regexp.MatchString("\\.(css|png)$", r.URL.String())
	if err != nil {
		log.Fatal(err)
	}

	if !matches {
		log.Printf("%v %v", r.Method, r.URL.String())
		go func() {
			PageLoadedTime = time.Now()
		}()
	}
	if r.Method == "GET" {
		_, err := LoginUser(r)
		if err != nil {
			log.Fatal(err)
		}
	}

	http.DefaultServeMux.ServeHTTP(w, r)
}
