package llm

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
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
		BaseURL: "https://api.anthropic.com/v1",
		APIKey:  os.Getenv("OPUS_API_KEY"),
		Model:   "claude-3-haiku-20240307",
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

type anthropicRequest struct {
	Model     string    `json:"model"`
	MaxTokens int       `json:"max_tokens"`
	Messages  []message `json:"messages"`
}

type message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

func (c *OpusClient) GenerateCaption(ctx context.Context, text string) (string, error) {
	// ...本来のAPI呼び出し処理...
	return "ダミーキャプション", nil
}

// SplitSummaryByArea はsummaryをエリアごとに分割し、各エリア名をキー、内容を値としたmapで返す
func (c *OpusClient) SplitSummaryByArea(ctx context.Context, summary string) (map[string]string, error) {
	dummy := "{\"屋外\":\"屋外テキスト\",\"玄関\":\"玄関テキスト\",\"廊下\":\"廊下テキスト\",\"階段\":\"階段テキスト\",\"寝室\":\"寝室テキスト\",\"居室\":\"居室テキスト\",\"台所\":\"台所テキスト\",\"トイレ\":\"トイレテキスト\",\"浴室\":\"浴室テキスト\",\"脱衣所\":\"脱衣所テキスト\"}"
	var result map[string]string
	if err := json.Unmarshal([]byte(dummy), &result); err != nil {
		return nil, fmt.Errorf("JSON parse error: %v", err)
	}
	return result, nil
}

// GenerateSouhyou は、全写真の依頼内容と状況説明を踏まえて総評を1つ生成する。
func (c *OpusClient) GenerateSouhyou(ctx context.Context, texts []string, comment string, captions []string) (string, error) {
	var b strings.Builder
	b.WriteString("【工事依頼内容と各写真の状況説明】\n\n")
	for i := range captions {
		if i < len(texts) {
			b.WriteString(fmt.Sprintf("写真%d 依頼内容: %s\n", i+1, texts[i]))
		}
		b.WriteString(fmt.Sprintf("写真%d 状況説明: %s\n\n", i+1, captions[i]))
	}
	b.WriteString(fmt.Sprintf("全体の依頼内容: %s\n\n", comment))
	b.WriteString("\n上記を踏まえ、報告書の「総評」を200文字以内で作成してください。")
	b.WriteString("現地確認の結果をまとめる形で、断定調で簡潔に記載してください。")

	reqBody := anthropicRequest{
		Model:     c.Model,
		MaxTokens: 300,
		Messages: []message{
			{Role: "user", Content: b.String()},
		},
	}
	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.BaseURL+"/messages", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", c.APIKey)
	req.Header.Set("anthropic-version", "2023-06-01")
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	bodyBytes, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("Claude error: %s", string(bodyBytes))
	}
	var parsed struct {
		Content []struct {
			Text string `json:"text"`
		} `json:"content"`
	}
	if err := json.Unmarshal(bodyBytes, &parsed); err != nil {
		return "", err
	}
	if len(parsed.Content) == 0 {
		return "", fmt.Errorf("empty response from Claude")
	}
	return strings.TrimSpace(parsed.Content[0].Text), nil
}
