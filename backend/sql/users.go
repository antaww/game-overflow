package sql

import (
	"database/sql"
	"fmt"
	"main/utils"
	"net/http"
	"strings"
	"time"
)

type Role string

const (
	RoleAdmin     Role = "admin"
	RoleModerator Role = "moderator"
	RoleUser      Role = "user"
	RoleBan       Role = "banned"
)

type User struct {
	Id             int64          `db:"id_user" json:"id,omitempty"`
	Username       string         `db:"username" json:"username"`
	IsOnline       bool           `db:"is_online" json:"isOnline"`
	Password       string         `db:"password" json:"password,omitempty"`
	Email          string         `db:"email" json:"email,omitempty"`
	Locale         string         `db:"locale" json:"locale,omitempty"`
	ProfilePic     sql.NullString `db:"profile_pic" json:"profilePic,omitempty"`
	Description    sql.NullString `db:"description" json:"description,omitempty"`
	CreationDate   time.Time      `db:"created_at" json:"creationDate"`  //todo
	Role           Role           `db:"role_type" json:"role,omitempty"` //todo
	Color          int            `db:"color" json:"color,omitempty"`
	CookiesEnabled sql.NullBool   `db:"cookies_enabled" json:"cookiesEnabled"`
	DefaultColor   int
}

type UserWithConnectedUser struct {
	*User
	ConnectedUser *User
}

type Like struct {
	IdMessage int64 `db:"id_message" json:"idMessage"`
	IdUser    int64 `db:"id_user" json:"idUser"`
	Like      bool  `db:"like" json:"like"`
}

func (user *User) IsFollowingUser(userId int64) bool {
	result, err := DB.Query("SELECT * FROM follow WHERE id_user_followed = ? AND id_user_follower = ? LIMIT 1", userId, user.Id)
	if err != nil {
		utils.SQLError(err)
	}

	return result.Next()
}

func (user *User) DisplayRole() string {
	switch user.Role {
	case RoleAdmin:
		return "<i class=\"fa-solid fa-crown fa-fw\"></i>"
	case RoleModerator:
		return "<i class=\"fa-solid fa-gavel fa-fw\"></i>"
	case RoleUser:
		return "<i></i>"
	default:
		return ""
	}
}

func (user *User) CalculateDefaultColor() {
	if user.ProfilePic.Valid && user.ProfilePic.String != "" {
		img, err := utils.GetImageFromBase64(user.ProfilePic.String)
		if err != nil {
			utils.RouteError(err)
		}

		if img != nil {
			go func() {
				user.DefaultColor = utils.GetDominantColor(img)
			}()
		}
	} else {
		user.DefaultColor = 0xcccccc
	}
}

func (user *User) CountTopics() int {
	topic, err := GetUserTopics(user.Id)
	if err != nil {
		return 0
	}

	return len(topic)

}

func (user *User) GetFollowers() []User {
	followers, err := GetFollowers(user.Id)
	if err != nil {
		utils.RouteError(err)
	}

	return followers
}

func (user *User) GetFollowing() []User {
	following, err := GetFollowing(user.Id)
	if err != nil {
		utils.RouteError(err)
	}

	return following
}

func (user *User) GetLinks() []Link {
	links, err := GetLinks(user.Id)
	if err != nil {
		utils.RouteError(err)
	}

	return links
}

func (user *User) GetTopics() []Topic {
	topic, err := GetUserTopics(user.Id)
	if err != nil {
		return nil
	}

	return topic
}

func (user *User) GetHasCookieEnabled() sql.NullBool {
	enabled, err := GetUserHasCookiesEnabled(user.Id)
	if err != nil {
		utils.SQLError(err)
	}

	return enabled
}

func (user *User) SetCookiesEnabled(enabled bool) error {
	nullBool := sql.NullBool{Bool: enabled, Valid: true}
	err := SetUserCookiesEnabled(user.Id, nullBool)
	if err != nil {
		return err
	}

	user.CookiesEnabled = nullBool
	return nil
}

func (user *User) WithConnectedUser(connectedUser *User) *UserWithConnectedUser {
	return &UserWithConnectedUser{
		ConnectedUser: connectedUser,
		User:          user,
	}
}

// ConfirmPassword checks if the password is correct
func ConfirmPassword(userId int64, password string) bool {
	var user User
	rows, err := DB.Query("SELECT password FROM users WHERE id_user = ?", userId)
	if err != nil {
		utils.SQLError(err)
		return false
	}
	err = Results(rows, &user.Password)
	if err != nil {
		utils.SQLError(err)
		return false
	}

	return password == user.Password
}

