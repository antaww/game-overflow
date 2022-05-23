package sql

import (
	"database/sql"
	"main/utils"
	"net/http"
	"strconv"
	"time"
)

type LoginSession struct {
	IdUser    string `db:"id_user"`
	IdSession string `db:"id_session"`
}

var DB *sql.DB

func AdminEditUsername(oldUsername string, newUsername string) error {
	_, err := DB.Exec("UPDATE users SET username = (?) WHERE username = (?)", newUsername, oldUsername)
	return err
}

func AddSession(user *User) (string, error) {
	sessionId := utils.RandomString(32)
	session := LoginSession{
		IdSession: sessionId,
		IdUser:    strconv.FormatInt(user.Id, 10),
	}

	_, err := DB.Exec("INSERT INTO sessions VALUES (?, ?)", session.IdSession, session.IdUser)
	return sessionId, err
}

func AddSessionCookie(user *User, w http.ResponseWriter) error {
	sessionId, err := AddSession(user)
	if err != nil {
		return err
	}

	cookie := &http.Cookie{
		HttpOnly: false,
		Name:     "session",
		Value:    sessionId,
		MaxAge:   7 * 24 * 60 * 60,
	}
	http.SetCookie(w, cookie)

	return err
}

func DeleteSessionCookie(sessionId string, w http.ResponseWriter) error {
	cookie := &http.Cookie{
		HttpOnly: true,
		Name:     "session",
		Value:    "",
		Expires:  time.Unix(0, 0),
		MaxAge:   -1,
	}
	http.SetCookie(w, cookie)

	_, err := DB.Exec("DELETE FROM sessions WHERE sessions.id_session = ?", sessionId)
	return err
}

type FeedSortType string

const (
	FeedSortNewest  FeedSortType = "newest"
	FeedSortOldest  FeedSortType = "oldest"
	FeedSortPopular FeedSortType = "views"
	FeedSortPoints  FeedSortType = "points"
)

// GetFeedSortingTypes returns all feed sorting types
func GetFeedSortingTypes() []FeedSortType {
	return []FeedSortType{FeedSortNewest, FeedSortOldest, FeedSortPopular, FeedSortPoints}
}
