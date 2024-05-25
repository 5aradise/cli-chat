package chat

import (
	"errors"
	"strconv"
	"unicode/utf8"
)

var commands map[header]func(*server, *user, []byte) error = map[header]func(*server, *user, []byte) error{
	userMsg: (*server).msgToChat,
	create:  (*server).createChat,
	connect: (*server).connChat,
	exit:    (*server).exitChat,
}

func (s *server) msgToChat(user *user, args []byte) error {
	if user.currChat == nil {
		return errors.New("you are not connected to any chat")
	}
	if utf8.RuneCount(args) > maxMsgLen {
		return errors.New("your message is too long (maximum 106 characters)")
	}
	user.currChat.c <- &message{user, args}

	return nil
}

func (s *server) createChat(user *user, args []byte) error {
	if user.currChat != nil {
		return errors.New("to connect to chat you must leave current")
	}

	chatId, err := strconv.Atoi(string(args))
	if err != nil {
		return err
	}

	s.chatsMux.RLock()
	_, ok := s.chats[chatId]
	s.chatsMux.RUnlock()
	if ok {
		return errors.New("chat with this id already exist")
	}

	chat := s.newChat(chatId)
	err = chat.addUser(user)
	if err != nil {
		return err
	}

	user.write(connect, nil)
	return nil
}

func (s *server) connChat(user *user, args []byte) error {
	if user.currChat != nil {
		return errors.New("to connect to chat you must leave current")
	}

	chatId, err := strconv.Atoi(string(args))
	if err != nil {
		return err
	}

	s.chatsMux.RLock()
	chat, ok := s.chats[chatId]
	s.chatsMux.RUnlock()
	if !ok {
		return errors.New("wrong chat id")
	}

	err = chat.addUser(user)
	if err != nil {
		return err
	}

	user.write(connect, nil)
	return nil
}

func (s *server) exitChat(user *user, args []byte) error {
	if user.currChat == nil {
		return errors.New("you are not in the chat")
	}

	user.currChat.deleteUser(user.id)

	user.write(exit, nil)
	return nil
}
