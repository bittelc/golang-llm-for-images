package main

import (
	"fmt"
	input "golang-ai-server/input"
	"golang-ai-server/ollama"
	"os"
	"time"

	"github.com/davecgh/go-spew/spew"
)

const Url = "http://localhost:11434/api/generate"

func main() {
	prompt, images, err := input.GetUserInput()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	start := time.Now()
	req := ollama.Request{
		Images: images,
		Model:  "granite3.2-vision:2b",
		Stream: false,
		Prompt: prompt,
	}
	resp, err := ollama.ReqOllama(Url, req)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	if resp.Response != "" {
		fmt.Println(resp.Response)
	} else {
		spew.Dump(resp)
	}
	fmt.Printf("Completed in %v", time.Since(start))
}
