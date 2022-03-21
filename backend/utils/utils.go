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
