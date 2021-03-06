package main

import (
	"database/sql"
	"github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	. "main/routes"
	. "main/sql"
	"main/utils"
	"net/http"
	"os"
)

const port = ":8091"

func main() {
	err := godotenv.Load("../.env")
	if err != nil {
		utils.MainError(err)
	}

	server := &http.Server{
		Addr:         port,
		Handler:      http.HandlerFunc(LogHandler),
		IdleTimeout:  0,
		ReadTimeout:  0,
		WriteTimeout: 0,
	}

	css := http.FileServer(http.Dir("../client/style/"))
	http.Handle("/static/", http.StripPrefix("/static/", css))

	resources := http.FileServer(http.Dir("../backend/resources/"))
	http.Handle("/resources/", http.StripPrefix("/resources/", resources))

	scripts := http.FileServer(http.Dir("../client/scripts/"))
	http.Handle("/scripts/", http.StripPrefix("/scripts/", scripts))

	http.HandleFunc("/", IndexRoute)
	http.HandleFunc("/admin/edit-username", AdminEditUsernameRoute)
	http.HandleFunc("/ban-user", BanRoute)
	http.HandleFunc("/change-category", ChangeCategoryRoute)
	http.HandleFunc("/close-topic", CloseTopicRoute)
	http.HandleFunc("/confirm-password", ConfirmPasswordRoute)
	http.HandleFunc("/cookies", CookieRoute)
	http.HandleFunc("/create-topic", CreateTopicRoute)
	http.HandleFunc("/delete-message", DeleteMessageRoute)
	http.HandleFunc("/delete-topic", DeleteTopicRoute)
	http.HandleFunc("/dislike", DislikeRoute)
	http.HandleFunc("/edit-message", EditMessageRoute)
	http.HandleFunc("/feed", FeedRoute)
	http.HandleFunc("/follow", FollowUserRoute)
	http.HandleFunc("/is-active", IsActiveRoute)
	http.HandleFunc("/legal-notice", LegalNoticeRoute)
	http.HandleFunc("/like", LikeRoute)
	http.HandleFunc("/likes", UserLikesRoute)
	http.HandleFunc("/login", LoginRoute)
	http.HandleFunc("/logout", LogoutRoute)
	http.HandleFunc("/open-topic", OpenTopicRoute)
	http.HandleFunc("/post-message", PostMessageRoute)
	http.HandleFunc("/posts", UserPostsRoute)
	http.HandleFunc("/profile", ProfileRoute)
	http.HandleFunc("/search", SearchRoute)
	http.HandleFunc("/settings", SettingsRoute)
	http.HandleFunc("/sign-up", SignUpRoute)
	http.HandleFunc("/topic", TopicRoute)
	http.HandleFunc("/unban-user", UnBanRoute)
	http.HandleFunc("/unfollow", UnFollowUserRoute)
	http.HandleFunc("/users", GetAllUsersRoute)
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
		utils.MainError(err)
	}

	// SetMaxIdleConns sets the maximum number of database transactions at the same time.
	// Setting this to 0 means that there is no limit and prevent random errors happening that stops the application.
	// It may lower performances a bit.
	DB.SetMaxIdleConns(0)

	pingErr := DB.Ping()
	if pingErr != nil {
		utils.MainError(pingErr)
	}

	utils.Logger.Printf("%vConnected!", utils.Reset)

	err = SetAllUsersOffline()
	if err != nil {
		utils.MainError(err)
	}

	utils.Logger.Printf("%vServer started at localhost: %v", utils.Reset, port)

	err = server.ListenAndServe()
	if err != nil {
		utils.MainError(err)
	}
}
