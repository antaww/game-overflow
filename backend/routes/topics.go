package routes

import (
	"fmt"
	"log"
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

// CreateTopicRoute is the route for creating a new topic
func CreateTopicRoute(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		if TemplatesData.ConnectedUser == nil {
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}
		err := utils.CallTemplate("create-topic", TemplatesData, w)
		if err != nil {
			log.Fatal(err)
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
		replace = strings.ReplaceAll(replace, ";", " ")
		fields := strings.Fields(replace)

		user, err := sql.GetUserByRequest(r)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(category)
		fmt.Println(title)
		fmt.Println(tags)
		idTopic, err := sql.CreateTopic(title, category, fields)
		if err != nil {
			log.Fatal(err)
		}

		idMessage, err := sql.AddMessage(user.Id, idTopic, content)
		if err != nil {
			log.Fatal(err)
		}

		_, err = sql.DB.Query("UPDATE topics SET id_first_message = ? WHERE id_topic = ? ", idMessage, idTopic)
		if err != nil {
			log.Fatal(err)
		}

		queriesCategory := url.Values{}
		queriesCategory.Add("category", category)

		http.Redirect(w, r, "/feed?"+queriesCategory.Encode(), http.StatusSeeOther)
	}
}

// DeleteMessageRoute is the route for deleting a message
func DeleteMessageRoute(w http.ResponseWriter, r *http.Request) {
	queries := r.URL.Query()

	if queries.Has("idMessage") {
		idMessage := queries.Get("idMessage")
		idTopic := queries.Get("id")

		Id, err := strconv.ParseInt(idMessage, 10, 64)
		if err != nil {
			log.Fatal(err)
		}

		err = sql.DeleteMessage(Id)
		fmt.Println("Delete message")
		if err != nil {
			log.Fatal(err)
		}

		queriesId := url.Values{}
		queriesId.Add("id", idTopic)

		http.Redirect(w, r, "/topic?"+queriesId.Encode(), http.StatusSeeOther)
	}
}

// DislikeRoute is the route for handling the dislike request
func DislikeRoute(w http.ResponseWriter, r *http.Request) {
	if r.Method == "PUT" {
		user, err := sql.GetUserByRequest(r)
		if err != nil {
			log.Fatal(err)
		}

		messageIdArg := r.URL.Query().Get("id")
		messageId, _ := strconv.ParseInt(messageIdArg, 10, 64)

		messageLike, err := sql.MessageGetLikeFrom(messageId, user.Id)
		if err != nil {
			log.Fatal(err)
		}

		if messageLike == nil {
			_, err = sql.DislikeMessage(messageId, user.Id)
			if err != nil {
				log.Fatal(err)
			}
		} else if !messageLike.IsLike {
			_, err = sql.DeleteDislikeMessage(messageId, user.Id)
			if err != nil {
				log.Fatal(err)
			}
		} else {
			_, err = sql.DeleteLikeMessage(messageId, user.Id)
			if err != nil {
				log.Fatal(err)
			}

			_, err = sql.DislikeMessage(messageId, user.Id)
			if err != nil {
				log.Fatal(err)
			}
		}

		response := LikeResponse{}
		message, err := sql.GetMessage(messageId)
		if err != nil {
			log.Fatal(err)
		}

		response.Points = message.CalculatePoints()

		err = utils.SendResponse(w, response)
		if err != nil {
			log.Fatal(err)
		}
	}
}

// EditMessageRoute is the route for editing a message
func EditMessageRoute(w http.ResponseWriter, r *http.Request) {
	queries := r.URL.Query()

	if queries.Has("idMessage") {
		idMessage := queries.Get("idMessage")
		idTopic := queries.Get("id")

		Id, err := strconv.ParseInt(idMessage, 10, 64)
		if err != nil {
			log.Fatal(err)
		}

		err = sql.DeleteMessage(Id)
		fmt.Println("Delete message")
		if err != nil {
			log.Fatal(err)
		}

		queriesId := url.Values{}
		queriesId.Add("id", idTopic)

		http.Redirect(w, r, "/topic?"+queriesId.Encode(), http.StatusSeeOther)
	}
}

// FeedRoute is the route for the feed
func FeedRoute(w http.ResponseWriter, r *http.Request) {
	queries := r.URL.Query()

	if queries.Has("category") {
		category := queries.Get("category")

		topics, err := sql.GetTopicsByCategory(category)
		if err != nil {
			log.Fatal(err)
		}

		for i := 0; i < len(topics); i++ {
			topics[i].Tags, err = sql.GetTags(topics[i].Id)
			if err != nil {
				log.Fatal(err)
			}
		}

		TemplatesData.ShownTopics = topics

		err = utils.CallTemplate("feed", TemplatesData, w)
		if err != nil {
			log.Fatal(err)
		}
	} else if queries.Has("tag") {
		FeedTagsRoute(w, r)
	}
}

// FeedTagsRoute is the route for the feed by tags
func FeedTagsRoute(w http.ResponseWriter, r *http.Request) {
	queries := r.URL.Query()

	if queries.Has("tag") {
		tag := queries.Get("tag")

		topics, err := sql.GetTopicsByTag(tag)
		if err != nil {
			log.Fatal(err)
		}

		TemplatesData.ShownTopics = topics

		err = utils.CallTemplate("feed", TemplatesData, w)
		if err != nil {
			log.Fatal(err)
		}
	}
}

// LikeRoute is the route for handling the like request
func LikeRoute(w http.ResponseWriter, r *http.Request) {
	if r.Method == "PUT" {
		user, err := sql.GetUserByRequest(r)
		if err != nil {
			log.Fatal(err)
		}

		messageIdArg := r.URL.Query().Get("id")
		messageId, err := strconv.ParseInt(messageIdArg, 10, 64)
		if err != nil {
			log.Fatal(err)
		}

		messageLike, err := sql.MessageGetLikeFrom(messageId, user.Id)
		if err != nil {
			log.Fatal(err)
		}

		if messageLike == nil {
			_, err = sql.LikeMessage(messageId, user.Id)
			if err != nil {
				log.Fatal(err)
			}
		} else if messageLike.IsLike {
			_, err = sql.DeleteLikeMessage(messageId, user.Id)
			if err != nil {
				log.Fatal(err)
			}
		} else {
			_, err = sql.DeleteDislikeMessage(messageId, user.Id)
			if err != nil {
				log.Fatal(err)
			}

			_, err = sql.LikeMessage(messageId, user.Id)
			if err != nil {
				log.Fatal(err)
			}
		}

		err = TemplatesData.ShownTopic.FetchMessages()
		if err != nil {
			log.Fatal(err)
		}

		response := LikeResponse{}
		message, err := sql.GetMessage(messageId)
		if err != nil {
			log.Fatal(err)
		}

		response.Points = message.CalculatePoints()

		err = utils.SendResponse(w, response)
		if err != nil {
			log.Fatal(err)
		}
	}
}

// PostMessageRoute is the route for posting a message
func PostMessageRoute(w http.ResponseWriter, r *http.Request) {
	queries := r.URL.Query()

	user, err := sql.GetUserByRequest(r)
	if err != nil {
		log.Fatal(err)
	}

	if queries.Has("id") {
		id := queries.Get("id")

		Id, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			log.Fatal(err)
		}

		_, err = sql.AddMessage(user.Id, Id, r.FormValue("post-text"))
		if err != nil {
			log.Fatal(err)
		}

		queriesId := url.Values{}
		queriesId.Add("id", id)

		http.Redirect(w, r, "/topic?"+queriesId.Encode(), http.StatusSeeOther)
	}
}

// TopicRoute is the route for showing a topic
func TopicRoute(w http.ResponseWriter, r *http.Request) {
	queries := r.URL.Query()

	if queries.Has("id") {
		id := queries.Get("id")

		Id, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			log.Fatal(err)
		}

		topic, err := sql.GetTopic(Id)
		if err != nil {
			log.Fatal(err)
		}

		err = topic.FetchMessages()
		if err != nil {
			log.Fatal(err)
		}

		topic.Tags, err = sql.GetTags(Id)
		if err != nil {
			log.Fatal(err)
		}

		TemplatesData.ShownTopic = *topic

		err = utils.CallTemplate("topic", TemplatesData, w)
		if err != nil {
			log.Fatal(err)
		}
	}
}
