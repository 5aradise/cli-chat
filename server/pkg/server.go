package chat

import (
	"fmt"
	"log"
	"net"
	"strconv"
	"sync"
	"unicode/utf8"
)

// TODO Update write with headers

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

func (s *server) authUser(conn net.Conn) (*user, error) {
	buf := make([]byte, bufferSize)
	var head header
	var username []byte
	for {
		l, err := conn.Read(buf)
		if err != nil {
			return nil, err
		}

		head, username = getHeader(buf[:l])
		if head != authAcc {
			_, err = conn.Write(authRej.setHeader([]byte("invalid request")))
			if err != nil {
				return nil, err
			}
			continue
		}

		if utf8.RuneCount(username) > maxUsernameLen {
			_, err = conn.Write(authRej.setHeader([]byte("username is too long (maximum 10 characters)")))
			if err != nil {
				return nil, err
			}
			continue
		}

		break
	}

	user := s.newUser(string(username), conn)
	_, err := conn.Write(authAcc.setHeader([]byte(strconv.Itoa(user.id))))
	if err != nil {
		return nil, err
	}

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
