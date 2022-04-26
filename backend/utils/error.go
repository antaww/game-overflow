package utils

import (
	"log"
	"runtime/debug"
)

func MainError(err error) {
	debug.PrintStack()
	log.Fatal(err)
}

func RouteError(err error) {
	debug.Stack()
	log.Fatal(err)
}

func SQLError(err error) {
	debug.Stack()
	log.Fatal(err)
}
