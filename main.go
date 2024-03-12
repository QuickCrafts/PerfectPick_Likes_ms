package main

import (
	"fmt"
	"log"
)

func main() {
	fmt.Println("Hello! Welcome to PerfectPick Likes Microservice")

	store, err := NewNeo4jStore()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%+v\n", store)
	server := NewAPIServer(":3000", store)
	server.Run()
}
