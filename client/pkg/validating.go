package client

import (
	"regexp"
	"slices"
	"unicode/utf8"
)

var (
	latinAndCyrillicLetters = regexp.MustCompile("[A-Za-zА-яІіЇїЄє]")
	reservedNames           = [][]byte{[]byte("You"), []byte("Chat"), []byte("System")}
)

const (
	maxUsernameLen = 10
	maxChatNameLen = 20
	maxMsgLen      = 106
)

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

func isValidChatName(name []byte) (bool, string) {
	if utf8.RuneCount(name) > maxChatNameLen {
		return false, "username is too long (maximum 20 characters)"
	}

	if !latinAndCyrillicLetters.Match(name) {
		return false, "chat name must contain at least 1 letter"
	}

	return true, ""
}

func isValidMsg(msg string) (bool, string) {
	if utf8.RuneCountInString(msg) > maxMsgLen {
		return false, "your message is too long (maximum 106 characters)"
	}
	return true, ""
}
