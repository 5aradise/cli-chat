package client

import (
	"errors"
	"strconv"
	"unicode/utf8"

	"github.com/5aradise/cli-chat/client/internal/cli"
)

var userCommands map[string]func(*client, []string) error = map[string]func(*client, []string) error{
	"help":   (*client).helpReq,
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

	c.write(create, []byte(strconv.Itoa(chatId)))
	return nil
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

	c.write(connect, []byte(strconv.Itoa(chatId)))
	return nil
}

func (c *client) chatExitReq(args []string) error {
	if !c.isInChat {
		c.printf(formatSystemMsg("you are not in the chat"))
		return nil
	}

	c.write(exit, nil)
	return nil
}

func (c *client) helpReq(args []string) error {
	c.printf(formatSystemMsg("create {chat id} - creates and connects to new chat room"))
	c.printf(cli.Colorize("            conn {chat id}   - connects to chat room", cli.RedS))
	c.printf(cli.Colorize("            exit             - exits current chat room", cli.RedS))
	return nil
}

func (c *client) sendMsg(msg string) error {
	if !c.isInChat {
		return errors.New("you are not connected to any chat")
	}
	if utf8.RuneCountInString(msg) > maxMsgLen {
		return errors.New("your message is too long (maximum 106 characters)")
	}

	c.write(userMsg, []byte(msg))
	c.printf(formatClientMsg(msg))
	return nil
}
