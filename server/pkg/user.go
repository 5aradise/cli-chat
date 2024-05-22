package chat

import (
	"log"
	"math/rand"
	"net"
)

const CommandSignal = 0

type user struct {
	conn     net.Conn
	id       int
	name     string
	currChat *chat
}

func (s *server) newUser(name string, conn net.Conn) *user {
	u := &user{
		conn:     conn,
		id:       rand.Intn(1000000),
		name:     name,
		currChat: nil,
	}

	s.usersMux.Lock()
	s.users[u.id] = u
	s.usersMux.Unlock()

	log.Printf("New user: %d (%v)\n", u.id, conn.RemoteAddr())

	return u
}

func (u *user) listenConn(s *server) {
	buf := make([]byte, 1024)
	for {
		l, err := u.conn.Read(buf)
		if err != nil {
			break
		}

		command, ok := commands[header(buf[0])]
		if !ok {
			u.writeSystemCall("unknown command")
			continue
		}

		err = command(s, u, buf[1:l])
		if err != nil {
			u.writeSystemCall(err.Error())
		}
	}
}

func (u *user) writeSystemCall(s string) error {
	_, err := u.conn.Write(systemMsg.setHeaderS(s))
	return err
}
