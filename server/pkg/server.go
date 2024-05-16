package chat

import (
	"fmt"
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
	defer s.Close()

	fmt.Println("Start listening on", s.Addr())
	for {
		conn, err := s.Accept()
		if err != nil {
			fmt.Println(err)
			continue
		}

		fmt.Println("Accept new connection:", conn.RemoteAddr())

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

	successMsg := fmt.Sprintf("User with id %d have been created", user.id)
	user.Write(append([]byte{0, 0}, []byte(successMsg)...))
}

func (s *Server) deleteUser(id int) error {
	s.usersMux.Lock()
	defer s.usersMux.Unlock()

	user, ok := s.users[id]
	if !ok {
		return fmt.Errorf("cannot find user with id: %d", id)
	}

	if user.currChat != nil {
		user.currChat.DeleteUser(id)
	}
	delete(s.users, id)
	user.Close()

	fmt.Printf("Delete user: %d (%v)\n", id, user.RemoteAddr())
	return nil
}
