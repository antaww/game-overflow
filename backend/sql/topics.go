package sql

import (
	"database/sql"
	"fmt"
	"main/utils"
	"net/http"
	"sort"
	"time"
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

type Topics []Topic

func (t Topics) SortBy(sortBy FeedSortType) {
	switch sortBy {
	case FeedSortNewest:
		sort.Slice(t, func(i, j int) bool {
			return t[i].GetDate().After(*t[j].GetDate())
		})
	case FeedSortOldest:
		sort.Slice(t, func(i, j int) bool {
			return t[i].GetDate().Before(*t[j].GetDate())
		})
	case FeedSortPopular:
		sort.Slice(t, func(i, j int) bool {
			return t[i].Views > t[j].Views
		})
	}
}

type Tags struct {
	Name string `db:"tag_name"`
}

func (topic Topic) GetDate() *time.Time {
	message, err := topic.GetFirstMessage()
	if err != nil {
		return nil
	}
	return &message.CreatedAt
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

// AddViews add a view to a topic
func AddViews(id int64) error {
	_, err := DB.Exec("UPDATE topics SET views = views + 1 WHERE id_topic = ?", id)
	if err != nil {
		return fmt.Errorf("AddViews error: %v", err)
	}
	return nil
}

// CloseTopic closes topic from user
func CloseTopic(topicId int64, userId int64) (bool, error) {
	editedLines, err := DB.Exec("UPDATE topics SET is_closed = 1 WHERE id_topic = ? AND id_first_message in (SELECT id_message FROM messages WHERE id_user = ?)", topicId, userId)
	if err != nil {
		return false, fmt.Errorf("CloseTopic error: %v", err)
	}

	affected, err := editedLines.RowsAffected()
	if err != nil {
		return false, fmt.Errorf("CloseTopic error: %v", err)
	}

	return affected > 0, nil
}

// GetShownTopic returns the actual shown topic on the URL
func GetShownTopic(r *http.Request) (*Topic, error) {
	queries := r.URL.Query()
	id := queries.Get("id")
	if id == "" {
		return nil, nil
	}

	topic := &Topic{}
	rows, err := DB.Query("SELECT * FROM topics WHERE id_topic = ?", id)
	if err != nil {
		return nil, err
	}

	err = Results(rows, &topic.Id, &topic.Title, &topic.IsClosed, &topic.Views, &topic.Category, &topic.IdFirstMessage)
	if err != nil {
		return nil, err
	}

	HandleSQLErrors(rows)

	err = topic.FetchMessages()
	if err != nil {
		return nil, err
	}

	err = topic.FetchTags()
	if err != nil {
		return nil, err
	}

	return topic, nil
}

// OpenTopic opens topic from user
func OpenTopic(topicId int64, userId int64) (bool, error) {
	editedLines, err := DB.Exec("UPDATE topics SET is_closed = 0 WHERE id_topic = ? AND id_first_message in (SELECT id_message FROM messages WHERE id_user = ?)", topicId, userId)
	if err != nil {
		return false, fmt.Errorf("CloseTopic error: %v", err)
	}

	affected, err := editedLines.RowsAffected()
	if err != nil {
		return false, fmt.Errorf("CloseTopic error: %v", err)
	}

	return affected > 0, nil
}

// CreateTopic create a new topic in db
func CreateTopic(title string, category string, tags []string) (int64, error) {
	id, err := utils.GenerateID()
	if err != nil {
		return 0, fmt.Errorf("CreateTopic error: %v", err)
	}
	_, err = DB.Exec("INSERT INTO topics VALUES (?, ?, ?, ?, ?, ?)", id, title, 0, 0, category, sql.NullInt64{})
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

// GetTopics returns topics from request
func GetTopics(rows *sql.Rows) ([]Topic, error) {
	var topics []Topic
	for rows.Next() {
		var topic Topic
		err := rows.Scan(&topic.Id, &topic.Title, &topic.IsClosed, &topic.Views, &topic.Category, &topic.IdFirstMessage)
		if err != nil {
			return nil, err
		}

		topics = append(topics, topic)
	}

	HandleSQLErrors(rows)

	return topics, nil
}

// GetTopicsByCategory returns topics by category
func GetTopicsByCategory(category string) ([]Topic, error) {
	rows, err := DB.Query("SELECT * FROM topics WHERE category_name = ?", category)
	if err != nil {
		return nil, err
	}

	return GetTopics(rows)
}

// GetTopicsNewest returns newest topics, limited by limit
func GetTopicsNewest(limit int) ([]Topic, error) {
	rows, err := DB.Query("SELECT * FROM topics ORDER BY id_topic DESC LIMIT ?", limit)
	if err != nil {
		return nil, err
	}

	return GetTopics(rows)
}

// GetTopicsOldest returns oldest topics, limited by limit
func GetTopicsOldest(limit int) ([]Topic, error) {
	rows, err := DB.Query("SELECT * FROM topics ORDER BY id_topic LIMIT ?", limit)
	if err != nil {
		return nil, err
	}

	return GetTopics(rows)
}

// GetTopicsPopular returns popular topics, limited by limit
func GetTopicsPopular(limit int) ([]Topic, error) {
	rows, err := DB.Query("SELECT * FROM topics ORDER BY views DESC LIMIT ?", limit)
	if err != nil {
		return nil, err
	}

	return GetTopics(rows)
}

func GetTopicsSortedByPoints(limit int) ([]Topic, error) {
	rows, err := DB.Query(`
SELECT topics.*
FROM topics
     INNER JOIN (SELECT SUM(IF(message_like.like = 0, -1, 1)) as likes, id_message
         FROM message_like
         GROUP BY id_message) as message_likes
    ON topics.id_first_message = message_likes.id_message
ORDER BY likes DESC LIMIT ?`, limit)
	if err != nil {
		return nil, err
	}

	return GetTopics(rows)
}

func GetTopicsFollowed(userId int64, limit int) ([]Topic, error) {
	//from follow, select id_user_followed where id_user_follower = userId
	//from messages, select id_message where id_user = id_user_followed and id_topic in (select id_topic from topics where id_first_message in (select id_message from messages where id_user = id_user_followed))
	rows, err := DB.Query("SELECT * FROM topics WHERE id_topic in (SELECT id_topic FROM messages WHERE id_user in (SELECT id_user_followed FROM follow WHERE id_user_follower = ?)) ORDER BY id_topic DESC LIMIT ?", userId, limit)
	if err != nil {
		return nil, err
	}
	fmt.Println(rows)
	fmt.Println(GetTopics(rows))
	return GetTopics(rows)
}

// GetTopicsByTag returns topics by tag
func GetTopicsByTag(tag string) ([]Topic, error) {
	rows, err := DB.Query("SELECT * FROM topics WHERE id_topic IN (SELECT id_topic FROM have WHERE tag_name = ?)", tag)
	if err != nil {
		return nil, err
	}

	return GetTopics(rows)
}

func ChangeCategory(id int64, category string) error {
	_, err := DB.Exec("UPDATE topics SET category_name = ? WHERE id_topic = ?", category, id)
	if err != nil {
		return fmt.Errorf("CloseTopic error: %v", err)
	}
	return nil
}

func DeleteTopic(id int64) error {
	_, err := DB.Exec("DELETE FROM message_like WHERE id_message IN (SELECT id_message FROM messages WHERE id_topic = ?)", id)
	if err != nil {
		return fmt.Errorf("DeleteTopic error: %v", err)
	}

	_, err = DB.Exec("DELETE FROM messages WHERE id_topic = ?", id)
	if err != nil {
		return fmt.Errorf("DeleteTopic error: %v", err)
	}

	_, err = DB.Exec("DELETE FROM have WHERE id_topic = ?", id)
	if err != nil {
		return fmt.Errorf("DeleteTopic error: %v", err)
	}

	_, err = DB.Exec("DELETE FROM topics WHERE id_topic = ?", id)
	if err != nil {
		return fmt.Errorf("DeleteTopic error: %v", err)
	}
	return nil
}

//SearchTopics returns every topic that contains the search string
func SearchTopics(search string) ([]Topic, error) {
	rows, err := DB.Query("SELECT * FROM topics WHERE title LIKE ?", "%"+search+"%")
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
