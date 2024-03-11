package llm

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/sirupsen/logrus"
)

// ClaudeRequest represents the request body for Claude API.
type ClaudeRequest struct {
	Model     string          `json:"model"`
	Messages  []ClaudeMessage `json:"messages"`
	MaxTokens int             `json:"max_tokens,omitempty"`
}

// ClaudeMessage represents a message in the request body for Claude API.
type ClaudeMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// ClaudeResponse represents the response from Claude API.
type ClaudeResponse struct {
	Content []struct {
		Text string `json:"text"`
		Type string `json:"type"`
	} `json:"content"`
	Error map[string]interface{} `json:"error"`
}

// ClaudeClient holds the information needed to make requests to the Claude API.
type ClaudeClient struct {
	APIKey string
	Model  string
}

// NewClaudeClient creates a new Claude API client.
func NewClaudeClient(apiKey, model string) *ClaudeClient {
	return &ClaudeClient{
		APIKey: apiKey,
		Model:  model,
	}
}

func NewClaudeClientFromEnv() *ClaudeClient {
	return NewClaudeClient(os.Getenv("CLAUDE_API_KEY"), os.Getenv("CLAUDE_MODEL"))
}

// Query sends a prompt to the Claude API and returns the response.
func (c *ClaudeClient) Query(prompt string, history []ChatMessage) (string, error) {
	systemPrompt := "The following is a conversation with an AI assistant."

	messages := []ClaudeMessage{{
		Role:    "system",
		Content: systemPrompt,
	}}
	for _, m := range history {
		messages = append(messages, ClaudeMessage{
			Role:    m.From,
			Content: m.Message,
		})
	}
	messages = append(messages, ClaudeMessage{
		Role:    "user",
		Content: prompt,
	})

	requestBody, err := json.Marshal(ClaudeRequest{
		Messages: messages,
		Model:    c.Model,
	})
	if err != nil {
		return "", err
	}
	logrus.Debugf("request: %s", prompt)

	req, err := http.NewRequest("POST", "https://api.openai.com/v1/chat/completions", bytes.NewBuffer(requestBody))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.APIKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var response ClaudeResponse
	logrus.Debugf("response: %s", string(body))
	if err := json.Unmarshal(body, &response); err != nil {
		return "", err
	}

	if response.Error != nil {
		logrus.Errorf("error from Claude: %v", response.Error)
		message, ok := response.Error["message"].(string)
		if ok {
			return "", errors.New(message)
		} else {
			return "", errors.New("unknown error from Claude")
		}
	}

	if len(response.Content) == 0 {
		return "", fmt.Errorf("no response from Claude")
	}
	out := response.Content[0].Text
	return out, nil
}
