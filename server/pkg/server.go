package chat

import (
	"fmt"
	"log"
	"net"
	"sync"
)

const (
	bufferSize = 256
)

type server struct {
	listener net.Listener
	chats    map[string]*chat
	chatsMux sync.RWMutex
	users    map[string]*user
	usersMux sync.RWMutex
}

func New(port string) (*server, error) {
	host := ""
	if port == "8080" {
		host = "127.0.0.1"
	}

	l, err := net.Listen("tcp4", net.JoinHostPort(host, port))
	if err != nil {
		return nil, err
	}
	return &server{
		listener: l,
		chats:    make(map[string]*chat),
		chatsMux: sync.RWMutex{},
		users:    make(map[string]*user),
		usersMux: sync.RWMutex{},
	}, nil
}

func (s *server) Run() {
	defer s.listener.Close()

	log.Println("Start listening on", s.listener.Addr())
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			log.Println(err)
			continue
		}

		log.Println("Accept new connection:", conn.RemoteAddr())

		go func() {
			user, err := s.authUser(conn)
			if err != nil {
				log.Println(err)
				return
			}
			user.listenConn(s)
			err = s.deleteUser(string(user.name))
			if err != nil {
				log.Println(err)
			}
		}()
	}
}

func (s *server) deleteUser(name string) error {
	s.usersMux.Lock()
	defer s.usersMux.Unlock()

	user, ok := s.users[name]
	if !ok {
		return fmt.Errorf("cannot find user with name: %s", name)
	}

	user.isActive = false
	user.conn.Close()
	delete(s.users, name)

	if user.currChat != nil {
		err := user.currChat.deleteUser(name)
		if err != nil {
			return err
		}
	}

	log.Printf("Delete user: %s (%v)\n", user.name, user.conn.RemoteAddr())

	return nil
}

func (s *server) deleteChat(name string) error {
	s.chatsMux.Lock()
	defer s.chatsMux.Unlock()

	chat, ok := s.chats[name]
	if !ok {
		return fmt.Errorf("cannot find chat with name: %s", name)
	}

	close(chat.c)
	chat.deleteTimer.Stop()
	for memberName, member := range chat.users {
		if member != chat.admin {
			chat.deleteUser(memberName)
		}
	}
	if chat.admin != nil {
		chat.deleteUser(chat.admin.name)
	}
	delete(s.chats, name)

	log.Printf("Delete chat: %s\n", chat.name)

	return nil
}
