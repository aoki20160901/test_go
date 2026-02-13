package llm

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type OpusClient struct {
	apiKey string
	client *http.Client
}

func NewOpusClient(apiKey string) *OpusClient {
	return &OpusClient{
		apiKey: apiKey,
		client: &http.Client{
			Timeout: 60 * time.Second,
		},
	}
}

type anthropicRequest struct {
	Model     string    `json:"model"`
	MaxTokens int       `json:"max_tokens"`
	Messages  []message `json:"messages"`
}

type message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type anthropicResponse struct {
	Content []struct {
		Text string `json:"text"`
	} `json:"content"`
}

func (c *OpusClient) RefineText(ctx context.Context, input string) (string, error) {

	reqBody := anthropicRequest{
		Model:     "claude-3-opus-20240229",
		MaxTokens: 500,
		Messages: []message{
			{
				Role:    "user",
				Content: input,
			},
		},
	}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		"https://api.anthropic.com/v1/messages",
		bytes.NewBuffer(jsonBody),
	)
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", c.apiKey)
	req.Header.Set("anthropic-version", "2023-06-01")

	resp, err := c.client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	// ❗ ステータスコードチェック（重要）
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("anthropic error: %s\n%s", resp.Status, string(body))
	}

	var parsed anthropicResponse
	if err := json.Unmarshal(body, &parsed); err != nil {
		return "", err
	}

	if len(parsed.Content) == 0 {
		return "", fmt.Errorf("empty response from anthropic")
	}

	return parsed.Content[0].Text, nil
}
