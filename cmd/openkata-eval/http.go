package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// HTTPCompleter calls an OpenAI-compatible chat completions API.
type HTTPCompleter struct {
	BaseURL string
	APIKey  string
	Model   string
	Timeout time.Duration
}

// NewHTTPCompleter creates an HTTP-based completer.
func NewHTTPCompleter(baseURL, apiKey, model string, timeout time.Duration) *HTTPCompleter {
	if timeout == 0 {
		timeout = 120 * time.Second
	}
	return &HTTPCompleter{
		BaseURL: baseURL,
		APIKey:  apiKey,
		Model:   model,
		Timeout: timeout,
	}
}

type chatRequest struct {
	Model    string        `json:"model"`
	Messages []chatMessage `json:"messages"`
}

type chatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type chatResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

// Complete sends a chat completion request and returns the response content.
func (h *HTTPCompleter) Complete(system, user string) (string, error) {
	var messages []chatMessage
	if system != "" {
		messages = append(messages, chatMessage{Role: "system", Content: system})
	}
	messages = append(messages, chatMessage{Role: "user", Content: user})

	body, err := json.Marshal(chatRequest{Model: h.Model, Messages: messages})
	if err != nil {
		return "", err
	}

	url := strings.TrimRight(h.BaseURL, "/") + "/chat/completions"
	req, err := http.NewRequest("POST", url, bytes.NewReader(body))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	if h.APIKey != "" {
		req.Header.Set("Authorization", "Bearer "+h.APIKey)
	}

	client := &http.Client{Timeout: h.Timeout}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("HTTP %d: %s", resp.StatusCode, string(respBody))
	}

	var cr chatResponse
	if err := json.Unmarshal(respBody, &cr); err != nil {
		return "", fmt.Errorf("parsing response: %w", err)
	}
	if len(cr.Choices) == 0 {
		return "", fmt.Errorf("no choices in response")
	}
	return cr.Choices[0].Message.Content, nil
}
