package sql

import (
	"database/sql"
	"fmt"
	"log"
	"main/utils"
	"net/http"
	"strconv"
	"strings"
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

// EditUser edits the user with the given id
// can edit Description, Locale, ProfilePic, Username, Email
func EditUser(idUser int64, newUser User) (bool, error) {
	request := "UPDATE users SET "
	var requestEdits []string
	var arguments []interface{}
	if newUser.Email != "" {
		requestEdits = append(requestEdits, "email = ?")
		arguments = append(arguments, newUser.Email)
	}
	if newUser.Locale != "" {
		requestEdits = append(requestEdits, "locale = ?")
		arguments = append(arguments, newUser.Locale)
	}
	if newUser.ProfilePic != "" {
		requestEdits = append(requestEdits, "profile_pic = ?")
		arguments = append(arguments, newUser.ProfilePic)
	}
	if newUser.Description != "" {
		requestEdits = append(requestEdits, "description = ?")
		arguments = append(arguments, newUser.Description)
	}
	if newUser.Username != "" {
		requestEdits = append(requestEdits, "username = ?")
		arguments = append(arguments, newUser.Username)
	}

	if requestEdits == nil {
		return false, nil
	}
	request += strings.Join(requestEdits, ",") + " WHERE id_user = ?"
	arguments = append(arguments, idUser)

	fmt.Println(request)
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

func GetUserByRequest(r *http.Request) (*User, error) {
	cookie, err := r.Cookie("session")
	if err != nil {
		return nil, fmt.Errorf("GetUserByRequest error: %v", err)
	}

	return GetUserBySession(cookie.Value)
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

func DeleteDislikeMessage(messageId, userId int64) (bool, error) {
	_, err := DB.Exec("DELETE FROM message_like WHERE id_message = ? AND id_user = ?", messageId, userId)
	if err != nil {
		return false, fmt.Errorf("DeleteDislikeMessage error: %v", err)
	}

	return true, nil
}

func DeleteLikeMessage(messageId, userId int64) (bool, error) {
	_, err := DB.Exec("DELETE FROM message_like WHERE id_message = ? AND id_user = ?", messageId, userId)
	if err != nil {
		return false, fmt.Errorf("DeleteLikeMessage error: %v", err)
	}

	return true, nil
}

func DislikeMessage(messageId, userId int64) (bool, error) {
	_, err := DB.Exec("INSERT INTO message_like VALUES (?, ?, ?)", messageId, userId, false)
	if err != nil {
		return false, fmt.Errorf("DislikeMessage error: %v", err)
	}

	return true, nil
}

func LikeMessage(messageId, userId int64) (bool, error) {
	_, err := DB.Exec("INSERT INTO message_like VALUES (?, ?, ?)", messageId, userId, true)
	if err != nil {
		return false, fmt.Errorf("LikeMessage error: %v", err)
	}

	return true, nil
}

func MessageGetLikeFrom(messageId, userId int64) (*MessageLike, error) {
	result, err := DB.Query("SELECT * FROM message_like WHERE id_message = ? AND id_user = ?", messageId, userId)
	if err != nil {
		return nil, fmt.Errorf("MessageGetLikeFrom error: %v", err)
	}

	messageLike := &MessageLike{}
	if result.Next() {
		err = result.Scan(&messageLike.IdMessage, &messageLike.IdUser, &messageLike.IsLike)
		HandleSQLErrors(result)
		return messageLike, nil
	} else {
		HandleSQLErrors(result)
		return nil, nil
	}
}

func AddMessage(idUser int64, idTopic int64, message string) (int64, error) {
	id := utils.GenerateID()
	_, err := DB.Exec("INSERT INTO messages VALUES (?, ?, ?, ?, ?)", id, message, time.Now(), idTopic, idUser)
	if err != nil {
		return 0, fmt.Errorf("CreateTopic error: %v", err)
	}
	fmt.Printf("message added to the topic %v\n", strconv.FormatInt(idTopic, 10))
	return id, nil
}

func CreateTopic(title string, category string) (int64, error) {
	id := utils.GenerateID()
	_, err := DB.Exec("INSERT INTO topics VALUES (?, ?, ?, ?, ?, ?)", id, title, 0, 0, category, sql.NullInt64{})
	if err != nil {
		return 0, fmt.Errorf("CreateTopic error: %v", err)
	}
	fmt.Printf("topic '%v' added to the category '%v'", title, category)
	return id, nil
}

//create a function that will set the user online (is_online = 1)
func SetUserOnline(idUser int64, isOnline bool) error {
	_, err := DB.Exec("UPDATE users SET is_online = ? WHERE id_user = ?", isOnline, idUser)
	if err != nil {
		return fmt.Errorf("SaveUser error: %v", err)
	}
	return nil
}
