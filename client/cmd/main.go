package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

func main() {
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		fmt.Println("Cannot connect to server")
		return
	}

	buf := make([]byte, 1024)
	go func() {
		for {
			l, err := conn.Read(buf)
			if err != nil {
				continue
			}
			fmt.Println(string(buf[:l]))
		}
	}()

	scanner := bufio.NewScanner(os.Stdin)
	for {
		scanner.Scan()
		_, err := conn.Write([]byte(scanner.Text()))
		if err != nil {
			fmt.Println(err)
		}
	}
}
