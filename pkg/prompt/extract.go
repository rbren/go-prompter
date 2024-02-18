package prompt

import (
	"errors"
	"strings"

	"github.com/sirupsen/logrus"
)

func extractDelimiters(body, startDelim, endDelim string) (string, error) {
	firstIndex := strings.Index(body, startDelim)
	lastIndex := strings.LastIndex(body, endDelim)
	if firstIndex == -1 || lastIndex == -1 {
		return "", errors.New("invalid response")
	}
	return body[firstIndex+len(startDelim) : lastIndex], nil
}

func extractTitle(body string) (string, error) {
	lines := strings.Split(body, "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "# ") {
			return strings.TrimPrefix(line, "# "), nil
		}
	}
	return "", nil
}

func extractJSONObject(body string) (string, error) {
	json, err := extractDelimiters(body, "{", "}")
	if err != nil {
		logrus.Errorf("invalid JSON response from LLM: %s", body)
		return "", errors.New("invalid JSON")
	}
	return "{" + json + "}", nil
}

func extractJavaScript(body string) (string, error) {
	blocks := strings.Split("\n"+body+"\n", "```")
	// blocks 0 is preamble
	// block 1 is first code
	// block 2 is another preamble
	// block 3 is second code
	// block 4 is postamble
	// basically we want the odd blocks
	longestBody := ""
	for i := 1; i < len(blocks); i += 2 {
		body := blocks[i]
		if strings.HasPrefix(body, "javascript") {
			body = body[len("javascript"):]
		}
		if strings.HasPrefix(body, "JavaScript") {
			body = body[len("JavaScript"):]
		}
		if strings.HasPrefix(body, "js") {
			body = body[len("js"):]
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
