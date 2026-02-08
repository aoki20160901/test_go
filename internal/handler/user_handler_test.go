package handler_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"myapi/internal/model"
	"myapi/internal/repository"
	"myapi/internal/router"
	"myapi/internal/service"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// -----------------------------
// mock repository
// -----------------------------
type mockUserRepository struct{}

func (m *mockUserRepository) Create(ctx context.Context, u *model.User) (*model.User, error) {
	u.ID = "mock-id-123"
	return u, nil
}

func (m *mockUserRepository) FindByID(ctx context.Context, id string) (*model.User, error) {
	return &model.User{ID: "mock-id-123", Name: "Mock User"}, nil
}

// -----------------------------
// Test handler
// -----------------------------
func TestCreateUserHandler(t *testing.T) {
	// ctx := context.Background()
	repo := &mockUserRepository{}
	userService := service.NewUserService(repo)
	r := router.Setup(userService)

	// HTTP Request
	payload := map[string]string{"name": "Hanako"}
	b, _ := json.Marshal(payload)
	req := httptest.NewRequest(http.MethodPost, "/v1/users", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// Check status code
	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d, body: %s", w.Code, w.Body.String())
	}

	// Check response body
	var resp map[string]interface{}
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode body: %v", err)
	}
	if resp["id"] == "" {
		t.Fatalf("expected non-empty id")
	}
	if resp["name"] != "Hanako" {
		t.Fatalf("expected name 'Hanako', got '%v'", resp["name"])
	}
}

func setupDB(t *testing.T) *gorm.DB {
	t.Helper()

	dsn := "host=localhost user=testuser password=testpass dbname=testdb port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to connect to DB: %v", err)
	}

	if err := db.AutoMigrate(&model.User{}); err != nil {
		t.Fatalf("failed to migrate: %v", err)
	}

	return db
}

func TestCreateUserHandler_Integration(t *testing.T) {
	db := setupDB(t)
	repo := repository.NewUserRepository(db)
	userService := service.NewUserService(repo)
	r := router.Setup(userService) // router を service 経由で初期化

	// リクエストボディ
	payload := map[string]string{"name": "Hanako"}
	b, _ := json.Marshal(payload)
	req := httptest.NewRequest(http.MethodPost, "/v1/users", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// レスポンスをパース
	var resp map[string]interface{}
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	// 作成されたレコードをテスト終了時に削除
	t.Cleanup(func() {
		if id, ok := resp["id"].(string); ok {
			db.Delete(&model.User{}, "id = ?", id)
		}
	})

	// 簡単なチェック
	if w.Code != http.StatusCreated {
		t.Fatalf("expected status 201, got %d", w.Code)
	}
	if resp["name"] != "Hanako" {
		t.Fatalf("expected name 'Hanako', got '%v'", resp["name"])
	}
	if resp["id"] == "" {
		t.Fatalf("expected non-empty id")
	}
}
