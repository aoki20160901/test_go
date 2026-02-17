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
	pdf.SetAutoPageBreak(true, 20)

	// 日本語フォント
	pdf.AddUTF8Font("NotoSans", "", "./fonts/NotoSansJP-Regular.ttf")
	pdf.SetFont("NotoSans", "", 12)

	pageWidth, pageHeight := 210.0, 297.0
	bottomLimit := pageHeight - 20

	for i, imagePath := range imagePaths {

		pdf.AddPage()

		// =========================
		// ① 作成日（右上）
		// =========================
		pdf.SetXY(140, 15)
		pdf.CellFormat(
			50, 10,
			fmt.Sprintf("作成日: %s", time.Now().Format("2006-01-02")),
			"", 0, "R", false, 0, "",
		)

		pdf.SetY(30)

		// =========================
		// ② タイトル
		// =========================
		pdf.SetFont("NotoSans", "", 14)
		pdf.Cell(0, 10, "工事報告書")
		pdf.Ln(15)

		pdf.SetFont("NotoSans", "", 12)

		// =========================
		// ③ 画像
		// =========================
		imageWidth := 150.0
		x := (pageWidth - imageWidth) / 2.0

		startY := pdf.GetY()

		pdf.ImageOptions(
			imagePath,
			x,
			startY,
			imageWidth,
			0, // 高さ自動
			false,
			gofpdf.ImageOptions{ImageType: "", ReadDpi: true},
			0,
			"",
		)

		// 実際の画像高さ計算
		info := pdf.GetImageInfo(imagePath)
		ratio := info.Height() / info.Width()
		imageHeight := imageWidth * ratio

		nextY := startY + imageHeight + 10

		// 画像がページ下限を超える場合は縮小
		if nextY > bottomLimit {
			maxHeight := bottomLimit - startY - 10
			scale := maxHeight / imageHeight
			imageWidth = imageWidth * scale
			imageHeight = imageHeight * scale
			x = (pageWidth - imageWidth) / 2.0

			pdf.ImageOptions(
				imagePath,
				x,
				startY,
				imageWidth,
				imageHeight,
				false,
				gofpdf.ImageOptions{ImageType: "", ReadDpi: true},
				0,
				"",
			)

			nextY = startY + imageHeight + 10
		}

		pdf.SetY(nextY)

		// =========================
		// ④ キャプション
		// =========================
		if pdf.GetY() > bottomLimit {
			pdf.AddPage()
		}

		pdf.MultiCell(0, 8, captions[i], "", "L", false)
	}

	var buf bytes.Buffer
	err := pdf.Output(&buf)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
