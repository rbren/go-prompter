package chat

import (
	"encoding/json"
)

// Prompt sends a prompt to the LLM and returns its response.
func (s *Session) Prompt(prompt string) (string, error) {
	return s.PromptWithID("default", prompt)
}

// PromptWithID sends a prompt to the LLM with a custom ID and returns its response.
func (s *Session) PromptWithID(id, prompt string) (string, error) {
	s.writeDebugRequest(id, prompt)
	s.AddUserMessage(prompt)
	resp, err := s.LLM.Query(prompt, s.History)
	s.AddBotMessage(resp)
	s.writeDebugResponse(id, resp)
	return resp, err
}

// PromptWithTemplate sends a templated prompt to the LLM and returns its response.
func (s *Session) PromptWithTemplate(template string, data map[string]any) (string, error) {
	prompt, err := s.fillTemplate(template, data)
	if err != nil {
		return "", err
	}
	return s.PromptWithID(template, prompt)
}

// PromptWithTemplateToJSON sends a templated prompt to the LLM, receives a response, and attempts to unmarshal it into a provided structure.
func (s *Session) PromptWithTemplateToJSON(template string, data map[string]any, dest any) error {
	resp, err := s.PromptWithTemplate(template, data)
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
