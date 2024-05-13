package chat

import (
	"fmt"
	"log"
	"net"
	"sync"
)

type Server struct {
	net.Listener
	chats    map[int]*Chat
	chatsMux sync.RWMutex
	users    map[int]*User
	usersMux sync.RWMutex
}

func New(port string) (*Server, error) {
	host := ""
	if port == "8080" {
		host = "127.0.0.1"
	}

	l, err := net.Listen("tcp", host+":"+port)
	if err != nil {
		return nil, err
	}
	return &Server{
		Listener: l,
		chats:    make(map[int]*Chat),
		chatsMux: sync.RWMutex{},
		users:    make(map[int]*User),
		usersMux: sync.RWMutex{},
	}, nil
}

func (s *Server) Run() {
	for {
		conn, err := s.Accept()
		if err != nil {
			log.Println(err)
			continue
		}

		go s.authUser(conn)
	}
}

func (s *Server) authUser(conn net.Conn) {
	_, err := conn.Write(append([]byte{0, 0}, []byte("Enter name")...))
	if err != nil {
		return
	}

	buf := make([]byte, 1024)
	l, err := conn.Read(buf)
	if err != nil {
		return
	}

	user := s.NewUser(string(buf[:l]), conn)

	conn.Write(append([]byte{0, 0}, []byte(fmt.Sprintf("User with id %d have been created", user.id))...))
}
