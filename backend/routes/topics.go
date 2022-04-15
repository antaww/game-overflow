package routes

import (
	"fmt"
	"log"
	"main/sql"
	"main/utils"
	"net/http"
	"net/url"
	"runtime/debug"
	"strconv"
	"strings"
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

		//get cookie from browser
		cookie, err := r.Cookie("session")
		if err != nil {
			debug.PrintStack()
			log.Fatal(err)
		}

		category := r.Form["category"][0]
		title := r.Form["title"][0]
		content := r.Form["content"][0]
		tags := r.Form["tags"][0]
		replace := strings.ReplaceAll(tags, ",", " ")
		replace = strings.ReplaceAll(replace, ";", " ")
		fields := strings.Fields(replace)
		//select user from session
		result, err := sql.DB.Query("SELECT id_user FROM sessions WHERE id_session = ?", cookie.Value)
		if err != nil {
			log.Fatal(err)
		}

		//get result from query
		var idUser int64
		if result.Next() {
			err = result.Scan(&idUser)
		}

		//Handle sql errors, close the query to avoid memory leaks
		sql.HandleSQLErrors(result)

		// Get User, save for TemplatesData (to show user logged in in templates)
		userConnected := sql.GetUserById(idUser)
		TemplatesData.ConnectedUser = userConnected

		fmt.Println(category)
		fmt.Println(title)
		fmt.Println(tags)
		idTopic, err := sql.CreateTopic(title, category, fields)
		if err != nil {
			log.Fatal(err)
		}

		idMessage, err := sql.AddMessage(userConnected.Id, idTopic, content)
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
		cookie, err := r.Cookie("session")
		if err != nil {
			if err == http.ErrNoCookie {
				return
			}

			log.Fatal(err)
		}

		user, err := sql.GetUserBySession(cookie.Value)
		if err != nil {
			log.Fatal(err)
		}

		TemplatesData.ConnectedUser = user

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
	queries := r.URL.Query()

	if queries.Has("category") {
		category := queries.Get("category")

		topics, err := sql.GetTopicsByCategories(category)
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
		FeedRouteTags(w, r)
	}

	return
}

func LikeRoute(w http.ResponseWriter, r *http.Request) {
	if r.Method == "PUT" {
		cookie, err := r.Cookie("session")
		if err != nil {
			if err == http.ErrNoCookie {
				return
			}

			log.Fatal(err)
		}

		user, err := sql.GetUserBySession(cookie.Value)
		if err != nil {
			log.Fatal(err)
		}

		TemplatesData.ConnectedUser = user

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
	//get cookie from browser
	cookie, err := r.Cookie("session")
	if err != nil {
		log.Fatal(err)
	}

	//select user from session
	result, err := sql.DB.Query("SELECT id_user FROM sessions WHERE id_session = ?", cookie.Value)
	if err != nil {
		log.Fatal(err)
	}

	//get result from query
	var idUser int64
	if result.Next() {
		err = result.Scan(&idUser)
	}

	//Handle sql errors, close the query to avoid memory leaks
	sql.HandleSQLErrors(result)

	// Get User, save for TemplatesData (to show user logged in templates)
	userConnected := sql.GetUserById(idUser)
	TemplatesData.ConnectedUser = userConnected

	// Get topic id from url
	queries := r.URL.Query()

	if queries.Has("id") {
		id := queries.Get("id")

		Id, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			log.Fatal(err)
		}

		_, err = sql.AddMessage(idUser, Id, r.FormValue("post-text"))
		if err != nil {
			log.Fatal(err)
		}

		queriesId := url.Values{}
		queriesId.Add("id", id)

		http.Redirect(w, r, "/topic?"+queriesId.Encode(), http.StatusSeeOther)
	}

	return
}

func DeleteMessageRoute(w http.ResponseWriter, r *http.Request) {
	queries := r.URL.Query()

	if queries.Has("idMessage") {
		idMessage := queries.Get("idMessage")

		Id, err := strconv.ParseInt(idMessage, 10, 64)
		if err != nil {
			log.Fatal(err)
		}

		err = sql.DeleteMessage(Id)
		fmt.Println("Delete message")
		if err != nil {
			log.Fatal(err)
		}
		//
		//queriesId := url.Values{}
		//queriesId.Add("id", id)

		http.Redirect(w, r, "/", http.StatusSeeOther)
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

	return
}

//If the url is /feed?tag=tagName, display every topics with the tag tagName
func FeedRouteTags(w http.ResponseWriter, r *http.Request) {
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
