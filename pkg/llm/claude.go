package llm

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"

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
	APIKey    string
	Model     string
	Version   string
	MaxTokens int
}

// NewClaudeClient creates a new Claude API client.
func NewClaudeClient(apiKey, model, version string) *ClaudeClient {
	return &ClaudeClient{
		APIKey:  apiKey,
		Model:   model,
		Version: version,
	}
}

func NewClaudeClientFromEnv() *ClaudeClient {
	client := NewClaudeClient(
		os.Getenv("CLAUDE_API_KEY"),
		os.Getenv("CLAUDE_MODEL"),
		os.Getenv("CLAUDE_VERSION"))
	client.MaxTokens = 1024
	if maxTokens, err := strconv.Atoi(os.Getenv("OPENAI_MAX_TOKENS")); err == nil {
		client.MaxTokens = maxTokens
	}

	return client
}

// Query sends a prompt to the Claude API and returns the response.
func (c *ClaudeClient) Query(prompt string, history []ChatMessage) (string, error) {
	messages := []ClaudeMessage{}
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
		Messages:  messages,
		Model:     c.Model,
		MaxTokens: c.MaxTokens,
	})
	if err != nil {
		return "", err
	}
	logrus.Debugf("request: %s", prompt)

	req, err := http.NewRequest("POST", "https://api.anthropic.com/v1/messages", bytes.NewBuffer(requestBody))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", c.APIKey)
	req.Header.Set("anthropic-version", c.Version)

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
