package main

import (
	"fmt"
	input "golang-ai-server/input"
	"golang-ai-server/ollama"
	"os"
	"time"

	"github.com/davecgh/go-spew/spew"
)

const generateUrl = "http://localhost:11434/api/generate"
const chatUrl = "http://localhost:11434/api/chat"

func main() {
	prompt, images, err := input.GetUserInput()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	start := time.Now()
	msg := ollama.Message{
		Role:    "user",
		Content: prompt,
		Images:  images,
	}
	req := ollama.Request{
		Model:    "granite3.2-vision:2b",
		Stream:   false,
		Messages: []ollama.Message{msg},
	}
	resp, err := ollama.ReqOllama(chatUrl, req)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	if resp.Message.Content != "" {
		fmt.Println(resp.Message.Content)
	} else {
		spew.Dump(resp)
	}
	fmt.Printf("Completed in %v", time.Since(start))
}
