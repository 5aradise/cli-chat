package chat

import "strconv"

func (s *Server) connToChat(user *User, args []byte) error {
	chatId, err := strconv.Atoi(string(args))
	if err != nil {
		return err
	}

	s.chatsMux.RLock()
	chat, ok := s.chats[chatId]
	s.chatsMux.RUnlock()
	if !ok {
		chat = s.NewChat(chatId)
	}

	return chat.AddUser(user)
}
