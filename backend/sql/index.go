package sql

import (
	"database/sql"
	"fmt"
	"main/utils"
	"net/http"
	"strconv"
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

func AddSessionCookie(user User, w http.ResponseWriter) error {
	SessionID := utils.RandomString(32)
	fmt.Println("session id:", SessionID)

	session := LoginSession{
		IdSession: SessionID,
		IdUser:    strconv.FormatInt(user.Id, 10),
	}
	cookie := &http.Cookie{
		HttpOnly: false,
		Name:     "session",
		Value:    session.IdSession,
		MaxAge:   7 * 24 * 60 * 60,
	}
	http.SetCookie(w, cookie)

	_, err := DB.Exec("INSERT INTO sessions VALUES (?, ?)", session.IdSession, session.IdUser)
	return err
}

func CookieLogout(getCookie http.Cookie, w http.ResponseWriter) error {
	cookie := &http.Cookie{
		HttpOnly: false,
		Name:     "session",
		Value:    "",
		MaxAge:   -1,
	}
	http.SetCookie(w, cookie)

	_, err := DB.Exec("DELETE FROM sessions WHERE sessions.id_session = (?)", getCookie.Value)
	return err
}
