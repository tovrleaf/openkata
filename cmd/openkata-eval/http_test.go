package main

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHTTPCompleterComplete(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("method = %s, want POST", r.Method)
		}
		if r.URL.Path != "/chat/completions" {
			t.Errorf("path = %s, want /chat/completions", r.URL.Path)
		}
		if r.Header.Get("Authorization") != "Bearer test-key" {
			t.Errorf("auth = %q, want %q", r.Header.Get("Authorization"), "Bearer test-key")
		}

		body, _ := io.ReadAll(r.Body)
		var req chatRequest
		json.Unmarshal(body, &req)

		if req.Model != "gpt-4o" {
			t.Errorf("model = %q, want %q", req.Model, "gpt-4o")
		}
		if len(req.Messages) != 2 {
			t.Errorf("messages count = %d, want 2", len(req.Messages))
		}
		if req.Messages[0].Role != "system" {
			t.Errorf("messages[0].role = %q, want %q", req.Messages[0].Role, "system")
		}

		resp := chatResponse{
			Choices: []struct {
				Message struct {
					Content string `json:"content"`
				} `json:"message"`
			}{
				{Message: struct {
					Content string `json:"content"`
				}{Content: "test response"}},
			},
		}
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	c := NewHTTPCompleter(server.URL, "test-key", "gpt-4o", 0)
	got, err := c.Complete("system prompt", "user prompt")
	if err != nil {
		t.Fatalf("Complete() error: %v", err)
	}
	if got != "test response" {
		t.Errorf("Complete() = %q, want %q", got, "test response")
	}
}

func TestHTTPCompleterNoSystem(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		var req chatRequest
		json.Unmarshal(body, &req)

		if len(req.Messages) != 1 {
			t.Errorf("messages count = %d, want 1", len(req.Messages))
		}
		if req.Messages[0].Role != "user" {
			t.Errorf("messages[0].role = %q, want %q", req.Messages[0].Role, "user")
		}

		resp := chatResponse{
			Choices: []struct {
				Message struct {
					Content string `json:"content"`
				} `json:"message"`
			}{
				{Message: struct {
					Content string `json:"content"`
				}{Content: "ok"}},
			},
		}
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	c := NewHTTPCompleter(server.URL, "", "model", 0)
	got, err := c.Complete("", "hello")
	if err != nil {
		t.Fatalf("Complete() error: %v", err)
	}
	if got != "ok" {
		t.Errorf("Complete() = %q, want %q", got, "ok")
	}
}

func TestHTTPCompleterError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusTooManyRequests)
		w.Write([]byte("rate limited"))
	}))
	defer server.Close()

	c := NewHTTPCompleter(server.URL, "key", "model", 0)
	_, err := c.Complete("sys", "user")
	if err == nil {
		t.Fatal("Complete() expected error, got nil")
	}
	if got := err.Error(); got != "HTTP 429: rate limited" {
		t.Errorf("error = %q, want %q", got, "HTTP 429: rate limited")
	}
}
