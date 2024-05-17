package chat

import (
	"errors"
	"strconv"
)

var commands map[byte]func(*Server, *User, []byte) error = map[byte]func(*Server, *User, []byte) error{
	userMsgCode: (*Server).msgToChat,
	createCode:  (*Server).createChat,
	connCode:    (*Server).connChat,
	exitCode:    (*Server).exitChat,
}

func (s *Server) msgToChat(user *User, args []byte) error {
	if user.currChat == nil {
		return errors.New("you are not connected to any chat")
	}
	user.currChat.Write(user, args)

	return nil
}

func (s *Server) createChat(user *User, args []byte) error {
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

	chat := s.NewChat(chatId)
	err = chat.AddUser(user)
	if err != nil {
		return err
	}

	user.Write([]byte{connCode, 0x00})
	return nil
}

func (s *Server) connChat(user *User, args []byte) error {
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

	err = chat.AddUser(user)
	if err != nil {
		return err
	}

	user.Write([]byte{connCode, 0x00})
	return nil
}

func (s *Server) exitChat(user *User, args []byte) error {
	if user.currChat == nil {
		return errors.New("you are not in the chat")
	}

	user.currChat.DeleteUser(user.id)

	user.Write([]byte{exitCode, 0x00})
	return nil
}
