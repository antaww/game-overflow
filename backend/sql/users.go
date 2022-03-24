package sql

import (
	"fmt"
	"log"
	"main/utils"
	"time"
)

type Role string

const (
	RoleAdmin     Role = "admin"
	RoleModerator Role = "moderator"
	RoleUser      Role = "user"
)

type User struct {
	Id           int64     `db:"id_user"`
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

// EditUsername edits the username of the user
func EditUsername(idUser int64, newUsername string) (bool, error) {
	_, err := DB.Exec("UPDATE users SET username = (?) WHERE id_user = (?)", newUsername, idUser)
	if err != nil {
		return false, fmt.Errorf("SaveUser error: %v", err)
	}

	fmt.Println(idUser, "nick =>", newUsername)
	return true, nil
}

func EditPassword(idUser int64, oldPassword string, newPassword string) (bool, error) {
	r, err := DB.Exec("UPDATE users SET password = (?) WHERE id_user = (?) AND password = (?)", newPassword, idUser, oldPassword)
	if err != nil {
		return false, fmt.Errorf("SaveUser error: %v", err)
	}
	affected, err := r.RowsAffected()
	if err != nil {
		return false, fmt.Errorf("SaveUser error: %v", err)
	}
	if affected > 0 {
		fmt.Println(idUser, "password : ", oldPassword, "=>", newPassword)
	} else {
		fmt.Println(idUser, "tried to change his password but `", oldPassword, "` is incorrect")
	}
	return true, nil
}

// CreateUser creates a new user with generated Id, creation date to now and locale to english
func CreateUser(username, password, email string) User {
	return User{
		Id:           utils.GenerateID(),
		Username:     username,
		Password:     password,
		Email:        email,
		Locale:       "en",
		CreationDate: time.Now(),
		Role:         RoleUser,
	}
}

// GetUserById finds a user by id, returns nil if not found
func GetUserById(id int64) *User {
	result, err := DB.Query("SELECT * FROM users WHERE id_user = ?", id)
	if err != nil {
		return nil
	}

	user := &User{}
	err = Results(result, &user.Id, &user.Username, &user.IsOnline, &user.Password, &user.Email, &user.Locale, &user.ProfilePic, &user.Description, &user.CreationDate, &user.Role)
	if err != nil {
		log.Fatal(err)
	}
	HandleSQLErrors(result)
	return user
}

// GetUserByUsername finds a user by username, returns nil if not found
func GetUserByUsername(username string) *User {
	result, err := DB.Query("SELECT * FROM users WHERE username = ?", username)
	if err != nil {
		log.Fatal(err)
	}

	user := &User{}
	err = Results(result, &user.Id, &user.Username, &user.IsOnline, &user.Password, &user.Email, &user.Locale, &user.ProfilePic, &user.Description, &user.CreationDate, &user.Role)
	if err != nil {
		log.Fatal(err)
	}
	HandleSQLErrors(result)
	return user
}

// LoginByIdentifiants logs in a user by username and password, return true if success
func LoginByIdentifiants(username, password string) (bool, error) {
	result, err := DB.Query("SELECT username, password FROM users WHERE username = ? AND password = ?", username, password)
	if err != nil {
		return false, fmt.Errorf("SaveUser error: %v", err)
	}

	if result.Next() {
		fmt.Println("result OK => login page loading")
		HandleSQLErrors(result)
		return true, nil
	} else {
		fmt.Println("result inexistant")
		HandleSQLErrors(result)
		return false, nil
	}
}

// LoginBySession logs in a user by sessionId (cookie), returns user found if success, else nil
func LoginBySession(sessionId string) (*User, error) {
	result, err := DB.Query("SELECT id_user FROM sessions WHERE id_session = ?", sessionId)
	if err != nil {
		return nil, fmt.Errorf("SaveUser error: %v", err)
	}

	var idUser int64
	err = Results(result, &idUser)
	if err != nil {
		return nil, fmt.Errorf("SaveUser error: %v", err)
	}

	HandleSQLErrors(result)

	return GetUserById(idUser), nil
}

// SaveUser saves a user in the database
func SaveUser(user User) (bool, error) {
	_, err := DB.Exec("INSERT INTO users VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)", user.Id, user.Username, user.IsOnline, user.Password, user.Email, user.Locale, user.ProfilePic, user.Description, user.CreationDate, user.Role)
	if err != nil {
		return false, fmt.Errorf("SaveUser error: %v", err)
	}

	return true, nil
}
