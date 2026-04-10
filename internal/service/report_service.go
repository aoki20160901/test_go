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
	entranceImages []string,
	entranceCaptions []string,
	hallwayImages []string,
	hallwayCaptions []string,
	outdoorImages []string,
	outdoorCaptions []string,
	bedroomImages []string,
	bedroomCaptions []string,
	livingImages []string,
	livingCaptions []string,
	toiletImages []string,
	toiletCaptions []string,
	bathroomImages []string,
	bathroomCaptions []string,
	summary string,
) ([]byte, error) {
	// --- ステップ1: 全エリアの画像・キャプションを1つの配列にまとめる ---
	var allImages []string
	var allCaptions []string
	allImages = append(allImages, entranceImages...)
	allCaptions = append(allCaptions, entranceCaptions...)
	allImages = append(allImages, hallwayImages...)
	allCaptions = append(allCaptions, hallwayCaptions...)
	allImages = append(allImages, outdoorImages...)
	allCaptions = append(allCaptions, outdoorCaptions...)
	allImages = append(allImages, bedroomImages...)
	allCaptions = append(allCaptions, bedroomCaptions...)
	allImages = append(allImages, livingImages...)
	allCaptions = append(allCaptions, livingCaptions...)
	allImages = append(allImages, toiletImages...)
	allCaptions = append(allCaptions, toiletCaptions...)
	allImages = append(allImages, bathroomImages...)
	allCaptions = append(allCaptions, bathroomCaptions...)

	if len(entranceImages) != len(entranceCaptions) {
		return nil, fmt.Errorf("玄関のimageとcaptionの数が一致しません")
	}
	if len(hallwayImages) != len(hallwayCaptions) {
		return nil, fmt.Errorf("廊下のimageとcaptionの数が一致しません")
	}
	if len(outdoorImages) != len(outdoorCaptions) {
		return nil, fmt.Errorf("屋外のimageとcaptionの数が一致しません")
	}
	pdf := gofpdf.New("L", "mm", "A4", "")
	pdf.SetMargins(20, 20, 20)
	pdf.SetAutoPageBreak(false, 20)

	var leftMargin, rightMargin, bottomLimit float64
	leftMargin = 20.0
	rightMargin = 20.0
	bottomLimit = 297.0 - 20

	// ===== ヘッダー =====
	pdf.AddPage()
	pdf.AddUTF8Font("NotoSans", "", "./fonts/NotoSansJP-Regular.ttf")
	pdf.SetFont("NotoSans", "", 12)
	pdf.SetXY(140, 15)
	pdf.CellFormat(
		50, 10,
		fmt.Sprintf("作成日: %s", time.Now().Format("2006-01-02")),
		"", 0, "R", false, 0, "",
	)
	pdf.SetY(30)
	pdf.SetFont("NotoSans", "", 14)
	pdf.CellFormat(0, 10, "テストユーザー様報告書", "", 0, "C", false, 0, "")
	pdf.Ln(15)
	pdf.SetFont("NotoSans", "", 11)

	// pageWidth, pageHeightは未使用のため削除

	// 1ページ目: タイトル・日付のみ
	pdf.SetY(30)
	pdf.SetFont("NotoSans", "", 14)
	pdf.CellFormat(0, 10, "テストユーザー様報告書", "", 0, "C", false, 0, "")
	pdf.Ln(15)
	pdf.SetFont("NotoSans", "", 11)

	// currentYは未使用のため削除

	// エリアリスト
	// areaListは未使用のため削除

	// --- ステップ2: 4枚ごとにページ分割し2×2で出力 ---
	for i := 0; i < len(allImages); i += 4 {
		end := i + 4
		if end > len(allImages) {
			end = len(allImages)
		}
		imgs := allImages[i:end]
		caps := allCaptions[i:end]
		if i != 0 {
			pdf.AddPage()
		}
		renderFourImagesOnePage(pdf, imgs, caps, leftMargin, rightMargin, bottomLimit)
	}

	// ===== 総評（summary）を最後に出力 =====
	// if strings.TrimSpace(summary) != "" {
	// 	// texts, captionsを合成
	// 	var allTexts []string
	// 	allTexts = append(allTexts, entranceCaptions...)
	// 	allTexts = append(allTexts, hallwayCaptions...)
	// 	var allCaptions []string
	// 	allCaptions = append(allCaptions, entranceCaptions...)
	// 	allCaptions = append(allCaptions, hallwayCaptions...)

	// 	// LLMで総評生成
	// 	souhyou, err := s.llm.GenerateSouhyou(ctx, allTexts, summary, allCaptions)
	// 	if err != nil {
	// 		souhyou = summary + "（総評生成エラー: " + err.Error() + ")"
	// 	}
	// 	pdf.AddPage()
	// 	pdf.SetFont("NotoSans", "", 14)
	// 	pdf.Cell(0, 10, "総評")
	// 	pdf.Ln(12)
	// 	pdf.SetFont("NotoSans", "", 11)
	// 	// 複数行対応
	//

	var buf bytes.Buffer
	var err error
	err = pdf.Output(&buf)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
	// 4枚の画像＋コメントを1ページに2×2でレイアウトする
}
func renderFourImagesOnePage(pdf *gofpdf.Fpdf, imagePaths []string, captions []string, leftMargin, rightMargin, bottomLimit float64) float64 {
	fmt.Println("[DEBUG] --- renderFourImagesOnePage START ---")
	fmt.Println("[DEBUG] imagePaths:", imagePaths)
	fmt.Println("[DEBUG] captions:", captions)
	// A4横: 297mm, 写真間余白10mm, 画像枠20%縮小
	pageWidth := 297.0
	gap := 10.0
	colWidth := ((pageWidth - 20.0 - 20.0 - gap) / 2) * 0.7
	rowHeight := (((bottomLimit-40)/2 - 5) * 0.7) * 0.8
	startY := 55.0 // ヘッダー・タイトル分の余白
	// 画像2枚＋gap分の幅を中央に揃えるための左端基準を計算
	totalContentWidth := colWidth*2 + gap
	leftX := (pageWidth - totalContentWidth) / 2

	for i := 0; i < 2; i++ {
		for j := 0; j < 2; j++ {
			idx := i*2 + j
			if idx >= len(imagePaths) {
				continue
			}
			// 画像アスペクト比取得
			ratio, err := getImageRatio(imagePaths[idx])
			if err != nil || ratio <= 0 {
				fmt.Printf("[DEBUG] getImageRatio error for %s: %v\n", imagePaths[idx], err)
				ratio = 0.75 // デフォルト比率
			}
			imgW := colWidth
			imgH := imgW * ratio
			if imgH > rowHeight-20 {
				imgH = rowHeight - 20
				imgW = imgH / ratio
			}
			// x座標をページ中央基準で中央揃えに調整
			x := leftX + float64(j)*(colWidth+gap) + (colWidth-imgW)/2
			y := startY + float64(i)*(rowHeight+10)
			fmt.Printf("[DEBUG] Drawing image idx=%d at x=%.2f y=%.2f\n", idx, x, y)

			// 画像描画
			pdf.ImageOptions(
				imagePaths[idx],
				x,
				y,
				imgW,
				imgH,
				false,
				gofpdf.ImageOptions{ImageType: "", ReadDpi: true},
				0,
				"",
			)

			// キャプション描画
			captionText := strings.TrimSpace(strings.TrimPrefix(captions[idx], "# 写真状況説明文"))
			pdf.SetFont("NotoSans", "", 12)
			lines := wrapText(captionText, 18)
			lineHeight := 6.0
			textY := y + imgH + 3
			for lineIdx, line := range lines {
				pdf.SetXY(x, textY+float64(lineIdx)*lineHeight)
				pdf.CellFormat(colWidth, lineHeight, line, "", 0, "L", false, 0, "")
				fmt.Printf("[DEBUG] Caption line %d: %s\n", lineIdx, line)
			}
		}
	}
	fmt.Println("[DEBUG] --- renderFourImagesOnePage END ---")
	return startY + 2*(rowHeight+10)
}
