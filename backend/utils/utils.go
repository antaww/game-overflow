package utils

import (
	"encoding/json"
	"log"
	"main/routes"
	"main/sql"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

func GenerateID() int64 {
	CurrentTime := time.Now().Unix()
	RandomNumbers := GenerateNumbers(4)
	result, err := strconv.Atoi(strconv.Itoa(int(CurrentTime)) + strconv.Itoa(int(RandomNumbers)))
	if err != nil {
		log.Fatal(err)
	}
	return int64(result)
}

func GenerateNumbers(length int) int64 {
	return rand.Int63n(int64(length))
}

func SendResponse(w http.ResponseWriter, data interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	bytes, err := json.Marshal(data)
	if err != nil {
		return err
	}

	_, err = w.Write(bytes)
	if err != nil {
		return err
	}

	return nil
}

func RandomString(n int) string {
	rand.Seed(time.Now().UnixNano())
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

	s := make([]rune, n)
	for i := range s {
		s[i] = letters[rand.Intn(len(letters))]
	}
	return string(s)
}

func CheckUserStatus() {
	for {
		time.Sleep(time.Minute)
	}
}

func LoginUser(r *http.Request) (*sql.User, error) {
	cookie, err := r.Cookie("session")

	if err != nil {
		routes.TemplatesData.ConnectedUser = nil
	} else {
		user, err := sql.GetUserBySession(cookie.Value)
		if err != nil {
			return nil, err
		}
		routes.TemplatesData.ConnectedUser = user

		err = sql.SetUserOnline(routes.TemplatesData.ConnectedUser.Id, true)
		if err != nil {
			return nil, err
		}

		return user, nil
	}

	return nil, nil
}
