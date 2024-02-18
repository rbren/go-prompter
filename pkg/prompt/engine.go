package prompt

import (
	"embed"

	"github.com/google/uuid"

	"github.com/rbren/go-prompter/pkg/llm"
)

var templateFS *embed.FS

func SetFS(f *embed.FS) {
	templateFS = f
}

type Engine struct {
	LLM llm.Client
	SessionID string
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

func (c *Engine) QueryWithTemplate(template string, data map[string]interface{}) (string, error) {
	prompt, err := fillTemplate(template, data)
	if err != nil {
		return "", err
	}
	go writeDebugRequest(c.SessionID, template, prompt)
	resp, err := c.LLM.Query(template, prompt)
	go writeDebugResponse(c.SessionID, template, resp)
	return resp, err
}
