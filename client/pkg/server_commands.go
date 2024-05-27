package client

import (
	"github.com/5aradise/cli-chat/client/internal/cli"
)

var serverCommands map[header]func(*client, []byte) = map[header]func(*client, []byte){
	systemMsg:   (*client).systemMsg,
	chatMsg:     (*client).chatMsg,
	userMsg:     (*client).userMsg,
	connectChat: (*client).chatConnResp,
	passAdmin:   (*client).passAdminResp,
	exitChat:    (*client).chatExitResp,
}

func (c *client) systemMsg(args []byte) {
	c.printf(formatSystemMsg(args))
}

func (c *client) chatMsg(args []byte) {
	c.printf(formatChatMsg(args))
}

func (c *client) userMsg(args []byte) {
	msg, err := c.formatUserMsg(args)
	if err != nil {
		return
	}
	c.printf(msg)
}

func (c *client) chatConnResp(args []byte) {
	c.updateScreen()
	c.isInChat = true
	c.chatColors = make(map[string]cli.Color)
}

func (c *client) passAdminResp(args []byte) {
	switch args[0] {
	case 0:
		if c.isAdmin {
			c.printf(formatChatMsg("you're no longer a chat room administrator"))
			c.isAdmin = false
		}
	case 1:
		if !c.isAdmin {
			c.printf(formatChatMsg("you are new admin of this chat"))
			c.isAdmin = true
		}
	}
}

func (c *client) chatExitResp(args []byte) {
	c.updateScreen()
	c.printf(formatSystemMsg("you have been deleted from the chat room"))
	c.isInChat = false
	c.chatColors = nil
}
