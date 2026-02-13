package service

import (
	"context"
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
	return []byte("ok"), nil
}
