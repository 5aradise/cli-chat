package client

import (
	"errors"
	"strconv"
)

var userCommands map[string]func(*client, []string) error = map[string]func(*client, []string) error{
	"create": (*client).chatCreateReq,
	"conn":   (*client).chatConnReq,
	"exit":   (*client).chatExitReq,
}

func (c *client) chatCreateReq(args []string) error {
	if c.isInChat {
		return errors.New("to create chat you must leave current")
	}

	if len(args) == 0 {
		return errors.New("to many arguments")
	}

	chatId, err := strconv.Atoi(args[0])
	if err != nil {
		return errors.New("wrong command (must be int)")
	}

	req := create.setHeaderS(strconv.Itoa(chatId))

	_, err = c.Write(req)
	return err
}

func (c *client) chatConnReq(args []string) error {
	if c.isInChat {
		return errors.New("to connect to chat you must leave current")
	}

	if len(args) == 0 {
		return errors.New("to many arguments")
	}

	chatId, err := strconv.Atoi(args[0])
	if err != nil {
		return errors.New("wrong command (must be int)")
	}

	req := connect.setHeaderS(strconv.Itoa(chatId))

	_, err = c.Write(req)
	return err
}

func (c *client) chatExitReq(args []string) error {
	if !c.isInChat {
		c.printf(formatSystemMsg("you are not in the chat"))
		return nil
	}

	req := exit.setHeaderB([]byte{0})

	_, err := c.Write(req)
	return err
}

func (c *client) sendMsg(msg string) error {
	if !c.isInChat {
		return errors.New("you are not connected to any chat")
	}

	c.printf(formatClientMsg(msg))

	req := userMsg.setHeaderS(msg)

	_, err := c.Write(req)
	return err
}
