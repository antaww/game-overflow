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


func SaveUser(user User) (int64, error) {

}

