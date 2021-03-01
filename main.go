package main

import (
	"github.com/codeliger/book-recommend/db"
	v1 "github.com/codeliger/book-recommend/server/rest/v1"
)

func main() {
	err := db.Init()
	if err != nil {
		panic(err.Error())
	}
	err = v1.Init()
	if err != nil {
		panic(err.Error())
	}
}
