package main

import (
	"log"
	"os"

	chat "github.com/5aradise/cli-chat/server/pkg"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	server, err := chat.NewServer("localhost:8080")
	if err != nil {
		log.Fatal(err)
	}

	server.Run()
}
