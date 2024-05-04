package chat

import (
	"errors"
	"log"
	"sync"
)

type Chat struct {
	id    int
	users map[int]*User
	mux   sync.RWMutex
}

func NewChat(id int) *Chat {
	chat := &Chat{
		id:    id,
		users: make(map[int]*User),
		mux:   sync.RWMutex{},
	}
	return chat
}

func (ch *Chat) AddUser(u *User) error {
	ch.mux.Lock()

	if _, ok := ch.users[u.id]; ok {
		ch.mux.Unlock()
		return errors.New("User with this id already exist")
	}
	ch.users[u.id] = u
	ch.mux.Unlock()

	ch.SystemCall(u.name + " has been added")
	go func() {
		for {
			msg, err := u.Read()
			if err != nil {
				break
			}
			ch.Write(u, msg)
		}
		u.conn.Close()
		ch.DeleteUser(u.id)
		ch.SystemCall(u.name + " leaved")
	}()
	return nil
}

func (ch *Chat) DeleteUser(id int) {
	ch.mux.Lock()
	defer ch.mux.Unlock()

	delete(ch.users, id)
}

func (ch *Chat) Write(src *User, msg string) {
	ch.mux.RLock()
	defer ch.mux.RUnlock()

	toSend := src.name + ": " + msg

	for _, dst := range ch.users {
		if dst != src {
			err := dst.Write(toSend)
			if err != nil {
				log.Println(err)
			}
		}
	}
}

func (ch *Chat) SystemCall(msg string) {
	ch.Write(&User{0, "System", nil}, msg)
}
