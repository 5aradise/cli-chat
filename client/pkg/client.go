package client

import (
	"fmt"
	"net"
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
	var input string
	fmt.Scan()

	buf := make([]byte, 1024)
	go func() {
		for {
			l, err := c.Read(buf)
			if err != nil {
				continue
			}
			fmt.Println(string(buf[:l]))
		}
	}()

	for {
		fmt.Scan(&input)
		_, err := c.Write([]byte(input))
		if err != nil {
			fmt.Println(err)
		}
	}
}
