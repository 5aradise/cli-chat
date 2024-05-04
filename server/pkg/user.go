package chat

import (
	"math/rand"
	"net"
)

type User struct {
	id   int
	name string
	conn net.Conn
}

func NewUser(name string, conn net.Conn) *User {
	return &User{
		id:   rand.Intn(1000000),
		name: name,
		conn: conn,
	}
}

func (u *User) Write(text string) error {
	_, err := u.conn.Write([]byte(text))
	return err
}

func (u *User) Read() (string, error) {
	buf := make([]byte, 1024)

	l, err := u.conn.Read(buf)
	if err != nil {
		return "", err
	}

	return string(buf[:l]), nil
}
