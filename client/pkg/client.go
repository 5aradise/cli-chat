package client

import (
	"bufio"
	"net"
	"os"

	"github.com/5aradise/cli-chat/client/internal/cli"
)

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

func (c *client) listenClient() {
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
		err := c.processReq(input)
		if err != nil {
			c.printf(cli.Colorize("System: "+err.Error(), cli.RedS))
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
