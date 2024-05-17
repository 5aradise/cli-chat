package client

import (
	"bufio"
	"net"
	"os"

	"github.com/5aradise/cli-chat/client/internal/cli"
)

const (
	systemMsgCode byte = 0x00
	chatMsgCode   byte = 0x10
	userMsgCode   byte = 0x20
	createCode    byte = 0x01
	connCode      byte = 0x02
	exitCode      byte = 0x03
)

type Client struct {
	net.Conn
	printLn    *int
	isInChat   bool
	chatColors map[string]cli.Color
}

func New(address string) (*Client, error) {
	conn, err := net.Dial("tcp", address)
	if err != nil {
		return nil, err
	}
	printLn := 1
	return &Client{conn, &printLn, false, nil}, nil
}

func (c *Client) updateScreen() {
	cli.ClearConsole()
	cli.PrintInputFrame()
	cli.MoveToInput()
	*c.printLn = 1
}

func (c *Client) printf(format string, a ...any) {
	cli.SafePrintf(c.printLn, format, a...)
}

func (c *Client) Run() error {
	c.updateScreen()
	go c.listenServer()
	return c.listenClient()
}

func (c *Client) listenServer() {
	buf := make([]byte, 1024)

	for {
		l, err := c.Read(buf)
		if err != nil {
			break
		}
		c.processResp(buf[:l])
	}

	cli.ClearConsole()
	c.printf(cli.Colorize("You've been disconnected from the server", cli.RedS))
}

func (c *Client) listenClient() error {
	var err error
	scanner := bufio.NewScanner(os.Stdin)
	cli.PrintInputFrame()
	cli.MoveToInput()
	var input string
	for len(input) == 0 {
		input = cli.Scan(scanner)
	}
	c.Write([]byte(input))
	for {
		input = cli.Scan(scanner)
		if len(input) == 0 {
			continue
		}
		err = c.processReq(input)
		if err != nil {
			break
		}
	}
	return err
}
