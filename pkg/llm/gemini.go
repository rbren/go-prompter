package llm

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/sirupsen/logrus"
)

const URL = "https://generativelanguage.googleapis.com/v1beta/models/gemini-pro:generateContent"

type GeminiRequest struct {
	Contents []struct {
		Parts []struct {
			Text string `json:"text"`
		} `json:"parts"`
	} `json:"contents"`
}

type GeminiResponse struct {
	Candidates []struct {
		Content struct {
			Parts []struct {
				Text string `json:"text"`
			} `json:"parts"`
		} `json:"content"`
	}
}

type Gemini struct {
	APIKey string
}

// NewGeminiClient initializes a new Gemini client with the provided API key.
func NewGeminiClient(apiKey string) *Gemini {
	return &Gemini{
		APIKey: apiKey,
	}
}

// NewGeminiClientFromEnv initializes a new Gemini client with the API key from the environment.
func NewGeminiClientFromEnv() *Gemini {
	return NewGeminiClient(os.Getenv("GEMINI_API_KEY"))
}

// Query sends a prompt and chat history to the Gemini API and returns the API's text response.
func (g *Gemini) Query(prompt string, history []ChatMessage) (string, error) {
	request := GeminiRequest{
		Contents: []struct {
			Parts []struct {
				Text string `json:"text"`
			} `json:"parts"`
		}{
			{Parts: []struct {
				Text string `json:"text"`
			}{{prompt}}},
		},
	}

	reqBody, err := json.Marshal(request)
	if err != nil {
		return "", err
	}

	url := URL + "?key=" + g.APIKey
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(reqBody))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")

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
	if resp.StatusCode != 200 {
		logrus.Errorf("Gemini Response: %s", resp.Status)
		return "", errors.New("non-200 status code")
	}

	var response GeminiResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return "", err
	}
	if len(response.Candidates) == 0 {
		return "", errors.New("no candidates")
	}
	if len(response.Candidates[0].Content.Parts) == 0 {
		return "", errors.New("no parts")
	}
	text := response.Candidates[0].Content.Parts[0].Text
	return text, nil
}
