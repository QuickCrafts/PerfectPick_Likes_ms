package main

import "fmt"

func main() {
	fmt.Println("Hello! Welcome to PerfectPick Likes Microservice")

	server := NewAPIServer(":3000")
	server.Run()
}
