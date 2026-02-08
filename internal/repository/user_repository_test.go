package repository_test

import (
	"context"
	"testing"

	"myapi/internal/model"
)

// -----------------------------
// mock repository for unit test
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
// Unit test with mock
// -----------------------------
func TestUserRepository_Mock(t *testing.T) {
	ctx := context.Background()
	repo := &mockUserRepository{}

	// Create
	u := &model.User{Name: "Taro"}
	created, err := repo.Create(ctx, u)
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}
	if created.ID == "" {
		t.Fatalf("expected non-empty ID")
	}

	// FindByID
	got, err := repo.FindByID(ctx, created.ID)
	if err != nil {
		t.Fatalf("FindByID failed: %v", err)
	}
	if got.Name != "Mock User" {
		t.Fatalf("expected name 'Mock User', got %s", got.Name)
	}
}
