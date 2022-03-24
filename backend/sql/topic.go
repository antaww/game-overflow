package sql

type Post struct {
	Id           int64  `db:"id_topic"`
	Title        string `db:"title"`
	IsClosed     bool   `db:"is_closed"`
	Views        int    `db:"views"`
	CategoryName string `db:"category_name"`
	Messages     []Message
}

// GetPost returns topic by id
func GetPost(id int64) (*Post, error) {
	rows, err := DB.Query("SELECT * FROM topic WHERE id_topic = ?", id)
	if err != nil {
		return nil, err
	}

	post := &Post{}
	err = Results(rows, &post.Id, &post.Title, &post.IsClosed, &post.Views, &post.CategoryName)
	if err != nil {
		return nil, err
	}

	HandleSQLErrors(rows)

	return post, nil
}

// FetchMessages get messages into topic from db using post id
func (post *Post) FetchMessages() error {
	message, err := GetMessages(post.Id)
	if err != nil {
		return err
	}

	post.Messages = message
	return nil
}
