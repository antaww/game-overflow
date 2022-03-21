package sql

import (
	"database/sql"
	"time"
)

var db *sql.DB


type User struct {
	ID int64
	Username string
	Password string
	Email string
	Locale string
	ProfilePic string
	Description string
	Creation time.Time
	Role string
}


func SaveUser(SavedUser User) (int64, error) {
	result, err := db.Exec("INSERT INTO users VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)", SavedUser.ID, SavedUser.Username, SavedUser.Password, SavedUser.Email, SavedUser.Locale, SavedUser.ProfilePic, SavedUser.Description, SavedUser.Creation, SavedUser.Role)

}

