package chat

import (
	"strings"
	"testing"
)

// TestExtractCodeSuccess tests the successful Extraction of Code code
func TestExtractCodeSuccess(t *testing.T) {
	tests := []struct {
		name     string
		body     string
		expected string
	}{
		{
			name:     "Extract with ```javascript",
			body:     "```javascript\nconsole.log('Hello, world!');\n```",
			expected: "console.log('Hello, world!');",
		},
		{
			name:     "Extract with ```js",
			body:     "```js\nconsole.log('Hello, world!');\n```",
			expected: "console.log('Hello, world!');",
		},
		{
			name:     "Extract with ```",
			body:     "```\nconsole.log('Hello, world!');\n```",
			expected: "console.log('Hello, world!');",
		},
		{
			name:     "Double block",
			body:     "we will use\n```javascript\nconsole.log()\n```\nhere's the full implementation:\n```\nconsole.log('Hello, world!');\n```",
			expected: "console.log('Hello, world!');",
		},
		{
			name:     "Double block 2",
			body:     "we will use\n```javascript\nconsole.log()\n```\n```\nconsole.log('Hello, world!');\n```",
			expected: "console.log('Hello, world!');",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			js, err := ExtractCode(tc.body)
			if err != nil {
				t.Errorf("Expected no error, got %v", err)
			}
			js = strings.TrimSpace(js)
			if js != tc.expected {
				t.Errorf("Expected %s, got %s", tc.expected, js)
			}
		})
	}
}

// TestExtractCodeFailure tests the function's handling of inputs that should result in errors
func TestExtractCodeFailure(t *testing.T) {
	body := "This is not a code block."
	_, err := ExtractCode(body)
	if err == nil {
		t.Error("Expected an error, but got nil")
	}
}
