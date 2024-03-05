package chat

import (
	"github.com/rbren/go-prompter/pkg/llm"
)

func (s *Session) AddUserMessage(msg string) {
	if !s.SaveHistory {
		return
	}
	s.History = append(s.History, llm.ChatMessage{
		From:    "user",
		Message: msg,
	})
}

func (s *Session) AddBotMessage(msg string) {
	if !s.SaveHistory {
		return
	}
	s.History = append(s.History, llm.ChatMessage{
		From:    "assistant",
		Message: msg,
	})
}
