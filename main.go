package main

import (
	"log"
)

func main() {

	store, err := NewPostgresStore()

	if err != nil {
		log.Fatal(err)
	}


	server := NewApiServer(":5000", store)

	server.Run()
}