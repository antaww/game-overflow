package sql

import (
	"fmt"
	"log"
	"main/utils"
	"time"
)

type Role string

const RoleUser Role = "user"

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

func NewUser(username, password, email string) User {
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

func SaveUser(user User) error {
	_, err := DB.Exec("INSERT INTO users VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)", user.Id, user.Username, user.IsOnline, user.Password, user.Email, user.Locale, user.ProfilePic, user.Description, user.CreationDate, user.Role)
	if err != nil {
		return fmt.Errorf("SaveUser error: %v", err)
	}

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
	result, err := Select("SELECT username, password FROM users WHERE username = ? AND password = ?", username, password)
	if err != nil {
		return false, fmt.Errorf("SaveUser error: %v", err)
	}

	if result.Next() {
		fmt.Println("result OK => login page loading")
		return true, nil
	} else {
		fmt.Println("result inexistant")
		return false, nil
	}
}

func GetUserById(id int64) *User {
	result, err := Select("SELECT * FROM users WHERE id_user = ?", id)
	if err != nil {
		log.Fatal(err)
	}

	user := &User{}
	err = Results(result, &user.Id, &user.Username, &user.IsOnline, &user.Password, &user.Email, &user.Locale, &user.ProfilePic, &user.Description, &user.CreationDate, &user.Role)
	if err != nil {
		log.Fatal(err)
	}
	return user
}

func GetUserByUsername(username string) *User {
	result, err := Select("SELECT * FROM users WHERE id_user = ?", username)
	if err != nil {
		log.Fatal(err)
	}

	user := &User{}
	err = Results(result, &user.Id, &user.Username, &user.IsOnline, &user.Password, &user.Email, &user.Locale, &user.ProfilePic, &user.Description, &user.CreationDate, &user.Role)
	if err != nil {
		log.Fatal(err)
	}
	return user
}

func UserLoginBySession(sessionId string) (*User, error) {
	result, err := Select("SELECT * FROM sessions WHERE id_session = ?", sessionId)
	if err != nil {
		return nil, fmt.Errorf("SaveUser error: %v", err)
	}

	var idUser int64
	err = Results(result, &idUser)
	if err != nil {
		return nil, fmt.Errorf("SaveUser error: %v", err)
	}

	return GetUserById(idUser), nil
}
