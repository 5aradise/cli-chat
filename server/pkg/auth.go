package chat

import (
	"net"
	"regexp"
	"slices"
	"strconv"
	"unicode/utf8"
)

var latinAndCyrillicLetters = regexp.MustCompile("[A-Za-zА-яІіЇїЄє]")
var reservedNames = [][]byte{[]byte("You"), []byte("Chat"), []byte("System")}

func (s *server) authUser(conn net.Conn) (*user, error) {
	buf := make([]byte, bufferSize)
	var head header
	var username []byte
	for {
		l, err := conn.Read(buf)
		if err != nil {
			return nil, err
		}

		head, username = getHeader(buf[:l])
		if head != authAcc {
			_, err = conn.Write(authRej.setHeader([]byte("invalid request")))
			if err != nil {
				return nil, err
			}
			continue
		}

		isValid, reas := isValidUsername(username)
		if !isValid {
			_, err = conn.Write(authRej.setHeader([]byte(reas)))
			if err != nil {
				return nil, err
			}
			continue
		}

		break
	}

	user := s.newUser(username, conn)
	_, err := conn.Write(authAcc.setHeader([]byte(strconv.Itoa(user.id))))
	if err != nil {
		return nil, err
	}

	return user, nil
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
