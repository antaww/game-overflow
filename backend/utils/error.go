package utils

import (
	"runtime/debug"
)

func MainError(err error) {
	debug.PrintStack()
	Logger.Fatal(err)
}

func RouteError(err error) {
	debug.PrintStack()
	Logger.Fatal(err)
}

func SQLError(err error) {
	debug.PrintStack()
	Logger.Fatal(err)
}
