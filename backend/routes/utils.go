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

	return &TemplatesDataType{
		ConnectedUser: connectedUser,
		FeedSort:      sql.FeedSortNewest,
		Locales:       locales,
		ShownTopics:   nil,
		ShownTopic:    shownTopic,
		ShownMessages: nil,
		ShownUser:     nil,
	}, nil
}
