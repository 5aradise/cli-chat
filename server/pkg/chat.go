package chat

import (
	"errors"
	"log"
	"sync"
	"time"
)

const chatDeleteDelay = time.Minute

type chat struct {
	name        string
	c           chan *message
	users       map[string]*user
	mux         sync.RWMutex
	admin       *user
	deleteTimer *time.Timer
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
		name:        name,
		c:           make(chan *message, 16),
		users:       make(map[string]*user),
		mux:         sync.RWMutex{},
		deleteTimer: time.NewTimer(chatDeleteDelay),
	}

	s.chats[name] = chat
	s.chatsMux.Unlock()

	go chat.broadcast()
	go func() {
		<-chat.deleteTimer.C
		s.deleteChat(chat.name)
	}()
	return chat, nil
}

func (ch *chat) addUser(u *user) error {
	ch.mux.Lock()
	if _, ok := ch.users[u.name]; ok {
		ch.mux.Unlock()
		return errors.New("user with this name already in chat")
	}
	ch.users[u.name] = u
	u.currChat = ch
	ch.mux.Unlock()

	u.write(connectChat, []byte(ch.name))
	ch.chatCall(append([]byte(u.name), addMsg...))

	if ch.admin == nil {
		err := ch.setAdmin(u.name)
		if err != nil {
			return err
		}
	}

	ch.deleteTimer.Stop()

	return nil
}

func (ch *chat) deleteUser(name string) error {
	ch.mux.Lock()
	u, ok := ch.users[name]
	if !ok {
		ch.mux.Unlock()
		return errors.New("there is no user by that name in the chat room")
	}
	delete(ch.users, name)
	if len(ch.users) == 0 {
		ch.deleteTimer.Reset(chatDeleteDelay)
	}
	ch.mux.Unlock()

	u.currChat = nil
	u.write(exitChat, nil)
	ch.chatCall(append([]byte(u.name), deleteMsg...))

	if u == ch.admin {
		return ch.setAdmin()
	}
	return nil
}

func (ch *chat) chatCall(msg []byte) {
	ch.mux.RLock()
	defer ch.mux.RUnlock()

	for _, dst := range ch.users {
		dst.write(chatMsg, msg)
	}
}

func (ch *chat) setAdmin(name ...string) error {
	ch.mux.RLock()
	defer ch.mux.RUnlock()

	oldAdmin := ch.admin
	if len(name) == 0 {
		membersCount := len(ch.users)
		if membersCount == 0 {
			ch.admin = nil
		}
		if oldAdmin != nil {
			if _, ok := ch.users[oldAdmin.name]; ok && membersCount == 1 {
				return errors.New("you already admin")
			}
		}
		for _, randMember := range ch.users {
			if randMember == oldAdmin {
				continue
			}
			ch.admin = randMember
			break
		}
	} else {
		newAdminName := name[0]
		if oldAdmin != nil && newAdminName == oldAdmin.name {
			return errors.New("you already admin")
		}
		newAdmin, ok := ch.users[newAdminName]
		if !ok {
			return errors.New("there is no user by that name in the chat room")
		}
		ch.admin = newAdmin
	}

	if oldAdmin != nil {
		oldAdmin.write(passAdmin, []byte{0})
	}
	if ch.admin != nil {
		ch.admin.write(passAdmin, []byte{1})
	}
	return nil
}

func (ch *chat) broadcast() {
	log.Printf("New chat: %s\n", ch.name)
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
	log.Printf("Delete chat: %s\n", ch.name)
}
