package service_test

import (
	"context"
	"testing"

	"myapi/internal/service"
)

func TestGetUser_EmptyID(t *testing.T) {
	_, err := service.GetUser(context.Background(), "")
	if err == nil {
		t.Fatalf("expected error for empty id, got nil")
	}
}

func TestGetUser_Success(t *testing.T) {
	u, err := service.GetUser(context.Background(), "123")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if u == nil {
		t.Fatalf("expected user, got nil")
	}
	if u.ID != "123" {
		t.Fatalf("expected id '123', got '%s'", u.ID)
	}
}
