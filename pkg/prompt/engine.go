package prompt

import (
	"embed"
	"encoding/json"
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

func (c *Engine) PromptWithTemplate(template string, data map[string]any) (string, error) {
	prompt, err := fillTemplate(template, data)
	if err != nil {
		return "", err
	}
	go writeDebugRequest(c.SessionID, template, prompt)
	resp, err := c.LLM.Query(template, prompt)
	go writeDebugResponse(c.SessionID, template, resp)
	return resp, err
}

func (c *Engine) PromptJSONWithTemplate(template string, data map[string]any, dest any) (error) {
	resp, err := c.PromptWithTemplate(template, data)
	if err != nil {
		return err
	}
	// TODO: handle arrays
	jsonString, err := ExtractJSONObject(resp)
	if err != nil {
		return err
	}
	return json.Unmarshal([]byte(jsonString), dest)
}
