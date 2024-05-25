package client

type header byte

const (
	systemMsg header = iota + 1
	chatMsg
	userMsg
	create
	connect
	exit
	authAcc
)

func (h header) setHeader(body []byte) []byte {
	return append(body, byte(h))
}

func getHeader(element []byte) (header, []byte) {
	return header(element[len(element)-1]), element[:len(element)-1]
}
