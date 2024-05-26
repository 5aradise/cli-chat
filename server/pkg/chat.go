package chat

import (
	"errors"
	"log"
	"sync"
)

type chat struct {
	name  string
	c     chan *message
	users map[string]*user
	mux   sync.RWMutex
}

type message struct {
	sender *user
	text   []byte
}

var (
	addMsg    = []byte(" has been added")
	deleteMsg = []byte(" left the chat room")
)

func (s *server) newChat(name string) (*chat, error) {
	isValid, reas := isValidChatName(name)
	if !isValid {
		return nil, errors.New(reas)
	}

	s.chatsMux.Lock()
	if _, ok := s.chats[name]; ok {
		s.chatsMux.Unlock()
		return nil, errors.New("chat with this name already exist")
	}

	chat := &chat{
		name:  name,
		c:     make(chan *message, 16),
		users: make(map[string]*user),
		mux:   sync.RWMutex{},
	}

	s.chats[name] = chat
	s.chatsMux.Unlock()

	go chat.broadcast()

	log.Printf("New chat: %s\n", name)

	return chat, nil
}

func (ch *chat) addUser(u *user) error {
	ch.mux.Lock()
	if _, ok := ch.users[u.name]; ok {
		ch.mux.Unlock()
		return errors.New("user with this id already in chat")
	}
	ch.users[u.name] = u
	ch.mux.Unlock()
	u.currChat = ch

	ch.chatCall(append([]byte(u.name), addMsg...))
	return nil
}

func (ch *chat) deleteUser(name string) {
	ch.mux.Lock()
	u := ch.users[name]
	delete(ch.users, name)
	ch.mux.Unlock()
	u.currChat = nil

	ch.chatCall(append([]byte(u.name), deleteMsg...))
}

func (ch *chat) chatCall(msg []byte) {
	ch.mux.RLock()
	defer ch.mux.RUnlock()

	for _, dst := range ch.users {
		dst.write(chatMsg, msg)
	}
}

func (ch *chat) broadcast() {
	const userMsgDiv byte = 0x00
	for msg := range ch.c {
		toSend := append([]byte(msg.sender.name), userMsgDiv)
		toSend = append(toSend, msg.text...)

		ch.mux.RLock()
		for _, dst := range ch.users {
			if dst != msg.sender {
				dst.write(userMsg, toSend)
			}
		}
		ch.mux.RUnlock()
	}
}
