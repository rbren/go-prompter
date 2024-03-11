package llm

import (
	"os"
)

type ChatMessage struct {
	From    string
	Message string
}

type Client interface {
	Query(string, []ChatMessage) (string, error)
}

// NewClientFromEnv selects an LLM backend based on environment settings and initializes the appropriate client.
func NewClientFromEnv() Client {
	var llmClient Client
	if os.Getenv("LLM_BACKEND") == "HUGGING_FACE" {
		llmClient = NewHuggingFaceClientFromEnv()
	} else if os.Getenv("LLM_BACKEND") == "GEMINI" {
		llmClient = NewGeminiClientFromEnv()
	} else if os.Getenv("LLM_BACKEND") == "CLAUDE" {
		llmClient = NewClaudeClientFromEnv()
	} else {
		llmClient = NewOpenAIClientFromEnv()
	}
	return llmClient
}
