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
	u := ch.users[id]
	delete(ch.users, id)
	ch.mux.Unlock()
	u.currChat = nil

	ch.ChatCall(u.name + " left the chat room")
}

func (ch *Chat) ChatCall(msg string) {
	ch.mux.RLock()
	defer ch.mux.RUnlock()

	toSend := append([]byte{chatMsgCode}, []byte(msg)...)

	for _, dst := range ch.users {
		_, err := dst.Write(toSend)
		if err != nil {
			log.Println(err)
		}
	}
}

func (ch *Chat) Write(src *User, msg []byte) {
	const userMsgDiv byte = 0x00

	ch.mux.RLock()
	defer ch.mux.RUnlock()

	toSend := append([]byte{userMsgCode}, []byte(src.name)...)
	toSend = append(toSend, userMsgDiv)
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
