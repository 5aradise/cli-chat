package client

import (
	"slices"

	"github.com/5aradise/cli-chat/client/internal/cli"
)

func (c *Client) formatUserMsg(b []byte) string {
	div := slices.Index(b, 0x00)
	if div == -1 {
		return ""
	}
	var reservedColors []cli.Color = []cli.Color{cli.Red, cli.RedS, cli.White, cli.WhiteS}

	user, msg := string(b[:div]), string(b[div+1:])

	userColor, ok := c.chatColors[user]
	if !ok {
		randColor := cli.RandColor()
		for slices.Index(reservedColors, randColor) != -1 {
			randColor = cli.RandColor()
		}
		c.chatColors[user] = randColor
		userColor = randColor
	}

	user = cli.Colorize(user, userColor)

	return user + ": " + msg
}

func formatSystemMsg(b []byte) string {
	return cli.Colorize("System: "+string(b), cli.Red)
}

func formatChatMsg(b []byte) string {
	return cli.Colorize("Chat: "+string(b), cli.RedS)
}
