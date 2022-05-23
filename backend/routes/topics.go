package routes

import (
	"io"
	"main/sql"
	"main/utils"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

type LikeResponse struct {
	Points int `json:"points"`
}

type EditMessageResponse struct {
	Message string `json:"message"`
}

// ChangeCategoryRoute handles the request to change the category of a topic
func ChangeCategoryRoute(w http.ResponseWriter, r *http.Request) {
	queries := r.URL.Query()

	if queries.Has("id") {
		id := queries.Get("id")
		category := queries.Get("category")

		user, err := sql.GetUserByRequest(r)
		if err != nil {
			utils.RouteError(err)
		}

		if user.Role == "admin" || user.Role == "moderator" {
			Id, err := strconv.ParseInt(id, 10, 64)
			if err != nil {
				utils.RouteError(err)
			}

			err = sql.ChangeCategory(Id, category)
		}
		http.Redirect(w, r, "/topic?id="+id, http.StatusSeeOther)
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// CloseTopicRoute closes a topic
func CloseTopicRoute(w http.ResponseWriter, r *http.Request) {
	queries := r.URL.Query()

	if queries.Has("id") {
		id := queries.Get("id")

		Id, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			utils.RouteError(err)
		}

		user, err := sql.GetUserByRequest(r)
		if err != nil {
			utils.RouteError(err)
		}

		if user.Role == "admin" || user.Role == "moderator" {
			Id, err := strconv.ParseInt(id, 10, 64)
			if err != nil {
				utils.RouteError(err)
			}

			Topic, err := sql.GetTopic(Id)
			if err != nil {
				utils.RouteError(err)
			}

			TopicFirstMessage, err := Topic.GetFirstMessage()
			if err != nil {
				utils.RouteError(err)
			}

			IsClosed, err := sql.CloseTopic(Id, TopicFirstMessage.AuthorId)
			if err != nil {
				utils.RouteError(err)
			}

			if !IsClosed {
				http.Redirect(w, r, "/topic?id="+id, http.StatusSeeOther)
			}
		} else {
			IsClosed, err := sql.CloseTopic(Id, user.Id)

			if !IsClosed {
				http.Redirect(w, r, "/topic?id="+id, http.StatusSeeOther)
			}
			if err != nil {
				utils.RouteError(err)
			}
		}

		http.Redirect(w, r, "/topic?id="+id, http.StatusSeeOther)
	}
}

// CreateTopicRoute is the route for creating a new topic
func CreateTopicRoute(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		templateData, err := GetTemplateDataFromRoute(w, r)
		if err != nil {
			utils.RouteError(err)
		}

		if templateData.ConnectedUser == nil {
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		err = utils.CallTemplate("create-topic", templateData, w)
		if err != nil {
			utils.RouteError(err)
		}
	}

	if r.Method == "POST" {
		err := r.ParseForm()
		if err != nil {
			return
		}

		category := r.Form["category"][0]
		title := r.Form["title"][0]
		content := r.Form["content"][0]
		tags := r.Form["tags"][0]
		replace := strings.ReplaceAll(tags, ",", " ")
		fields := strings.Fields(replace)

		user, err := sql.GetUserByRequest(r)
		if err != nil {
			utils.RouteError(err)
		}

		idTopic, err := sql.CreateTopic(title, category, fields)
		if err != nil {
			utils.RouteError(err)
		}

		idMessage, err := sql.AddMessage(user.Id, idTopic, content)
		if err != nil {
			utils.RouteError(err)
		}

		_, err = sql.DB.Query("UPDATE topics SET id_first_message = ? WHERE id_topic = ? ", idMessage, idTopic)
		if err != nil {
			utils.RouteError(err)
		}

		queriesCategory := url.Values{}
		queriesCategory.Add("category", category)

		http.Redirect(w, r, "/feed?"+queriesCategory.Encode(), http.StatusSeeOther)
	}
}

// DeleteTopicRoute is the route for deleting a topic
func DeleteTopicRoute(w http.ResponseWriter, r *http.Request) {
	queries := r.URL.Query()

	if queries.Has("id") {
		idTopic := queries.Get("id")
		Id, err := strconv.ParseInt(idTopic, 10, 64)

		topic, err := sql.GetTopic(Id)
		if err != nil {
			utils.RouteError(err)
		}

		topicFirstMsg, err := topic.GetFirstMessage()
		if err != nil {
			utils.RouteError(err)
		}

		user, err := sql.GetUserByRequest(r)
		if err != nil {
			utils.RouteError(err)
		}

		if user.Role == "admin" || user.Role == "moderator" || user.Id == topicFirstMsg.AuthorId {
			err := sql.DeleteTopic(Id)
			if err != nil {
				utils.RouteError(err)
			}
		}
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// DeleteMessageRoute is the route for deleting a message
func DeleteMessageRoute(w http.ResponseWriter, r *http.Request) {
	queries := r.URL.Query()

	if queries.Has("message-id") {
		idMessage := queries.Get("message-id")

		id, err := strconv.ParseInt(idMessage, 10, 64)
		if err != nil {
			utils.RouteError(err)
		}

		err = sql.DeleteMessage(id)
		if err != nil {
			utils.RouteError(err)
		}

		http.Redirect(w, r, r.Header.Get("Referer"), http.StatusFound)
	}
}

// DislikeRoute is the route for handling the dislike request
func DislikeRoute(w http.ResponseWriter, r *http.Request) {
	if r.Method == "PUT" {
		user, err := sql.GetUserByRequest(r)
		if err != nil {
			utils.RouteError(err)
		}

		messageIdArg := r.URL.Query().Get("id")
		messageId, _ := strconv.ParseInt(messageIdArg, 10, 64)

		messageLike, err := sql.MessageGetLikeFrom(messageId, user.Id)
		if err != nil {
			utils.RouteError(err)
		}

		if messageLike == nil {
			_, err = sql.DislikeMessage(messageId, user.Id)
			if err != nil {
				utils.RouteError(err)
			}
		} else if !messageLike.IsLike {
			_, err = sql.DeleteDislikeMessage(messageId, user.Id)
			if err != nil {
				utils.RouteError(err)
			}
		} else {
			_, err = sql.DeleteLikeMessage(messageId, user.Id)
			if err != nil {
				utils.RouteError(err)
			}

			_, err = sql.DislikeMessage(messageId, user.Id)
			if err != nil {
				utils.RouteError(err)
			}
		}

		response := LikeResponse{}
		message, err := sql.GetMessage(messageId)
		if err != nil {
			utils.RouteError(err)
		}

		response.Points = message.CalculatePoints()

		err = utils.SendResponse(w, response)
		if err != nil {
			utils.RouteError(err)
		}
	}
}

// EditMessageRoute is the route for editing a message
func EditMessageRoute(w http.ResponseWriter, r *http.Request) {
	queries := r.URL.Query()

	if queries.Has("idMessage") {
		idMessage := queries.Get("idMessage")
		idTopic := queries.Get("id")
		contentMessage, _ := io.ReadAll(r.Body)

		Id, err := strconv.ParseInt(idMessage, 10, 64)
		if err != nil {
			utils.RouteError(err)
		}

		// make a variable message taking the content from the posts-content class in the html
		err = sql.EditMessage(Id, string(contentMessage))
		if err != nil {
			utils.RouteError(err)
		}

		queriesId := url.Values{}
		queriesId.Add("id", idTopic)

		http.Redirect(w, r, "/topic?"+queriesId.Encode(), http.StatusSeeOther)
	}
}

// FeedRoute is the route for the feed
func FeedRoute(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		queries := r.URL.Query()

		if queries.Has("category") {
			category := queries.Get("category")

			topics, err := sql.GetTopicsByCategory(category)
			if err != nil {
				utils.RouteError(err)
			}

			for i := 0; i < len(topics); i++ {
				err = topics[i].FetchTags()
				if err != nil {
					utils.RouteError(err)
				}
			}

			templateData, err := GetTemplateDataFromRoute(w, r)
			if err != nil {
				utils.RouteError(err)
			}

			templateData.ShownTopics = topics
			templateData.SortTopics()

			err = utils.CallTemplate("feed", templateData, w)
			if err != nil {
				utils.RouteError(err)
			}
		} else if queries.Has("tag") {
			FeedTagsRoute(w, r)
		}
	}
}

// FeedTagsRoute is the route for the feed by tags
func FeedTagsRoute(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		queries := r.URL.Query()

		if queries.Has("tag") {
			tag := queries.Get("tag")

			topics, err := sql.GetTopicsByTag(tag)
			if err != nil {
				utils.RouteError(err)
			}

			templateData, err := GetTemplateDataFromRoute(w, r)
			if err != nil {
				utils.RouteError(err)
			}

			templateData.ShownTopics = topics

			err = utils.CallTemplate("feed", templateData, w)
			if err != nil {
				utils.RouteError(err)
			}
		}
	}
}

// LikeRoute is the route for handling the like request
func LikeRoute(w http.ResponseWriter, r *http.Request) {
	if r.Method == "PUT" {
		user, err := sql.GetUserByRequest(r)
		if err != nil {
			utils.RouteError(err)
		}

		messageIdArg := r.URL.Query().Get("id")
		messageId, err := strconv.ParseInt(messageIdArg, 10, 64)
		if err != nil {
			utils.RouteError(err)
		}

		messageLike, err := sql.MessageGetLikeFrom(messageId, user.Id)
		if err != nil {
			utils.RouteError(err)
		}

		if messageLike == nil {
			_, err = sql.LikeMessage(messageId, user.Id)
			if err != nil {
				utils.RouteError(err)
			}
		} else if messageLike.IsLike {
			_, err = sql.DeleteLikeMessage(messageId, user.Id)
			if err != nil {
				utils.RouteError(err)
			}
		} else {
			_, err = sql.DeleteDislikeMessage(messageId, user.Id)
			if err != nil {
				utils.RouteError(err)
			}

			_, err = sql.LikeMessage(messageId, user.Id)
			if err != nil {
				utils.RouteError(err)
			}
		}

		response := LikeResponse{}
		message, err := sql.GetMessage(messageId)
		if err != nil {
			utils.RouteError(err)
		}

		response.Points = message.CalculatePoints()

		err = utils.SendResponse(w, response)
		if err != nil {
			utils.RouteError(err)
		}
	}
}

// OpenTopicRoute is the route for opening a topic
func OpenTopicRoute(w http.ResponseWriter, r *http.Request) {
	queries := r.URL.Query()

	if queries.Has("id") {
		id := queries.Get("id")

		user, err := sql.GetUserByRequest(r)
		if err != nil {
			utils.RouteError(err)
		}

		if user.Role == "admin" || user.Role == "moderator" {
			Id, err := strconv.ParseInt(id, 10, 64)
			if err != nil {
				utils.RouteError(err)
			}

			Topic, err := sql.GetTopic(Id)
			if err != nil {
				utils.RouteError(err)
			}

			TopicFirstMessage, err := Topic.GetFirstMessage()
			if err != nil {
				utils.RouteError(err)
			}

			IsOpen, err := sql.OpenTopic(Id, TopicFirstMessage.AuthorId)
			if err != nil {
				utils.RouteError(err)
			}

			if !IsOpen {
				http.Redirect(w, r, "/topic?id="+id, http.StatusSeeOther)
			}
		}

		http.Redirect(w, r, "/topic?id="+id, http.StatusSeeOther)
	}
}

// PostMessageRoute is the route for posting a message
func PostMessageRoute(w http.ResponseWriter, r *http.Request) {
	queries := r.URL.Query()

	user, err := sql.GetUserByRequest(r)
	if err != nil {
		utils.RouteError(err)
	}

	if queries.Has("id") {
		id := queries.Get("id")

		Id, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			utils.RouteError(err)
		}

		_, err = sql.AddMessage(user.Id, Id, r.FormValue("post-text"))
		if err != nil {
			utils.RouteError(err)
		}

		queriesId := url.Values{}
		queriesId.Add("id", id)

		http.Redirect(w, r, "/topic?"+queriesId.Encode(), http.StatusSeeOther)
	}
}

// TopicRoute is the route for showing a topic
func TopicRoute(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		queries := r.URL.Query()

		if queries.Has("id") {
			id := queries.Get("id")

			Id, err := strconv.ParseInt(id, 10, 64)
			if err != nil {
				utils.RouteError(err)
			}

			topic, err := sql.GetTopic(Id)
			if err != nil {
				utils.RouteError(err)
			}

			err = sql.AddViews(Id)
			if err != nil {
				utils.RouteError(err)
			}

			err = topic.FetchMessages()
			if err != nil {
				utils.RouteError(err)
			}

			topic.Tags, err = sql.GetTags(Id)
			if err != nil {
				utils.RouteError(err)
			}

			templateData, err := GetTemplateDataFromRoute(w, r)
			if err != nil {
				utils.RouteError(err)
			}

			templateData.ShownTopic = topic

			err = utils.CallTemplate("topic", templateData, w)
			if err != nil {
				utils.RouteError(err)
			}
		}
	}
}

// SearchRoute is the route for searching topics
func SearchRoute(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		templateData, err := GetTemplateDataFromRoute(w, r)
		if err != nil {
			utils.RouteError(err)
		}

		search := r.FormValue("search")

		if search == "" {
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		} else {
			topics, err := sql.SearchTopics(search)
			if err != nil {
				utils.RouteError(err)
			}

			templateData.ShownTopics = topics
			if len(templateData.ShownTopics) == 0 {
				http.Redirect(w, r, "/", http.StatusSeeOther)
				return
			}

			err = utils.CallTemplate("feed", templateData, w)
			if err != nil {
				utils.RouteError(err)
			}
		}
	}
}
