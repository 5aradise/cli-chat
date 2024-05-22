package client

type header byte

const (
	systemMsg header = iota
	chatMsg
	userMsg
	create
	connect
	exit
)

func (h header) setHeaderB(body []byte) []byte {
	return append([]byte{byte(h)}, body...)
}

func (h header) setHeaderS(body string) []byte {
	return append([]byte{byte(h)}, body...)
}
