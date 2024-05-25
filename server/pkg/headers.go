package chat

type header byte

const (
	systemMsg header = iota
	chatMsg
	userMsg
	create
	connect
	exit
	authAcc
	authRej
)

func (h header) setHeader(body []byte) []byte {
	return append(body, byte(h))
}

func getHeader(element []byte) (header, []byte) {
	return header(element[len(element)-1]), element[:len(element)-1]
}
