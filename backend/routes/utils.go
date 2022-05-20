package routes

import (
	"main/sql"
	"net/http"
)

func GetTemplatesDataFromRoute(w http.ResponseWriter, r *http.Request) (*TemplatesDataType, error) {
	connectedUser, err := sql.GetUserByRequest(r)
	if err != nil {
		return nil, err
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
