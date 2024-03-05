package chat

import (
	"errors"
	"strings"

	"github.com/sirupsen/logrus"
)

// ExtractDelimiters extracts content between specified delimiters.
func ExtractDelimiters(body, startDelim, endDelim string) (string, error) {
	firstIndex := strings.Index(body, startDelim)
	lastIndex := strings.LastIndex(body, endDelim)
	if firstIndex == -1 || lastIndex == -1 {
		return "", errors.New("invalid response")
	}
	return body[firstIndex+len(startDelim) : lastIndex], nil
}

// ExtractTitle finds and returns the title from the provided text.
func ExtractTitle(body string) (string, error) {
	lines := strings.Split(body, "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "# ") {
			return strings.TrimPrefix(line, "# "), nil
		}
	}
	return "", nil
}

// ExtractJSONObject attempts to extract a JSON object from the given text.
func ExtractJSONObject(body string) (string, error) {
	json, err := ExtractDelimiters(body, "{", "}")
	if err != nil {
		logrus.Errorf("invalid JSON response from LLM: %s", body)
		return "", errors.New("invalid JSON")
	}
	return "{" + json + "}", nil
}

// ExtractCode attempts to find and return the longest code block in the given text.
func ExtractCode(body string) (string, error) {
	blocks := strings.Split("\n"+body+"\n", "```")
	// blocks 0 is preamble
	// block 1 is first code
	// block 2 is more text
	// block 3 is second code
	// block 4 is more text
	// basically, we want the odd blocks
	longestBody := ""
	for i := 1; i < len(blocks); i += 2 {
		body := blocks[i]
		newline := strings.Index(body, "\n")
		if newline != -1 {
			// strip the first line, which is either empty or the language name, e.g. ```javascript
			body = body[newline:]
		}
		if len(body) > len(longestBody) {
			longestBody = body
		}
	}
	if len(longestBody) == 0 {
		return "", errors.New("no JavaScript code blocks found")
	}
	return longestBody, nil
}

