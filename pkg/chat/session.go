package chat

import (
	"embed"
	"text/template"

	"github.com/google/uuid"

	"github.com/rbren/go-prompter/pkg/files"
	"github.com/rbren/go-prompter/pkg/llm"
)

// Session structures the session of a chat interaction, holding information necessary for communication and history.
type Session struct {
	LLM              llm.Client
	History          []llm.ChatMessage
	SessionID        string
	SaveHistory      bool
	templateFS       *embed.FS
	templateFuncMap  template.FuncMap
	debugFileManager files.FileManager
}

// NewSession creates and initializes a new chat session.
func NewSession() *Session {
	return &Session{
		LLM:       llm.New(),
		SessionID: uuid.New().String(),
	}
}

// SetFS sets the filesystem for templates.
func (s *Session) SetFS(f *embed.FS) {
	s.templateFS = f
}

// SetTemplateFuncMap sets the template function map for templating.
func (s *Session) SetTemplateFuncMap(funcMap template.FuncMap) {
	s.templateFuncMap = funcMap
}

