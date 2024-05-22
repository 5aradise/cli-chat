package chat

import (
	"fmt"
	"log"
	"net"
	"sync"
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

	l, err := net.Listen("tcp", host+":"+port)
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
				return
			}
			user.listenConn(s)
			s.deleteUser(user.id)
		}()
	}
}

func (s *server) authUser(conn net.Conn) (*user, error) {
	_, err := conn.Write(systemMsg.setHeaderS("Enter name"))
	if err != nil {
		return nil, err
	}

	buf := make([]byte, 1024)
	l, err := conn.Read(buf)
	if err != nil {
		return nil, err
	}

	user := s.newUser(string(buf[:l]), conn)

	createMsg := fmt.Sprintf("User with id %d have been created", user.id)
	user.conn.Write(systemMsg.setHeaderS(createMsg))
	return user, nil
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
