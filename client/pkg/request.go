package client

import (
	"errors"
	"strconv"
	"strings"

	"github.com/5aradise/cli-chat/client/internal/cli"
)

func (c *Client) processReq(input string) error {
	input = strings.Trim(input, " ")
	if input[:1] == "/" {
		if len(input) == 1 {
			c.printf(cli.Colorize("System: unknown command", cli.RedS))
			return nil
		}
		splited := strings.Split(input[1:], " ")
		commandName := splited[0]
		args := make([]string, 0)
		if len(splited) != 1 {
			args = splited[1:]
		}
		command, ok := userCommands[commandName]
		if !ok {
			c.printf(cli.Colorize("System: unknown command", cli.RedS))
			return nil
		}
		err := command(c, args)
		if err != nil {
			c.printf(cli.Colorize("System: "+err.Error(), cli.RedS))
		}
		return nil
	}
	err := c.sendMsg(input)
	if err != nil {
		c.printf(cli.Colorize("System: "+err.Error(), cli.RedS))
	}
	return nil
}

func (c *Client) sendMsg(msg string) error {
	if !c.isInChat {
		return errors.New("you are not connected to any chat")
	}

	c.printf(cli.Colorize("You: "+msg, cli.WhiteS))

	req := append([]byte{userMsgCode}, []byte(msg)...)

	_, err := c.Write(req)
	return err
}

var userCommands map[string]func(*Client, []string) error = map[string]func(*Client, []string) error{
	"create": (*Client).chatCreateReq,
	"conn":   (*Client).chatConnReq,
	"exit":   (*Client).chatExitReq,
}

func (c *Client) chatCreateReq(args []string) error {
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

	req := strconv.AppendInt([]byte{createCode}, int64(chatId), 10)

	_, err = c.Write(req)
	return err
}

func (c *Client) chatConnReq(args []string) error {
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

	req := strconv.AppendInt([]byte{connCode}, int64(chatId), 10)

	_, err = c.Write(req)
	return err
}

func (c *Client) chatExitReq(args []string) error {
	if !c.isInChat {
		c.printf(cli.Colorize("System: you are not in the chat", cli.RedS))
		return nil
	}

	req := []byte{exitCode}

	_, err := c.Write(req)
	return err
}
