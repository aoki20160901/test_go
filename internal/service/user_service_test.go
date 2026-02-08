package service_test

import (
	"context"
	"testing"

	"myapi/internal/model"
	"myapi/internal/repository"
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
	return &model.User{ID: id, Name: "Mock User"}, nil
}

// -----------------------------
// Test service
// -----------------------------
func TestUserService_CreateAndGet(t *testing.T) {
	ctx := context.Background()
	repo := &mockUserRepository{}
	svc := service.NewUserService(repo)

	// Create
	u := &model.User{Name: "Taro"}
	created, err := svc.CreateUser(ctx, u)
	if err != nil {
		t.Fatalf("CreateUser failed: %v", err)
	}
	if created.ID == "" {
		t.Fatalf("expected non-empty ID")
	}

	// GetUser
	got, err := svc.GetUser(ctx, created.ID)
	if err != nil {
		t.Fatalf("GetUser failed: %v", err)
	}
	if got.Name != "Mock User" {
		t.Fatalf("expected name 'Mock User', got %s", got.Name)
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

func TestUserService_Integration(t *testing.T) {
	ctx := context.Background()
	db := setupDB(t)
	repo := repository.NewUserRepository(db)
	svc := service.NewUserService(repo)

	// Create user
	u := &model.User{Name: "Taro"}
	created, err := svc.CreateUser(ctx, u)
	if err != nil {
		t.Fatalf("CreateUser failed: %v", err)
	}

	// テスト終了時に削除
	t.Cleanup(func() {
		db.Delete(&model.User{}, "id = ?", created.ID)
	})

	if created.ID == "" {
		t.Fatalf("expected non-empty ID")
	}

	// GetUser
	got, err := svc.GetUser(ctx, created.ID)
	if err != nil {
		t.Fatalf("GetUser failed: %v", err)
	}
	if got.Name != "Taro" {
		t.Fatalf("expected name 'Taro', got '%s'", got.Name)
	}
}
