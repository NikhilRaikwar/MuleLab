package modelgateway

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestZeroCostBlocksPaidModelBeforeNetwork(t *testing.T) {
	g := OpenRouter{APIKey: "test", BaseURL: "http://unused", Model: "openai/gpt-paid", ZeroCostMode: true, Client: &http.Client{Timeout: time.Second}}
	_, err := g.Generate(context.Background(), Request{})
	if !errors.Is(err, ErrPaidBlocked) {
		t.Fatalf("expected paid block, got %v", err)
	}
}

func TestStructuredRequestAndResolvedModel(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != "Bearer test" {
			t.Fatal("missing server-side auth")
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"model":"free/resolved:free","choices":[{"message":{"content":"{\"agentType\":\"STORE\"}"}}],"usage":{"prompt_tokens":10,"completion_tokens":4,"cost":0}}`))
	}))
	defer server.Close()
	g := OpenRouter{APIKey: "test", BaseURL: server.URL, Model: "openrouter/free", ZeroCostMode: true, MaxOutputTokens: 1200, Client: &http.Client{Timeout: time.Second}}
	r, err := g.Generate(context.Background(), Request{SchemaName: "proposal", Schema: map[string]any{"type": "object"}})
	if err != nil {
		t.Fatal(err)
	}
	if r.ResolvedModel != "free/resolved:free" || r.CostUSD != 0 || r.PromptTokens != 10 {
		t.Fatalf("metadata not captured: %#v", r)
	}
}

func TestInvalidStructuredOutputNeverSucceeds(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(`{"model":"x:free","choices":[{"message":{"content":"not json"}}]}`))
	}))
	defer server.Close()
	g := OpenRouter{APIKey: "test", BaseURL: server.URL, Model: "x:free", ZeroCostMode: true, Client: &http.Client{Timeout: time.Second}}
	_, err := g.Generate(context.Background(), Request{})
	if !errors.Is(err, ErrInvalidJSON) {
		t.Fatalf("invalid output accepted: %v", err)
	}
}
