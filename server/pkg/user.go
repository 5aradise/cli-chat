package chat

import (
	"errors"
	"log"
	"net"
)

const CommandSignal = 0

type user struct {
	isActive bool
	conn     net.Conn
	readBuf  []byte
	name     string
	currChat *chat
}

func (s *server) newUser(name string, conn net.Conn) (*user, error) {
	isValid, reas := isValidUsername(name)
	if !isValid {
		return nil, errors.New(reas)
	}

	s.usersMux.Lock()
	if _, ok := s.users[name]; ok {
		s.usersMux.Unlock()
		return nil, errors.New("user with this name already exist")
	}

	u := &user{
		isActive: true,
		conn:     conn,
		readBuf:  make([]byte, bufferSize),
		name:     name,
		currChat: nil,
	}

	s.users[u.name] = u
	s.usersMux.Unlock()

	log.Printf("New user: %s (%v)\n", u.name, conn.RemoteAddr())

	return u, nil
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
