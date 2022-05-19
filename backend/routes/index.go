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
	FeedSort      sql.FeedSortType
	GetAllTags    []string
	Locales       map[string]string
	ShownTopic    *sql.Topic
	ShownTopics   sql.Topics
	ShownMessages []sql.Message
	ShownUser     *sql.User
}

// GetCategories returns all categories
func (t TemplatesDataType) GetCategories() []sql.Category {
	categories, err := sql.GetCategories()
	if err != nil {
		log.Println(err)
	}
	return categories
}

// GetFeedSortingTypes returns all feed sorting types
func (t TemplatesDataType) GetFeedSortingTypes() []sql.FeedSortType {
	return []sql.FeedSortType{sql.FeedSortNewest, sql.FeedSortOldest, sql.FeedSortPopular, sql.FeedSortFollow}
}

// GetLocales returns all locales
func (t TemplatesDataType) GetLocales() map[string]string {
	locales, err := sql.GetLocales()
	if err != nil {
		utils.RouteError(err)
	}

	return locales
}

// GetTags returns all tags
func (t TemplatesDataType) GetTags() []sql.Tags {
	tags, err := sql.GetAllTags()
	if err != nil {
		utils.RouteError(err)
	}
	fmt.Println("tags:", tags)
	return tags
}

// GetTrendingTags returns the trending tags, limited by the limit
func (t TemplatesDataType) GetTrendingTags(limit int) []sql.TagListItem {
	tags, err := sql.GetTrendingTags(limit)
	if err != nil {
		utils.RouteError(err)
	}

	return tags
}

// GetTopicsDependingSort returns all topics depending on the sort type, limited by limit
func (t TemplatesDataType) GetTopicsDependingSort(sortType sql.FeedSortType, limit int) ([]sql.Topic, error) {
	switch sortType {
	case sql.FeedSortNewest:
		return sql.GetNewestTopics(limit)
	case sql.FeedSortOldest:
		return sql.GetOldestTopics(limit)
	case sql.FeedSortPopular:
		return sql.GetPopularTopics(limit)
		//case sql.FeedSortFollow:
		//	return sql.GetFollowedTopics(t.ConnectedUser.Id, limit)
	}

	return nil, nil
}

func (t TemplatesDataType) SortTopics(sortType string) {
	sortTypes := sql.GetFeedSortingTypes()

	var isValid bool
	for _, sortType := range sortTypes {
		if sortType == sortType {
			isValid = true
			break
		}
	}

	if isValid {
		t.FeedSort = sql.FeedSortType(sortType)
		t.ShownTopics.SortBy(t.FeedSort)
	}
}

// IndexRoute is the route for the home page
func IndexRoute(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/")

	if path == "" {
		templateData, err := GetTemplatesDataFromRoute(w, r)
		if err != nil {
			utils.RouteError(err)
		}

		topics, err := templateData.GetTopicsDependingSort(templateData.FeedSort, 20)
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

		queries := r.URL.Query()

		if queries.Has("s") {
			sortType := queries.Get("s")

			templateData.SortTopics(sortType)
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
