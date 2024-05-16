package chat

import (
	"fmt"
	"math/rand"
	"net"
)

const CommandSignal = 0

type User struct {
	net.Conn
	id       int
	name     string
	currChat *Chat
}

func (s *Server) NewUser(name string, conn net.Conn) *User {
	u := &User{
		Conn:     conn,
		id:       rand.Intn(1000000),
		name:     name,
		currChat: nil,
	}

	s.usersMux.Lock()
	s.users[u.id] = u
	s.usersMux.Unlock()

	fmt.Printf("New user: %d (%v)\n", u.id, conn.RemoteAddr())

	go u.listenConn(s)

	return u
}

func (u *User) listenConn(s *Server) {
	buf := make([]byte, 1024)
	for {
		l, err := u.Read(buf)
		if err != nil {
			break
		}
		if buf[0] == CommandSignal {
			command, ok := commands[buf[1]]
			if !ok {
				u.WriteSystemCall("unknown command")
				continue
			}
			err := command(s, u, buf[2:l])
			if err != nil {
				u.WriteSystemCall(err.Error())
			}
			continue
		}
		if u.currChat == nil {
			u.WriteSystemCall("you are not connected to any chat")
			continue
		}
		u.currChat.Write(u, buf[:l])
	}
	s.deleteUser(u.id)
}

func (u *User) WriteSystemCall(s string) error {
	_, err := u.Write(append([]byte{0, 0}, []byte(s)...))
	return err
}
