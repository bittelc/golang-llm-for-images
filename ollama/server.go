package ollama

import (
	"bytes"
	"encoding/json"
	"net/http"
	"time"

	"github.com/davecgh/go-spew/spew"
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
	js, err := json.Marshal(&ollamaReq)
	if err != nil {
		return nil, err
	}
	client := http.Client{}
	httpReq, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(js))
	if err != nil {
		return nil, err
	}
	httpResp, err := client.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer httpResp.Body.Close()
	ollamaResp := Response{}
	err = json.NewDecoder(httpResp.Body).Decode(&ollamaResp)
	spew.Dump(httpResp.Body)
	return &ollamaResp, err
}
