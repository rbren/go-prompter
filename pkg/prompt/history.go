package prompt

import (
	"github.com/rbren/go-prompter/pkg/llm"
)

func (c *Engine) AddUserMessage(msg string) {
	if (!c.SaveHistory) {
		return
	}
	c.History = append(c.History, llm.ChatMessage{
		From: "user",
		Message: msg,
	})
}

func (c *Engine) AddBotMessage(msg string) {
	if (!c.SaveHistory) {
		return
	}
	c.History = append(c.History, llm.ChatMessage{
		From: "assistant",
		Message: msg,
	})
}

