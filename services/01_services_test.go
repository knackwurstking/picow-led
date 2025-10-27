package services

import (
	"database/sql"
	"testing"
)

const TestDatabaseName = "picow-led.test.db"

var registry *Registry

func init() {
	// Create test database
	db, err := sql.Open("sqlite3", TestDatabaseName)
	if err != nil {
		panic(err)
	}

	registry = NewRegistry(db)
	if err := registry.CreateTables(); err != nil {
		panic(err)
	}
}

func openDB(t *testing.T) *sql.DB {
	db, err := sql.Open("sqlite3", TestDatabaseName)
	if err != nil {
		t.Fatalf("Failed to open database: %v", err)
	}
	return db
}
