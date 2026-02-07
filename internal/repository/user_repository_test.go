package repository_test

import (
	"context"
	"testing"

	"myapi/internal/model"
	"myapi/internal/repository"
)

func TestCreateUser_Nil(t *testing.T) {
	_, err := repository.CreateUser(context.Background(), nil)
	if err == nil {
		t.Fatalf("expected error when creating nil user, got nil")
	}
}

func TestCreateUser_Success(t *testing.T) {
	in := &model.User{Name: "Alice"}
	out, err := repository.CreateUser(context.Background(), in)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out.ID == "" {
		t.Fatalf("expected non-empty ID, got empty")
	}
	if out.Name != "Alice" {
		t.Fatalf("expected name 'Alice', got '%s'", out.Name)
	}
}

func TestFindUserByID(t *testing.T) {
	u, err := repository.FindUserByID(context.Background(), "42")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if u == nil {
		t.Fatalf("expected user, got nil")
	}
	if u.ID != "42" {
		t.Fatalf("expected id '42', got '%s'", u.ID)
	}
}
