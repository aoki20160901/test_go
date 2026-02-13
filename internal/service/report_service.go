package service

import (
	"context"
	"fmt"
	"io"

	"myapi/internal/llm"
	"myapi/internal/pdf"
)

type ReportService interface {
	GenerateReport(ctx context.Context, text string, image io.Reader) ([]byte, error)
}

type reportService struct {
	llm llm.LLM
	pdf pdf.Generator
}

func NewReportService(
	llm llm.LLM,
	pdf pdf.Generator,
) ReportService {
	return &reportService{
		llm: llm,
		pdf: pdf,
	}
}

func (s *reportService) GenerateReport(
	ctx context.Context,
	text string,
	image io.Reader,
) ([]byte, error) {

	fmt.Println("SERVICE CALLED")

	refined, err := s.llm.RefineText(ctx, text)
	if err != nil {
		return nil, err
	}

	fmt.Println("REFINED:", refined)

	// 今は仮でそのまま返す
	return []byte(refined), nil
}
