package prompt

import (
	"errors"
	"strings"
	"encoding/json"

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

func extractJSONObject(body string, dest any) (error) {
	str, err := extractDelimiters(body, "{", "}")
	if err != nil {
		logrus.Errorf("invalid JSON response from LLM: %s", body)
		return errors.New("invalid JSON")
	}
	str = "{" + str + "}"
	return json.Unmarshal([]byte(str), dest)
}

