package chat

import (
	"errors"
	"fmt"
	"log"
	"sync"
)

type Chat struct {
	id    int
	users map[int]*User
	mux   sync.RWMutex
}

func (s *Server) NewChat(id int) *Chat {
	chat := &Chat{
		id:    id,
		users: make(map[int]*User),
		mux:   sync.RWMutex{},
	}

	s.chatsMux.Lock()
	s.chats[id] = chat
	s.chatsMux.Unlock()

	fmt.Printf("New chat: %d\n", id)

	return chat
}

func (ch *Chat) AddUser(u *User) error {
	ch.mux.Lock()

	if _, ok := ch.users[u.id]; ok {
		ch.mux.Unlock()
		return errors.New("User with this id already exist")
	}
	ch.users[u.id] = u
	u.currChat = ch
	ch.mux.Unlock()

	ch.ChatCall(u.name + " has been added")
	return nil
}

func (ch *Chat) DeleteUser(id int) {
	ch.mux.Lock()
	defer ch.mux.Unlock()

	u := ch.users[id]
	ch.ChatCall(u.name + " left the chat room")
	u.currChat = nil
	delete(ch.users, id)
}

func (ch *Chat) ChatCall(msg string) {
	ch.mux.RLock()
	defer ch.mux.RUnlock()

	toSend := append([]byte{0, 1}, []byte(msg)...)

	for _, dst := range ch.users {
		_, err := dst.Write(toSend)
		if err != nil {
			log.Println(err)
		}
	}
}

func (ch *Chat) Write(src *User, msg []byte) {
	ch.mux.RLock()
	defer ch.mux.RUnlock()

	const userMsgDiv byte = 0

	toSend := append([]byte(src.name), userMsgDiv)
	toSend = append(toSend, msg...)

	for _, dst := range ch.users {
		if dst != src {
			_, err := dst.Write(toSend)
			if err != nil {
				log.Println(err)
			}
		}
	}
}
