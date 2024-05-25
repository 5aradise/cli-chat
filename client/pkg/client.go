package client

import (
	"bufio"
	"fmt"
	"net"
	"os"

	"github.com/5aradise/cli-chat/client/internal/cli"
)

const (
	bufferSize     = 256
	maxUsernameLen = 10
	maxMsgLen      = 106
)

type client struct {
	conn       net.Conn
	readBuf    []byte
	printLn    *int
	isInChat   bool
	chatColors map[string]cli.Color
}

func New(address string) (*client, error) {
	conn, err := net.Dial("tcp4", address)
	if err != nil {
		return nil, err
	}
	printLn := 1
	readBuf := make([]byte, bufferSize)
	return &client{conn, readBuf, &printLn, false, nil}, nil
}

func (c *client) Run() {
	c.updateScreen()

	scanner := bufio.NewScanner(os.Stdin)
	c.authClient(scanner)

	go c.listenServer()
	c.listenClient(scanner)
}

func (c *client) authClient(scanner *bufio.Scanner) {
	c.printf(formatSystemMsg("Enter name"))
	for {
		input, inputLen := cli.Scan(scanner)
		if inputLen == 0 {
			continue
		}
		if inputLen > maxUsernameLen {
			c.printf(formatSystemMsg("username is too long (maximum 10 characters)"))
			continue
		}
		c.write(authAcc, []byte(input))
		
		head, body := c.read()
		if head == authAcc {
			c.printf(formatSystemMsg("User with id " + string(body) + " have been created"))
			break
		}
		c.printf(formatSystemMsg(body))
	}
	c.printf(formatSystemMsg("Type /help to see all available commands"))
}

func (c *client) listenClient(scanner *bufio.Scanner) {
	for {
		input, inputLen := cli.Scan(scanner)
		if inputLen == 0 {
			continue
		}
		if inputLen > cli.MaxInputLen {
			c.printf(formatSystemMsg("your message is too long (maximum 106 characters)"))
			continue
		}
		err := c.processReq(input)
		if err != nil {
			c.printf(formatSystemMsg(err.Error()))
		}
	}
}

func (c *client) listenServer() {
	var head header
	var body []byte
	for {
		head, body = c.read()
		c.processResp(head, body)
	}
}

func (c *client) shutDown(msg string) {
	c.conn.Close()
	cli.ClearConsole()
	fmt.Println(msg)
	fmt.Scanln()
	os.Exit(0)
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

func (c *client) write(h header, b []byte) {
	_, err := c.conn.Write(h.setHeader(b))
	if err != nil {
		c.shutDown("you've been disconnected from the server")
	}
}

func (c *client) read() (header, []byte) {
	l, err := c.conn.Read(c.readBuf)
	if err != nil {
		c.shutDown("you've been disconnected from the server")
	}
	return getHeader(c.readBuf[:l])
}
