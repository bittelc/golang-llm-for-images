package main

import (
	"bufio"
	"fmt"
	"golang-ai-server/ollama"
	"log"
	"os"
	"time"
)

const generateUrl = "http://localhost:11434/api/generate"
const chatUrl = "http://localhost:11434/api/chat"

func main() {
	var input string
	fmt.Print("User: ")
	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		log.Fatal(err)
	}
	start := time.Now()
	msg := ollama.Message{
		Role:    "user",
		Content: input,
	}
	req := ollama.Request{
		Model:    "llama3.2:latest",
		Stream:   false,
		Messages: []ollama.Message{msg},
	}
	// resp, err := ollama.ReqOllama(generateUrl, req)
	resp, err := ollama.ReqOllama(chatUrl, req)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println(resp.Message.Content)
	fmt.Printf("Completed in %v", time.Since(start))
}
