package client

import (
	"errors"
	"fmt"
	"slices"
	"unicode/utf8"

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

	user, msg := b[:div], b[div+1:]

	if !isValidUserMsg(user, msg) {
		return "", errors.New("invalid user or message length")
	}

	userS, msgS := string(user), string(msg)

	userColor, ok := c.chatColors[userS]
	if !ok {
		randColor := cli.RandColor()
		for slices.Index(reservedColors, randColor) != -1 {
			randColor = cli.RandColor()
		}
		c.chatColors[userS] = randColor
		userColor = randColor
	}

	userS = fmt.Sprintf("%-10s", userS)
	userS = cli.Colorize(userS, userColor)

	return userS + ": " + msgS, nil
}

func formatSystemMsg(a any) string {
	switch msg := a.(type) {
	case string:
		return cli.Colorize("System    : "+msg, cli.RedS)
	case []byte:
		return cli.Colorize("System    : "+string(msg), cli.RedS)
	}
	return ""
}

func formatChatMsg(b []byte) string {
	return cli.Colorize("Chat      : "+string(b), cli.Red)
}

func formatClientMsg(s string) string {
	return cli.Colorize("You       : "+s, cli.WhiteS)
}

func isValidUserMsg(user, msg []byte) bool {
	return utf8.RuneCount(user) <= maxUsernameLen && utf8.RuneCount(msg) <= maxMsgLen
}