// CreateUser creates a new user with generated id, creation date to now and locale to english
func CreateUser(username, password, email string) User {
	id, err := utils.GenerateID()
	if err != nil {
		utils.SQLError(err)
	}

	return User{
		Id:           id,
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
// returns true if the user was edited successfully
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
	if newUser.ProfilePic.String != "" {
		requestEdits = append(requestEdits, "profile_pic = ?")
		arguments = append(arguments, newUser.ProfilePic)
	}
	if newUser.Description.String != "" {
		requestEdits = append(requestEdits, "description = ?")
		arguments = append(arguments, newUser.Description.String)
	}
	if newUser.Username != "" {
		requestEdits = append(requestEdits, "username = ?")
		arguments = append(arguments, newUser.Username)
	}
	if newUser.Color != 0 {
		requestEdits = append(requestEdits, "color = ?")
		arguments = append(arguments, newUser.Color)
	}
	if newUser.CookiesEnabled.Valid {
		requestEdits = append(requestEdits, "cookies_enabled = ?")
		arguments = append(arguments, newUser.CookiesEnabled.Bool)
	}

	if requestEdits == nil {
		return false, nil
	}
	request += strings.Join(requestEdits, ", ") + " WHERE id_user = ?"
	arguments = append(arguments, idUser)

	r, err := DB.Exec(request, arguments...)
	if err != nil {
		return false, fmt.Errorf("EditUser error: %v", err)
	}
	_, err = r.RowsAffected()
	if err != nil {
		return false, fmt.Errorf("EditUser error: %v", err)
	}

	return true, nil
}

// GetAllUsers returns all users
func GetAllUsers() ([]User, error) {
	rows, err := DB.Query("SELECT * FROM users")
	if err != nil {
		return nil, err
	}

	return GetUsers(rows)
}

// GetFollowers returns the followers of the user with the given id
func GetFollowers(idUser int64) ([]User, error) {
	rows, err := DB.Query("SELECT * from users where id_user in (select id_user_follower from follow where id_user_followed = ?)", idUser)
	if err != nil {
		return nil, fmt.Errorf("GetFollowers error: %v", err)
	}
	return GetUsers(rows)
}

// GetFollowing returns the users followed by the user with the given id
func GetFollowing(idUser int64) ([]User, error) {
	rows, err := DB.Query("SELECT * from users where id_user in (select id_user_followed from follow where id_user_follower = ?)", idUser)
	if err != nil {
		return nil, fmt.Errorf("GetFollowing error: %v", err)
	}
	return GetUsers(rows)
}

// GetSessionId returns the session id of the user with the given id
func GetSessionId(r *http.Request) (string, error) {
	ip := r.Header.Get("X-Forwarded-For")
	sessionWithoutCookie := utils.ConnectedUsersWithoutCookies[ip]
	if sessionWithoutCookie != "" {
		return sessionWithoutCookie, nil
	}

	sessionId, err := r.Cookie("session")
	if err == http.ErrNoCookie {
		return "", nil
	}

	if sessionId == nil {
		return "", nil
	}

	return sessionId.Value, err
}

// GetUserById finds a user by id, returns nil if not found
func GetUserById(id int64) (*User, error) {
	result, err := DB.Query("SELECT * FROM users WHERE id_user = ?", id)
	if err != nil {
		return nil, fmt.Errorf("GetUserById error: %v", err)
	}
	var profilePicture []byte

	user := &User{}
	err = Results(result, &user.Id, &user.Username, &user.IsOnline, &user.Password, &user.Email, &user.Locale, &profilePicture, &user.Description, &user.CreationDate, &user.Role, &user.Color, &user.CookiesEnabled)
	if err != nil {
		return nil, fmt.Errorf("GetUserById error: %v", err)
	}
	user.ProfilePic = sql.NullString{Valid: true, String: string(profilePicture)}

	HandleSQLErrors(result)
	return user, nil
}

// GetUserByRequest gets a user by request, returns nil if not found
func GetUserByRequest(r *http.Request) (*User, error) {
	cookie, err := r.Cookie("session")

	if err == http.ErrNoCookie {
		ip := r.Header.Get("X-Forwarded-For")
		connectedSessionWithoutCookie := utils.ConnectedUsersWithoutCookies[ip]
		if connectedSessionWithoutCookie != "" {
			return GetUserBySession(connectedSessionWithoutCookie)
		}

		return nil, err
	}

	if err != nil {
		return nil, fmt.Errorf("GetUserByRequest error: %v", err)
	}

	return GetUserBySession(cookie.Value)
}

// GetUserBySession logs in a user by sessionId (cookie), returns user found if success, else nil
func GetUserBySession(sessionId string) (*User, error) {
	result, err := DB.Query("SELECT id_user FROM sessions WHERE id_session = ?", sessionId)
	if err != nil {
		return nil, fmt.Errorf("GetUserBySession error: %v", err)
	}

	var idUser int64
	err = Results(result, &idUser)
	if err != nil {
		return nil, fmt.Errorf("GetUserBySession error: %v", err)
	}

	HandleSQLErrors(result)

	return GetUserById(idUser)
}

// GetUserByUsername finds a user by username, returns nil if not found
func GetUserByUsername(username string) (*User, error) {
	result, err := DB.Query("SELECT * FROM users WHERE username = ?", username)
	if err != nil {
		return nil, fmt.Errorf("GetUserByUsername error: %v", err)
	}

	var profilePicture []byte

	user := &User{}
	err = Results(result, &user.Id, &user.Username, &user.IsOnline, &user.Password, &user.Email, &user.Locale, &profilePicture, &user.Description, &user.CreationDate, &user.Role, &user.Color, &user.CookiesEnabled)
	if err != nil {
		return nil, fmt.Errorf("GetUserByUsername error: %v", err)
	}
	user.ProfilePic = sql.NullString{Valid: true, String: string(profilePicture)}

	HandleSQLErrors(result)
	return user, nil
}

// GetUserHasCookiesEnabled returns true if the user has cookies enabled
func GetUserHasCookiesEnabled(id int64) (sql.NullBool, error) {
	enabled := sql.NullBool{Bool: false, Valid: false}
	result, err := DB.Query("SELECT cookies_enabled FROM users WHERE id_user = ?", id)
	if err != nil {
		return enabled, fmt.Errorf("GetUserHasCookiesEnabled error: %v", err)
	}

	err = Results(result, &enabled)
	if err != nil {
		return enabled, fmt.Errorf("GetUserHasCookiesEnabled error: %v", err)
	}

	HandleSQLErrors(result)
	return enabled, nil
}

// GetUsers get all users from the database
func GetUsers(rows *sql.Rows) ([]User, error) {
	var users []User
	for rows.Next() {
		var profilePicture []byte
		user := &User{}
		err := rows.Scan(&user.Id, &user.Username, &user.IsOnline, &user.Password, &user.Email, &user.Locale, &profilePicture, &user.Description, &user.CreationDate, &user.Role, &user.Color, &user.CookiesEnabled)
		if err != nil {
			return nil, fmt.Errorf("GetUsers error: %v", err)
		}
		user.ProfilePic.String = string(profilePicture)
		users = append(users, *user)
	}

	return users, nil
}

// GetUsersStatus returns an array of users with their status
func GetUsersStatus(users []string) ([]*User, error) {
	var usersOnline []*User
	for _, user := range users {
		result, err := DB.Query("SELECT is_online, username FROM users WHERE username = ?", user)
		if err != nil {
			return nil, fmt.Errorf("GetUsersStatus error: %v", err)
		}
		if result.Next() {
			user := &User{}
			err = result.Scan(&user.IsOnline, &user.Username)
			HandleSQLErrors(result)
			usersOnline = append(usersOnline, user)
		} else {
			HandleSQLErrors(result)
		}
	}
	return usersOnline, nil
}

// GetUserTopics returns an array of topics of a user
func GetUserTopics(id int64) ([]Topic, error) {
	var topics []Topic
	result, err := DB.Query("SELECT * FROM topics WHERE id_first_message in (SELECT id_message FROM messages WHERE id_user = ?)", id)
	if err != nil {
		return nil, fmt.Errorf("GetUserTopics error: %v", err)
	}

	for result.Next() {
		topic := Topic{}
		err = result.Scan(&topic.Id, &topic.Title, &topic.IsClosed, &topic.Views, &topic.Category, &topic.IdFirstMessage)
		topics = append(topics, topic)
	}
	HandleSQLErrors(result)

	return topics, nil
}

func GetUserLikedMessages(id int64) ([]Message, error) {
	var messages []Message
	result, err := DB.Query("SELECT * FROM messages WHERE id_message in (SELECT id_message FROM message_like WHERE id_user = ? and `like` = 1) AND id_message not in (SELECT id_first_message FROM topics)", id)
	if err != nil {
		return nil, fmt.Errorf("GetUserLikes error: %v", err)
	}

	for result.Next() {
		message := Message{}
		err = result.Scan(&message.Id, &message.Content, &message.CreatedAt, &message.IdTopic, &message.AuthorId)
		messages = append(messages, message)
	}
	HandleSQLErrors(result)

	return messages, nil
}

func GetUserLikesTopics(id int64) ([]Topic, error) {
	var topics []Topic
	result, err := DB.Query("SELECT * FROM topics WHERE id_first_message in (SELECT id_message FROM message_like WHERE id_user = ? AND `like` = 1)", id)
	if err != nil {
		return nil, fmt.Errorf("GetUserLikesTopics error: %v", err)
	}

	for result.Next() {
		topic := Topic{}
		err = result.Scan(&topic.Id, &topic.Title, &topic.IsClosed, &topic.Views, &topic.Category, &topic.IdFirstMessage)
		topics = append(topics, topic)
	}
	HandleSQLErrors(result)

	return topics, nil
}

// LoginByIdentifiants logs in a user by username and password, return true if success
func LoginByIdentifiants(username, password string) (bool, error) {
	result, err := DB.Query("SELECT username, password FROM users WHERE username = ? AND password = ?", username, password)
	if err != nil {
		return false, fmt.Errorf("LoginByIdentifiants error: %v", err)
	}

	valid := result.Next()
	HandleSQLErrors(result)

	return valid, nil
}

// SaveUser saves a user in the database
func SaveUser(user User) (bool, error) {
	_, err := DB.Exec("INSERT INTO users VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)", user.Id, user.Username, user.IsOnline, user.Password, user.Email, user.Locale, user.ProfilePic, user.Description, user.CreationDate, user.Role, user.Color, user.CookiesEnabled)
	if err != nil {
		return false, fmt.Errorf("SaveUser error: %v", err)
	}

	return true, nil
}

// SetUserOnline sets a user online
func SetUserOnline(idUser int64, isOnline bool) error {
	_, err := DB.Exec("UPDATE users SET is_online = ? WHERE id_user = ?", isOnline, idUser)
	if err != nil {
		return fmt.Errorf("SetUserOnline error: %v", err)
	}
	return nil
}

// SetAllUsersOffline sets all users offline
func SetAllUsersOffline() error {
	_, err := DB.Exec("UPDATE users SET is_online = ?", false)
	if err != nil {
		return fmt.Errorf("SetAllUsersOffline error: %v", err)
	}
	utils.Logger.Printf("%v All users have been set offline", utils.Reset)
	return nil
}

// SetUserCookiesEnabled sets a user's cookies enabled
func SetUserCookiesEnabled(idUser int64, cookiesEnabled sql.NullBool) error {
	_, err := DB.Exec("UPDATE users SET cookies_enabled = ? WHERE id_user = ?", cookiesEnabled, idUser)
	if err != nil {
		return fmt.Errorf("SetUserCookiesEnabled error: %v", err)
	}
	return nil
}

func FollowUser(idUserFollowed, idUserFollower int64) error {
	result, err := DB.Exec("INSERT INTO follow VALUES (?, ?)", idUserFollowed, idUserFollower)
	if err != nil {
		return fmt.Errorf("FollowUser error: %v", err)
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("FollowUser error: %v", err)
	}
	if affected == 0 {
		return fmt.Errorf("FollowUser error: no rows affected")
	}
	return nil
}

func UnfollowUser(idUserFollowed, idUserFollower int64) error {
	result, err := DB.Exec("DELETE FROM follow WHERE id_user_followed = ? AND id_user_follower = ?", idUserFollowed, idUserFollower)
	if err != nil {
		return fmt.Errorf("UnfollowUser error: %v", err)
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("UnfollowUser error: %v", err)
	}
	if affected == 0 {
		return fmt.Errorf("UnfollowUser error: no rows affected")
	}
	return nil
}

func BanUser(idUser int64) error {
	_, err := DB.Exec("UPDATE users SET role_type = 'banned', profile_pic = '', description = '' WHERE id_user = ?", idUser)
	if err != nil {
		return fmt.Errorf("BanUser error: %v", err)
	}

	_, err = DB.Exec("UPDATE users SET profile_pic = '', description = '' WHERE id_user = ?", idUser)
	if err != nil {
		return fmt.Errorf("BanUser error: %v", err)
	}
	return nil
}

func UnbanUser(idUser int64) error {
	_, err := DB.Exec("UPDATE users SET role_type = 'user' WHERE id_user = ?", idUser)
	if err != nil {
		return fmt.Errorf("Unban error: %v", err)
	}
	return nil
}
