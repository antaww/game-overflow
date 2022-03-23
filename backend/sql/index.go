package sql

import (
	"database/sql"
	"fmt"
	"log"
	"main/utils"
	"time"
)

var DB *sql.DB

type User struct {
	ID           int64
	Username     string
	IsOnline     bool
	Password     string
	Email        string
	Locale       string
	ProfilePic   string
	Description  string
	CreationDate time.Time
	Role         Role
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

func EditUsername(user User, newUsername string) error {
	_, err := DB.Exec("UPDATE users SET username = (?) WHERE username = (?)", newUsername, user.Username)
	if err != nil {
		return fmt.Errorf("SaveUser error: %v", err)
	}
	fmt.Println(user.Username, "nick =>", newUsername)
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
