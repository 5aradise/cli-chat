package client

import (
	"errors"
	"strings"
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

func (c *client) processResp(h header, b []byte) {
	command, ok := serverCommands[h]
	if !ok {
		return
	}
	command(c, b)
}
