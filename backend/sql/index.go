package sql

import (
	"database/sql"
	"fmt"
	"log"
	"main/utils"
	"net/http"
	"strconv"
	"time"
)

var DB *sql.DB

type User struct {
	ID           int64     `db:"id_user"`
	Username     string    `db:"username"`
	IsOnline     bool      `db:"is_online"`
	Password     string    `db:"password"`
	Email        string    `db:"email"`
	Locale       string    `db:"locale"`
	ProfilePic   string    `db:"profile_pic"`
	Description  string    `db:"description"`
	CreationDate time.Time `db:"created_at"`
	Role         Role      `db:"role_type"`
}

type LoginSession struct {
	IdUser    string `db:"id_user"`
	IdSession string `db:"id_session"`
}

type Role string

const (
	RoleAdmin     Role = "admin"
	RoleModerator Role = "moderator"
	RoleUser      Role = "user"
)

func NewUser(username, password, email string) User {
	return User{
		ID:           utils.GenerateID(),
		Username:     username,
		Password:     password,
		Email:        email,
		Locale:       "en",
		CreationDate: time.Now(),
		Role:         RoleUser,
	}
}

func SaveUser(user User) error {
	_, err := DB.Exec("INSERT INTO users VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)", user.ID, user.Username, user.IsOnline, user.Password, user.Email, user.Locale, user.ProfilePic, user.Description, user.CreationDate, user.Role)
	if err != nil {
		return fmt.Errorf("SaveUser error: %v", err)
	}
	return nil
}

func AdminEditUsername(oldUsername string, newUsername string) error {
	_, err := DB.Exec("UPDATE users SET username = (?) WHERE username = (?)", newUsername, oldUsername)
	if err != nil {
		return fmt.Errorf("SaveUser error: %v", err)
	}
	fmt.Println(oldUsername, "nick =>", newUsername)
	return nil
}

func EditUsername(idUser int64, newUsername string) error {
	_, err := DB.Exec("UPDATE users SET username = (?) WHERE id_user = (?)", newUsername, idUser)
	if err != nil {
		return fmt.Errorf("SaveUser error: %v", err)
	}
	fmt.Println(idUser, "nick =>", newUsername)
	return nil
}

func UserLogin(username string, password string) (bool, error) {
	fmt.Println("username :", username)
	fmt.Println("password : ", password)
	result, err := DB.Query("SELECT username, password FROM users WHERE username = ? AND password = ?", username, password)

	if err != nil {
		return false, fmt.Errorf("SaveUser error: %v", err)
	}

	defer func(result *sql.Rows) {
		err := result.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(result)

	err = result.Err()
	if err != nil {
		log.Fatal(err)
	}

	if result.Next() {
		fmt.Println("result OK => login page loading")
		return true, nil
	} else {
		fmt.Println("result inexistant")
		return false, nil
	}

	return true, nil
}

func SessionID(user User, w http.ResponseWriter) {
	SessionID := utils.RandomString(32)
	fmt.Println("session id:", SessionID)
	session := LoginSession{IdUser: strconv.FormatInt(user.ID, 10), IdSession: SessionID}
	cookie1 := &http.Cookie{Name: "session", Value: session.IdSession, HttpOnly: false}
	http.SetCookie(w, cookie1)
	_, err := DB.Exec("INSERT INTO sessions VALUES (?, ?)", session.IdSession, session.IdUser)
	if err != nil {
		log.Fatal(err)
	}
}

func GetUser(username string) *User {
	result, err := DB.Query("SELECT * FROM users WHERE username = ?", username)
	if err != nil {
		log.Fatal(err)
	}
	defer func(result *sql.Rows) {
		err := result.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(result)
	user := &User{}
	if result.Next() {
		err = result.Scan(
			&user.ID,
			&user.Username,
			&user.IsOnline,
			&user.Password,
			&user.Email,
			&user.Locale,
			&user.ProfilePic,
			&user.Description,
			&user.CreationDate,
			&user.Role,
		)
		if err != nil {
			log.Fatal(err)
		}
		return user
	}
	return nil
}
