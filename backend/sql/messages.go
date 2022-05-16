package sql

import (
	"fmt"
	"main/utils"
	"strconv"
	"time"
)

type Message struct {
	Id        int64     `db:"id_message"`
	Content   string    `db:"content"`
	CreatedAt time.Time `db:"created_at"`
	IdTopic   string    `db:"id_topic"`
	AuthorId  int64     `db:"id_user" `
	Likes     []MessageLike
}

type MessageWithConnectedUser struct {
	*Message
	ConnectedUser *User
}

type MessageLike struct {
	IdMessage int64 `db:"id_message"`
	IdUser    int64 `db:"id_user"`
	IsLike    bool  `db:"like"`
}

// CalculatePoints CalculateLikes returns the number of likes for a message
func (message *Message) CalculatePoints() int {
	err := message.FetchLikes()
	if err != nil {
		return 0
	}

	var count int
	for _, like := range message.Likes {
		if like.IsLike {
			count++
		} else {
			count--
		}
	}

	return count
}

// GetUser returns the user who wrote the message
func (message *Message) GetUser() (*User, error) {
	return GetUserById(message.AuthorId)
}

// FetchLikes FetchMessages get messages into topic from db using post id
func (message *Message) FetchLikes() error {
	messageLike, err := GetLikes(message.Id)
	if err != nil {
		return err
	}

	message.Likes = messageLike
	return nil
}

// WithConnectedUser returns a message with the connected user
func (message *Message) WithConnectedUser(user *User) MessageWithConnectedUser {
	return MessageWithConnectedUser{
		Message:       message,
		ConnectedUser: user,
	}
}

// AddMessage add a message into a topic
func AddMessage(idUser int64, idTopic int64, message string) (int64, error) {
	id, err := utils.GenerateID()
	if err != nil {
		return 0, fmt.Errorf("error while generating id: %v", err)
	}
	_, err = DB.Exec("INSERT INTO messages VALUES (?, ?, ?, ?, ?)", id, message, time.Now(), idTopic, idUser)
	if err != nil {
		return 0, fmt.Errorf("CreateTopic error: %v", err)
	}
	fmt.Printf("message added to the topic %v\n", strconv.FormatInt(idTopic, 10))
	return id, nil
}

// DeleteMessage delete a message from a topic
func DeleteMessage(idMessage int64) error {
	_, err := DB.Exec("UPDATE messages SET content = 'Message deleted' WHERE id_message = ?", idMessage)
	if err != nil {
		return fmt.Errorf("DeleteMessage error: %v", err)
	}
	return nil
}

// DeleteDislikeMessage delete a dislike from a message
func DeleteDislikeMessage(messageId, userId int64) (bool, error) {
	_, err := DB.Exec("DELETE FROM message_like WHERE id_message = ? AND id_user = ?", messageId, userId)
	if err != nil {
		return false, fmt.Errorf("DeleteDislikeMessage error: %v", err)
	}

	return true, nil
}

// DeleteLikeMessage delete a like from a message
func DeleteLikeMessage(messageId, userId int64) (bool, error) {
	_, err := DB.Exec("DELETE FROM message_like WHERE id_message = ? AND id_user = ?", messageId, userId)
	if err != nil {
		return false, fmt.Errorf("DeleteLikeMessage error: %v", err)
	}

	return true, nil
}

// DislikeMessage dislike a message
func DislikeMessage(messageId, userId int64) (bool, error) {
	_, err := DB.Exec("INSERT INTO message_like VALUES (?, ?, ?)", messageId, userId, false)
	if err != nil {
		return false, fmt.Errorf("DislikeMessage error: %v", err)
	}

	return true, nil
}

func EditMessage(messageId int64, message string) error {
	_, err := DB.Exec("UPDATE messages SET content = ? WHERE id_message = ?", message, messageId)
	if err != nil {
		return fmt.Errorf("EditMessage error: %v", err)
	}
	return nil
}

// GetMessages returns all messages from a topic id
func GetMessages(postId int64) ([]Message, error) {
	rows, err := DB.Query("SELECT * FROM messages WHERE id_topic = ? ORDER BY created_at", postId)
	if err != nil {
		return nil, err
	}

	var messages []Message
	for rows.Next() {
		var message Message
		err = rows.Scan(&message.Id, &message.Content, &message.CreatedAt, &message.IdTopic, &message.AuthorId)
		if err != nil {
			return nil, err
		}

		messages = append(messages, message)
	}

	HandleSQLErrors(rows)

	return messages, nil
}

// GetMessage returns a message from a message id
func GetMessage(messageId int64) (*Message, error) {
	var message Message
	row := DB.QueryRow("SELECT * FROM messages WHERE id_message = ?", messageId)
	err := row.Scan(&message.Id, &message.Content, &message.CreatedAt, &message.IdTopic, &message.AuthorId)
	if err != nil {
		return nil, err
	}
	return &message, nil
}

// GetLikes returns all likes from a message id
func GetLikes(messageId int64) ([]MessageLike, error) {
	rows, err := DB.Query("SELECT * FROM message_like WHERE id_message = ?", messageId)
	if err != nil {
		return nil, err
	}

	var likes []MessageLike
	for rows.Next() {
		var like MessageLike
		err = rows.Scan(&like.IdMessage, &like.IdUser, &like.IsLike)
		if err != nil {
			return nil, err
		}

		likes = append(likes, like)
	}

	HandleSQLErrors(rows)

	return likes, nil
}

// GetTags returns all tags from a topic id
func GetTags(topicId int64) ([]string, error) {
	rows, err := DB.Query("SELECT tag_name FROM have WHERE id_topic = ?", topicId)
	if err != nil {
		return nil, err
	}

	var tags []string
	for rows.Next() {
		var tag string
		err = rows.Scan(&tag)
		if err != nil {
			return nil, err
		}

		tags = append(tags, tag)
	}

	HandleSQLErrors(rows)

	return tags, nil
}

// LikeMessage like a message
func LikeMessage(messageId, userId int64) (bool, error) {
	_, err := DB.Exec("INSERT INTO message_like VALUES (?, ?, ?)", messageId, userId, true)
	if err != nil {
		return false, fmt.Errorf("LikeMessage error: %v", err)
	}

	return true, nil
}

// MessageGetLikeFrom returns a like from a message id and user id
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

func GetAllTags() ([]Tags, error) {
	rows, err := DB.Query("SELECT tag_name FROM have")
	if err != nil {
		return nil, fmt.Errorf("GetAllTags error: %v", err)
	}
	//save all tags in a slice
	var tags []Tags
	for rows.Next() {
		var tag string
		err = rows.Scan(&tag)
		if err != nil {
			return nil, fmt.Errorf("GetAllTags error: %v", err)
		}
		if !contains(tags, Tags{tag}) {
			tags = append(tags, Tags{tag})
		}
	}
	HandleSQLErrors(rows)
	return tags, nil
}

func (message *Message) IsLiked(userId int64) (bool, error) {
	like, err := MessageGetLikeFrom(message.Id, userId)
	if err != nil {
		return false, err
	}
	if like != nil {
		return like.IsLike, nil
	}
	return false, nil
}

func (message *Message) IsDisliked(userId int64) (bool, error) {
	like, err := MessageGetLikeFrom(message.Id, userId)
	if err != nil {
		return false, err
	}
	if like != nil {
		return !like.IsLike, nil
	}
	return false, nil
}

