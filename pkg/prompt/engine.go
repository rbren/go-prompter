package prompt

import (
	"embed"
	"encoding/json"

	"github.com/rbren/go-prompter/pkg/llm"
)

var templateFS *embed.FS

func SetFS(f *embed.FS) {
	templateFS = f
}

type Engine struct {
	LLM llm.Client
}

func New() *Engine {
	return &Engine{
		LLM: llm.New(),
	}
}

func (c *Engine) QueryWithTemplate(template string, data map[string]interface{}) (string, error) {
	prompt, err := fillTemplate(template, data)
	if err != nil {
		return "", err
	}
	return c.LLM.Query(template, prompt)
}

