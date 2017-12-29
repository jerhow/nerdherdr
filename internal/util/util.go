package util

import (
	"log"
)

func ErrChk(err error) {
	if err != nil {
		log.Fatal(err) // panic(err.Error)
	}
}

func Hi(name string) string {
	return "Hello, " + name
}

func AddTwoInts(x int, y int) int {
	return x + y
}
