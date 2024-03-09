package llm

import (
	"os"

	"github.com/sirupsen/logrus"
)

type ChatMessage struct {
	From    string
	Message string
}

type Client interface {
	Query(string, []ChatMessage) (string, error)
}

// New selects an LLM backend based on environment settings and initializes the appropriate client.
func New() Client {
	var llmClient Client
	if os.Getenv("LLM_BACKEND") == "HUGGING_FACE" {
		logrus.Infof("Using Hugging Face API")
		llmClient = NewHuggingFaceClient(os.Getenv("HUGGING_FACE_API_KEY"), os.Getenv("HUGGING_FACE_URL"))
	} else if os.Getenv("LLM_BACKEND") == "GEMINI" {
		gClient := NewGeminiClient(os.Getenv("GEMINI_API_KEY"))
		llmClient = gClient
	} else {
		oaiClient := NewOpenAIClient(os.Getenv("OPENAI_API_KEY"), os.Getenv("OPENAI_MODEL"))
		oaiClient.Seed = 42
		llmClient = oaiClient
	}
	return llmClient
}

