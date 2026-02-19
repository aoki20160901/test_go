package main

import (
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"

	"myapi/internal/handler"
	"myapi/internal/llm"
	// "myapi/internal/pdf"
	"myapi/internal/service"
)

func main() {
	// ==========================
	// 1. 環境変数取得
	// ==========================
	apiKey := os.Getenv("ANTHROPIC_API_KEY")
	if apiKey == "" {
		log.Fatal("ANTHROPIC_API_KEY is required")
	}

	// ==========================
	// 2. Infrastructure生成
	// ==========================
	llmClient := llm.NewOpusClient()
	// pdfGenerator := pdf.NewTemplateGenerator("templates/report_template.pdf")

	// ==========================
	// 3. Service生成
	// ==========================
	reportService := service.NewReportService(llmClient)

	// ==========================
	// 4. Handler生成
	// ==========================
	reportHandler := handler.NewReportHandler(reportService)

	// ==========================
	// 5. Router設定
	// ==========================
	r := chi.NewRouter()

	r.Post("/report", reportHandler.GenerateReport)  // Generate → GenerateReport

	// ==========================
	// 6. Server起動
	// ==========================
	log.Println("Server started at :8080")
	err := http.ListenAndServe(":8080", r)
	if err != nil {
		log.Fatal(err)
	}
}
