package client

import (
	"bufio"
	"net"
	"os"
	"strings"

	"github.com/5aradise/cli-chat/client/internal/cli"
)

type Client struct {
	net.Conn
	chatColors map[string]cli.Color
	printLn    *int
}

func New(address string) (*Client, error) {
	conn, err := net.Dial("tcp", address)
	if err != nil {
		return nil, err
	}
	printLn := 1
	return &Client{conn, make(map[string]cli.Color), &printLn}, nil
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
		if buf[0] == 0 {
			switch buf[1] {
			case 0:
				cli.SafePrint(c.printLn, formatSystemMsg(buf[2:l]))
			case 1:
				cli.SafePrint(c.printLn, formatChatMsg(buf[2:l]))
			}
			continue
		}
		cli.SafePrint(c.printLn, c.formatUserMsg(buf[:l]))
	}
}

func (c *Client) listenClient() error {
	var err error
	scanner := bufio.NewScanner(os.Stdin)
	cli.PrintInputFrame()
	cli.MoveToInput()
	for {
		input := cli.Scan(scanner)
		if len(input) == 0 {
			continue
		}
		if input[:1] == "/" {
			args := strings.Split(input[1:], " ")
			command, ok := commands[args[0]]
			if !ok {
				cli.SafePrint(c.printLn, cli.Colorize("System: unknown command", cli.RedS))
				continue
			}
			err := command.fn(c, args[1:])
			if err != nil {
				cli.SafePrint(c.printLn, err.Error())
			}
			continue
		}
		_, err := c.Write([]byte(input))
		if err != nil {
			cli.SafePrint(c.printLn, err.Error())
			break
		}
	}
	return err
}
