package main

import (
	"fmt"
	"io/ioutil"

	"github.com/rbren/go-prompter/pkg/llm"
)

const baseDir = "./test/regression/"

func main() {
	// iterate over directories in ./test/regression/
	dirs, err := ioutil.ReadDir(baseDir)
	if err != nil {
		panic(err)
	}
	for _, dir := range dirs {
		if !dir.IsDir() {
			continue
		}
		runRegressionTest(dir.Name())
	}
}

func runRegressionTest(dir string) {
	runRegressionTestWithClient(dir, "openai", llm.NewOpenAIClientFromEnv())
	runRegressionTestWithClient(dir, "claude", llm.NewClaudeClientFromEnv())
	runRegressionTestWithClient(dir, "huggingface", llm.NewHuggingFaceClientFromEnv())
	runRegressionTestWithClient(dir, "gemini", llm.NewGeminiClientFromEnv())
}

func runRegressionTestWithClient(dir string, name string, client llm.Client) {
	fmt.Println("testing " + name + " for case " + dir)
	prompt, err := ioutil.ReadFile(baseDir + dir + "/prompt.txt")
	if err != nil {
		panic(err)
	}

	resp, err := client.Query(string(prompt), nil)
	if err != nil {
		panic(err)
	}

	err = ioutil.WriteFile(baseDir+dir+"/"+name+".txt", []byte(resp), 0644)
	if err != nil {
		panic(err)
	}
}
