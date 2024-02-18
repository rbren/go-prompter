# go-prompter

This is a set of utilities for prompting LLMs.

This is a very early version of this library, and the API is likely to change.

## Supported Backends
* OpenAI's GPT
* Google Gemini
* Hugging Face (experimental)

## Usage
### Minimal Example
```go
package main

import (
	"github.com/rbren/go-prompter/pkg/prompt"
)

func main() {
	engine := prompt.New()
    resp, err := engine.Prompt("Who was the 44th president of the US?")
    if err != nil {
        panic(err)
    }
    fmt.Println(resp) # "Barack Obama was the 44th president."
}
```

### With Templates
You can use Go's text templating engine to create more powerful and dynamic prompts.
Using markdown to generate queries is recommended.

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
