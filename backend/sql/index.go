package sql

import (
	"database/sql"
	"fmt"
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
