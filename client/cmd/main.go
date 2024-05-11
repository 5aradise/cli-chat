package main

import (
	"fmt"

	client "github.com/5aradise/cli-chat/client/pkg"
)

func main() {
	defer fmt.Scanln()
	
	client, err := client.New()
	if err != nil {
		fmt.Println(err)
		return
	}

	err = client.Run()
	if err != nil {
		fmt.Println(err)
	}
}
