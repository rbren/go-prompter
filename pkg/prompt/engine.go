package prompt

import (
	"embed"
	"text/template"


	"github.com/google/uuid"

	"github.com/rbren/go-prompter/pkg/llm"
)

var templateFS *embed.FS
var templateFuncMap template.FuncMap

func SetFS(f *embed.FS) {
	templateFS = f
}

func SetTemplateFuncMap(funcMap template.FuncMap) {
	templateFuncMap = funcMap
}

type Engine struct {
	LLM llm.Client
	History []llm.ChatMessage
	SessionID string
	SaveHistory bool
}

func New() *Engine {
	return &Engine{
		LLM: llm.New(),
	}
}

func (c *Engine) WithSession(sessionID string) *Engine {
	if sessionID == "" {
		sessionID = uuid.New().String()
	}
	return &Engine{
		SessionID: sessionID,
		LLM: c.LLM,
	}
}

