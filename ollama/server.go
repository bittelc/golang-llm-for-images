package ollama

import (
	"bytes"
	"encoding/json"
	"golang-ai-server/logger"
	"net/http"
	"time"
)

type Request struct {
	Model     string   `json:"model"`
	Prompt    string   `json:"prompt"`
	Stream    bool     `json:"stream"`
	Think     bool     `json:"think"`
	Images    []string `json:"images"`
	Content   string   `json:"content"`
	ToolCalls []string `json:"tool_calls"`
}

type Response struct {
	Model      string    `json:"model"`
	CreatedAt  time.Time `json:"created_at"`
	Response   string    `json:"response"`
	Done       bool      `json:"done"`
	DoneReason string    `json:"done_reason"`
}

func ReqOllama(url string, ollamaReq Request) (*Response, error) {
	logger.LogProcessingStep("start_ollama_request", map[string]interface{}{
		"url":           url,
		"model":         ollamaReq.Model,
		"prompt_length": len(ollamaReq.Prompt),
		"images_count":  len(ollamaReq.Images),
		"stream":        ollamaReq.Stream,
	})

	js, err := json.Marshal(&ollamaReq)
	if err != nil {
		logger.LogError("marshal_ollama_request", err, map[string]interface{}{
			"url":   url,
			"model": ollamaReq.Model,
		})
		return nil, err
	}
	logger.LogProcessingStep("marshal_request", map[string]interface{}{
		"json_size_bytes": len(js),
	})

	client := http.Client{}
	logger.LogProcessingStep("create_http_client", map[string]interface{}{})

	httpReq, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(js))
	if err != nil {
		logger.LogError("create_http_request", err, map[string]interface{}{
			"url":    url,
			"method": http.MethodPost,
		})
		return nil, err
	}
	logger.LogRequest(http.MethodPost, url, len(js), map[string]string{
		"Content-Type": "application/json",
	})

	// Set appropriate headers
	httpReq.Header.Set("Content-Type", "application/json")
	logger.LogProcessingStep("set_headers", map[string]interface{}{
		"content_type": "application/json",
	})

	// Perform the request
	start := time.Now()
	httpResp, err := client.Do(httpReq)
	requestDuration := time.Since(start)

	if err != nil {
		logger.LogError("http_request_failed", err, map[string]interface{}{
			"url":      url,
			"duration": requestDuration,
		})
		return nil, err
	}
	defer httpResp.Body.Close()

	// Convert headers to map[string]string for logging
	responseHeaders := make(map[string]string)
	for key, values := range httpResp.Header {
		if len(values) > 0 {
			responseHeaders[key] = values[0]
		}
	}

	logger.LogResponse(httpResp.StatusCode, int(httpResp.ContentLength), requestDuration.String(), responseHeaders)

	ollamaResp := Response{}
	err = json.NewDecoder(httpResp.Body).Decode(&ollamaResp)
	if err != nil {
		logger.LogError("decode_ollama_response", err, map[string]interface{}{
			"status_code": httpResp.StatusCode,
			"url":         url,
		})
		return nil, err
	}

	logger.LogProcessingStep("decode_response", map[string]interface{}{
		"model":           ollamaResp.Model,
		"response_length": len(ollamaResp.Response),
		"done":            ollamaResp.Done,
		"done_reason":     ollamaResp.DoneReason,
		"created_at":      ollamaResp.CreatedAt,
	})

	// Only dump the response body if there's an issue or if debug level is very high
	if ollamaResp.Response == "" && !ollamaResp.Done {
		logger.LogError("incomplete_response", nil, map[string]interface{}{
			"response_empty": ollamaResp.Response == "",
			"done":           ollamaResp.Done,
			"done_reason":    ollamaResp.DoneReason,
		})
	}

	return &ollamaResp, err
}
