package utils

import (
	"html/template"
	"net/http"
	"time"
)

var TemplateMap = template.FuncMap{
	"safeURL": func(u string) template.URL {
		return template.URL(u)
	},
	"formatDate": func(t time.Time) string {
		return t.Format("2006-01-02 15:04:05")
	},
	"MinusOne": func(i int) int {
		return i - 1
	},
	/*"ConnectedUserMessage": func(r *http.Request, id int64) bool {
		// recreate database connection
		DatabaseConfig := mysql.Config{
			User:                 os.Getenv("DB_USER"),
			Passwd:               os.Getenv("DB_PASSWORD"),
			Net:                  "tcp",
			Addr:                 os.Getenv("DB_ADDRESS"),
			DBName:               "forum",
			AllowNativePasswords: true,
			ParseTime:            true,
		}

		// connect
		db, err := sql.Open("mysql", DatabaseConfig.FormatDSN())
		if err != nil {
			log.Fatal(err)
		}

		// test ping
		pingErr := db.Ping()
		if pingErr != nil {
			log.Fatal(pingErr)
		}

		// get cookie
		cookie, err := r.Cookie("session")
		if err != nil {
			return false
		}

		// get user id by session id
		result, err := db.Query("SELECT id_user FROM sessions WHERE id_session = ?", cookie.Value)
		if err != nil {
			return false
		}

		var idUser int64

		if result.Next() {
			err := result.Scan(&idUser)
			if err != nil {
				return false
			}
		} else {
			return false
		}

		defer func(rows *sql.Rows) {
			err := rows.Close()
			if err != nil {
				log.Fatal(err)
			}
		}(result)

		defer func(rows *sql.Rows) {
			err := rows.Err()
			if err != nil {
				log.Fatal(err)
			}
		}(result)

		return idUser == id
	},*/
}

func CallTemplate(templateName string, data interface{}, w http.ResponseWriter) error {
	templates := template.New("").Funcs(TemplateMap)
	templates, err := templates.ParseFiles("../client/templates/main.gohtml", "../client/templates/"+templateName+".gohtml")
	if err != nil {
		return err
	}
	templates, err = templates.ParseGlob("../client/templates/components/*.gohtml")
	if err != nil {
		return err
	}

	err = templates.ExecuteTemplate(w, "main", data)
	if err != nil {
		return err
	}

	_ = templateName
	return nil
}
