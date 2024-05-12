package client

func (c *Client) connToChat(args []string) error {
	const CommandCode byte = 1

	req := []byte{commandCode, CommandCode}

	req = append(req, []byte(args[0])...)

	_, err := c.Write(req)
	return err
}
