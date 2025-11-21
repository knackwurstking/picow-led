package services

import (
	"database/sql"
	"os"
	"strings"
	"testing"
)

const TestDatabaseName = "picow-led.test.db"

func openDB(t *testing.T, clean bool) *Registry {
	if clean {
		err := os.Remove(TestDatabaseName)
		if err != nil && !strings.Contains(err.Error(), "no such file or directory") {
			t.Fatalf("Failed to remove database file: %[1]s (%#[1]v)", err)
		}
	}

	db, err := sql.Open("sqlite3", TestDatabaseName)
	if err != nil {
		t.Fatalf("Failed to open database: %v", err)
	}

	registry, err := NewRegistry(db)
	if err != nil {
		t.Fatalf("Failed to create registry: %v", err)
	}

	return registry
}
