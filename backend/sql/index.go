package sql

import (
	"database/sql"
	"fmt"
	"log"
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
	if err != nil {
		return fmt.Errorf("SaveUser error: %v", err)
	}
	fmt.Println(oldUsername, "nick =>", newUsername)
	return nil
}

func SessionID(user User, w http.ResponseWriter) {
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
		MaxAge:   1000,
	}
	http.SetCookie(w, cookie)

	_, err := DB.Exec("INSERT INTO sessions VALUES (?, ?)", session.IdSession, session.IdUser)
	if err != nil {
		log.Fatal(err)
	}
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
	if err != nil {
		return err
	}
	return nil
}
