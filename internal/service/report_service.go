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
	pdf.SetAutoPageBreak(false, 20)

	pdf.AddPage()

	pdf.AddUTF8Font("NotoSans", "", "./fonts/NotoSansJP-Regular.ttf")
	pdf.SetFont("NotoSans", "", 12)

	pageWidth, pageHeight := 210.0, 297.0
	leftMargin := 20.0
	rightMargin := 20.0
	bottomLimit := pageHeight - 20

	contentWidth := pageWidth - leftMargin - rightMargin
	columnWidth := contentWidth/2 - 5

	// ===== ヘッダー =====
	pdf.SetXY(140, 15)
	pdf.CellFormat(
		50, 10,
		fmt.Sprintf("作成日: %s", time.Now().Format("2006-01-02")),
		"", 0, "R", false, 0, "",
	)

	pdf.SetY(30)
	pdf.SetFont("NotoSans", "", 14)
	pdf.Cell(0, 10, "工事報告書")
	pdf.Ln(15)
	pdf.SetFont("NotoSans", "", 11)

	currentY := pdf.GetY()

	for i := 0; i < len(imagePaths); i += 2 {

		rowStartY := currentY
		maxRowHeight := 0.0

		for col := 0; col < 2; col++ {

			index := i + col
			if index >= len(imagePaths) {
				break
			}

			x := leftMargin + float64(col)*(columnWidth+10)
			y := rowStartY

			// 画像サイズ計算
			imageWidth := columnWidth
			info := pdf.GetImageInfo(imagePaths[index])
			ratio := info.Height() / info.Width()
			imageHeight := imageWidth * ratio

			// ページ下限チェック
			if y+imageHeight+30 > bottomLimit {
				pdf.AddPage()
				currentY = 30
				rowStartY = currentY
				y = rowStartY
			}

			// 画像描画
			pdf.ImageOptions(
				imagePaths[index],
				x,
				y,
				imageWidth,
				imageHeight,
				false,
				gofpdf.ImageOptions{ImageType: "", ReadDpi: true},
				0,
				"",
			)

			// キャプション位置
			textY := y + imageHeight + 3
			pdf.SetXY(x, textY)
			pdf.MultiCell(columnWidth, 6, captions[index], "", "L", false)

			// この列の総高さ
			colHeight := imageHeight + 3 + pdf.GetY() - textY
			if colHeight > maxRowHeight {
				maxRowHeight = colHeight
			}
		}

		currentY = rowStartY + maxRowHeight + 10
	}

	var buf bytes.Buffer
	err := pdf.Output(&buf)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
