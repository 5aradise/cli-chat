package client

import (
	"errors"
	"strings"
)

func (c *client) processReq(input string) error {
	input = strings.Trim(input, " ")
	if len(input) == 0 {
		return nil
	}
	if input[:1] == "/" {
		if len(input) == 1 {
			return errors.New("unknown command")
		}

		splited := strings.Split(input[1:], " ")
		commandName := splited[0]
		args := splited[1:]

		command, ok := userCommands[commandName]
		if !ok {
			return errors.New("unknown command")
		}

		return command(c, args)
	}

	return c.sendMsg(input)
}

func (c *client) processResp(h header, b []byte) {
	command, ok := serverCommands[h]
	if !ok {
		return
	}

	command(c, b)
}
