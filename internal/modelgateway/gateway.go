package modelgateway

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

var (
	ErrNoKey       = errors.New("openrouter key is not configured")
	ErrPaidBlocked = errors.New("paid model blocked by zero-cost policy")
	ErrInvalidJSON = errors.New("model response is not valid structured JSON")
)

type Request struct {
	System     string
	User       string
	SchemaName string
	Schema     map[string]any
}

type Result struct {
	JSON             json.RawMessage
	RequestedModel   string
	ResolvedModel    string
	Source           string
	Latency          time.Duration
	PromptTokens     int
	CompletionTokens int
	CostUSD          float64
	Retries          int
}

type Gateway interface {
	Generate(context.Context, Request) (Result, error)
}

type OpenRouter struct {
	APIKey, BaseURL, Model, FallbackModel, Referer, AppName string
	ZeroCostMode                                            bool
	MaxRetries                                              int
	MaxOutputTokens                                         int
	Client                                                  *http.Client
}

type chatRequest struct {
	Model          string              `json:"model"`
	Messages       []map[string]string `json:"messages"`
	MaxTokens      int                 `json:"max_tokens"`
	ResponseFormat map[string]any      `json:"response_format"`
	Provider       map[string]any      `json:"provider"`
}

type chatResponse struct {
	Model   string `json:"model"`
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
	Usage struct {
		PromptTokens     int     `json:"prompt_tokens"`
		CompletionTokens int     `json:"completion_tokens"`
		Cost             float64 `json:"cost"`
	} `json:"usage"`
}

func (g *OpenRouter) Generate(ctx context.Context, req Request) (Result, error) {
	if g.APIKey == "" {
		return Result{}, ErrNoKey
	}
	models := []string{g.Model}
	if g.FallbackModel != "" && g.FallbackModel != g.Model {
		models = append(models, g.FallbackModel)
	}
	var last error
	for modelIndex, model := range models {
		if err := enforceFree(model, g.ZeroCostMode); err != nil {
			return Result{}, err
		}
		attempts := 1 + g.MaxRetries
		for attempt := 0; attempt < attempts; attempt++ {
			result, err := g.call(ctx, model, req)
			result.Retries = attempt + modelIndex*attempts
			if err == nil {
				return result, nil
			}
			last = err
			if !retryable(err) || attempt+1 == attempts {
				break
			}
			select {
			case <-ctx.Done():
				return Result{}, ctx.Err()
			case <-time.After(time.Duration(50*(1<<attempt)) * time.Millisecond):
			}
		}
	}
	return Result{}, fmt.Errorf("model and fallback exhausted: %w", last)
}

func (g *OpenRouter) call(ctx context.Context, model string, req Request) (Result, error) {
	payload := chatRequest{Model: model, Messages: []map[string]string{{"role": "system", "content": req.System}, {"role": "user", "content": req.User}}, MaxTokens: g.MaxOutputTokens, ResponseFormat: map[string]any{"type": "json_schema", "json_schema": map[string]any{"name": req.SchemaName, "strict": true, "schema": req.Schema}}, Provider: map[string]any{"require_parameters": true}}
	body, err := json.Marshal(payload)
	if err != nil {
		return Result{}, err
	}
	url := strings.TrimRight(g.BaseURL, "/") + "/chat/completions"
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return Result{}, err
	}
	httpReq.Header.Set("Authorization", "Bearer "+g.APIKey)
	httpReq.Header.Set("Content-Type", "application/json")
	if g.Referer != "" {
		httpReq.Header.Set("HTTP-Referer", g.Referer)
	}
	if g.AppName != "" {
		httpReq.Header.Set("X-Title", g.AppName)
	}
	start := time.Now()
	resp, err := g.Client.Do(httpReq)
	latency := time.Since(start)
	if err != nil {
		return Result{}, err
	}
	defer resp.Body.Close()
	raw, err := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	if err != nil {
		return Result{}, err
	}
	if resp.StatusCode == http.StatusTooManyRequests || resp.StatusCode >= 500 {
		return Result{}, temporaryError{fmt.Errorf("openrouter status %d", resp.StatusCode)}
	}
	if resp.StatusCode >= 400 {
		return Result{}, fmt.Errorf("openrouter status %d", resp.StatusCode)
	}
	var decoded chatResponse
	if err := json.Unmarshal(raw, &decoded); err != nil {
		return Result{}, err
	}
	if len(decoded.Choices) == 0 {
		return Result{}, errors.New("openrouter returned no choices")
	}
	content := json.RawMessage(decoded.Choices[0].Message.Content)
	if !json.Valid(content) {
		return Result{}, ErrInvalidJSON
	}
	return Result{JSON: content, RequestedModel: model, ResolvedModel: decoded.Model, Source: "LIVE", Latency: latency, PromptTokens: decoded.Usage.PromptTokens, CompletionTokens: decoded.Usage.CompletionTokens, CostUSD: decoded.Usage.Cost}, nil
}

func enforceFree(model string, enabled bool) error {
	if !enabled {
		return nil
	}
	if model == "openrouter/free" || strings.HasSuffix(model, ":free") || strings.HasPrefix(model, "deterministic/") {
		return nil
	}
	return ErrPaidBlocked
}

type temporaryError struct{ error }

func retryable(err error) bool { var temp temporaryError; return errors.As(err, &temp) }

type Deterministic struct{ Output json.RawMessage }

func (d Deterministic) Generate(_ context.Context, _ Request) (Result, error) {
	if !json.Valid(d.Output) {
		return Result{}, ErrInvalidJSON
	}
	return Result{JSON: d.Output, RequestedModel: "deterministic/demo-v1", ResolvedModel: "deterministic/demo-v1", Source: "DETERMINISTIC_FALLBACK"}, nil
}
