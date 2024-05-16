package main

import (
	"fmt"
	"os"

	client "github.com/5aradise/cli-chat/client/pkg"
)

func main() {
	defer fmt.Scanln()

	host := os.Getenv("HOST")
	if host == "" {
		host = "127.0.0.1"
	}
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	client, err := client.New(host + ":" + port)
	if err != nil {
		fmt.Println(err)
		return
	}

	err = client.Run()
	if err != nil {
		fmt.Println(err)
	}
}
