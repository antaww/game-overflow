package routes

import (
	"log"
	"main/sql"
	"main/utils"
	"net/http"
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
	Session       string
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

// GetTags returns all tags
func (t TemplatesDataType) GetTags() []sql.Tags {
	tags, err := sql.GetAllTags()
	if err != nil {
		utils.RouteError(err)
	}

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

// GetTopicsSortedBy returns all topics depending on the sort type, limited by limit
func (t TemplatesDataType) GetTopicsSortedBy(sortType sql.FeedSortType, limit int) ([]sql.Topic, error) {
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

func (t TemplatesDataType) SortTopics() {
	t.ShownTopics.SortBy(t.FeedSort)
}

func GetTemplatesDataFromRoute(w http.ResponseWriter, r *http.Request) (*TemplatesDataType, error) {
	connectedUser, err := sql.GetUserByRequest(r)
	if err != nil {
		if err != http.ErrNoCookie {
			return nil, err
		}
	}

	locales, err := sql.GetLocales()
	if err != nil {
		return nil, err
	}

	shownTopic, err := sql.GetShownTopic(r)
	if err != nil {
		return nil, err
	}

	queries := r.URL.Query()

	feedSort := sql.FeedSortNewest
	if queries.Has("s") {
		sortType := queries.Get("s")

		sortTypes := sql.GetFeedSortingTypes()

		var isValid bool
		for _, sortType := range sortTypes {
			if sortType == sortType {
				isValid = true
				break
			}
		}

		if isValid {
			feedSort = sql.FeedSortType(sortType)
		}
	}

	return &TemplatesDataType{
		ConnectedUser: connectedUser,
		FeedSort:      feedSort,
		Locales:       locales,
		ShownTopics:   nil,
		ShownTopic:    shownTopic,
		ShownMessages: nil,
		ShownUser:     nil,
	}, nil
}
