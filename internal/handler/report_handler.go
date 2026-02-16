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
	GeneratePDF(ctx context.Context, imagePaths []string, captions []string) ([]byte, error)
}

type ReportHandler struct {
	service ReportService
}

func NewReportHandler(s ReportService) *ReportHandler {
	return &ReportHandler{service: s}
}

func (h *ReportHandler) GenerateReport(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// 最大20MB
	err := r.ParseMultipartForm(20 << 20)
	if err != nil {
		http.Error(w, "invalid multipart form", http.StatusBadRequest)
		return
	}

	texts := r.MultipartForm.Value["text"]
	files := r.MultipartForm.File["image"]

	if len(texts) == 0 || len(files) == 0 {
		http.Error(w, "textとimageは必須です", http.StatusBadRequest)
		return
	}

	if len(texts) != len(files) {
		http.Error(w, "textとimageの数が一致しません", http.StatusBadRequest)
		return
	}

	// アップロード保存先
	uploadDir := "./uploads"
	os.MkdirAll(uploadDir, os.ModePerm)

	var imagePaths []string
	var captions []string

	for i, fileHeader := range files {

		// =====================
		// ① 画像保存
		// =====================
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

		imagePaths = append(imagePaths, savePath)

		// =====================
		// ② LLM説明生成（Visionなし）
		// =====================
		caption, err := h.service.GenerateCaption(ctx, texts[i])
		if err != nil {
			http.Error(w, "LLM error: "+err.Error(), http.StatusInternalServerError)
			return
		}

		captions = append(captions, caption)
	}

	// =====================
	// ③ PDF生成
	// =====================
	pdfBytes, err := h.service.GeneratePDF(ctx, imagePaths, captions)
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
