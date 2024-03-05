package chat

import (
	"github.com/rbren/go-prompter/pkg/llm"
)

// AddUserMessage adds a user message to the session history.
func (s *Session) AddUserMessage(msg string) {
	if !s.SaveHistory {
		return
	}
	s.History = append(s.History, llm.ChatMessage{
		From:    "user",
		Message: msg,
	})
}

// AddBotMessage adds a bot (assistant) message to the session history.
func (s *Session) AddBotMessage(msg string) {
	if !s.SaveHistory {
		return
	}
	s.History = append(s.History, llm.ChatMessage{
		From:    "assistant",
		Message: msg,
	})
}

