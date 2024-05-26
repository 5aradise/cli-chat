package client

import "github.com/5aradise/cli-chat/client/internal/cli"

var serverCommands map[header]func(*client, []byte) = map[header]func(*client, []byte){
	systemMsg:   (*client).systemMsg,
	chatMsg:     (*client).chatMsg,
	userMsg:     (*client).userMsg,
	connectChat: (*client).chatConnResp,
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
	c.printf(formatSystemMsg("you have been added to the chat room `" + string(args) + "`"))
	c.isInChat = true
	c.chatColors = make(map[string]cli.Color)
}

func (c *client) chatExitResp(args []byte) {
	c.updateScreen()
	c.printf(formatSystemMsg("you have been deleted from the chat room"))
	c.isInChat = false
	c.chatColors = nil
}
