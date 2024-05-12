package client

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"

	"github.com/5aradise/cli-chat/client/internal/cli"
)

type Client struct {
	net.Conn
}

func New() (*Client, error) {
	conn, err := net.Dial("tcp", "127.0.0.1:8080")
	if err != nil {
		return nil, err
	}
	return &Client{conn}, nil
}

func (c *Client) Run() error {
	go c.listenServer()
	return c.listenClient()
}

func (c *Client) listenServer() {
	buf := make([]byte, 1024)

	for {
		l, err := c.Read(buf)
		if err != nil {
			continue
		}
		fmt.Println(string(buf[:l]))
	}
}

func (c *Client) listenClient() error {
	scanner := bufio.NewScanner(os.Stdin)
	var err error
	for {
		scanner.Scan()
		input := scanner.Text()
		if input[:1] == "/" {
			args := strings.Split(input[1:], " ")
			command, ok := commands[args[0]]
			if !ok {
				fmt.Println(cli.Color("System: unknown command", cli.Red))
				continue
			}
			err := command.fn(c, args[1:])
			if err != nil {
				fmt.Println(err)
			}
			continue
		}
		_, err := c.Write([]byte(input))
		if err != nil {
			fmt.Println(err)
			break
		}
	}
	return err
}
