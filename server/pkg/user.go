package chat

import (
	"log"
	"math/rand"
	"net"
)

const CommandSignal = 0

type user struct {
	isActive bool
	conn     net.Conn
	readBuf  []byte
	id       int
	name     []byte
	currChat *chat
}

func (s *server) newUser(name []byte, conn net.Conn) *user {
	u := &user{
		isActive: true,
		conn:     conn,
		readBuf:  make([]byte, bufferSize),
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
	for u.isActive {
		head, body := u.read()
		command, ok := commands[head]
		if !ok {
			u.write(systemMsg, []byte("unknown command"))
			continue
		}

		err := command(s, u, body)
		if err != nil {
			u.write(systemMsg, []byte(err.Error()))
		}
	}
}

func (u *user) write(h header, b []byte) {
	_, err := u.conn.Write(h.setHeader(b))
	if err != nil {
		u.isActive = false
	}
}

func (u *user) read() (header, []byte) {
	l, err := u.conn.Read(u.readBuf)
	if err != nil {
		u.isActive = false
		return 0, nil
	}
	return getHeader(u.readBuf[:l])
}
