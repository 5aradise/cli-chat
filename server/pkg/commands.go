package chat

import (
	"errors"
)

var commands map[header]func(*server, *user, []byte) error = map[header]func(*server, *user, []byte) error{
	userMsg:     (*server).msgToChat,
	createChat:  (*server).createChat,
	connectChat: (*server).connChat,
	exitChat:    (*server).exitChat,
}

func (s *server) msgToChat(user *user, args []byte) error {
	if user.currChat == nil {
		return errors.New("you are not connected to any chat")
	}
	isValid, reas := isValidMsg(string(args))
	if !isValid {
		return errors.New(reas)
	}

	user.currChat.c <- &message{user, args}

	return nil
}

func (s *server) createChat(user *user, args []byte) error {
	if user.currChat != nil {
		return errors.New("to connect to chat you must leave current")
	}

	chat, err := s.newChat(string(args))
	if err != nil {
		return err
	}
	chat.addUser(user)

	user.write(connectChat, args)
	return nil
}

func (s *server) connChat(user *user, args []byte) error {
	if user.currChat != nil {
		return errors.New("to connect to chat you must leave current")
	}

	s.chatsMux.RLock()
	chat, ok := s.chats[string(args)]
	s.chatsMux.RUnlock()
	if !ok {
		return errors.New("wrong chat name")
	}

	err := chat.addUser(user)
	if err != nil {
		return err
	}

	user.write(connectChat, args)
	return nil
}

func (s *server) exitChat(user *user, args []byte) error {
	if user.currChat == nil {
		return errors.New("you are not in the chat")
	}

	user.currChat.deleteUser(string(user.name))

	user.write(exitChat, nil)
	return nil
}
