package handler

import (
	"context"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

type ReportService interface {
	GenerateCaption(ctx context.Context, text string) (string, error)
	GeneratePDF(ctx context.Context,
		entranceImages []string, entranceCaptions []string,
		hallwayImages []string, hallwayCaptions []string,
		outdoorImages []string, outdoorCaptions []string,
		bedroomImages []string, bedroomCaptions []string,
		livingImages []string, livingCaptions []string,
		toiletImages []string, toiletCaptions []string,
		bathroomImages []string, bathroomCaptions []string,
		summary string) ([]byte, error)
}

type ReportHandler struct {
	service ReportService
}

func NewReportHandler(s ReportService) *ReportHandler {
	return &ReportHandler{service: s}
}

func (h *ReportHandler) GenerateReport(w http.ResponseWriter, r *http.Request) {
	// ...既存の変数宣言・初期化・フォーム値取得の後に移動...

	ctx := r.Context()

	// 最大20MB
	err := r.ParseMultipartForm(20 << 20)
	if err != nil {
		http.Error(w, "invalid multipart form", http.StatusBadRequest)
		return
	}

	entranceFiles := r.MultipartForm.File["entrance_images"]
	hallwayFiles := r.MultipartForm.File["hallway_images"]

	// 総評(summary)を取得
	summary := ""
	if arr, ok := r.MultipartForm.Value["summary"]; ok && len(arr) > 0 {
		summary = arr[0]
	}

	// ステップ1: summaryをAIでエリアごとに分割
	areaSummaries := map[string]string{}
	if summary != "" {
		if s, ok := h.service.(interface {
			SplitSummaryByAreaWithAI(ctx context.Context, summary string) (map[string]string, error)
		}); ok {
			m, err := s.SplitSummaryByAreaWithAI(ctx, summary)
			if err == nil {
				areaSummaries = m
			}
		}
	}

	// if len(entranceFiles) == 0 && len(hallwayFiles) == 0 {
	// 	http.Error(w, "玄関または廊下の画像が必要です", http.StatusBadRequest)
	// 	return
	// }

	// if len(entranceTexts) != len(entranceFiles) {
	// 	http.Error(w, "玄関のテキストと画像の数が一致しません", http.StatusBadRequest)
	// 	return
	// }

	// if len(hallwayTexts) != len(hallwayFiles) {
	// 	http.Error(w, "廊下のテキストと画像の数が一致しません", http.StatusBadRequest)
	// 	return
	// }

	// アップロード保存先
	uploadDir := "./uploads"
	os.MkdirAll(uploadDir, os.ModePerm)

	var entranceImages []string
	var entranceCaptions []string
	var hallwayImages []string
	var hallwayCaptions []string
	var outdoorImages []string
	var outdoorCaptions []string
	var bedroomImages []string
	var bedroomCaptions []string
	var livingImages []string
	var livingCaptions []string
	var toiletImages []string
	var toiletCaptions []string
	var bathroomImages []string
	var bathroomCaptions []string

	// =====================
	// 玄関の画像処理
	// =====================
	for _, fileHeader := range entranceFiles {
		// ① 画像保存
		src, err := fileHeader.Open()
		if err != nil {
			http.Error(w, "file open error", http.StatusInternalServerError)
			return
		}

		filename := time.Now().Format("20060102150405") + "_" + fileHeader.Filename
		savePath := filepath.Join(uploadDir, filename)

		dst, err := os.Create(savePath)
		if err != nil {
			src.Close()
			http.Error(w, "file save error", http.StatusInternalServerError)
			return
		}

		_, err = io.Copy(dst, src)
		src.Close()
		dst.Close()

		if err != nil {
			http.Error(w, "file copy error", http.StatusInternalServerError)
			return
		}

		entranceImages = append(entranceImages, savePath)

		// ② LLM説明生成（areaSummaries["玄関"]を使う）
		text := ""
		if v, ok := areaSummaries["玄関"]; ok {
			text = v
		}
		caption, err := h.service.GenerateCaption(ctx, text)
		if err != nil {
			http.Error(w, "LLM error: "+err.Error(), http.StatusInternalServerError)
			return
		}
		entranceCaptions = append(entranceCaptions, caption)
	}

	// =====================
	// 廊下の画像処理
	// =====================
	for _, fileHeader := range hallwayFiles {
		// ① 画像保存
		src, err := fileHeader.Open()
		if err != nil {
			http.Error(w, "file open error", http.StatusInternalServerError)
			return
		}

		filename := time.Now().Format("20060102150405") + "_" + fileHeader.Filename
		savePath := filepath.Join(uploadDir, filename)

		dst, err := os.Create(savePath)
		if err != nil {
			src.Close()
			http.Error(w, "file save error", http.StatusInternalServerError)
			return
		}

		_, err = io.Copy(dst, src)
		src.Close()
		dst.Close()

		if err != nil {
			http.Error(w, "file copy error", http.StatusInternalServerError)
			return
		}

		hallwayImages = append(hallwayImages, savePath)

		// ② LLM説明生成（areaSummaries["廊下"]を使う）
		text := ""
		if v, ok := areaSummaries["廊下"]; ok {
			text = v
		}
		caption, err := h.service.GenerateCaption(ctx, text)
		if err != nil {
			http.Error(w, "LLM error: "+err.Error(), http.StatusInternalServerError)
			return
		}
		hallwayCaptions = append(hallwayCaptions, caption)
	}

	// =====================
	// ③ PDF生成
	// =====================
	pdfBytes, err := h.service.GeneratePDF(
		ctx,
		entranceImages, entranceCaptions,
		hallwayImages, hallwayCaptions,
		outdoorImages, outdoorCaptions,
		bedroomImages, bedroomCaptions,
		livingImages, livingCaptions,
		toiletImages, toiletCaptions,
		bathroomImages, bathroomCaptions,
		summary,
	)
	if err != nil {
		http.Error(w, "PDF生成失敗: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// =====================
	// ④ PDF返却
	// =====================
	w.Header().Set("Content-Type", "application/pdf")
	w.Header().Set("Content-Disposition", "attachment; filename=report.pdf")
	w.WriteHeader(http.StatusOK)
	w.Write(pdfBytes)
}
