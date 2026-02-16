package service

import (
	"context"
	"errors"
	"strings"

	"myapi/internal/llm"
)

const systemPrompt = `
あなたは住宅改修会社の社内報告書作成アシスタントです。
入力内容を基に、社内提出用の報告書文章を作成してください。

条件：
・200文字以内
・簡潔で丁寧な文章
・主観的表現は禁止
・箇条書きは禁止
・一文または二文でまとめる
・200文字を超える出力は禁止
・余計な説明は出力しない
`

type ReportService interface {
	GenerateReport(ctx context.Context, text string) (string, error)
}

type reportService struct {
	llm llm.Client
}

func NewReportService(llmClient llm.Client) ReportService {
	return &reportService{
		llm: llmClient,
	}
}

func (s *reportService) GenerateReport(
	ctx context.Context,
	text string,
) (string, error) {

	if strings.TrimSpace(text) == "" {
		return "", errors.New("input text is empty")
	}

	req := llm.Request{
		System: systemPrompt,
		User:   text,
	}

	result, err := s.llm.Generate(ctx, req)
	if err != nil {
		return "", err
	}

	result = cleanOutput(result)
	result = trimTo200(result)

	if result == "" {
		return "", errors.New("llm returned empty result")
	}

	return result, nil
}

// -------------------------
// private helper functions
// -------------------------

func trimTo200(s string) string {
	r := []rune(s)
	if len(r) > 200 {
		return string(r[:200])
	}
	return s
}

func cleanOutput(s string) string {
	s = strings.TrimSpace(s)
	s = strings.ReplaceAll(s, "\n", "")
	return s
}
