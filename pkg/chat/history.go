package chat

import (
	"github.com/rbren/go-prompter/pkg/llm"
)

// AddUserMessage adds a user message to the session history.
func (s *Session) AddUserMessage(msg string) {
	s.addMessage("user", msg)
}

// AddBotMessage adds a bot (assistant) message to the session history.
func (s *Session) AddBotMessage(msg string) {
	s.addMessage("assistant", msg)
}

func (s *Session) addMessage(from, msg string) {
	if s.MaxHistory == 0 {
		return
	}
	s.History = append(s.History, llm.ChatMessage{
		From:    "user",
		Message: msg,
	})
	for len(s.History) > s.MaxHistory {
		s.History = s.History[1:]
	}
}
