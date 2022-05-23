package utils

import (
	"log"
	"os"
)

const Grey = "\033[90m"
const Blue = "\033[94m"
const Reset = "\033[0m"

var Logger = log.New(os.Stdout, Grey, log.LstdFlags)
