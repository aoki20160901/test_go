package handler

import (
	"io"
	"myapi/internal/service"
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

	// multipart解析
	err := r.ParseMultipartForm(10 << 20) // 10MB
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	text := r.FormValue("text")

	file, _, err := r.FormFile("image")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer file.Close()

	// io.Reader に渡す
	pdfBytes, err := h.service.GenerateReport(ctx, text, file.(io.Reader))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/pdf")
	w.Write(pdfBytes)
}
