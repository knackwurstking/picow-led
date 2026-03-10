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
			t.Fatalf("remove database file: %[1]s (%#[1]v)", err)
		}
	}

	db, err := sql.Open("sqlite3", TestDatabaseName)
	if err != nil {
		t.Fatalf("open database: %v", err)
	}

	registry, err := NewRegistry(db)
	if err != nil {
		t.Fatalf("create registry: %v", err)
	}

	return registry
}
