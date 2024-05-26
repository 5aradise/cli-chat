package chat

import (
	"net"
)

func (s *server) authUser(conn net.Conn) (*user, error) {
	buf := make([]byte, bufferSize)
	var head header
	var username []byte
	var user *user
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

		user, err = s.newUser(string(username), conn)
		if err != nil {
			_, err = conn.Write(authRej.setHeader([]byte(err.Error())))
			if err != nil {
				return nil, err
			}
			continue
		}

		break
	}

	_, err := conn.Write(authAcc.setHeader([]byte(user.name)))
	if err != nil {
		return nil, err
	}

	return user, nil
}
