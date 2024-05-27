package client

import (
	"errors"
	"fmt"
	"slices"

	"github.com/5aradise/cli-chat/client/internal/cli"
)

var reservedColors []cli.Color = []cli.Color{cli.Red, cli.RedS, cli.WhiteS}

func (c *client) formatUserMsg(b []byte) (string, error) {
	const userMsgDiv byte = 0x00

	if !c.isInChat {
		return "", errors.New("you are not connected to any chat")
	}

	div := slices.Index(b, userMsgDiv)
	if div == -1 {
		return "", errors.New("invalid user message")
	}

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

	user = fmt.Sprintf("%-10s", user)
	user = cli.Colorize(user, userColor)

	return user + ": " + msg, nil
}

func formatSystemMsg(a any) string {
	switch msg := a.(type) {
	case string:
		return cli.Colorize("System    : "+msg, cli.Red)
	case []byte:
		return cli.Colorize("System    : "+string(msg), cli.Red)
	}
	return ""
}

func formatChatMsg(a any) string {
	switch msg := a.(type) {
	case string:
		return cli.Colorize("Chat      : "+msg, cli.RedS)
	case []byte:
		return cli.Colorize("Chat      : "+string(msg), cli.RedS)
	}
	return ""
}

func formatClientMsg(a any) string {
	switch msg := a.(type) {
	case string:
		return cli.Colorize("You       : "+msg, cli.WhiteS)
	case []byte:
		return cli.Colorize("You       : "+string(msg), cli.WhiteS)
	}
	return ""
}
