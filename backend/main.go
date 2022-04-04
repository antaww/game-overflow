package main

import (
	"database/sql"
	"fmt"
	"github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"log"
	. "main/routes"
	. "main/sql"
	"net/http"
	"os"
	"time"
)

const inactiveTime = 60 * time.Second
const printDelay = 10 * time.Second

func main() {
	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatal(err)
	}

	css := http.FileServer(http.Dir("../client/style/"))
	http.Handle("/static/", http.StripPrefix("/static/", css))

	resources := http.FileServer(http.Dir("../backend/resources/"))
	http.Handle("/resources/", http.StripPrefix("/resources/", resources))

	scripts := http.FileServer(http.Dir("../client/scripts/"))
	http.Handle("/scripts/", http.StripPrefix("/scripts/", scripts))

	http.HandleFunc("/", IndexRoute)

	http.HandleFunc("/login", LoginRoute)

	http.HandleFunc("/sign-up", SignUpRoute)

	http.HandleFunc("/confirm-password", ConfirmPasswordRoute)

	http.HandleFunc("/admin/edit-username", AdminEditUsernameRoute)

	http.HandleFunc("/post-message", PostMessageRoute)

	http.HandleFunc("/create-topic", CreateTopicRoute)

	http.HandleFunc("/logout", LogoutRoute)

	http.HandleFunc("/settings", SettingsRoute)

	http.HandleFunc("/feed", FeedRoute)

	http.HandleFunc("/topic", TopicsRoute)

	http.HandleFunc("/like", LikeRoute)
	//http.HandleFunc("/like", DislikeRoute)

	// Capture connection properties.
	cfg := mysql.Config{
		User:                 os.Getenv("DB_USER"),
		Passwd:               os.Getenv("DB_PASSWORD"),
		Net:                  "tcp",
		Addr:                 os.Getenv("DB_ADDRESS"),
		DBName:               "forum",
		AllowNativePasswords: true,
		ParseTime:            true,
	}
	// Get a database handle.
	DB, err = sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		log.Fatal(err)
	}

	pingErr := DB.Ping()
	if pingErr != nil {
		log.Fatal(pingErr)
	}
	fmt.Println("Connected!")

	go IsUserActive()

	fmt.Println("Server started at localhost:8091")
	err = http.ListenAndServe(":8091", http.HandlerFunc(LogHandler))
	if err != nil {
		log.Fatal(err)
	}
}

func IsUserActive() {
	for {
		if TemplatesData.ConnectedUser == nil {
			continue
		}
		if PageLoadedTime.Add(inactiveTime).After(time.Now()) {
			fmt.Println("User is active")
			if !TemplatesData.ConnectedUser.IsOnline {
				err := SetUserOnline(TemplatesData.ConnectedUser.Id, true)
				if err != nil {
					log.Fatal(err)
				}
				TemplatesData.ConnectedUser = GetUserById(TemplatesData.ConnectedUser.Id)
			}

		} else {
			fmt.Println("User is inactive")
			if TemplatesData.ConnectedUser.IsOnline {
				err := SetUserOnline(TemplatesData.ConnectedUser.Id, false)
				if err != nil {
					log.Fatal(err)
				}
				TemplatesData.ConnectedUser = GetUserById(TemplatesData.ConnectedUser.Id)
			}

		}
		time.Sleep(printDelay)
	}
}
