package service

import (
	"bytes"
	"context"
	"fmt"
	"time"

	"your_project/llm"

	"github.com/jung-kurt/gofpdf"
)

type ReportService struct {
	llm *llm.OpusClient
}

func NewReportService(llmClient *llm.OpusClient) *ReportService {
	return &ReportService{
		llm: llmClient,
	}
}

func (s *ReportService) GenerateCaption(ctx context.Context, text string) (string, error) {
	return s.llm.GenerateCaption(ctx, text)
}

func (s *ReportService) GeneratePDF(
	ctx context.Context,
	imagePaths []string,
	captions []string,
) ([]byte, error) {

	if len(imagePaths) != len(captions) {
		return nil, fmt.Errorf("imageとcaptionの数が一致しません")
	}

	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.SetMargins(20, 20, 20)

	// 日本語フォント（事前にttf配置）
	pdf.AddUTF8Font("NotoSans", "", "./fonts/NotoSansJP-Regular.ttf")
	pdf.SetFont("NotoSans", "", 12)

	for i, imagePath := range imagePaths {

		pdf.AddPage()

		// =========================
		// ① 右上：作成日
		// =========================
		pdf.SetXY(140, 15)
		pdf.CellFormat(50, 10,
			fmt.Sprintf("作成日: %s", time.Now().Format("2006-01-02")),
			"", 0, "R", false, 0, "")

		pdf.SetY(30)

		// =========================
		// ② タイトル
		// =========================
		pdf.SetFont("NotoSans", "", 14)
		pdf.Cell(0, 10, "工事報告書")
		pdf.Ln(15)

		// =========================
		// ③ 画像
		// =========================
		imageWidth := 150.0
		x := (210.0 - imageWidth) / 2.0

		pdf.ImageOptions(
			imagePath,
			x,
			pdf.GetY(),
			imageWidth,
			0,
			false,
			gofpdf.ImageOptions{ImageType: "", ReadDpi: true},
			0,
			"",
		)

		pdf.Ln(95)

		// =========================
		// ④ キャプション（画像下）
		// =========================
		pdf.SetFont("NotoSans", "", 12)
		pdf.MultiCell(0, 8, captions[i], "", "L", false)
	}

	var buf bytes.Buffer
	err := pdf.Output(&buf)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
