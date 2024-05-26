package chat

import (
	"regexp"
	"slices"
	"strings"
	"unicode/utf8"
)

var (
	latinAndCyrillicLetters = regexp.MustCompile("[A-Za-zА-яІіЇїЄє]")
	reservedNames           = []string{"You", "Chat", "System"}
)

const (
	maxUsernameLen = 10
	maxChatNameLen = 20
	maxMsgLen      = 106
)

func isValidUsername(name string) (bool, string) {
	if utf8.RuneCountInString(name) > maxUsernameLen {
		return false, "username is too long (maximum 10 characters)"
	}

	if strings.Contains(name, " ") {
		return false, "username mustn't contain spaces"
	}

	if slices.Contains(reservedNames, name) {
		return false, "username is equal to reserved name"
	}

	if !latinAndCyrillicLetters.MatchString(name) {
		return false, "username must contain at least 1 letter"
	}

	return true, ""
}

func isValidChatName(name string) (bool, string) {
	if utf8.RuneCountInString(name) > maxChatNameLen {
		return false, "username is too long (maximum 20 characters)"
	}

	if strings.Contains(name, " ") {
		return false, "chat name mustn't contain spaces"
	}

	if !latinAndCyrillicLetters.MatchString(name) {
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
