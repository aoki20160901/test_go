package main

import (
	"log"
	"net/http"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"myapi/internal/model"
	"myapi/internal/repository"
	"myapi/internal/router"
	"myapi/internal/service"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println(".env not found (using OS env)")
	}

	// ① DB接続
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatal("DATABASE_URL is not set")
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	// ② マイグレーション（開発中のみ推奨）
	if err := db.AutoMigrate(&model.User{}); err != nil {
		log.Fatal(err)
	}

	// ③ repository / service 作成
	userRepo := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepo)

	// ④ router 作成（service 注入）
	r := router.Setup(userService)

	// ⑤ サーバ起動
	log.Println("server running on :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
