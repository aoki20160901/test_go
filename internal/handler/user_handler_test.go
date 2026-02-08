package handler_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"myapi/internal/router"
)

func TestHealth(t *testing.T) {
	r := router.Setup()
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}
	var body map[string]string
	if err := json.NewDecoder(w.Body).Decode(&body); err != nil {
		t.Fatalf("failed to decode body: %v", err)
	}
	if body["status"] != "ok" {
		t.Fatalf("expected status 'ok', got '%v'", body["status"])
	}
}

func TestCreateUser(t *testing.T) {
	r := router.Setup()
	payload := map[string]string{"name": "Hanako"}
	b, _ := json.Marshal(payload)
	req := httptest.NewRequest(http.MethodPost, "/v1/users", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected status 201, got %d, body: %s", w.Code, w.Body.String())
	}
	var got map[string]interface{}
	if err := json.NewDecoder(w.Body).Decode(&got); err != nil {
		t.Fatalf("failed to decode body: %v", err)
	}
	if got["id"] == "" || got["id"] == nil {
		t.Fatalf("expected non-empty id, got %v", got["id"])
	}
	if got["name"] != "Hanako" {
		t.Fatalf("expected name 'Hanako', got '%v'", got["name"])
	}
}

func TestGetUser(t *testing.T) {
	r := router.Setup()
	req := httptest.NewRequest(http.MethodGet, "/v1/users/123", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d, body: %s", w.Code, w.Body.String())
	}
	var got map[string]interface{}
	if err := json.NewDecoder(w.Body).Decode(&got); err != nil {
		t.Fatalf("failed to decode body: %v", err)
	}
	if got["id"] != "123" {
		t.Fatalf("expected id '123', got '%v'", got["id"])
	}
}
