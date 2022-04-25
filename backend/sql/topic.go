package sql

import (
	"database/sql"
	"fmt"
	"main/utils"
)

type Topic struct {
	Id             int64  `db:"id_topic"`
	Title          string `db:"title"`
	IsClosed       bool   `db:"is_closed"`
	Views          int    `db:"views"`
	Category       string `db:"category_name"`
	IdFirstMessage int64  `db:"id_first_message"`
	Messages       []Message
	Tags           []string
}

// GetAnswersNumber returns the number of answers for a topic
func (topic Topic) GetAnswersNumber() int {
	err := topic.FetchMessages()
	if err != nil {
		return 0
	}
	return len(topic.Messages) - 1
}

// GetFirstMessage returns the first message of the topic
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

// FetchTags get tags into topic from db using post id
func (topic *Topic) FetchTags() error {
	tags, err := GetTags(topic.Id)
	if err != nil {
		return err
	}

	topic.Tags = tags
	return nil
}

// CreateTopic create a new topic in db
func CreateTopic(title string, category string, tags []string) (int64, error) {
	id := utils.GenerateID()
	_, err := DB.Exec("INSERT INTO topics VALUES (?, ?, ?, ?, ?, ?)", id, title, 0, 0, category, sql.NullInt64{})
	if err != nil {
		return 0, fmt.Errorf("CreateTopic error: %v", err)
	}
	for _, tag := range tags {
		_, err := DB.Exec("INSERT INTO have VALUES (?, ?)", id, tag)
		if err != nil {
			return 0, fmt.Errorf("Tags error: %v", err)
		}
	}
	fmt.Printf("topic '%v' added to the category '%v'", title, category)
	return id, nil
}

// GetTopic returns topic by id
func GetTopic(id int64) (*Topic, error) {
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

// GetTopicsByCategory returns topics by category
func GetTopicsByCategory(category string) ([]Topic, error) {
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

// GetTopicsByTag returns topics by tag
func GetTopicsByTag(tag string) ([]Topic, error) {
	rows, err := DB.Query("SELECT * FROM topics WHERE id_topic IN (SELECT id_topic FROM have WHERE tag_name = ?)", tag)
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

		err := topic.FetchTags()
		if err != nil {
			return nil, err
		}
		topics = append(topics, topic)
	}

	HandleSQLErrors(rows)

	return topics, nil
}

//AddViews add views to topic
func AddViews(id int64) error {
	_, err := DB.Exec("UPDATE topics SET views = views + 1 WHERE id_topic = ?", id)
	if err != nil {
		return fmt.Errorf("AddViews error: %v", err)
	}
	return nil
}
