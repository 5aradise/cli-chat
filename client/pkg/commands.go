package client

const commandCode byte = 0

type command struct {
	desc string
	fn   func(*Client, []string) error
}

var commands map[string]command = map[string]command{
	"chat": {
		desc: "Lol",
		fn:   (*Client).connToChat,
	},
}
