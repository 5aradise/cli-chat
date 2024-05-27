package client

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

func isValidUsername(name string) (isValid bool, reason string) {
	if utf8.RuneCountInString(name) > maxUsernameLen {
		reason = "username is too long (maximum 10 characters)"
		return
	}

	if strings.Contains(name, " ") {
		reason = "username mustn't contain spaces"
		return
	}

	if slices.Contains(reservedNames, name) {
		reason = "username is equal to reserved name"
		return
	}

	if !latinAndCyrillicLetters.MatchString(name) {
		reason = "username must contain at least 1 letter"
		return
	}

	isValid = true
	return
}

func isValidChatName(name string) (isValid bool, reason string) {
	if utf8.RuneCountInString(name) > maxChatNameLen {
		reason = "username is too long (maximum 20 characters)"
		return
	}

	if strings.Contains(name, " ") {
		reason = "chat name mustn't contain spaces"
		return
	}

	if !latinAndCyrillicLetters.MatchString(name) {
		reason = "chat name must contain at least 1 letter"
		return
	}

	isValid = true
	return
}

func isValidMsg(msg string) (isValid bool, reason string) {
	if utf8.RuneCountInString(msg) > maxMsgLen {
		return false, "your message is too long (maximum 106 characters)"
	}
	return true, ""
}
