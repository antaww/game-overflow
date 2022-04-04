package sql

type Topic struct {
	Id             int64  `db:"id_topic"`
	Title          string `db:"title"`
	IsClosed       bool   `db:"is_closed"`
	Views          int    `db:"views"`
	Category       string `db:"category_name"`
	IdFirstMessage int64  `db:"id_first_message"`
	Messages       []Message
}

// GetPost returns topic by id
func GetPost(id int64) (*Topic, error) {
	rows, err := DB.Query("SELECT * FROM topics WHERE id_topic = ?", id)
	if err != nil {
		return nil, err
	}

	post := &Topic{}
	err = Results(rows, &post.Id, &post.Title, &post.IsClosed, &post.Views, &post.Category, &post.IdFirstMessage)
	if err != nil {
		return nil, err
	}

	HandleSQLErrors(rows)

	return post, nil
}

func (topic *Topic) GetFirstMessage() (*Message, error) {
	rows, err := DB.Query("SELECT * FROM messages WHERE id_message = ?", topic.IdFirstMessage)
	if err != nil {
		return nil, err
	}

	message := &Message{}
	err = Results(rows, &message.Id, &message.Content, &message.CreatedAt, &message.IdTopic, &message.AuthorId)
	if err != nil {
		return nil, err
	}

	HandleSQLErrors(rows)

	return message, nil
}

// FetchMessages get messages into topic from db using post id
func (topic *Topic) FetchMessages() error {
	message, err := GetMessages(topic.Id)
	if err != nil {
		return err
	}

	topic.Messages = message
	return nil
}

func GetTopicsByCategories(category string) ([]Topic, error) {
	rows, err := DB.Query("SELECT * FROM topics WHERE category_name = ?", category)
	if err != nil {
		return nil, err
	}
	var topics []Topic
	for rows.Next() {
		var topic Topic
		err = rows.Scan(&topic.Id, &topic.Title, &topic.IsClosed, &topic.Views, &topic.Category, &topic.IdFirstMessage)
		if err != nil {
			return nil, err
		}

		topics = append(topics, topic)
	}

	HandleSQLErrors(rows)

	return topics, nil
}
