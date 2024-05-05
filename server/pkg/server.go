package chat

import (
	"fmt"
	"log"
	"net"
	"strconv"
)

type Server struct {
	listener net.Listener
	chats    map[int]*Chat
	users    map[int]*User
}

func NewServer(port string) (*Server, error) {
	ip := ":"
	if port == "8080" {
		ip = "127.0.0.1:"
	}

	l, err := net.Listen("tcp", ip+port)
	if err != nil {
		return nil, err
	}
	return &Server{
		l,
		make(map[int]*Chat),
		make(map[int]*User),
	}, nil
}

func (s *Server) Run() {
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			log.Println(err)
			continue
		}

		go s.acceptConn(conn)
	}
}

func (s *Server) acceptConn(conn net.Conn) {
	user, err := s.authUser(conn)
	if err != nil {
		fmt.Println(err)
		return
	}

	err = s.connetToChat(user)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func (s *Server) authUser(conn net.Conn) (*User, error) {
	_, err := conn.Write([]byte("Enter name"))
	if err != nil {
		return nil, err
	}

	buf := make([]byte, 1024)
	l, err := conn.Read(buf)
	if err != nil {
		return nil, err
	}

	user := NewUser(string(buf[:l]), conn)
	s.users[user.id] = user

	_, err = conn.Write([]byte(fmt.Sprintf("User with id %d have been created", user.id)))
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *Server) connetToChat(user *User) error {
	err := user.Write("Enter chat id")
	if err != nil {
		return err
	}

	input, err := user.Read()
	if err != nil {
		return err
	}

	chatId, err := strconv.Atoi(input)
	for err == strconv.ErrSyntax {
		err = user.Write("Please write a number")
		if err != nil {
			return err
		}

		input, err = user.Read()
		if err != nil {
			return err
		}
		chatId, err = strconv.Atoi(input)
	}
	if err != nil {
		return err
	}

	chat, ok := s.chats[chatId]
	if !ok {
		chat = NewChat(chatId)
		s.chats[chatId] = chat
	}

	return chat.AddUser(user)
}
