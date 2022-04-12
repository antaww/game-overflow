package routes

import (
	"fmt"
	"log"
	"main/sql"
	"main/utils"
	"net/http"
	"net/url"
	"strconv"
)

type LikeResponse struct {
	Points int `json:"points"`
}

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

		user, err := sql.GetUserByRequest(r)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(category)
		fmt.Println(title)
		idTopic, err := sql.CreateTopic(title, category)
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

func DislikeRoute(w http.ResponseWriter, r *http.Request) {
	if r.Method == "PUT" {
		user, err := sql.GetUserByRequest(r)

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

func FeedRoute(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		queries := r.URL.Query()

		if queries.Has("category") {
			category := queries.Get("category")

			topics, err := sql.GetTopicsByCategories(category)
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
}

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

func PostMessageRoute(w http.ResponseWriter, r *http.Request) {
	user, err := sql.GetUserByRequest(r)
	if err != nil {
		log.Fatal(err)
	}

	// Get topic id from url
	queries := r.URL.Query()

	if queries.Has("id") {
		id := queries.Get("id")

		messageId, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			log.Fatal(err)
		}

		_, err = sql.AddMessage(user.Id, messageId, r.FormValue("post-text"))
		if err != nil {
			log.Fatal(err)
		}

		queriesId := url.Values{}
		queriesId.Add("id", id)

		http.Redirect(w, r, "/topic?"+queriesId.Encode(), http.StatusSeeOther)
	}

	return
}

func TopicRoute(w http.ResponseWriter, r *http.Request) {
	queries := r.URL.Query()

	if queries.Has("id") {
		id := queries.Get("id")

		Id, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			log.Fatal(err)
		}

		topic, err := sql.GetPost(Id)
		if err != nil {
			log.Fatal(err)
		}

		err = topic.FetchMessages()
		if err != nil {
			log.Fatal(err)
		}

		TemplatesData.ShownTopic = *topic

		err = utils.CallTemplate("topic", TemplatesData, w)
		if err != nil {
			log.Fatal(err)
		}
	}

	return
}
