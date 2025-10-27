package services

import (
	"database/sql"
	"os"
	"testing"
)

const TestDatabaseName = "picow-led.test.db"

var registry *Registry

func openDB(t *testing.T, clean bool) *sql.DB {
	if clean {
		if err := os.Remove(TestDatabaseName); err != nil {
			t.Fatalf("Failed to remove database file: %v", err)
		}
	}

	db, err := sql.Open("sqlite3", TestDatabaseName)
	if err != nil {
		t.Fatalf("Failed to open database: %v", err)
	}

	registry = NewRegistry(db)
	if err := registry.CreateTables(); err != nil {
		panic(err)
	}

	return db
}
