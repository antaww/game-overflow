package sql

import (
	"fmt"
	"log"
	"main/utils"
	"strconv"
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

func ConfirmPassword(userId int64, password string) bool {
	var user User
	rows, err := DB.Query("SELECT password FROM users WHERE id_user = ?", userId)
	if err != nil {
		log.Println(err)
		return false
	}
	err = Results(rows, &user.Password)
	if err != nil {
		log.Println(err)
		return false
	}

	return password == user.Password
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

// EditPassword edits the password of the user
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

// EditUser edits the user with the given id
// can edit Description, Locale, ProfilePic, Username, Email
func EditUser(idUser int64, newUser User) (bool, error) {
	request := "UPDATE users SET "
	var arguments []interface{}
	if newUser.Email != "" {
		request += "email = ?, "
		arguments = append(arguments, newUser.Email)
	}
	if newUser.Locale != "" {
		request += "locale = ?, "
		arguments = append(arguments, newUser.Locale)
	}
	if newUser.ProfilePic != "" {
		request += "profile_pic = ?, "
		arguments = append(arguments, newUser.ProfilePic)
	}
	if newUser.Description != "" {
		request += "description = ?, "
		arguments = append(arguments, newUser.Description)
	}
	if newUser.Username != "" {
		request += "username = ?"
		arguments = append(arguments, newUser.Username)
	}

	request += " WHERE id_user = ?"
	arguments = append(arguments, idUser)

	r, err := DB.Exec(request, arguments...)
	if err != nil {
		return false, fmt.Errorf("SaveUser error: %v", err)
	}
	_, err = r.RowsAffected()
	if err != nil {
		return false, fmt.Errorf("SaveUser error: %v", err)
	}

	return true, nil
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

// EditAvatar edits the username of the user
func EditAvatar(idUser int64, avatarUrl string) (bool, error) {
	_, err := DB.Exec("UPDATE users SET profile_pic = (?) WHERE id_user = (?)", avatarUrl, idUser)
	if err != nil {
		return false, fmt.Errorf("SaveUser error: %v", err)
	}

	fmt.Println(idUser, "avatar =>", avatarUrl)
	return true, nil
}

// GetUserById finds a user by id, returns nil if not found
func GetUserById(id int64) *User {
	result, err := DB.Query("SELECT * FROM users WHERE id_user = ?", id)
	if err != nil {
		return nil
	}
	var profilePicture []byte

	user := &User{}
	err = Results(result, &user.Id, &user.Username, &user.IsOnline, &user.Password, &user.Email, &user.Locale, &profilePicture, &user.Description, &user.CreationDate, &user.Role)
	if err != nil {
		log.Fatal(err)
	}
	user.ProfilePic = string(profilePicture)

	HandleSQLErrors(result)
	return user
}

// GetUserBySession logs in a user by sessionId (cookie), returns user found if success, else nil
func GetUserBySession(sessionId string) (*User, error) {
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

// GetUserByUsername finds a user by username, returns nil if not found
func GetUserByUsername(username string) *User {
	result, err := DB.Query("SELECT * FROM users WHERE username = ?", username)
	if err != nil {
		log.Fatal(err)
	}

	var profilePicture []byte

	user := &User{}
	err = Results(result, &user.Id, &user.Username, &user.IsOnline, &user.Password, &user.Email, &user.Locale, &profilePicture, &user.Description, &user.CreationDate, &user.Role)
	if err != nil {
		log.Fatal(err)
	}
	user.ProfilePic = string(profilePicture)

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
		fmt.Println("result not found")
		HandleSQLErrors(result)
		return false, nil
	}
}

// SaveUser saves a user in the database
func SaveUser(user User) (bool, error) {
	_, err := DB.Exec("INSERT INTO users VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)", user.Id, user.Username, user.IsOnline, user.Password, user.Email, user.Locale, user.ProfilePic, user.Description, user.CreationDate, user.Role)
	if err != nil {
		return false, fmt.Errorf("SaveUser error: %v", err)
	}

	return true, nil
}

func LikeMessage(idMessage int64) (bool, error) {
	_, err := DB.Exec("UPDATE messages SET likes = likes + 1 WHERE id_message = ?", idMessage)
	if err != nil {
		return false, fmt.Errorf("SaveUser error: %v", err)
	}

	return true, nil
}

func DislikeMessage(idMessage int64) (bool, error) {
	_, err := DB.Exec("UPDATE messages SET dislikes = dislikes + 1 WHERE id_message = ?", idMessage)
	if err != nil {
		return false, fmt.Errorf("SaveUser error: %v", err)
	}

	return true, nil
}

func AddMessage(idUser int64, idTopic int64, message string) error {
	_, err := DB.Exec("INSERT INTO messages VALUES (?, ?, ?, ?, ?, ?, ?)", utils.GenerateID(), message, time.Now(), 0, 0, idTopic, idUser)
	if err != nil {
		return fmt.Errorf("SaveUser error: %v", err)
	}
	fmt.Printf("message added to the topic %v\n", strconv.FormatInt(idTopic, 10))
	return nil
}

func CreateTopic(title string, category string) error {
	_, err := DB.Exec("INSERT INTO topics VALUES (?, ?, ?, ?, ?, ?)", utils.GenerateID(), title, 0, 0, category, nil)
	if err != nil {
		return fmt.Errorf("SaveUser error: %v", err)
	}
	fmt.Printf("topic '%v' added to the category '%v'", title, category)
	return nil
}
