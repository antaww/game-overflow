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
)

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

	http.HandleFunc("/admin/edit-username", AdminEditUsernameRoute)

	http.HandleFunc("/confirm-password", ConfirmPasswordRoute)

	http.HandleFunc("/create-topic", CreateTopicRoute)

	http.HandleFunc("/delete-message", DeleteMessageRoute)

	http.HandleFunc("/dislike", DislikeRoute)

	http.HandleFunc("/edit-message", EditMessageRoute)

	http.HandleFunc("/feed", FeedRoute)

	http.HandleFunc("/is-active", IsActiveRoute)

	http.HandleFunc("/like", LikeRoute)

	http.HandleFunc("/login", LoginRoute)

	http.HandleFunc("/logout", LogoutRoute)

	http.HandleFunc("/post-message", PostMessageRoute)

	http.HandleFunc("/settings", SettingsRoute)

	http.HandleFunc("/sign-up", SignUpRoute)

	http.HandleFunc("/topic", TopicRoute)

	http.HandleFunc("/users-active", UsersActive)

	// Capture connection properties.
	DatabaseConfig := mysql.Config{
		User:                 os.Getenv("DB_USER"),
		Passwd:               os.Getenv("DB_PASSWORD"),
		Net:                  "tcp",
		Addr:                 os.Getenv("DB_ADDRESS"),
		DBName:               "forum",
		AllowNativePasswords: true,
		ParseTime:            true,
	}
	// Get a database handle.
	DB, err = sql.Open("mysql", DatabaseConfig.FormatDSN())
	if err != nil {
		log.Fatal(err)
	}

	pingErr := DB.Ping()
	if pingErr != nil {
		log.Fatal(pingErr)
	}
	fmt.Println("Connected!")
	UsersOfflineRoute()

	const port = ":8091"

	fmt.Println("Server started at localhost:", port)
	err = http.ListenAndServe(port, http.HandlerFunc(LogHandler))
	if err != nil {
		log.Fatal(err)
	}
}
