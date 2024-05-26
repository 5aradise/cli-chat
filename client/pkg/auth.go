package client

import (
	"bufio"

	"github.com/5aradise/cli-chat/client/internal/cli"
)

func (c *client) authClient(scanner *bufio.Scanner) {
	c.printf(formatSystemMsg("Enter name"))
	for {
		username, inputLen := cli.Scan(scanner)
		if inputLen == 0 {
			continue
		}

		isValid, reas := isValidUsername([]byte(username))
		if !isValid {
			c.printf(formatSystemMsg(reas))
			continue
		}

		c.write(authAcc, []byte(username))
		head, body := c.read()
		if head == authAcc {
			c.printf(formatSystemMsg("User with id " + string(body) + " have been created"))
			break
		}
		c.printf(formatSystemMsg(body))
	}
	c.printf(formatSystemMsg("Type /help to see all available commands"))
}
