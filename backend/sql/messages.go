package sql

import "time"

type Message struct {
	Id        int64     `db:"id_message"`
	Content   string    `db:"content"`
	CreatedAt time.Time `db:"created_at"`
	IsFirst   bool      `db:"is_first"`
	Likes     int       `db:"likes"`
	Dislikes  int       `db:"dislikes"`
	IdTopic   string    `db:"id_topic"`
	AuthorId  int64     `db:"id_user" `
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
		err = rows.Scan(&message.Id, &message.Content, &message.CreatedAt, &message.IsFirst, &message.Likes, &message.Dislikes, &message.IdTopic, &message.AuthorId)
		if err != nil {
			return nil, err
		}

		messages = append(messages, message)
	}

	HandleSQLErrors(rows)

	return messages, nil
}

func (message *Message) CalculatePoints() int {
	return message.Likes - message.Dislikes
}

func (message *Message) GetUser() *User {
	return GetUserById(message.AuthorId)
}
