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
)

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

		http.Redirect(w, r, "/topic?" + queriesId.Encode(), http.StatusSeeOther)
	}

	return
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
		idTopic, err := sql.CreateTopic(title, category)
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

		http.Redirect(w, r, "/feed?" + queriesCategory.Encode(), http.StatusSeeOther)
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

		TemplatesData.ShownTopics = topics

		err = utils.CallTemplate("feed", TemplatesData, w)
		if err != nil {
			log.Fatal(err)
		}
	}

	return
}

func TopicsRoute(w http.ResponseWriter, r *http.Request) {
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
