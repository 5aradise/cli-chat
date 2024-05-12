package chat

var commands map[byte]func(*Server, *User, []byte) error = map[byte]func(*Server, *User, []byte) error{
	// 0: (*Server).createChat,
	1: (*Server).connToChat,
	// 2: (*Server).exitChat,
}
