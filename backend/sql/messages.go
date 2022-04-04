package sql

import (
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

type MessageLike struct {
	idMessage int64 `db:"id_message"`
	idUser    int64 `db:"id_user"`
	like      bool  `db:"like"`
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

func (message *Message) GetUser() *User {
	return GetUserById(message.AuthorId)
}

func GetLikes(messageId int64) ([]MessageLike, error) {
	rows, err := DB.Query("SELECT * FROM message_like WHERE id_message = ?", messageId)
	if err != nil {
		return nil, err
	}

	var likes []MessageLike
	for rows.Next() {
		var like MessageLike
		err = rows.Scan(&like.idMessage, &like.idUser, &like.like)
		if err != nil {
			return nil, err
		}

		likes = append(likes, like)
	}

	HandleSQLErrors(rows)

	return likes, nil
}

// FetchMessages get messages into topic from db using post id
func (message *Message) FetchLikes() error {
	messageLike, err := GetLikes(message.Id)
	if err != nil {
		return err
	}

	message.Likes = messageLike
	return nil
}

// CalculateLikes returns the number of likes for a message
func (message *Message) CalculatePoints() int {
	err := message.FetchLikes()
	if err != nil {
		return 0
	}

	var count int
	for _, like := range message.Likes {
		if like.like {
			count++
		} else {
			count--
		}
	}

	return count
}
