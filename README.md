# go-prompter

This is a set of utilities for prompting LLMs in Go.

This is a very early version of this library, and the API is likely to change.

## Supported Backends
* OpenAI's GPT
* Anthropic's Claude
* Google's Gemini
* Hugging Face models (experimental)

## Features
* Consistent interface for different models
* Session management (currently only OpenAI and Claude)
* Craft prompts using Go's text templating engine
* Extract JSON, Markdown, and code from responses
* Save prompts and responses to local files or S3 for debugging and analysis

## Usage

### Minimal Example
```go
package main

import (
    "github.com/rbren/go-prompter/pkg/chat"
)

func main() {
    session := chat.NewSession()
    resp, _ := session.Prompt("Who was the 44th president of the US?")
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
    "github.com/rbren/go-prompter/pkg/chat"
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
    "github.com/rbren/go-prompter/pkg/chat"
)

//go:embed prompts/*.md
var templateFS embed.FS

func main() {
    session := chat.NewSession()
    session.SetFS(&templateFS)
    resp, err := session.PromptWithTemplate("polite", map[string]any{
        user_query: "How tall is Barack Obama?",
    })
    if err != nil {
        panic(err)
    }
    fmt.Println(resp)
}
```

### Extract JSON, Markdown, and Code
You can extract JSON, Markdown, and code from the response.

```go
package main

import (
    "github.com/rbren/go-prompter/pkg/chat"
)

type Person struct {
  Height int `json:"height"`
  Age    int `json:"age"`
}

func main() {
    session := chat.NewSession()
    resp, _ := session.Prompt("Please tell me Obama's height in inches and age in years. Respond in JSON format.")
    p := Person{}
    _ = chat.ExtractJSONAndUnmarshal(resp, &p)

    resp, _ := session.Prompt("Write a bash script that prints Obama's height and age.")
    code := chat.ExtractCode(resp)

    resp, _ := session.Prompt("Write an essay in Markdown about Obama")
    title := chat.ExtractMarkdownTitle(resp)


    fmt.Println(resp) // "Barack Obama was the 44th president."
}
```

### Send Chat History as Context
You can optionally send the entire session history to the model as context.
Be sure to start a new session when you want to clear the context, and don't
share sessions across users.

```go
package main

import (
    "github.com/rbren/go-prompter/pkg/chat"
)

func main() {
    session := chat.NewSession()
    session.SaveHistory = true
    resp, err := session.Prompt("Who was the 44th president of the US?")
    resp, err = session.Prompt("How tall is he?")
}
```

### Save Debug Prompts and Responses
You can save prompts and responses to local files or S3 for debugging and analysis.

```go
package main

import (
    "github.com/rbren/go-prompter/pkg/chat"
    "github.com/rbren/go-prompter/pkg/files"
)

func main() {
    session := chat.NewSession()
    session.SessionID = "presidents" // this will be a random UUID otherwise
    session.SetDebugFileManager(files.LocalFileManager{
      BasePath: "./debug/",
    })
    resp, err := session.PromptWithID("44", "Who was the 44th president of the US?")
    // creates ./debug/presidents/44/prompt.md and ./debug/presidents/44/response.md
}
```


# Example Projects
* https://github.com/rbren/vizzy
