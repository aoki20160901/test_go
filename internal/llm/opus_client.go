package llm

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"time"
)

type OpusClient struct {
	apiKey string
	client *http.Client
}

func NewOpusClient() *OpusClient {
	return &OpusClient{
		apiKey: os.Getenv("ANTHROPIC_API_KEY"),
		client: &http.Client{
			Timeout: 60 * time.Second,
		},
	}
}

func (c *OpusClient) Generate(ctx context.Context, req Request) (string, error) {
	if c.apiKey == "" {
		return "", errors.New("ANTHROPIC_API_KEY is not set")
	}

	body := map[string]any{
		"model":      "claude-3-opus-20240229",
		"max_tokens": 512,
		"system":     req.System,
		"messages": []map[string]string{
			{
				"role":    "user",
				"content": req.User,
			},
		},
	}

	jsonBody, err := json.Marshal(body)
	if err != nil {
		return "", err
	}

	httpReq, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		"https://api.anthropic.com/v1/messages",
		bytes.NewBuffer(jsonBody),
	)
	if err != nil {
		return "", err
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("x-api-key", c.apiKey)
	httpReq.Header.Set("anthropic-version", "2023-06-01")

	resp, err := c.client.Do(httpReq)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", errors.New("anthropic error: " + resp.Status)
	}

	var result struct {
		Content []struct {
			Text string `json:"text"`
		} `json:"content"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	if len(result.Content) == 0 {
		return "", errors.New("empty response from anthropic")
	}

	return result.Content[0].Text, nil
}
