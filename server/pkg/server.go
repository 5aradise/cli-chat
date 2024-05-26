package chat

import (
	"fmt"
	"log"
	"net"
	"sync"
)

const (
	bufferSize     = 256
	maxUsernameLen = 10
	maxMsgLen      = 106
)

type server struct {
	net.Listener
	chats    map[int]*chat
	chatsMux sync.RWMutex
	users    map[int]*user
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
		Listener: l,
		chats:    make(map[int]*chat),
		chatsMux: sync.RWMutex{},
		users:    make(map[int]*user),
		usersMux: sync.RWMutex{},
	}, nil
}

func (s *server) Run() {
	defer s.Close()

	log.Println("Start listening on", s.Addr())
	for {
		conn, err := s.Accept()
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
			s.deleteUser(user.id)
		}()
	}
}

func (s *server) deleteUser(id int) error {
	s.usersMux.Lock()
	defer s.usersMux.Unlock()

	user, ok := s.users[id]
	if !ok {
		return fmt.Errorf("cannot find user with id: %d", id)
	}

	if user.currChat != nil {
		user.currChat.deleteUser(id)
	}
	delete(s.users, id)
	user.conn.Close()

	log.Printf("Delete user: %d (%v)\n", id, user.conn.RemoteAddr())
	return nil
}
