package prompt

import (
	"encoding/json"
)

func (c *Engine) Prompt(prompt string) (string, error) {
	return c.PromptWithID("default", prompt)
}

func (c *Engine) PromptWithID(id, prompt string) (string, error) {
	writeDebugRequest(c.SessionID, id, prompt)
	resp, err := c.LLM.Query(prompt)
	writeDebugResponse(c.SessionID, id, resp)
	return resp, err
}

func (c *Engine) PromptWithTemplate(template string, data map[string]any) (string, error) {
	prompt, err := fillTemplate(template, data)
	if err != nil {
		return "", err
	}
	return c.PromptWithID(template, prompt)
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

