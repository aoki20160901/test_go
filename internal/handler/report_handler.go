package handler

import (
	"io"
	"myapi/internal/service"
	"myapi/pkg/logger"
	"net/http"
)

type ReportHandler struct {
	service service.ReportService
}

func NewReportHandler(s service.ReportService) *ReportHandler {
	return &ReportHandler{service: s}
}

func (h *ReportHandler) Generate(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	logger.Info("レポート生成開始")

	// multipart解析
	err := r.ParseMultipartForm(10 << 20) // 10MB
	if err != nil {
		logger.Error("multipart解析エラー", "error", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	logger.Info("multipart解析成功")

	text := r.FormValue("text")
	logger.Info("テキスト取得", "text_length", len(text))

	file, _, err := r.FormFile("image")
	if err != nil {
		logger.Error("画像ファイル取得エラー", "error", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer file.Close()
	logger.Info("画像ファイル取得成功")

	// io.Reader に渡す
	pdfBytes, err := h.service.GenerateReport(ctx, text, file.(io.Reader))
	if err != nil {
		logger.Error("レポート生成エラー", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	logger.Info("レポート生成成功", "pdf_size", len(pdfBytes))

	w.Header().Set("Content-Type", "application/pdf")
	w.Write(pdfBytes)
	logger.Info("レポート生成完了")
}
