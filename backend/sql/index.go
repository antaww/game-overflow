package sql

import (
	"database/sql"
	"fmt"
	"time"
)

var DB *sql.DB


type User struct {
	ID int64
	Username string
	IsOnline bool
	Password string
	Email string
	Locale string
	ProfilePic string
	Description string
	CreationDate time.Time
	Role string
}


func SaveUser(SavedUser User) (int64, error) {
	result, err := DB.Exec("INSERT INTO users VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)", SavedUser.ID,  SavedUser.Username, SavedUser.IsOnline, SavedUser.Password, SavedUser.Email, SavedUser.Locale, SavedUser.ProfilePic, SavedUser.Description, SavedUser.CreationDate, SavedUser.Role)
	if err != nil {
		return 0, fmt.Errorf("SaveUser error: %v", err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("SaveUser error: %v", err)
	}
	return id, nil
}

