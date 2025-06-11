package main

import (
	"fmt"
	input "golang-ai-server/input"
	"golang-ai-server/logger"
	"golang-ai-server/ollama"
	"log/slog"
	"os"
	"time"

	"github.com/davecgh/go-spew/spew"
)

const Url = "http://localhost:11434/api/generate"

func main() {
	// Initialize logging with configuration from environment or defaults
	logConfig := logger.NewDefaultConfig()
	logger.InitLogger(logConfig)
	
	slog.Info("Starting golang-ai-server application", 
		"version", "1.0.0",
		"log_level", logConfig.Level,
		"log_format", logConfig.Format)
	slog.Debug("Configuration", "ollama_url", Url)
	
	// Get user input
	slog.Info("Requesting user input")
	prompt, images, err := input.GetUserInput()
	if err != nil {
		logger.LogError("get_user_input", err, map[string]interface{}{
			"operation": "input_collection",
		})
		fmt.Println(err)
		os.Exit(1)
	}
	logger.LogUserInput("prompt_and_images", len(prompt)+len(images))
	slog.Info("Successfully received user input", "prompt_length", len(prompt), "images_count", len(images))

	// Prepare Ollama request
	start := time.Now()
	logger.LogProcessingStep("prepare_ollama_request", map[string]interface{}{
		"start_time": start,
		"model": "granite3.2-vision:2b",
	})
	
	req := ollama.Request{
		Images: images,
		Model:  "granite3.2-vision:2b",
		Stream: false,
		Prompt: prompt,
	}
	slog.Info("Created Ollama request", 
		"model", req.Model, 
		"stream", req.Stream, 
		"prompt_length", len(req.Prompt),
		"images_count", len(req.Images))
	
	// Make request to Ollama
	slog.Info("Sending request to Ollama", "url", Url)
	resp, err := ollama.ReqOllama(Url, req)
	if err != nil {
		logger.LogError("ollama_request", err, map[string]interface{}{
			"url": Url,
			"duration": time.Since(start),
			"model": req.Model,
		})
		fmt.Println(err)
		os.Exit(1)
	}
	
	totalDuration := time.Since(start)
	slog.Info("Ollama request completed", "duration", totalDuration)
	
	// Process response
	if resp.Response != "" {
		logger.LogProcessingStep("process_response", map[string]interface{}{
			"response_length": len(resp.Response),
			"model": resp.Model,
			"done": resp.Done,
		})
		fmt.Println(resp.Response)
	} else {
		slog.Warn("Received empty response, dumping full response object", 
			"done", resp.Done,
			"done_reason", resp.DoneReason)
		spew.Dump(resp)
	}
	
	slog.Info("Application completed successfully", "total_duration", totalDuration)
	fmt.Printf("Completed in %v", totalDuration)
}
