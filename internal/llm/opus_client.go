package llm

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
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

	prompt := fmt.Sprintf(`
あなたは住宅改修および福祉用具選定の専門職です。
添付写真および入力テキストをもとに、
ケアマネジャー提出用の環境評価および改善提案文を作成してください。

【文章条件】
・専門職記録用の簡潔な文章（敬体）
・客観的表現
・50文字以内
・文章は1文で簡潔に作成する

【評価の観点】
以下の流れを参考に簡潔にまとめる
1. 住宅環境の状況
2. 動作への影響
3. 安全性（転倒等）
4. 改善提案

【重要ルール】
・入力テキストおよび写真から確認できる事実のみ記載する
・入力情報にない数値や寸法は記載しない
・入力内容を補完して新しい状況を作らない
・住宅環境の評価を優先して記述する

【用語統一】
「立ち上がり手すり」「据置手すり」「縦手すり」などの表現は
すべて「手すり」に置き換える。
出力では必ず「手すり」を使用する。

例
立ち上がり手すり → 手すり
縦手すり → 手すり
横手すり → 手すり

【利用者情報の扱い】
・利用者情報は動作能力の説明に必要な場合のみ使用する
・年齢、性別、独居などの属性情報は原則記載しない
・住宅環境評価に関係しない情報は記載しない

【フィラー処理】
入力文に含まれる「えー」「あのー」「その」「まあ」等の
フィラーは除去する

【環境特定ルール】
・評価場所は入力テキストを基準とする
・入力にない場所（浴室、トイレ等）は記載しない
・写真の環境認識が不確実な場合は入力テキストを優先する

【結論表現】
・改善提案は「〜が必要と判断される」で表現する
・「提案する」「望ましい」などの表現は使用しない

【文章形式】
以下の形式で文章を作成する

「〇〇リスクがある事から、〇〇設置を提案します。」

文末は必ず「〜設置を提案します。」で終える。
「提案」「提案と考えられる」「必要」「必要と判断される」「必要と考えられる」などは使用しない。

例
段差からの転倒リスクがある事から、手すりの設置を提案します。
滑りによる転倒リスクがある事から、手すりの設置を提案します。
動作時の不安定さから転倒リスクがある事から、手すりの設置を提案します。

依頼内容:
%s
`, text)

	reqBody := anthropicRequest{
		Model:     c.Model,
		MaxTokens: 300,
		Messages: []message{
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
		http.MethodPost,
		c.BaseURL+"/messages",
		bytes.NewBuffer(jsonData),
	)
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

	return parsed.Content[0].Text, nil
}

// SplitSummaryByArea はsummaryをエリアごとに分割し、各エリア名をキー、内容を値としたmapで返す
func (c *OpusClient) SplitSummaryByArea(ctx context.Context, summary string) (map[string]string, error) {
	// エリア名リスト
	// areas := []string{"屋外", "玄関", "廊下", "階段", "寝室", "居室", "台所", "トイレ", "浴室", "脱衣所"}
	result := make(map[string]string)
	areaNames := []string{"屋外", "玄関", "廊下", "階段", "寝室", "居室", "台所", "トイレ", "浴室", "脱衣所"}
	// area名でsplitし、各エリア名＋内容を再構成
	re := regexp.MustCompile("(" + strings.Join(areaNames, "|") + ")")
	parts := re.Split(summary, -1)
	indices := re.FindAllStringIndex(summary, -1)
	for i, idx := range indices {
		area := summary[idx[0]:idx[1]]
		var text string
		if i+1 < len(parts) {
			text = strings.TrimSpace(parts[i+1])
		} else {
			text = ""
		}
		result[area] = text
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
