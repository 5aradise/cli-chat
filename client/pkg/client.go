package client

import (
	"bufio"
	"net"
	"os"

	"github.com/5aradise/cli-chat/client/internal/cli"
)

const bufferSize = 512

type client struct {
	net.Conn
	printLn    *int
	isInChat   bool
	chatColors map[string]cli.Color
}

func New(address string) (*client, error) {
	conn, err := net.Dial("tcp", address)
	if err != nil {
		return nil, err
	}
	printLn := 1
	return &client{conn, &printLn, false, nil}, nil
}

func (c *client) Run() {
	c.updateScreen()

	go c.listenClient()
	c.listenServer()
}

func (c *client) listenServer() {
	buf := make([]byte, bufferSize)

	for {
		l, err := c.Read(buf)
		if err != nil {
			break
		}
		c.processResp(buf[:l])
	}

	cli.ClearConsole()
	c.printf(formatSystemMsg("you've been disconnected from the server"))
}

func (c *client) listenClient() {
	scanner := bufio.NewScanner(os.Stdin)
	var input string
	for len(input) == 0 {
		input, _ = cli.Scan(scanner)
	}
	c.Write([]byte(input))
	for {
		input, inputLen := cli.Scan(scanner)
		if inputLen == 0 {
			continue
		}
		if inputLen > cli.MaxInputLen {
			c.printf(formatSystemMsg("your message is too long"))
			continue
		}
		err := c.processReq(input)
		if err != nil {
			c.printf(formatSystemMsg(err.Error()))
		}
	}
}

func (c *client) updateScreen() {
	cli.ClearConsole()
	cli.PrintInputFrame()
	cli.MoveToInput()
	*c.printLn = 1
}

func (c *client) printf(format string, a ...any) {
	cli.SafePrintf(c.printLn, format, a...)
}
