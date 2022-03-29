package utils

import (
	"log"
	"math/rand"
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
