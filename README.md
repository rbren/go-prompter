# go-prompter

This is a set of utilities for prompting LLMs in Go.

This is a very early version of this library, and the API is likely to change.

## Supported Backends
* OpenAI's GPT
* Google Gemini
* Hugging Face (experimental)

## Features
* Consistent interface for different models
* Session management (only OpenAI currently)
* Craft prompts using Go's text templating engine
* Extract JSON, Markdown, and code from responses

## Usage

### Minimal Example
```go
package main

import (
    "github.com/rbren/go-prompter/pkg/prompt"
)

func main() {
    resp, err := engine.Prompt("Who was the 44th president of the US?")
    if err != nil {
        panic(err)
    }
    fmt.Println(resp) // "Barack Obama was the 44th president."
}
```

### Select a Model
You can use env vars to choose a backend and supply credentials.
```
export HUGGING_FACE_API_KEY="hf_..."
export HUGGING_FACE_URL="https://api-inference.huggingface.co/models/codellama/CodeLlama-70b-Instruct-hf"

export OPENAI_API_KEY="sk-..."
export OPENAI_MODEL="gpt-4-0125-preview"

export GEMINI_API_KEY="AI..."

export LLM_BACKEND="OPENAI"
```

The model will be automatically selected using the `LLM_BACKEND` env var when you run:
```go
engine := prompt.New()
```

You can also instantiate a specific model directly:
```go
package main

import (
    "github.com/rbren/go-prompter/pkg/prompt"
)

func main() {
    model := llm.NewHuggingFaceClient(apiKey, url)
    // or
    model := llm.NewOpenAIClient(apiKey, model)
    // or
    model := llm.NewGeminiClient(apiKey)
    engine := prompt.Engine{LLM: model}
}
```

### With Templates
You can use Go's text templating engine to create more powerful and dynamic prompts.

Note that you can include one template file inside of another with `{{ template "example" . }}`.
This is helpful for e.g. adding a boilerplate preamble to all your prompts.

##### prompts/polite.md
```markdown
# Task
Your task is to respond to the user's query below. Please do so
as politely as possible. You MUST always refer to the user as "Sir or Madam".

## User Query
{{ .user_query }}
```

##### main.go
```go
package main

import (
    "github.com/rbren/go-prompter/pkg/prompt"
)

//go:embed prompts/*.md
var templateFS embed.FS

func init() {
    prompt.SetFS(&templateFS)
}

func main() {
    engine := prompt.New()
    resp, err := engine.PromptWithTemplate("polite", map[string]any{
        user_query: "How tall is Barack Obama?",
    })
    if err != nil {
        panic(err)
    }
    fmt.Println(resp)
}
```

# Example Projects
* https://github.com/rbren/vizzy
