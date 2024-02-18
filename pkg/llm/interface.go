package llm

import (
	"os"

	"github.com/sirupsen/logrus"

	"github.com/rbren/vizzy/pkg/files"
)

type Client interface {
	Query(string, string) (string, error)
	Copy() Client
	SetDebugFileManager(files.FileManager)
}

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
