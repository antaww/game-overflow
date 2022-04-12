package routes

import (
	"main/sql"
	"net/http"
)

func LoginUser(r *http.Request) (*sql.User, error) {
	cookie, err := r.Cookie("session")

	if err != nil {
		TemplatesData.ConnectedUser = nil
	} else {
		user, err := sql.GetUserBySession(cookie.Value)
		if err != nil {
			return nil, err
		}
		TemplatesData.ConnectedUser = user

		err = sql.SetUserOnline(TemplatesData.ConnectedUser.Id, true)
		if err != nil {
			return nil, err
		}

		return user, nil
	}

	return nil, nil
}
