package chatgpt

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

const (
	// baseURL  = "https://api.openai.com/v1/"
	proxyURL           = "https://open.aiproxy.xyz/v1/"
	modelGPT3Dot5Turbo = "gpt-3.5-turbo"
	wfurl              = "https://api.gpt.wf/v3/completions"
)

/*
type chatkeymessage struct {
	Code           int     `json:"code"`
	Msg            string  `json:"msg"`
	TotalGranted   float64 `json:"total_granted"`
	TotalUsed      float64 `json:"total_used"`
	TotalAvailable float64 `json:"total_available"`
	EffectiveAt    int64   `json:"effective_at"`
	ExpiresAt      int64   `json:"expires_at"`
}*/

// chatGPTResponseBody 响应体
type chatGPTResponseBody struct {
	ID      string       `json:"id"`
	Object  string       `json:"object"`
	Created int          `json:"created"`
	Model   string       `json:"model"`
	Choices []chatChoice `json:"choices"`
	Usage   chatUsage    `json:"usage"`
}

// chatGPTRequestBody 请求体
type chatGPTRequestBody struct {
	Model       string        `json:"model,omitempty"` // gpt3.5-turbo
	Messages    []chatMessage `json:"messages,omitempty"`
	Temperature float64       `json:"temperature,omitempty"`
	N           int           `json:"n,omitempty"`
	MaxTokens   int           `json:"max_tokens,omitempty"`
}

// chatMessage 消息
type chatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type chatChoice struct {
	Index        int `json:"index"`
	Message      chatMessage
	FinishReason string `json:"finish_reason"`
}

type chatUsage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

type wfmess struct {
	Keys   string `json:"keys"`
	Prompt string `json:"prompt"`
}

var client = &http.Client{
	Transport: &http.Transport{
		Proxy: http.ProxyFromEnvironment,
	},
	Timeout: time.Minute * 5,
}

// completions gtp3.5文本模型回复
// curl https://api.openai.com/v1/chat/completions
// -H "Content-Type: application/json"
// -H "Authorization: Bearer YOUR_API_KEY"
// -d '{ "model": "gpt-3.5-turbo",  "messages": [{"role": "user", "content": "Hello!"}]}'
func completions(messages []chatMessage, apiKey string) (*chatGPTResponseBody, error) {
	com := chatGPTRequestBody{
		Messages: messages,
	}
	// default model
	if com.Model == "" {
		com.Model = modelGPT3Dot5Turbo
	}

	body, err := json.Marshal(com)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, proxyURL+"chat/completions", bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", "application/json; charset=utf-8")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", apiKey))
	req.Header.Set("Content-Type", "application/json; charset=utf-8")

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode < http.StatusOK || res.StatusCode >= http.StatusBadRequest {
		return nil, errors.New("response error: " + strconv.Itoa(res.StatusCode))
	}

	v := new(chatGPTResponseBody)
	if err = json.NewDecoder(res.Body).Decode(&v); err != nil {
		return nil, err
	}
	return v, nil
}

func completionsWF(message, key string) (string, error) {
	body, err := json.Marshal(wfmess{
		Prompt: message,
		Keys:   key,
	})
	if err != nil {
		return "", err
	}

	payload := bytes.NewReader(body)
	client := &http.Client{
		Timeout: time.Minute * 5,
	}
	req, err := http.NewRequest(http.MethodPost, wfurl, payload)

	if err != nil {
		return "", err
	}
	req.Header.Add("User-Agent", "Apifox/1.0.0 (https://www.apifox.cn)")
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	body, err = ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}
	return string(body), nil
}
