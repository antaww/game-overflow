package routes

import (
	"log"
	"main/sql"
	"main/utils"
	"net/http"
	"regexp"
	"strings"
)

// IndexRoute is the route for the home page
func IndexRoute(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/")

	if path == "" {
		templateData, err := GetTemplateDataFromRoute(w, r)
		if err != nil {
			utils.RouteError(err)
		}

		topics, err := templateData.GetTopicsSortedBy(templateData.FeedSort, 20)
		if err != nil {
			utils.RouteError(err)
		}

		for i := 0; i < len(topics); i++ {
			err = topics[i].FetchTags()
			if err != nil {
				utils.RouteError(err)
			}
		}

		templateData.ShownTopics = topics

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
			if err != http.ErrNoCookie {
				utils.RouteError(err)
			}

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
