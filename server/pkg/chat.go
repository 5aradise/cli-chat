package chat

import (
	"errors"
	"log"
	"sync"
)

type chat struct {
	id    int
	c     chan *message
	users map[int]*user
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

func (s *server) newChat(id int) *chat {
	chat := &chat{
		id:    id,
		c:     make(chan *message, 16),
		users: make(map[int]*user),
		mux:   sync.RWMutex{},
	}

	s.chatsMux.Lock()
	s.chats[id] = chat
	s.chatsMux.Unlock()

	go chat.broadcast()

	log.Printf("New chat: %d\n", id)

	return chat
}

func (ch *chat) addUser(u *user) error {
	ch.mux.Lock()
	if _, ok := ch.users[u.id]; ok {
		ch.mux.Unlock()
		return errors.New("user with this id already exist")
	}
	ch.users[u.id] = u
	u.currChat = ch
	ch.mux.Unlock()

	ch.chatCall(append(u.name, addMsg...))
	return nil
}

func (ch *chat) deleteUser(id int) {
	ch.mux.Lock()
	u := ch.users[id]
	delete(ch.users, id)
	ch.mux.Unlock()
	u.currChat = nil

	ch.chatCall(append(u.name, deleteMsg...))
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
		toSend := append(msg.sender.name, userMsgDiv)
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
