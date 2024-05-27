package chat

import (
	"errors"
)

var commands map[header]func(*server, *user, []byte) error = map[header]func(*server, *user, []byte) error{
	userMsg:     (*server).msgToChat,
	createChat:  (*server).createChat,
	connectChat: (*server).connChat,
	passAdmin:   (*server).passAdmin,
	kickUser:    (*server).kickUser,
	exitChat:    (*server).exitChat,
	deleteChat:  (*server).deleteChatCommmand,
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

	return chat.addUser(user)
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

	return chat.addUser(user)
}

func (s *server) kickUser(user *user, args []byte) error {
	if user.currChat == nil {
		return errors.New("you are not in the chat")
	}

	if user != user.currChat.admin {
		return errors.New("you do not have permission")
	}

	userToDelete := string(args)
	return user.currChat.deleteUser(userToDelete)
}

func (s *server) exitChat(user *user, args []byte) error {
	if user.currChat == nil {
		return errors.New("you are not in the chat")
	}

	return user.currChat.deleteUser(user.name)
}

func (s *server) passAdmin(user *user, args []byte) error {
	if user.currChat == nil {
		return errors.New("you are not in the chat")
	}

	if user != user.currChat.admin {
		return errors.New("you do not have permission")
	}

	if len(args) == 0 {
		return errors.New("to many arguments")
	}

	newAdmin := string(args)
	return user.currChat.setAdmin(newAdmin)
}

func (s *server) deleteChatCommmand(user *user, args []byte) error {
	if user.currChat == nil {
		return errors.New("you are not in the chat")
	}

	if user != user.currChat.admin {
		return errors.New("you do not have permission")
	}

	return s.deleteChat(user.currChat.name)
}
