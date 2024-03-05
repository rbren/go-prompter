package chat

import (
	"embed"
	"text/template"

	"github.com/google/uuid"

	"github.com/rbren/go-prompter/pkg/files"
	"github.com/rbren/go-prompter/pkg/llm"
)

type Session struct {
	LLM              llm.Client
	History          []llm.ChatMessage
	SessionID        string
	SaveHistory      bool
	templateFS       *embed.FS
	templateFuncMap  template.FuncMap
	debugFileManager files.FileManager
}

func NewSession() *Session {
	return &Session{
		LLM:       llm.New(),
		SessionID: uuid.New().String(),
	}
}

func (s *Session) SetFS(f *embed.FS) {
	s.templateFS = f
}

func (s *Session) SetTemplateFuncMap(funcMap template.FuncMap) {
	s.templateFuncMap = funcMap
}
