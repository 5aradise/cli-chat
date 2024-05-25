package main

import (
	"fmt"
	"net"
	"os"

	client "github.com/5aradise/cli-chat/client/pkg"
)

func main() {
	host := os.Getenv("HOST")
	if host == "" {
		host = "127.0.0.1"
	}
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	client, err := client.New(net.JoinHostPort(host, port))
	if err != nil {
		fmt.Println(err)
		return
	}

	client.Run()
}
