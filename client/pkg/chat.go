package client

func (c *Client) connToChat(args []string) error {
	const connectionCode byte = 1

	req := []byte{commandsCode, connectionCode}

	req = append(req, []byte(args[0])...)

	_, err := c.Write(req)
	return err
}
