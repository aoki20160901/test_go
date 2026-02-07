package logger_test

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"myapi/pkg/logger"
)

func TestLogger_InfoJSON(t *testing.T) {
	var buf bytes.Buffer
	logger.InitWithWriter(&buf)

	logger.Info("user created", "id", 123, "name", "Alice")

	out := strings.TrimSpace(buf.String())
	if out == "" {
		t.Fatalf("expected log output, got empty")
	}

	var m map[string]interface{}
	if err := json.Unmarshal([]byte(out), &m); err != nil {
		t.Fatalf("failed to unmarshal json log: %v; raw=%s", err, out)
	}

	if m["level"] != "INFO" {
		t.Fatalf("expected level INFO, got %v", m["level"])
	}
	if m["msg"] != "user created" {
		t.Fatalf("expected msg 'user created', got %v", m["msg"])
	}
	if m["id"] == nil {
		t.Fatalf("expected id field, got %v", m)
	}
	if m["name"] != "Alice" {
		t.Fatalf("expected name Alice, got %v", m["name"])
	}
}
