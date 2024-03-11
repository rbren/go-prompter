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

// OpenAIRequest represents the request body for OpenAI API.
type OpenAIRequest struct {
	Model    string          `json:"model"`
	Messages []OpenAIMessage `json:"messages"`
	Seed     int             `json:"seed,omitempty"`
}

// OpenAIMessage represents a message in the request body for OpenAI API.
type OpenAIMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// OpenAIResponse represents the response from OpenAI API.
type OpenAIResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
	Error map[string]interface{} `json:"error"`
}

// OpenAIClient holds the information needed to make requests to the OpenAI API.
type OpenAIClient struct {
	APIKey string
	Model  string
	Seed   int
}

// NewOpenAIClient creates a new OpenAI API client.
func NewOpenAIClient(apiKey, model string) *OpenAIClient {
	return &OpenAIClient{
		APIKey: apiKey,
		Model:  model,
	}
}

func NewOpenAIClientFromEnv() *OpenAIClient {
	client := NewOpenAIClient(os.Getenv("OPENAI_API_KEY"), os.Getenv("OPENAI_MODEL"))
	if seed, err := strconv.Atoi(os.Getenv("OPENAI_SEED")); err == nil {
		client.Seed = seed
	}
	return client
}

// Query sends a prompt to the OpenAI API and returns the response.
func (c *OpenAIClient) Query(prompt string, history []ChatMessage) (string, error) {
	systemPrompt := "The following is a conversation with an AI assistant."

	messages := []OpenAIMessage{{
		Role:    "system",
		Content: systemPrompt,
	}}
	for _, m := range history {
		messages = append(messages, OpenAIMessage{
			Role:    m.From,
			Content: m.Message,
		})
	}
	messages = append(messages, OpenAIMessage{
		Role:    "user",
		Content: prompt,
	})

	requestBody, err := json.Marshal(OpenAIRequest{
		Seed:     c.Seed,
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

	var response OpenAIResponse
	logrus.Debugf("response: %s", string(body))
	if err := json.Unmarshal(body, &response); err != nil {
		return "", err
	}

	if response.Error != nil {
		logrus.Errorf("error from OpenAI: %v", response.Error)
		message, ok := response.Error["message"].(string)
		if ok {
			return "", errors.New(message)
		} else {
			return "", errors.New("unknown error from OpenAI")
		}
	}

	if len(response.Choices) == 0 {
		return "", fmt.Errorf("no response from OpenAI")
	}
	out := response.Choices[0].Message.Content
	return out, nil
}
