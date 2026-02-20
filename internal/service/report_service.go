package service

import (
	"bytes"
	"context"
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"os"
	"strings"
	"time"

	"myapi/internal/llm"

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

// getImageRatio は画像ファイルの縦横比（高さ/幅）を返す。gofpdf に登録前でも使える。
func getImageRatio(path string) (float64, error) {
    f, err := os.Open(path)
    if err != nil {
        return 0, err
    }
    defer f.Close()
    img, _, err := image.Decode(f)
    if err != nil {
        return 0, err
    }
    b := img.Bounds()
    w := float64(b.Dx())
    if w == 0 {
        return 0, fmt.Errorf("image width is 0: %s", path)
    }
    return float64(b.Dy()) / w, nil
}

// wrapText は日本語テキストを指定文字数ごとに分割して文字列配列を返す
func wrapText(text string, maxChars int) []string {
    // まず改行コードで分割
    originalLines := strings.Split(text, "\n")
    var result []string
    
    // 各行について、長ければ折り返す
    for _, line := range originalLines {
        if line == "" {
            // 空行はそのまま保持
            result = append(result, "")
            continue
        }
        
        runes := []rune(line)
        if len(runes) <= maxChars {
            // 短い行はそのまま
            result = append(result, line)
        } else {
            // 長い行は指定文字数で折り返す
            for i := 0; i < len(runes); i += maxChars {
                end := i + maxChars
                if end > len(runes) {
                    end = len(runes)
                }
                result = append(result, string(runes[i:end]))
            }
        }
    }
    
    return result
}

func (s *ReportService) GeneratePDF(
    ctx context.Context,
    imagePaths []string,
    captions []string,
    souhyou string,
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
    pdf.Cell(0, 10, "XX報告書")
    pdf.Ln(15)
    pdf.SetFont("NotoSans", "", 11)

    currentY := pdf.GetY()

    for i := 0; i < len(imagePaths); i += 2 {

        rowStartY := currentY
        maxRowHeight := 0.0

        // 各列のキャプション行数を事前に計算
        var captionLinesLeft, captionLinesRight []string
        
        for col := 0; col < 2; col++ {
            index := i + col
            if index >= len(imagePaths) {
                break
            }
            captionText := strings.TrimSpace(strings.TrimPrefix(captions[index], "# 写真状況説明文"))
            lines := wrapText(captionText, 25)
            if col == 0 {
                captionLinesLeft = lines
            } else {
                captionLinesRight = lines
            }
        }

        for col := 0; col < 2; col++ {

            index := i + col
            if index >= len(imagePaths) {
                break
            }

            // 描画前に標準ライブラリでアスペクト比を取得（GetImageInfo は登録後のみ有効）
            ratio, err := getImageRatio(imagePaths[index])
            if err != nil {
                return nil, err
            }

            x := leftMargin + float64(col)*(columnWidth+10)
            y := rowStartY

            imageWidth := columnWidth
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

			// キャプション描画（Y座標を完全に手動制御）
			textY := y + imageHeight + 3
			pdf.SetFont("NotoSans", "", 9)
			lineHeight := 4.5
			
			var lines []string
			if col == 0 {
				lines = captionLinesLeft
			} else {
				lines = captionLinesRight
			}
			
			// 1行ずつ描画（Y座標を明示的に指定）
			for lineIdx, line := range lines {
				currentLineY := textY + float64(lineIdx)*lineHeight
				pdf.SetXY(x, currentLineY)
				// Cellではなく、Text系メソッドで描画してY座標を変更しない
				pdf.CellFormat(columnWidth, lineHeight, line, "", 0, "L", false, 0, "")
			}
			
			pdf.SetFont("NotoSans", "", 11)
			
			// 列の高さ計算
			textHeight := float64(len(lines)) * lineHeight
			colHeight := imageHeight + 3 + textHeight
            if colHeight > maxRowHeight {
                maxRowHeight = colHeight
            }
        }

        currentY = rowStartY + maxRowHeight + 10
    }

    // ===== 総評（最後のページにのみ） =====
    if souhyou != "" {
        if pdf.GetY()+40 > bottomLimit {
            pdf.AddPage()
        }
        pdf.Ln(10)
        
        // 本文
        pdf.SetFont("NotoSans", "", 11)
        // 55文字ごとに分割して右端に余裕を持たせる
        souhyouLines := wrapText(souhyou, 40)
        lineHeight := 6.0
        startY := pdf.GetY()
        
        for lineIdx, line := range souhyouLines {
            currentLineY := startY + float64(lineIdx)*lineHeight
            pdf.SetXY(leftMargin, currentLineY)
            pdf.CellFormat(contentWidth, lineHeight, line, "", 0, "L", false, 0, "")
        }
    }

    var buf bytes.Buffer
    err := pdf.Output(&buf)
    if err != nil {
        return nil, err
    }

    return buf.Bytes(), nil
}

// GenerateSouhyou は全写真の依頼内容とキャプションを踏まえ、LLMで総評を1つ生成する。
func (s *ReportService) GenerateSouhyou(ctx context.Context, texts []string, comment string, captions []string) (string, error) {
    return s.llm.GenerateSouhyou(ctx, texts, comment, captions)
}
