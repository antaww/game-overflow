package utils

import (
	"encoding/json"
	"golang.org/x/crypto/bcrypt"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

func GenerateID() (int64, error) {
	CurrentTime := time.Now().Unix()
	RandomNumbers := GenerateNumbers(4)
	result, err := strconv.Atoi(strconv.Itoa(int(CurrentTime)) + strconv.Itoa(int(RandomNumbers)))
	return int64(result), err
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

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
