package client

import (
	"bufio"
	"regexp"
	"slices"
	"unicode/utf8"

	"github.com/5aradise/cli-chat/client/internal/cli"
)

var latinAndCyrillicLetters = regexp.MustCompile("[A-Za-zА-яІіЇїЄє]")
var reservedNames = [][]byte{[]byte("You"), []byte("Chat"), []byte("System")}

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

func isValidUsername(name []byte) (bool, string) {
	if utf8.RuneCount(name) > maxUsernameLen {
		return false, "username is too long (maximum 10 characters)"
	}

	if slices.Contains(name, 0x20) {
		return false, "username mustn't contain spaces"
	}

	for _, reservedName := range reservedNames {
		if slices.Equal(reservedName, name) {
			return false, "username is equal to reserved name"
		}
	}

	if !latinAndCyrillicLetters.Match(name) {
		return false, "username must contain at least 1 letter"
	}

	return true, ""
}
