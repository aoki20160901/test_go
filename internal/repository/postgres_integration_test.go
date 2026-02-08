package repository_test

import (
	"database/sql"
	"os"
	"testing"

	_ "github.com/lib/pq"
)

func dsnFromEnv() string {
	if s := os.Getenv("DATABASE_URL"); s != "" {
		return s
	}
	return "postgres://pguser:pgpass@localhost:5432/go_test2?sslmode=disable"
}

func TestPostgresCRUD(t *testing.T) {
	dsn := dsnFromEnv()
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		t.Fatalf("open db: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		t.Fatalf("ping db: %v", err)
	}

	// Ensure table exists
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS users (id TEXT PRIMARY KEY, name TEXT NOT NULL)`)
	if err != nil {
		t.Fatalf("create table: %v", err)
	}

	id := "integration-test-1"
	name := "IntegrationUser"

	// Upsert
	_, err = db.Exec(`INSERT INTO users (id, name) VALUES ($1, $2) ON CONFLICT (id) DO UPDATE SET name = EXCLUDED.name`, id, name)
	if err != nil {
		t.Fatalf("insert: %v", err)
	}

	var got string
	if err := db.QueryRow(`SELECT name FROM users WHERE id = $1`, id).Scan(&got); err != nil {
		t.Fatalf("select: %v", err)
	}
	if got != name {
		t.Fatalf("expected name %q, got %q", name, got)
	}

	// cleanup
	_, _ = db.Exec(`DELETE FROM users WHERE id = $1`, id)
}
