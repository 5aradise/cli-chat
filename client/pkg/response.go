package client

import "github.com/5aradise/cli-chat/client/internal/cli"

var serverCommands map[byte]func(*Client, []byte) = map[byte]func(*Client, []byte){
	systemMsgCode: (*Client).systemMsg,
	chatMsgCode:   (*Client).chatMsg,
	userMsgCode:   (*Client).userMsg,
	connCode:      (*Client).chatConnResp,
	exitCode:      (*Client).chatExitResp,
}

func (c *Client) processResp(b []byte) {
	header, args := b[0], b[1:]
	command := serverCommands[header]
	command(c, args)
}

func (c *Client) systemMsg(args []byte) {
	c.printf(formatSystemMsg(args))
}

func (c *Client) chatMsg(args []byte) {
	c.printf(formatChatMsg(args))
}

func (c *Client) userMsg(args []byte) {
	c.printf(c.formatUserMsg(args))
}

func (c *Client) chatConnResp(args []byte) {
	c.updateScreen()
	c.isInChat = true
	c.chatColors = make(map[string]cli.Color)
}

func (c *Client) chatExitResp(args []byte) {
	c.updateScreen()
	c.isInChat = false
	c.chatColors = nil
}
