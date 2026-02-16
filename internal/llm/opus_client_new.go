package llm

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

type OpusClient struct {
	BaseURL    string
	APIKey     string
	Model      string
	HTTPClient *http.Client
}

func NewOpusClient() *OpusClient {
	return &OpusClient{
		BaseURL: "https://api.openai.com/v1", // Opus互換エンドポイントに変更可
		APIKey:  os.Getenv("OPUS_API_KEY"),
		Model:   "gpt-4o-mini", // 任意モデルに変更
		HTTPClient: &http.Client{
			Timeout: 60 * time.Second,
		},
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
		Message chatMessage `json:"message"`
	} `json:"choices"`
}

func (c *OpusClient) GenerateCaption(ctx context.Context, text string) (string, error) {

	prompt := fmt.Sprintf(`
あなたは社内提出用報告書作成AIです。

以下の工事依頼内容をもとに、
写真の説明文として100文字以内で生成してください。
断定調で簡潔に記載してください。

依頼内容:
%s
`, text)

	reqBody := chatRequest{
		Model: c.Model,
		Messages: []chatMessage{
			{
				Role:    "user",
				Content: prompt,
			},
		},
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequestWithContext(
		ctx,
		"POST",
		c.BaseURL+"/chat/completions",
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		return "", err
	}

	req.Header.Set("Authorization", "Bearer "+c.APIKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("LLM error: %s", string(bodyBytes))
	}

	var result chatResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	if len(result.Choices) == 0 {
		return "", fmt.Errorf("no response from LLM")
	}

	return result.Choices[0].Message.Content, nil
}
