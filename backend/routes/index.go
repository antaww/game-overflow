package routes

import (
	"fmt"
	"log"
	"main/sql"
	"main/utils"
	"net/http"
	"regexp"
	"strings"
)

type TemplatesDataType struct {
	ConnectedUser *sql.User
	Locales       map[string]string
	ShownTopics   []sql.Topic
	ShownTopic    *sql.Topic
	ShownMessages []sql.Message
	ShownUser     *sql.User
	GetAllTags    []string
}

// GetCategories returns all categories
func (t TemplatesDataType) GetCategories() []sql.Category {
	categories, err := sql.GetCategories()
	if err != nil {
		log.Println(err)
	}
	return categories
}

func (t TemplatesDataType) GetLocales() map[string]string {
	locales, err := sql.GetLocales()
	if err != nil {
		utils.RouteError(err)
	}

	return locales
}

// IndexRoute is the route for the home page
func IndexRoute(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/")

	if path == "" {
		templateData, err := GetTemplatesDataFromRoute(w, r)
		if err != nil {
			utils.RouteError(err)
		}

		err = utils.CallTemplate("main", templateData, w)
		if err != nil {
			utils.RouteError(err)
		}
	}
}

// LogHandler is a middleware that logs the request and connects the user using the session cookie
func LogHandler(w http.ResponseWriter, r *http.Request) {
	matches, err := regexp.MatchString("\\.(css|png)$", r.URL.String())
	if err != nil {
		utils.RouteError(err)
	}

	if !matches {
		log.Printf("%v %v", r.Method, r.URL.RequestURI())
	}

	if r.Method == "GET" {
		user, err := sql.GetUserByRequest(r)
		if err != nil {
			utils.RouteError(err)
		}

		if user != nil {
			err = sql.SetUserOnline(user.Id, true)
			if err != nil {
				utils.RouteError(err)
			}
		}
	}

	http.DefaultServeMux.ServeHTTP(w, r)
}

func (t TemplatesDataType) GetTags() []sql.Tags {
	tags, err := sql.GetAllTags()
	if err != nil {
		log.Println(err)
	}
	fmt.Println("tags:", tags)
	return tags
}
