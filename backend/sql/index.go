package sql

import (
	"database/sql"
	"fmt"
	"log"
	"main/utils"
	"net/http"
	"strconv"
)

const (
	RoleAdmin     Role = "admin"
	RoleModerator Role = "moderator"
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
	session := LoginSession{IdUser: strconv.FormatInt(user.Id, 10), IdSession: SessionID}
	cookie1 := &http.Cookie{Name: "session", Value: session.IdSession, HttpOnly: false}
	http.SetCookie(w, cookie1)
	_, err := DB.Exec("INSERT INTO sessions VALUES (?, ?)", session.IdSession, session.IdUser)
	if err != nil {
		log.Fatal(err)
	}
}
