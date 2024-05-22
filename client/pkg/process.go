package client

import (
	"errors"
	"strings"

	"github.com/5aradise/cli-chat/client/internal/cli"
)

func (c *client) processReq(input string) error {
	input = strings.Trim(input, " ")
	if input[:1] == "/" {
		if len(input) == 1 {
			return errors.New("unknown command")
		}
		splited := strings.Split(input[1:], " ")
		commandName := splited[0]
		args := make([]string, 0)
		if len(splited) != 1 {
			args = splited[1:]
		}
		command, ok := userCommands[commandName]
		if !ok {
			return errors.New("unknown command")
		}
		err := command(c, args)
		if err != nil {
			return err
		}
		return nil
	}
	return c.sendMsg(input)
}

func (c *client) sendMsg(msg string) error {
	if !c.isInChat {
		return errors.New("you are not connected to any chat")
	}

	c.printf(cli.Colorize("You: "+msg, cli.WhiteS))

	req := userMsg.setHeaderS(msg)

	_, err := c.Write(req)
	return err
}

func (c *client) processResp(b []byte) {
	header, args := header(b[0]), b[1:]
	command := serverCommands[header]
	command(c, args)
}
