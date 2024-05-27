package client

import (
	"errors"

	"github.com/5aradise/cli-chat/client/internal/cli"
)

var userCommands map[string]func(*client, []string) error = map[string]func(*client, []string) error{
	"help":   (*client).helpReq,
	"create": (*client).chatCreateReq,
	"conn":   (*client).chatConnReq,
	"admin":  (*client).passAdminReq,
	"exit":   (*client).chatExitReq,
	"delete": (*client).chatDeleteReq,
}

func (c *client) chatCreateReq(args []string) error {
	if c.isInChat {
		return errors.New("to create chat you must leave current")
	}

	if len(args) == 0 {
		return errors.New("to many arguments")
	}

	chatName := args[0]

	isValid, reas := isValidChatName(chatName)
	if !isValid {
		return errors.New(reas)
	}

	c.write(createChat, []byte(chatName))
	return nil
}

func (c *client) chatConnReq(args []string) error {
	if c.isInChat {
		return errors.New("to connect to chat you must leave current")
	}

	if len(args) == 0 {
		return errors.New("to many arguments")
	}

	chatName := args[0]

	isValid, reas := isValidChatName(chatName)
	if !isValid {
		return errors.New(reas)
	}

	c.write(connectChat, []byte(chatName))
	return nil
}

func (c *client) chatExitReq(args []string) error {
	if !c.isInChat {
		return errors.New("you are not in the chat")
	}

	c.write(exitChat, nil)
	return nil
}

func (c *client) passAdminReq(args []string) error {
	if !c.isInChat {
		return errors.New("you are not in the chat")
	}

	if !c.isAdmin {
		return errors.New("you do not have permission")
	}

	if len(args) == 0 {
		return errors.New("to many arguments")
	}

	newAdmin := args[0]

	isValid, reas := isValidUsername(newAdmin)
	if !isValid {
		return errors.New(reas)
	}

	c.write(passAdmin, []byte(newAdmin))
	return nil
}

func (c *client) chatDeleteReq(args []string) error {
	if !c.isInChat {
		return errors.New("you are not in the chat")
	}

	if !c.isAdmin {
		return errors.New("you do not have permission")
	}

	c.write(deleteChat, nil)
	return nil
}

func (c *client) helpReq(args []string) error {
	c.printf(formatSystemMsg("/create {chat name} - creates and connects to new chat room"))
	c.printf(cli.Colorize("            /conn {chat name}   - connects to chat room", cli.Red))
	c.printf(cli.Colorize("            /exit               - exits current chat room", cli.Red))
	c.printf(cli.Colorize("            /help               - shows a list of commands", cli.Red))
	return nil
}

func (c *client) sendMsg(msg string) error {
	if !c.isInChat {
		return errors.New("you are not connected to any chat")
	}
	isValid, reas := isValidMsg(msg)
	if !isValid {
		return errors.New(reas)
	}

	c.write(userMsg, []byte(msg))
	c.printf(formatClientMsg(msg))
	return nil
}
